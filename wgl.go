package glfw

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	opengl32 = windows.NewLazySystemDLL("opengl32.dll")
)

const (
	wgl_NUMBER_PIXEL_FORMATS_ARB                = 0x2000
	wgl_DRAW_TO_WINDOW_ARB                      = 0x2001
	wgl_ACCELERATION_ARB                        = 0x2003
	wgl_SUPPORT_OPENGL_ARB                      = 0x2010
	wgl_DOUBLE_BUFFER_ARB                       = 0x2011
	wgl_STEREO_ARB                              = 0x2012
	wgl_PIXEL_TYPE_ARB                          = 0x2013
	wgl_RED_BITS_ARB                            = 0x2015
	wgl_RED_SHIFT_ARB                           = 0x2016
	wgl_GREEN_BITS_ARB                          = 0x2017
	wgl_GREEN_SHIFT_ARB                         = 0x2018
	wgl_BLUE_BITS_ARB                           = 0x2019
	wgl_BLUE_SHIFT_ARB                          = 0x201A
	wgl_ALPHA_BITS_ARB                          = 0x201B
	wgl_ALPHA_SHIFT_ARB                         = 0x201C
	wgl_ACCUM_BITS_ARB                          = 0x201D
	wgl_ACCUM_RED_BITS_ARB                      = 0x201E
	wgl_ACCUM_GREEN_BITS_ARB                    = 0x201F
	wgl_ACCUM_BLUE_BITS_ARB                     = 0x2020
	wgl_ACCUM_ALPHA_BITS_ARB                    = 0x2021
	wgl_DEPTH_BITS_ARB                          = 0x2022
	wgl_STENCIL_BITS_ARB                        = 0x2023
	wgl_AUX_BUFFERS_ARB                         = 0x2024
	wgl_NO_ACCELERATION_ARB                     = 0x2025
	wgl_TYPE_RGBA_ARB                           = 0x202B
	wgl_SAMPLES_ARB                             = 0x2042
	wgl_FRAMEBUFFER_SRGB_CAPABLE_ARB            = 0x20a9
	wgl_COLORSPACE_EXT                          = 0x309d
	wgl_COLORSPACE_SRGB_EXT                     = 0x3089
	wgl_CONTEXT_DEBUG_BIT_ARB                   = 0x00000001
	wgl_CONTEXT_FORWARD_COMPATIBLE_BIT_ARB      = 0x00000002
	wgl_CONTEXT_PROFILE_MASK_ARB                = 0x9126
	wgl_CONTEXT_CORE_PROFILE_BIT_ARB            = 0x00000001
	wgl_CONTEXT_COMPATIBILITY_PROFILE_BIT_ARB   = 0x00000002
	wgl_CONTEXT_MAJOR_VERSION_ARB               = 0x2091
	wgl_CONTEXT_MINOR_VERSION_ARB               = 0x2092
	wgl_CONTEXT_FLAGS_ARB                       = 0x2094
	wgl_CONTEXT_ES2_PROFILE_BIT_EXT             = 0x00000004
	wgl_CONTEXT_ROBUST_ACCESS_BIT_ARB           = 0x00000004
	wgl_LOSE_CONTEXT_ON_RESET_ARB               = 0x8252
	wgl_CONTEXT_RESET_NOTIFICATION_STRATEGY_ARB = 0x8256
	wgl_NO_RESET_NOTIFICATION_ARB               = 0x8261
	wgl_CONTEXT_RELEASE_BEHAVIOR_ARB            = 0x2097
	wgl_CONTEXT_RELEASE_BEHAVIOR_NONE_ARB       = 0
	wgl_CONTEXT_OPENGL_NO_ERROR_ARB             = 0x31b3
	wgl_CONTEXT_RELEASE_BEHAVIOR_FLUSH_ARB      = 0x2098
)

func addAttrib(w *Window, a int32) {
	w.attribs[w.attribCount] = a
	w.attribCount++
}

func findAttrib(w *Window, a int32) int32 {
	for i := 0; i < w.attribCount; i++ {
		if w.attribs[i] == a {
			return w.values[i]
		}
	}
	panic("WGL: Unknown pixel format attribute requested")
}

