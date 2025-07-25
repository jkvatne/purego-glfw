package glfw

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows"
	"log/slog"
	"syscall"
	"unsafe"
)

var opengl32 *windows.LazyDLL

func swapBuffersWGL(window *_GLFWwindow) {
	if window.monitor != nil {
		if isWindowsVistaOrGreater() {
			/*
				// DWM Composition is always enabled on Win8+
				enabled := IsWindows8OrGreater()
				// HACK: Use DwmFlush when desktop composition is enabled
					if enabled || (SUCCEEDED(DwmIsCompositionEnabled(&enabled)) && enabled) {
						for count := abs(window.context.wgl.interval); count > 0; count-- {
							DwmFlush()
						}
					}
			*/
		}
	}
	r, _, err := _glfw.wgl.wglSwapBuffers.Call(uintptr(window.context.wgl.dc))
	if !errors.Is(err, syscall.Errno(0)) {
		panic(err)
	}
	if r == 0 {
		err = syscall.GetLastError()
		panic(err)
	}
}

func createContext(dc HDC) HANDLE {
	r1, _, err := _glfw.wgl.wglCreateContext.Call(uintptr(dc))
	if !errors.Is(err, syscall.Errno(0)) {
		// panic("Could not create context, " + err.Error())
	}
	return HANDLE(r1)
}

func deleteContext(handle HANDLE) {
	_, _, err := _glfw.wgl.wglDeleteContext.Call(uintptr(handle))
	if !errors.Is(err, syscall.Errno(0)) {
		panic(err)
	}
}

func getCurrentDC() HDC {
	r1, _, err := _glfw.wgl.wglGetCurrentDC.Call()
	if !errors.Is(err, syscall.Errno(0)) {
		panic("getCurrentDC failed, " + err.Error())
	}
	return HDC(r1)
}

func getCurrentContext() HANDLE {
	r1, _, err := _glfw.wgl.wglGetCurrentContext.Call()
	if !errors.Is(err, syscall.Errno(0)) {
		panic("getCurrentDC failed, " + err.Error())
	}
	return HANDLE(r1)
}

func makeCurrent(dc HDC, handle HANDLE) bool {
	r1, _, err := _glfw.wgl.wglMakeCurrent.Call(uintptr(dc), uintptr(handle))
	if !errors.Is(err, syscall.Errno(0)) {
		slog.Error("makeCurrent failed", "Error", err.Error(), "dc", dc, "hMonitor", handle)
	}
	if r1 == 0 {
		err = syscall.GetLastError()
		if err != nil {
			panic(err)
		}
	}
	return r1 != 0
}

func glfwMakeContextCurrent(window *_GLFWwindow) error {
	// _GLFWwindow* window = (_GLFWwindow*) hMonitor;
	// previous := _glfwPlatformGetTls(&_glfw.contextSlot);
	// if previous!=nil && w, r1indow!=nil || window.context.source != previous.context.source)
	//		previous.context.makeCurrent(NULL);
	// }

	if window != nil {
		window.context.makeCurrent(window)
	}
	return nil
}
func setPixelFormat(dc HDC, iPixelFormat int, pfd *PIXELFORMATDESCRIPTOR) int {
	ret, _, err := _glfw.wgl.wglSetPixelFormat.Call(uintptr(unsafe.Pointer(dc)), uintptr(iPixelFormat), uintptr(unsafe.Pointer(pfd)))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("wglSetPixelFormat failed, " + err.Error())
	}
	if ret == 0 {
		err = syscall.GetLastError()
		panic("wglSetPixelFormat failed" + err.Error())
	}
	return int(ret)
}

func choosePixelFormat(dc HDC, pfd *PIXELFORMATDESCRIPTOR) int {
	ret, _, err := _glfw.wgl.wglChoosePixelFormat.Call(uintptr(unsafe.Pointer(dc)), uintptr(unsafe.Pointer(pfd)))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("wglSetPixewglChoosePixelFormatlFormat failed, " + err.Error())
	}
	if ret == 0 {
		err = syscall.GetLastError()
		panic("wglChoosePixelFormat failed, " + err.Error())
	}
	return int(ret)
}