func wglGetProcAddress(name string) uintptr {
	var b [64]byte
	copy(b[:], name)
	r, _, err := _glfw.wgl.wglGetProcAddress.Call(uintptr(unsafe.Pointer(&b)))
	if !errors.Is(err, syscall.Errno(0)) {
		// Ignore errors. The address is checked before calling the extensions.
		// This is done to enable use of older computers which does not have the extensions.
		return 0
	}
	return r
}

func swapBuffersWGL(window *_GLFWwindow) {
	_, _, _ = _SwapBuffers.Call(uintptr(window.context.wgl.dc))
	// Ignore errors becausee it sometimes fails without reason
}

func getCurrentWindow() *_GLFWwindow {
	p := glfwPlatformGetTls(&_glfw.contextSlot)
	window := _glfw.windowListHead
	for window != nil {
		if uintptr(unsafe.Pointer(window)) == p {
			return window
		}
		window = window.next
	}
	return nil
}

func swapIntervalWGL(interval int) {
	// p := glfwPlatformGetTls(&_glfw.contextSlot)
	// window := (*Window)(unsafe.Pointer(p))
	window := getCurrentWindow()
	if window == nil {
		panic("swapIntervalWGL failed, window is nil\n")
	}
	window.context.wgl.interval = interval
	if _glfw.wgl.EXT_swap_control {
		_, _, _ = syscall.SyscallN(_glfw.wgl.SwapIntervalEXT, uintptr(interval))
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
		slog.Error("wgl makeCurrent failed", "err", err, "dc", dc, "handle", handle)
	}
	return r1 != 0
}

func setPixelFormat(dc HDC, iPixelFormat int32, pfd *PIXELFORMATDESCRIPTOR) int {
	ret, _, err := _SetPixelFormat.Call(uintptr(dc), uintptr(iPixelFormat), uintptr(unsafe.Pointer(pfd)))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("wglSetPixelFormat failed, " + err.Error())
	}
	return int(ret)
}

func choosePixelFormat(dc HDC, pfd *PIXELFORMATDESCRIPTOR) int32 {
	ret, _, err := _ChoosePixelFormat.Call(uintptr(dc), uintptr(unsafe.Pointer(pfd)))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("choosePixelFormat failed, " + err.Error())
	}
	return int32(ret)
}

func wglCreateContextAttribsARB(dc HDC, share syscall.Handle, attribs *int32) HANDLE {
	ret, _, err := syscall.SyscallN(_glfw.wgl.wglCreateContextAttribsARB, uintptr(dc), uintptr(share), uintptr(unsafe.Pointer(attribs)))
	// We do not check err, as it seems to be 126 all the time, even when ok.
	if ret == 0 {
		panic("wglCreateContextAttribsARB failed " + err.Error())
	}
	return HANDLE(ret)
}

func shareLists(share syscall.Handle, handle HANDLE) bool {
	ret, _, err := _glfw.wgl.wglShareLists.Call(uintptr(share), uintptr(handle))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("wglShareLists failed, " + err.Error())
	}
	return ret == 0
}

func extensionSupportedWGL(extension string) bool {
	var extensions string
	if _glfw.wgl.GetExtensionsStringARB != 0 {
		r, _, err := syscall.SyscallN(_glfw.wgl.GetExtensionsStringARB, uintptr(getCurrentDC()))
		if !errors.Is(err, syscall.Errno(0)) {
			panic("GetExtensionsStringARB failed, " + err.Error())
		}
		extensions = GoStr((*uint8)(unsafe.Pointer(r)))
	} else if _glfw.wgl.GetExtensionsStringEXT != 0 {
		r, _, err := syscall.SyscallN(_glfw.wgl.GetExtensionsStringEXT, uintptr(getCurrentDC()))
		if !errors.Is(err, syscall.Errno(0)) {
			panic("GetExtensionsStringEXT failed, " + err.Error())
		}
		extensions = GoStr((*uint8)(unsafe.Pointer(r)))
	}
	if extensions == "" {
		return false
	}
	return strings.Contains(extensions, extension)
}

func getProcAddressWGL(procName string) uintptr {
	proc := wglGetProcAddress(procName)
	if proc != 0 {
		return proc
	}
	return opengl32.NewProc(procName).Addr()
}