// Initialize WGL
func _glfwInitWGL() error {
	var pfd PIXELFORMATDESCRIPTOR
	if _glfw.wgl.instance != nil {
		return nil
	}

	opengl32 = windows.NewLazySystemDLL("opengl32")
	_glfw.wgl.instance = opengl32
	_glfw.wgl.wglSwapBuffers = gdi32.NewProc("SwapBuffers")
	_glfw.wgl.wglSetPixelFormat = gdi32.NewProc("SetPixelFormat")
	_glfw.wgl.wglChoosePixelFormat = gdi32.NewProc("ChoosePixelFormat")
	_glfw.wgl.wglDescribePixelFormat = gdi32.NewProc("DescribePixelFormat")
	_glfw.wgl.GetDeviceCaps = gdi32.NewProc("GetDeviceCaps")
	_glfw.wgl.instance = windows.NewLazySystemDLL("opengl32.dll")
	_glfw.wgl.getProcAddress = opengl32.NewProc("wglGetProcAddress")
	_glfw.wgl.wglCreateContext = opengl32.NewProc("wglCreateContext")
	_glfw.wgl.wglDeleteContext = opengl32.NewProc("wglDeleteContext")
	_glfw.wgl.wglGetProcAddress = opengl32.NewProc("wglGetProcAddress")
	_glfw.wgl.wglGetCurrentDC = opengl32.NewProc("wglGetCurrentDC")
	_glfw.wgl.wglGetCurrentContext = opengl32.NewProc("wglGetCurrentContext")
	_glfw.wgl.wglMakeCurrent = opengl32.NewProc("wglMakeCurrent")
	_glfw.wgl.GetString = opengl32.NewProc("wglGetString")
	// NOTE: A dummy context has to be created for opengl32.dll to load the
	//       OpenGL ICD, from which we can then query WGL extensions
	dc := getDC(_glfw.win32.helperWindowHandle)
	pfd.nSize = uint16(unsafe.Sizeof(pfd))
	pfd.nVersion = 1
	pfd.dwFlags = PFD_DRAW_TO_WINDOW | PFD_SUPPORT_OPENGL | PFD_DOUBLEBUFFER
	pfd.iPixelType = PFD_TYPE_RGBA
	pfd.cColorBits = 24
	setPixelFormat(dc, choosePixelFormat(dc, &pfd), &pfd)
	err := syscall.GetLastError()
	if err != nil {
		panic(err)
	}
	rc := createContext(dc)
	if rc == 0 {
		panic("WGL: Failed to create dummy context, " + err.Error())
	}
	pdc := getCurrentDC()
	prc := getCurrentContext()
	if !makeCurrent(dc, HANDLE(rc)) {
		slog.Error("WGL: Failed to make dummy context current")
	}
	/* TODO
	ret, _, _ := _glfw.wgl.wglMakeCurrent.Call(pdc, prc)
	if ret == 0 {
		_, _, _ = _glfw.wgl.wglMakeCurrent.Call(pdc, prc)
		_, _, _ = _glfw.wgl.wglDeleteContext.Call(rc)
		return fmt.Errorf("WGL: Failed to make dummy context current")
	}
	*/
	// NOTE: Functions must be loaded first as they're needed to retrieve the
	//       extension string that tells us whether the functions are supported
	// _glfw.wgl.GetExtensionsStringEXT = _glfw.wgl.wglGetProcAddress.Call("wglGetExtensionsStringEXT")
	// _glfw.wgl.GetExtensionsStringARB = (PFNWGLGETEXTENSIONSSTRINGARBPROC)wglGetProcAddress("wglGetExtensionsStringARB")
	// _glfw.wgl.wglCreateContextAttribsARB, _, _ = wglGetProcAddress.Call("wglCreateContextAttribsARB")
	_glfw.wgl.GetPixelFormatAttribivARB = opengl32.NewProc("wglGetPixelFormatAttribivARB")

	// NOTE: WGL_ARB_extensions_string and WGL_EXT_extensions_string are not
	//       checked below as we are already using them
	/* TODO
	_glfw.wgl.ARB_multisample =
		extensionSupportedWGL("WGL_ARB_multisample")
	_glfw.wgl.ARB_framebuffer_sRGB =
		extensionSupportedWGL("WGL_ARB_framebuffer_sRGB")
	_glfw.wgl.EXT_framebuffer_sRGB =
		extensionSupportedWGL("WGL_EXT_framebuffer_sRGB")
	_glfw.wgl.ARB_create_context = extensionSupportedWGL("WGL_ARB_create_context")
	_glfw.wgl.ARB_create_context_profile =
		extensionSupportedWGL("WGL_ARB_create_context_profile")
	_glfw.wgl.EXT_create_context_es2_profile =
		extensionSupportedWGL("WGL_EXT_create_context_es2_profile")
	_glfw.wgl.ARB_create_context_robustness =
		extensionSupportedWGL("WGL_ARB_create_context_robustness")
	_glfw.wgl.ARB_create_context_no_error =
		extensionSupportedWGL("WGL_ARB_create_context_no_error")
	_glfw.wgl.EXT_swap_control =
		extensionSupportedWGL("WGL_EXT_swap_control")
	_glfw.wgl.EXT_colorspace =
		extensionSupportedWGL("WGL_EXT_colorspace")
	_glfw.wgl.ARB_pixel_format =
		extensionSupportedWGL("WGL_ARB_pixel_format")
	_glfw.wgl.ARB_context_flush_control =
		extensionSupportedWGL("WGL_ARB_context_flush_control")
	*/
	makeCurrent(pdc, prc)
	deleteContext(HANDLE(rc))
	return nil
}

func getDC(w syscall.Handle) HDC {
	r1, _, err := _GetDC.Call(uintptr(w))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("getDC failed, " + err.Error())
	}
	return HDC(r1)
}

func releaseDC(w syscall.Handle, dc HDC) {
	_, _, err := _ReleaseDC.Call(uintptr(w), uintptr(dc))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("getDC failed, " + err.Error())
	}
}

func describePixelFormat(dc HDC, iPixelFormat int, nBytes int, ppfd *PIXELFORMATDESCRIPTOR) int {
	r1, _, err := _glfw.wgl.wglDescribePixelFormat.Call(uintptr(dc), uintptr(iPixelFormat), uintptr(nBytes), uintptr(unsafe.Pointer(ppfd)))
	if r1 == 0 || !errors.Is(err, syscall.Errno(0)) {
		slog.Error("describePixelFormat failed, " + err.Error())
		r1 = 0
	}
	return int(r1)
}

var attribs [40]int
var values [40]int
var index int

func SetAttrib(a int, v int) {
	attribs[index] = a
	index++
	attribs[index] = v
	index++
}

var attribCount int

func ADD_ATTRIB(a int) {
	attribs[attribCount] = a
	attribCount++
}

func FIND_ATTRIB_VALUE(a int) int {
	return findPixelFormatAttribValueWGL(attribs, attribCount, values, a)
}