func _glfwInitWGL() error {
	var pfd PIXELFORMATDESCRIPTOR
	// Force the lazy handle to load now
	if err := opengl32.Load(); err != nil {
		return fmt.Errorf("your graphic card does not support OpenGl, 'opengl32.dll' not available: %w", err)
	}
	if _glfw.wgl.instance != nil {
		return nil
	}
	_glfw.wgl.instance = opengl32
	_glfw.wgl.wglCreateContext = opengl32.NewProc("wglCreateContext")
	_glfw.wgl.wglDeleteContext = opengl32.NewProc("wglDeleteContext")
	_glfw.wgl.wglGetProcAddress = opengl32.NewProc("wglGetProcAddress")
	_glfw.wgl.wglGetCurrentDC = opengl32.NewProc("wglGetCurrentDC")
	_glfw.wgl.wglGetCurrentContext = opengl32.NewProc("wglGetCurrentContext")
	_glfw.wgl.wglMakeCurrent = opengl32.NewProc("wglMakeCurrent")
	_glfw.wgl.wglShareLists = opengl32.NewProc("wglShareLists")

	if err := _glfw.wgl.wglGetProcAddress.Find(); err != nil {
		return fmt.Errorf("your graphic card does not support OpenGl, wglGetProcAddress not found: %w", err)
	}
	// NOTE: A dummy context has to be created for opengl32.dll to load the
	// OpenGL Installable Client Driver, from which we can then query WGL extensions
	dc := getDC(_glfw.win32.helperWindowHandle)
	pfd.nSize = uint16(unsafe.Sizeof(pfd))
	pfd.nVersion = 1
	pfd.dwFlags = PFD_DRAW_TO_WINDOW | PFD_SUPPORT_OPENGL | PFD_DOUBLEBUFFER
	pfd.iPixelType = PFD_TYPE_RGBA
	pfd.cColorBits = 24
	if setPixelFormat(dc, choosePixelFormat(dc, &pfd), &pfd) == 0 {
		err := syscall.GetLastError()
		if err != nil {
			panic(err)
		}
	}

	rc := createContext(dc)
	if rc == 0 {
		panic("WGL: Failed to create dummy context")
	}
	pdc := getCurrentDC()
	prc := getCurrentContext()
	if !makeCurrent(dc, rc) {
		slog.Error("WGL: Failed to make dummy context current")
		_, _, _ = _glfw.wgl.wglMakeCurrent.Call(uintptr(pdc), uintptr(prc))
		_, _, _ = _glfw.wgl.wglDeleteContext.Call(uintptr(rc))
		return fmt.Errorf("WGL: Failed to make dummy context current")
	}

	// NOTE: Functions must be loaded first as they're needed to retrieve the
	//       extension string that tells us whether the functions are supported
	_glfw.wgl.GetExtensionsStringARB = wglGetProcAddress("wglGetExtensionsStringARB")
	_glfw.wgl.GetExtensionsStringEXT = wglGetProcAddress("wglGetExtensionsStringEXT")
	_glfw.wgl.wglCreateContextAttribsARB = wglGetProcAddress("wglCreateContextAttribsARB")
	_glfw.wgl.SwapIntervalEXT = wglGetProcAddress("wglSwapIntervalEXT")
	_glfw.wgl.GetPixelFormatAttribivARB = wglGetProcAddress("wglGetPixelFormatAttribivARB")

	_glfw.wgl.ARB_multisample = extensionSupportedWGL("WGL_ARB_multisample")
	_glfw.wgl.ARB_framebuffer_sRGB = extensionSupportedWGL("WGL_ARB_framebuffer_sRGB")
	_glfw.wgl.EXT_framebuffer_sRGB = extensionSupportedWGL("WGL_EXT_framebuffer_sRGB")
	_glfw.wgl.ARB_create_context = extensionSupportedWGL("WGL_ARB_create_context")
	_glfw.wgl.ARB_create_context_profile = extensionSupportedWGL("WGL_ARB_create_context_profile")
	_glfw.wgl.EXT_create_context_es2_profile = extensionSupportedWGL("WGL_EXT_create_context_es2_profile")
	_glfw.wgl.ARB_create_context_robustness = extensionSupportedWGL("WGL_ARB_create_context_robustness")
	_glfw.wgl.ARB_create_context_no_error = extensionSupportedWGL("WGL_ARB_create_context_no_error")
	_glfw.wgl.EXT_swap_control = extensionSupportedWGL("WGL_EXT_swap_control")
	_glfw.wgl.EXT_colorspace = extensionSupportedWGL("WGL_EXT_colorspace")
	_glfw.wgl.ARB_pixel_format = extensionSupportedWGL("WGL_ARB_pixel_format")
	_glfw.wgl.ARB_context_flush_control = extensionSupportedWGL("WGL_ARB_context_flush_control")
	makeCurrent(pdc, prc)
	deleteContext(rc)
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

func describePixelFormat(dc HDC, iPixelFormat int32, nBytes int, ppfd *PIXELFORMATDESCRIPTOR) int32 {
	r1, _, err := _DescribePixelFormat.Call(uintptr(dc), uintptr(iPixelFormat), uintptr(nBytes), uintptr(unsafe.Pointer(ppfd)))
	if r1 == 0 || !errors.Is(err, syscall.Errno(0)) {
		slog.Error("describePixelFormat failed, " + err.Error())
		r1 = 0
	}
	return int32(r1)
}

func glfwCreateContextWGL(window *_GLFWwindow, ctxConfig *_GLFWctxconfig, fbConfig *_GLFWfbconfig) error {
	attribList := make([]int32, 0, 40)
	var pfd PIXELFORMATDESCRIPTOR
	hShare := syscall.Handle(0)
	if ctxConfig.share != nil {
		hShare = ctxConfig.share.Win32.Handle
	}
	share := ctxConfig.share
	window.context.wgl.dc = getDC(window.Win32.Handle)
	if window.context.wgl.dc == 0 {
		return fmt.Errorf("WGL: Failed to retrieve DC for window")
	}
	pixelFormat := choosePixelFormatWGL(window, ctxConfig, fbConfig) // 14
	if pixelFormat == 0 {
		return fmt.Errorf("WGL: Failed to retrieve PixelFormat for window")
	}
	if describePixelFormat(window.context.wgl.dc, pixelFormat, int(unsafe.Sizeof(pfd)), &pfd) == 0 {
		return fmt.Errorf("WGL: Failed to retrieve PFD for selected pixel format")
	}
	if setPixelFormat(window.context.wgl.dc, pixelFormat, &pfd) == 0 {
		return fmt.Errorf("WGL: Failed to set selected pixel format")
	}
	if ctxConfig.client == OpenGLAPI {
		if ctxConfig.forward && !_glfw.wgl.ARB_create_context {
			return fmt.Errorf("WGL: A forward compatible OpenGL context requested but wgl_ARB_create_context is unavailable")
		}
		if (ctxConfig.profile != 0) && !_glfw.wgl.ARB_create_context_profile {
			return fmt.Errorf("WGL: OpenGL profile requested but wgl_ARB_create_context_profile is unavailable")
		}
	} else {
		if !_glfw.wgl.ARB_create_context || !_glfw.wgl.ARB_create_context_profile || !_glfw.wgl.EXT_create_context_es2_profile {
			return fmt.Errorf("WGL: OpenGL ES requested but wgl_ARB_create_context_es2_profile is unavailable")
		}
	}
	if _glfw.wgl.ARB_create_context {
		mask := 0
		flags := 0
		if ctxConfig.client == OpenGLAPI {
			if ctxConfig.forward {
				flags |= wgl_CONTEXT_FORWARD_COMPATIBLE_BIT_ARB
			}
			if ctxConfig.profile == OpenGLCoreProfile {
				mask |= wgl_CONTEXT_CORE_PROFILE_BIT_ARB
			} else if ctxConfig.profile == OpenGLCompatProfile {
				mask |= wgl_CONTEXT_COMPATIBILITY_PROFILE_BIT_ARB
			}
		} else {
			mask |= wgl_CONTEXT_ES2_PROFILE_BIT_EXT
		}
		if ctxConfig.debug {
			flags |= wgl_CONTEXT_DEBUG_BIT_ARB
		}
		if ctxConfig.robustness != 0 {
			if _glfw.wgl.ARB_create_context_robustness {
				if ctxConfig.robustness == NoResetNotification {
					attribList = append(attribList, wgl_CONTEXT_RESET_NOTIFICATION_STRATEGY_ARB, wgl_NO_RESET_NOTIFICATION_ARB)
				}
			} else if ctxConfig.robustness == LoseContextOnReset {
				attribList = append(attribList, wgl_CONTEXT_RESET_NOTIFICATION_STRATEGY_ARB, wgl_LOSE_CONTEXT_ON_RESET_ARB)
			}
			flags |= wgl_CONTEXT_ROBUST_ACCESS_BIT_ARB
		}
		if ctxConfig.release != 0 {
			if _glfw.wgl.ARB_context_flush_control {
				if ctxConfig.release == ReleaseBehaviorNone {
					attribList = append(attribList, wgl_CONTEXT_RELEASE_BEHAVIOR_ARB, wgl_CONTEXT_RELEASE_BEHAVIOR_NONE_ARB)
				} else if ctxConfig.release == ReleaseBehaviorFlush {
					attribList = append(attribList, wgl_CONTEXT_RELEASE_BEHAVIOR_ARB, wgl_CONTEXT_RELEASE_BEHAVIOR_FLUSH_ARB)
				}
			}
		}
		if ctxConfig.noerror {
			if _glfw.wgl.ARB_create_context_no_error {
				attribList = append(attribList, wgl_CONTEXT_OPENGL_NO_ERROR_ARB, 1)
			}
		}
		// Only request an explicitly versioned context when necessary,
		if ctxConfig.major != 1 || ctxConfig.minor != 0 {
			attribList = append(attribList, wgl_CONTEXT_MAJOR_VERSION_ARB, ctxConfig.major)
			attribList = append(attribList, wgl_CONTEXT_MINOR_VERSION_ARB, ctxConfig.minor)
		}
		if flags != 0 {
			attribList = append(attribList, wgl_CONTEXT_FLAGS_ARB, int32(flags))
		}
		if mask != 0 {
			attribList = append(attribList, wgl_CONTEXT_PROFILE_MASK_ARB, int32(mask))
		}
		// Add sentinel for end-of-list
		attribList = append(attribList, 0, 0)
		window.context.wgl.handle = wglCreateContextAttribsARB(window.context.wgl.dc, hShare, &attribList[0])
		if window.context.wgl.handle == 0 {
			return fmt.Errorf("WGL: Driver does not support OpenGL version %d.%d", ctxConfig.major, ctxConfig.minor)
		}
	} else {
		window.context.wgl.handle = createContext(window.context.wgl.dc)
		if window.context.wgl.handle == 0 {
			return fmt.Errorf("WGL: Failed to create OpenGL context")
		}
		if share != nil {
			if shareLists(share.Win32.Handle, window.context.wgl.handle) {
				return fmt.Errorf("WGL: Failed to enable sharing with specified OpenGL context")
			}
		}
	}

	window.context.makeCurrent = makeContextCurrentWGL
	window.context.swapBuffers = swapBuffersWGL
	window.context.swapInterval = swapIntervalWGL
	window.context.extensionSupported = extensionSupportedWGL
	window.context.getProcAddress = getProcAddressWGL
	window.context.destroy = destroyContextWGL
	return nil
}

func wglGetPixelFormatAttribivARB(dc HDC, pixelFormat int32, layerPlane int, nAttrib int, attributes *int32, piValues *int32) {
	r, _, err := syscall.SyscallN(_glfw.wgl.GetPixelFormatAttribivARB, uintptr(dc), uintptr(pixelFormat), uintptr(layerPlane),
		uintptr(nAttrib), uintptr(unsafe.Pointer(attributes)), uintptr(unsafe.Pointer(piValues)))
	if !errors.Is(err, syscall.Errno(0)) || r == 0 {
		panic("WGL: GetPixelFormatAttribivARB failed, " + err.Error())
	}
}

func choosePixelFormatWGL(w *_GLFWwindow, ctxConfig *_GLFWctxconfig, fbConfig *_GLFWfbconfig) int32 {
	var (
		closest                               *_GLFWfbconfig
		pixelFormat, nativeCount, usableCount int32
		pfd                                   PIXELFORMATDESCRIPTOR
	)
	nativeCount = describePixelFormat(w.context.wgl.dc, 1, int(unsafe.Sizeof(pfd)), nil)
	w.attribCount = 0
	if _glfw.wgl.ARB_pixel_format {
		addAttrib(w, wgl_SUPPORT_OPENGL_ARB)
		addAttrib(w, wgl_DRAW_TO_WINDOW_ARB)
		addAttrib(w, wgl_PIXEL_TYPE_ARB)
		addAttrib(w, wgl_ACCELERATION_ARB)
		addAttrib(w, wgl_RED_BITS_ARB)
		addAttrib(w, wgl_RED_SHIFT_ARB)
		addAttrib(w, wgl_GREEN_BITS_ARB)
		addAttrib(w, wgl_GREEN_SHIFT_ARB)
		addAttrib(w, wgl_BLUE_BITS_ARB)
		addAttrib(w, wgl_BLUE_SHIFT_ARB)
		addAttrib(w, wgl_ALPHA_BITS_ARB)
		addAttrib(w, wgl_ALPHA_SHIFT_ARB)
		addAttrib(w, wgl_DEPTH_BITS_ARB)
		addAttrib(w, wgl_STENCIL_BITS_ARB)
		addAttrib(w, wgl_ACCUM_BITS_ARB)
		addAttrib(w, wgl_ACCUM_RED_BITS_ARB)
		addAttrib(w, wgl_ACCUM_GREEN_BITS_ARB)
		addAttrib(w, wgl_ACCUM_BLUE_BITS_ARB)
		addAttrib(w, wgl_ACCUM_ALPHA_BITS_ARB)
		addAttrib(w, wgl_AUX_BUFFERS_ARB)
		addAttrib(w, wgl_STEREO_ARB)
		addAttrib(w, wgl_DOUBLE_BUFFER_ARB)

		if _glfw.wgl.ARB_multisample {
			addAttrib(w, wgl_SAMPLES_ARB)
		}
		if ctxConfig.client == OpenGLAPI {
			if _glfw.wgl.ARB_framebuffer_sRGB || _glfw.wgl.EXT_framebuffer_sRGB {
				addAttrib(w, wgl_FRAMEBUFFER_SRGB_CAPABLE_ARB)
			}
		} else {
			if _glfw.wgl.EXT_colorspace {
				addAttrib(w, wgl_COLORSPACE_EXT)
			}
		}
		attrib := int32(wgl_NUMBER_PIXEL_FORMATS_ARB)
		var extensionCount int32
		wglGetPixelFormatAttribivARB(w.context.wgl.dc, 1, 0, 1, &attrib, &extensionCount)
		nativeCount = min(nativeCount, extensionCount)
	}
	usableConfigs := make([]_GLFWfbconfig, nativeCount)
	for i := 0; i < int(nativeCount); i++ {
		u := &usableConfigs[usableCount]
		pixelFormat = int32(i) + 1
		if _glfw.wgl.ARB_pixel_format {
			// Get pixel format attributes through "modern" extension
			w.values[0] = 0
			wglGetPixelFormatAttribivARB(w.context.wgl.dc, pixelFormat, 0, 1 /*attribCount*/, &w.attribs[2], &w.values[2])
			for j := 0; j < w.attribCount; j++ {
				wglGetPixelFormatAttribivARB(w.context.wgl.dc, pixelFormat, 0, 1 /*attribCount*/, &w.attribs[j], &w.values[j])
			}
			if findAttrib(w, wgl_SUPPORT_OPENGL_ARB) == 0 || findAttrib(w, wgl_DRAW_TO_WINDOW_ARB) == 0 {
				continue
			}
			if findAttrib(w, wgl_PIXEL_TYPE_ARB) != wgl_TYPE_RGBA_ARB {
				continue
			}
			if findAttrib(w, wgl_ACCELERATION_ARB) == wgl_NO_ACCELERATION_ARB {
				continue
			}
			if (findAttrib(w, wgl_DOUBLE_BUFFER_ARB) != 0) != fbConfig.doublebuffer {
				continue
			}
			u.redBits = findAttrib(w, wgl_RED_BITS_ARB)
			u.greenBits = findAttrib(w, wgl_GREEN_BITS_ARB)
			u.blueBits = findAttrib(w, wgl_BLUE_BITS_ARB)
			u.alphaBits = findAttrib(w, wgl_ALPHA_BITS_ARB)
			u.depthBits = findAttrib(w, wgl_DEPTH_BITS_ARB)
			u.stencilBits = findAttrib(w, wgl_STENCIL_BITS_ARB)
			u.accumRedBits = findAttrib(w, wgl_ACCUM_RED_BITS_ARB)
			u.accumGreenBits = findAttrib(w, wgl_ACCUM_GREEN_BITS_ARB)
			u.accumBlueBits = findAttrib(w, wgl_ACCUM_BLUE_BITS_ARB)
			u.accumAlphaBits = findAttrib(w, wgl_ACCUM_ALPHA_BITS_ARB)
			u.auxBuffers = findAttrib(w, wgl_AUX_BUFFERS_ARB)
			if findAttrib(w, wgl_STEREO_ARB) != 0 {
				u.stereo = true
			}
			if _glfw.wgl.ARB_multisample {
				u.samples = findAttrib(w, wgl_SAMPLES_ARB)
			}
			if ctxConfig.client == OpenGLAPI {
				if _glfw.wgl.ARB_framebuffer_sRGB || _glfw.wgl.EXT_framebuffer_sRGB {
					if findAttrib(w, wgl_FRAMEBUFFER_SRGB_CAPABLE_ARB) != 0 {
						u.sRGB = true
					}
				}
			} else {
				if _glfw.wgl.EXT_colorspace {
					if findAttrib(w, wgl_COLORSPACE_EXT) == wgl_COLORSPACE_SRGB_EXT {
						u.sRGB = true
					}
				}
			}
		} else {
			// Get pixel format attributes through legacy PFDs
			if describePixelFormat(w.context.wgl.dc, pixelFormat, int(unsafe.Sizeof(pfd)), &pfd) == 0 {
				panic("WGL: Failed to describe pixel format")
			}
			if (pfd.dwFlags&PFD_DRAW_TO_WINDOW) == 0 || (pfd.dwFlags&PFD_SUPPORT_OPENGL) == 0 {
				continue
			}
			if pfd.dwFlags&PFD_GENERIC_ACCELERATED == 0 && pfd.dwFlags&PFD_GENERIC_FORMAT == 0 {
				continue
			}
			if pfd.iPixelType != PFD_TYPE_RGBA {
				continue
			}
			if ((pfd.dwFlags & PFD_DOUBLEBUFFER) != 0) != fbConfig.doublebuffer {
				continue
			}
			u.redBits = int32(pfd.cRedBits)
			u.greenBits = int32(pfd.cGreenBits)
			u.blueBits = int32(pfd.cBlueBits)
			u.alphaBits = int32(pfd.cAlphaBits)
			u.depthBits = int32(pfd.cDepthBits)
			u.stencilBits = int32(pfd.cStencilBits)
			u.accumRedBits = int32(pfd.cAccumRedBits)
			u.accumGreenBits = int32(pfd.cAccumGreenBits)
			u.accumBlueBits = int32(pfd.cAccumBlueBits)
			u.accumAlphaBits = int32(pfd.cAccumAlphaBits)
			u.auxBuffers = int32(pfd.cAuxBuffers)
			if pfd.dwFlags&PFD_STEREO != 0 {
				u.stereo = true
			}
		}
		u.handle = uintptr(pixelFormat)
		usableCount++
	}
	if usableCount == 0 {
		panic("WGL: The driver does not appear to support OpenGL")
	}
	closest = glfwChooseFBConfig(fbConfig, usableConfigs, usableCount)
	if closest == nil {
		panic("WGL: Failed to find a suitable pixel format")
	}
	pixelFormat = int32(closest.handle)
	return pixelFormat
}

func makeContextCurrentWGL(window *_GLFWwindow) error {
	if window != nil {
		if !makeCurrent(window.context.wgl.dc, window.context.wgl.handle) {
			glfwPlatformSetTls(&_glfw.contextSlot, uintptr(unsafe.Pointer(window)))
			return fmt.Errorf("WGL: Failed to make context current")
		}
	} else {
		if !makeCurrent(0, 0) {
			glfwPlatformSetTls(&_glfw.contextSlot, uintptr(unsafe.Pointer(window)))
			return fmt.Errorf("WGL: Failed to clear current context")
		}
	}
	glfwPlatformSetTls(&_glfw.contextSlot, uintptr(unsafe.Pointer(window)))
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
func GoStr(cStr *uint8) string {
	str := ""
	if cStr == nil {
		return str
	}
	for {
		if *cStr == 0 {
			break
		}
		str += string(*cStr)
		cStr = (*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(cStr)) + 1))
	}
	return str
}