func glfwCreateContextWGL(window *_GLFWwindow, ctxconfig *_GLFWctxconfig, fbconfig *_GLFWfbconfig) error {
	index = 0
	attribs[0] = 0
	attribs[1] = 0
	var pfd PIXELFORMATDESCRIPTOR
	window.context.wgl.dc = getDC(window.Win32.handle)
	if window.context.wgl.dc == 0 {
		return fmt.Errorf("WGL: Failed to retrieve DC for window")
	}
	pixelFormat := choosePixelFormatWGL(window, ctxconfig, fbconfig) // 14
	if pixelFormat == 0 {
		return fmt.Errorf("WGL: Failed to retrieve PixelFormat for window")
	}

	if describePixelFormat(window.context.wgl.dc, pixelFormat, int(unsafe.Sizeof(pfd)), &pfd) == 0 {
		return fmt.Errorf("WGL: Failed to retrieve PFD for selected pixel format")
	}
	if setPixelFormat(window.context.wgl.dc, pixelFormat, &pfd) == 0 {
		return fmt.Errorf("WGL: Failed to set selected pixel format")
	}

	/* TODO
	if ctxconfig.client == GLFW_OPENGL_API {
		if ctxconfig.forward && (_glfw.wgl.ARB_create_context == nil) {
			return fmt.Errorf("WGL: A forward compatible OpenGL context requested but WGL_ARB_create_context is unavailable")
		}
		if (ctxconfig.profile == 0) && _glfw.wgl.ARB_create_context_profile != nil {
			return fmt.Errorf("WGL: OpenGL profile requested but WGL_ARB_create_context_profile is unavailable")
		}
	} else {
		if !_glfw.wgl.ARB_create_context || !_glfw.wgl.ARB_create_context_profile || !_glfw.wgl.EXT_create_context_es2_profile {
			return fmt.Errorf("WGL: OpenGL ES requested but WGL_ARB_create_context_es2_profile is unavailable")
		}
	}
	*/

	// if _glfw.wgl.ARB_create_context != nil {
	/*
		mask := 0
		flags := 0;
		if (ctxconfig.client == GLFW_OPENGL_API) {
			if (ctxconfig.forward) {
				flags |= WGL_CONTEXT_FORWARD_COMPATIBLE_BIT_ARB;
			}
			if (ctxconfig.profile == GLFW_OPENGL_CORE_PROFILE) {
				mask |= WGL_CONTEXT_CORE_PROFILE_BIT_ARB;
			} else if (ctxconfig.profile == GLFW_OPENGL_COMPAT_PROFILE) {
				mask |= WGL_CONTEXT_COMPATIBILITY_PROFILE_BIT_ARB;
			}
		} else {
			mask |= WGL_CONTEXT_ES2_PROFILE_BIT_EXT;
		}
		if ctxconfig.debug {
			flags |= WGL_CONTEXT_DEBUG_BIT_ARB;
		}
		if ctxconfig.robustness!=0 {
			if _glfw.wgl.ARB_create_context_robustness {
				if (ctxconfig.robustness == GLFW_NO_RESET_NOTIFICATION) {
					setAttrib(WGL_CONTEXT_RESET_NOTIFICATION_STRATEGY_ARB, WGL_NO_RESET_NOTIFICATION_ARB);
				}
			} else if (ctxconfig.robustness == GLFW_LOSE_CONTEXT_ON_RESET) {
				setAttrib(WGL_CONTEXT_RESET_NOTIFICATION_STRATEGY_ARB, WGL_LOSE_CONTEXT_ON_RESET_ARB);
			}
			flags |= WGL_CONTEXT_ROBUST_ACCESS_BIT_ARB;
		}

		if ctxconfig.release!=0 {
			if (_glfw.wgl.ARB_context_flush_control) {
				if (ctxconfig.release == GLFW_RELEASE_BEHAVIOR_NONE) {
					setAttrib(WGL_CONTEXT_RELEASE_BEHAVIOR_ARB, WGL_CONTEXT_RELEASE_BEHAVIOR_NONE_ARB);
				} else if ctxconfig.release == GLFW_RELEASE_BEHAVIOR_FLUSH {
					setAttrib(WGL_CONTEXT_RELEASE_BEHAVIOR_ARB,
						WGL_CONTEXT_RELEASE_BEHAVIOR_FLUSH_ARB);
				}
			}
		}
		if (ctxconfig.noerror) {
			if (_glfw.wgl.ARB_create_context_no_error) {
				setAttrib(WGL_CONTEXT_OPENGL_NO_ERROR_ARB, GLFW_TRUE);
			}
		}
	*/
	// NOTE: Only request an explicitly versioned context when necessary, as
	//       explicitly requesting version 1.0 does not always return the
	//       highest version supported by the driver
	/*
				if ctxconfig.major != 1 || ctxconfig.minor != 0 {
					setAttrib(WGL_CONTEXT_MAJOR_VERSION_ARB, ctxconfig.major);
					setAttrib(WGL_CONTEXT_MINOR_VERSION_ARB, ctxconfig.minor);
				}
				if (flags!=0) {
					setAttrib(WGL_CONTEXT_FLAGS_ARB, flags)
				}
				if (mask!=0) {
					setAttrib(WGL_CONTEXT_PROFILE_MASK_ARB, mask)
				}
			// setAttrib(0, 0);
			window.context.wgl.hMonitor = wglCreateContextAttribsARB(window.context.wgl.dc, share, &attribs)
			if window.context.wgl.hMonitor == 0 {
				return fmt.Errorf("WGL: Driver does not support OpenGL version %d.%d", ctxconfig.major, ctxconfig.minor)
			}
		} else {
	*/
	window.context.wgl.handle = createContext(window.context.wgl.dc)
	if window.context.wgl.handle == 0 {
		return fmt.Errorf("WGL: Failed to create OpenGL context")
	}
	window.context.makeCurrent = makeContextCurrentWGL
	window.context.swapBuffers = swapBuffersWGL
	window.context.destroy = destroyContextWGL
	return nil
}

func findPixelFormatAttribValueWGL(attribs [40]int, attribCount int, values [40]int, attrib int) int {
	for i := 0; i < attribCount; i++ {
		if attribs[i] == attrib {
			return values[i]
		}
	}
	panic("WGL: Unknown pixel format attribute requested")
	return 0
}

func wglGetPixelFormatAttribivARB(dc HDC, pixelFormat int, layerPlane int, nAttrib int, attributes *int, piValues *int) int {
	r1, _, err := _glfw.wgl.GetPixelFormatAttribivARB.Call(uintptr(dc), uintptr(pixelFormat),
		uintptr(layerPlane), uintptr(nAttrib), uintptr(unsafe.Pointer(attributes)), uintptr(unsafe.Pointer(piValues)))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("WGL: GetPixelFormatAttribivARB failed, " + err.Error())
	}
	if r1 == 0 {
		panic("WGL: GetPixelFormatAttribivARB failed")
	}
	return int(r1)
}

func choosePixelFormatWGL(window *_GLFWwindow, ctxconfig *_GLFWctxconfig, fbconfig *_GLFWfbconfig) int {
	var (
		closest                               *_GLFWfbconfig
		pixelFormat, nativeCount, usableCount int
		pfd                                   PIXELFORMATDESCRIPTOR
	)
	nativeCount = describePixelFormat(window.context.wgl.dc, 1, int(unsafe.Sizeof(pfd)), nil)
	if _glfw.wgl.ARB_pixel_format != 0 {
		ADD_ATTRIB(WGL_SUPPORT_OPENGL_ARB)
		ADD_ATTRIB(WGL_DRAW_TO_WINDOW_ARB)
		ADD_ATTRIB(WGL_PIXEL_TYPE_ARB)
		ADD_ATTRIB(WGL_ACCELERATION_ARB)
		ADD_ATTRIB(WGL_RED_BITS_ARB)
		ADD_ATTRIB(WGL_RED_SHIFT_ARB)
		ADD_ATTRIB(WGL_GREEN_BITS_ARB)
		ADD_ATTRIB(WGL_GREEN_SHIFT_ARB)
		ADD_ATTRIB(WGL_BLUE_BITS_ARB)
		ADD_ATTRIB(WGL_BLUE_SHIFT_ARB)
		ADD_ATTRIB(WGL_ALPHA_BITS_ARB)
		ADD_ATTRIB(WGL_ALPHA_SHIFT_ARB)
		ADD_ATTRIB(WGL_DEPTH_BITS_ARB)
		ADD_ATTRIB(WGL_STENCIL_BITS_ARB)
		ADD_ATTRIB(WGL_ACCUM_BITS_ARB)
		ADD_ATTRIB(WGL_ACCUM_RED_BITS_ARB)
		ADD_ATTRIB(WGL_ACCUM_GREEN_BITS_ARB)
		ADD_ATTRIB(WGL_ACCUM_BLUE_BITS_ARB)
		ADD_ATTRIB(WGL_ACCUM_ALPHA_BITS_ARB)
		ADD_ATTRIB(WGL_AUX_BUFFERS_ARB)
		ADD_ATTRIB(WGL_DOUBLE_BUFFER_ARB)

		if _glfw.wgl.ARB_multisample {
			ADD_ATTRIB(WGL_SAMPLES_ARB)
		}
		if ctxconfig.client == GLFW_OPENGL_API && (_glfw.wgl.ARB_framebuffer_sRGB || _glfw.wgl.EXT_framebuffer_sRGB) {
			ADD_ATTRIB(WGL_FRAMEBUFFER_SRGB_CAPABLE_ARB)
		} else if _glfw.wgl.EXT_colorspace {
			ADD_ATTRIB(WGL_COLORSPACE_EXT)
		}
		attrib := WGL_NUMBER_PIXEL_FORMATS_ARB
		var extensionCount int
		if wglGetPixelFormatAttribivARB(window.context.wgl.dc, 1, 0, 1, &attrib, &extensionCount) == 0 {
			panic("WGL: Failed to retrieve pixel format attribute")
		}
		nativeCount = min(nativeCount, extensionCount)
	}
	usableConfigs := make([]_GLFWfbconfig, nativeCount)
	for i := 0; i < nativeCount; i++ {
		u := &usableConfigs[usableCount]
		pixelFormat = i + 1
		if _glfw.wgl.ARB_pixel_format != 0 {
			// Get pixel format attributes through "modern" extension
			if wglGetPixelFormatAttribivARB(window.context.wgl.dc, pixelFormat, 0, attribCount, &attribs[0], &values[0]) == 0 {
				panic("WGL: Failed to retrieve pixel format attributes")
			}
			if FIND_ATTRIB_VALUE(WGL_SUPPORT_OPENGL_ARB) == 0 || FIND_ATTRIB_VALUE(WGL_DRAW_TO_WINDOW_ARB) == 0 {
				continue
			}
			if FIND_ATTRIB_VALUE(WGL_PIXEL_TYPE_ARB) != WGL_TYPE_RGBA_ARB {
				continue
			}
			if FIND_ATTRIB_VALUE(WGL_ACCELERATION_ARB) == WGL_NO_ACCELERATION_ARB {
				continue
			}
			if (FIND_ATTRIB_VALUE(WGL_DOUBLE_BUFFER_ARB) != 0) != fbconfig.doublebuffer {
				continue
			}
			u.redBits = FIND_ATTRIB_VALUE(WGL_RED_BITS_ARB)
			u.greenBits = FIND_ATTRIB_VALUE(WGL_GREEN_BITS_ARB)
			u.blueBits = FIND_ATTRIB_VALUE(WGL_BLUE_BITS_ARB)
			u.alphaBits = FIND_ATTRIB_VALUE(WGL_ALPHA_BITS_ARB)
			u.depthBits = FIND_ATTRIB_VALUE(WGL_DEPTH_BITS_ARB)
			u.stencilBits = FIND_ATTRIB_VALUE(WGL_STENCIL_BITS_ARB)
			u.accumRedBits = FIND_ATTRIB_VALUE(WGL_ACCUM_RED_BITS_ARB)
			u.accumGreenBits = FIND_ATTRIB_VALUE(WGL_ACCUM_GREEN_BITS_ARB)
			u.accumBlueBits = FIND_ATTRIB_VALUE(WGL_ACCUM_BLUE_BITS_ARB)
			u.accumAlphaBits = FIND_ATTRIB_VALUE(WGL_ACCUM_ALPHA_BITS_ARB)
			u.auxBuffers = FIND_ATTRIB_VALUE(WGL_AUX_BUFFERS_ARB)
			if _glfw.wgl.ARB_multisample {
				u.samples = FIND_ATTRIB_VALUE(WGL_SAMPLES_ARB)
			}
			if ctxconfig.client == GLFW_OPENGL_API {
				if _glfw.wgl.ARB_framebuffer_sRGB || _glfw.wgl.EXT_framebuffer_sRGB {
					if FIND_ATTRIB_VALUE(WGL_FRAMEBUFFER_SRGB_CAPABLE_ARB) != 0 {
						u.sRGB = true
					}
				}
			} else {
				if _glfw.wgl.EXT_colorspace {
					if FIND_ATTRIB_VALUE(WGL_COLORSPACE_EXT) == WGL_COLORSPACE_SRGB_EXT {
						u.sRGB = true
					}
				}
			}
		} else {
			// Get pixel format attributes through legacy PFDs
			if describePixelFormat(window.context.wgl.dc, pixelFormat, int(unsafe.Sizeof(pfd)), &pfd) == 0 {
				panic("WGL: Failed to describe pixel format")
			}
			if (pfd.dwFlags&PFD_DRAW_TO_WINDOW) == 0 || (pfd.dwFlags&PFD_SUPPORT_OPENGL) == 0 {
				continue
			}
			// if pfd.dwFlags&PFD_GENERIC_ACCELERATED == 0 && pfd.dwFlags&PFD_GENERIC_FORMAT == 0 {
			//	continue
			// }
			if pfd.iPixelType != PFD_TYPE_RGBA {
				continue
			}
			if ((pfd.dwFlags & PFD_DOUBLEBUFFER) != 0) != fbconfig.doublebuffer {
				continue
			}
			u.redBits = int(pfd.cRedBits)
			u.greenBits = int(pfd.cGreenBits)
			u.blueBits = int(pfd.cBlueBits)
			u.alphaBits = int(pfd.cAlphaBits)
			u.depthBits = int(pfd.cDepthBits)
			u.stencilBits = int(pfd.cStencilBits)
			u.accumRedBits = int(pfd.cAccumRedBits)
			u.accumGreenBits = int(pfd.cAccumGreenBits)
			u.accumBlueBits = int(pfd.cAccumBlueBits)
			u.accumAlphaBits = int(pfd.cAccumAlphaBits)
			u.auxBuffers = int(pfd.cAuxBuffers)
		}
		u.handle = uintptr(pixelFormat)
		usableCount++
	}
	if usableCount == 0 {
		panic("WGL: The driver does not appear to support OpenGL")
	}
	closest = glfwChooseFBConfig(fbconfig, usableConfigs, usableCount)
	if closest == nil {
		panic("WGL: Failed to find a suitable pixel format")
	}
	pixelFormat = int(closest.handle)
	return pixelFormat
}

func glfwChooseFBConfig(desired *_GLFWfbconfig, alternatives []_GLFWfbconfig, count int) *_GLFWfbconfig {
	var i int
	var missing, leastMissing = INT_MAX, INT_MAX
	var colorDiff, leastColorDiff = INT_MAX, INT_MAX
	var extraDiff, leastExtraDiff = INT_MAX, INT_MAX
	var closest *_GLFWfbconfig

	for i = 0; i < count; i++ {
		current := &alternatives[i]
		// Count number of missing buffers
		missing = 0
		if desired.alphaBits > 0 && current.alphaBits == 0 {
			missing++
		}
		if desired.depthBits > 0 && current.depthBits == 0 {
			missing++
		}

		if desired.stencilBits > 0 && current.stencilBits == 0 {
			missing++
		}
		if desired.auxBuffers > 0 && current.auxBuffers < desired.auxBuffers {
			missing += desired.auxBuffers - current.auxBuffers
		}
		if desired.samples > 0 && current.samples == 0 {
			// Technically, several multisampling buffers could be
			// involved, but that's a lower level implementation detail and
			// not important to us here, so we count them as one
			missing++
		}
		if desired.transparent != current.transparent {
			missing++
		}
		// These polynomials make many small channel size differences matter
		// less than one large channel size difference
		// Calculate color channel size difference value
		colorDiff = 0
		if desired.redBits != GLFW_DONT_CARE {
			colorDiff += (desired.redBits - current.redBits) * (desired.redBits - current.redBits)
		}
		if desired.greenBits != GLFW_DONT_CARE {
			colorDiff += (desired.greenBits - current.greenBits) * (desired.greenBits - current.greenBits)
		}
		if desired.blueBits != GLFW_DONT_CARE {
			colorDiff += (desired.blueBits - current.blueBits) * (desired.blueBits - current.blueBits)
		}

		// Calculate non-color channel size difference value
		extraDiff = 0
		if desired.alphaBits != GLFW_DONT_CARE {
			extraDiff += (desired.alphaBits - current.alphaBits) * (desired.alphaBits - current.alphaBits)
		}
		if desired.depthBits != GLFW_DONT_CARE {
			extraDiff += (desired.depthBits - current.depthBits) * (desired.depthBits - current.depthBits)
		}
		if desired.stencilBits != GLFW_DONT_CARE {
			extraDiff += (desired.stencilBits - current.stencilBits) * (desired.stencilBits - current.stencilBits)
		}
		if desired.accumRedBits != GLFW_DONT_CARE {
			extraDiff += (desired.accumRedBits - current.accumRedBits) * (desired.accumRedBits - current.accumRedBits)
		}
		if desired.accumGreenBits != GLFW_DONT_CARE {
			extraDiff += (desired.accumGreenBits - current.accumGreenBits) * (desired.accumGreenBits - current.accumGreenBits)
		}
		if desired.accumBlueBits != GLFW_DONT_CARE {
			extraDiff += (desired.accumBlueBits - current.accumBlueBits) * (desired.accumBlueBits - current.accumBlueBits)
		}
		if desired.accumAlphaBits != GLFW_DONT_CARE {
			extraDiff += (desired.accumAlphaBits - current.accumAlphaBits) * (desired.accumAlphaBits - current.accumAlphaBits)
		}
		if desired.samples != GLFW_DONT_CARE {
			extraDiff += (desired.samples - current.samples) * (desired.samples - current.samples)
		}
		if desired.sRGB && !current.sRGB {
			extraDiff++
		}

		// Figure out if the current one is better than the best one found so far
		// Least number of missing buffers is the most important heuristic,
		// then color buffer size match and lastly size match for other buffers
		if missing < leastMissing {
			closest = current
		} else if missing == leastMissing {
			if (colorDiff < leastColorDiff) || (colorDiff == leastColorDiff && extraDiff < leastExtraDiff) {
				closest = current
			}
		}

		if current == closest {
			leastMissing = missing
			leastColorDiff = colorDiff
			leastExtraDiff = extraDiff
		}
	}
	return closest
}

func wglMakeCurrent(g *_GLFWtls, w *_GLFWwindow) bool {
	return false
}

func makeContextCurrentWGL(window *_GLFWwindow) error {
	if window != nil {
		if !makeCurrent(window.context.wgl.dc, window.context.wgl.handle) {
			return fmt.Errorf("WGL: Failed to make context current")
		}
	} else {
		if !makeCurrent(0, 0) {
			return fmt.Errorf("WGL: Failed to clear current context")
		}
	}
	return nil
}

func destroyContextWGL(window *_GLFWwindow) {
	if window.context.wgl.handle != 0 {
		deleteContext(window.context.wgl.handle)
		window.context.wgl.handle = 0
	}
}

// GoStr takes a null-terminated string returned by OpenGL and constructs a
// corresponding Go string.
func GoStr(cstr *uint8) string {
	str := ""
	if cstr == nil {
		return str
	}
	for {
		if *cstr == 0 {
			break
		}
		str += string(*cstr)
		cstr = (*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(cstr)) + 1))
	}
	return str
}

func _glfwRefreshContextAttribs(window *_GLFWwindow, ctxconfig *_GLFWctxconfig) error {
	glfwMakeContextCurrent(window)
	window.context.major = 3
	window.context.minor = 3
	window.context.revision = 3
	if window.context.major == 0 {
		return fmt.Errorf("No version found in OpenGL version string")
	}
	if window.context.major < ctxconfig.major || window.context.major == ctxconfig.major && window.context.minor < ctxconfig.minor {
		// The desired OpenGL version is greater than the actual version
		// This only happens if the machine lacks {GLX|WGL}_ARB_create_context
		// /and/ the user has requested an OpenGL version greater than 1.0
		return fmt.Errorf("Requested OpenGL version %d.%d, got version %d.%d", ctxconfig.major, ctxconfig.minor, window.context.major, window.context.minor)
	}
	return nil
}
