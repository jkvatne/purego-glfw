package glfw

import (
	"errors"
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

func glfwIsValidContextConfig(ctxconfig *_GLFWctxconfig) error {
	if ctxconfig.source != NativeContextAPI && ctxconfig.source != EGLContextAPI && ctxconfig.source != OSMesaContextAPI {
		return errors.New("Invalid context creation API")
	}
	if ctxconfig.client != NoAPI && ctxconfig.client != OpenGLAPI && ctxconfig.client != OpenGLESAPI {
		return errors.New("Invalid cclient API")
	}
	if ctxconfig.share != nil {
		if ctxconfig.client == NoAPI || ctxconfig.share.context.client == OpenGLAPI {
			return errors.New("No context API")
		}
		if ctxconfig.client != ctxconfig.share.context.client {
			return errors.New("Context creation APIs do not match between contexts")
		}
	}
	if ctxconfig.client == OpenGLAPI {
		if (ctxconfig.major < 1 || ctxconfig.minor < 0) ||
			(ctxconfig.major == 1 && ctxconfig.minor > 5) ||
			(ctxconfig.major == 2 && ctxconfig.minor > 1) ||
			(ctxconfig.major == 3 && ctxconfig.minor > 3) {
			return fmt.Errorf("Invalid OpenGL version %d.%d", ctxconfig.major, ctxconfig.minor)
		}
		if ctxconfig.profile != 0 {
			if ctxconfig.profile != OpenGLCoreProfile && ctxconfig.profile != OpenGLCompatProfile {
				return fmt.Errorf("Invalid OpenGL profile")
			}
			if ctxconfig.major <= 2 || ctxconfig.major == 3 && ctxconfig.minor < 2 {
				// Desktop OpenGL context profiles are only defined for version 3.2 and above
				return fmt.Errorf("Context profiles are only defined for OpenGL version 3.2 and above")
			}
			if ctxconfig.forward && ctxconfig.major <= 2 {
				// Forward-compatible contexts are only defined for OpenGL version 3.0 and above
				return fmt.Errorf("Forward-compatibility is only defined for OpenGL version 3.0 and above")
			}
		}
	} else if ctxconfig.client == OpenGLESAPI {
		if ctxconfig.major < 1 || ctxconfig.minor < 0 || ctxconfig.major == 1 && ctxconfig.minor > 1 || ctxconfig.major == 2 && ctxconfig.minor > 0 {
			// OpenGL ES 1.0 is the smallest valid version
			// OpenGL ES 1.x series ended with version 1.1
			// OpenGL ES 2.x series ended with version 2.0
			// For now, let everything else through
			return fmt.Errorf("Invalid OpenGL ES version %i.%i", ctxconfig.major, ctxconfig.minor)
		}
	}

	if ctxconfig.robustness != 0 {
		if ctxconfig.robustness != NoResetNotification && ctxconfig.robustness != LoseContextOnReset {
			return fmt.Errorf("Invalid context robustness mode 0x%08X", ctxconfig.robustness)
		}
	}
	return nil
}

func glfwChooseFBConfig(desired *_GLFWfbconfig, alternatives []_GLFWfbconfig, count int32) *_GLFWfbconfig {
	var missing, leastMissing = int32(_INT_MAX), int32(_INT_MAX)
	var colorDiff, leastColorDiff = int32(_INT_MAX), int32(_INT_MAX)
	var extraDiff, leastExtraDiff = int32(_INT_MAX), int32(_INT_MAX)
	var closest *_GLFWfbconfig

	for i := int32(0); i < count; i++ {
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
			missing++
		}
		if desired.transparent != current.transparent {
			missing++
		}
		colorDiff = 0
		if desired.redBits != DontCare {
			colorDiff += (desired.redBits - current.redBits) * (desired.redBits - current.redBits)
		}
		if desired.greenBits != DontCare {
			colorDiff += (desired.greenBits - current.greenBits) * (desired.greenBits - current.greenBits)
		}
		if desired.blueBits != DontCare {
			colorDiff += (desired.blueBits - current.blueBits) * (desired.blueBits - current.blueBits)
		}

		// Calculate non-color channel size difference value
		extraDiff = 0
		if desired.alphaBits != DontCare {
			extraDiff += (desired.alphaBits - current.alphaBits) * (desired.alphaBits - current.alphaBits)
		}
		if desired.depthBits != DontCare {
			extraDiff += (desired.depthBits - current.depthBits) * (desired.depthBits - current.depthBits)
		}
		if desired.stencilBits != DontCare {
			extraDiff += (desired.stencilBits - current.stencilBits) * (desired.stencilBits - current.stencilBits)
		}
		if desired.accumRedBits != DontCare {
			extraDiff += (desired.accumRedBits - current.accumRedBits) * (desired.accumRedBits - current.accumRedBits)
		}
		if desired.accumGreenBits != DontCare {
			extraDiff += (desired.accumGreenBits - current.accumGreenBits) * (desired.accumGreenBits - current.accumGreenBits)
		}
		if desired.accumBlueBits != DontCare {
			extraDiff += (desired.accumBlueBits - current.accumBlueBits) * (desired.accumBlueBits - current.accumBlueBits)
		}
		if desired.accumAlphaBits != DontCare {
			extraDiff += (desired.accumAlphaBits - current.accumAlphaBits) * (desired.accumAlphaBits - current.accumAlphaBits)
		}
		if desired.samples != DontCare {
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

func glfwRefreshContextAttribs(window *_GLFWwindow, ctxconfig *_GLFWctxconfig) error {
	window.context.source = ctxconfig.source
	window.context.client = OpenGLAPI
	previous := (*Window)(unsafe.Pointer(glfwPlatformGetTls(&_glfw.contextSlot)))
	_ = glfwMakeContextCurrent(window)
	if glfwPlatformGetTls(&_glfw.contextSlot) != uintptr(unsafe.Pointer(window)) {
		return fmt.Errorf("glfwRefrechContext got Tls slot error")
	}

	window.context.GetIntegerv = window.context.getProcAddress("glGetIntegerv")
	window.context.GetString = window.context.getProcAddress("glGetString")
	if window.context.GetIntegerv == 0 || window.context.GetString == 0 {
		return fmt.Errorf("glfwRefreshContextAttribs: Entry point retrieval is broken")
	}
	r, _, err := syscall.SyscallN(window.context.GetString, uintptr(_GL_VERSION))
	if !errors.Is(err, syscall.Errno(0)) {
		return fmt.Errorf("String retrieval is broken, " + err.Error())
	}
	version := GoStr((*uint8)(unsafe.Pointer(uintptr(r))))
	prefixes := []string{"OpenGL ES-CM ", "OpenGL ES-CL ", "OpenGL ES ", ""}
	for s := range prefixes {
		if strings.HasPrefix(version, prefixes[s]) {
			version = strings.TrimPrefix(version, prefixes[s])
			window.context.client = OpenGLESAPI
		}
	}
	i := 0
	var v [3]int32
	for _, ch := range version {
		if ch >= '0' && ch <= '9' {
			v[i] = v[i]*10 + int32(ch) - int32('0')
		} else {
			i++
			if i >= 3 {
				break
			}
		}
	}
	window.context.major = v[0]
	window.context.minor = v[1]
	window.context.revision = v[2]
	if window.context.major == 0 {
		return fmt.Errorf("No version found in OpenGL version string")
	}
	if window.context.major < ctxconfig.major || window.context.major == ctxconfig.major && window.context.minor < ctxconfig.minor {
		// The desired OpenGL version is greater than the actual version
		// This only happens if the machine lacks {GLX|WGL}_ARB_create_context
		// /and/ the user has requested an OpenGL version greater than 1.0
		if window.context.client == OpenGLAPI {
			return fmt.Errorf("Requested OpenGL version %d.%d, got version %d.%d", ctxconfig.major, ctxconfig.minor, window.context.major, window.context.minor)
		} else {
			return fmt.Errorf("Requested OpenGL ES version %d.%d, got version %d.%d", ctxconfig.major, ctxconfig.minor, window.context.major, window.context.minor)
		}
		// makeContextCurrentWGL(previous)
	}
	if window.context.major >= 3 {
		// OpenGL 3.0+ uses a different function for extension string retrieval
		// We cache it here instead of in glfwExtensionSupported mostly to alert
		// users as early as possible that their build may be broken
		window.context.GetStringi = window.context.getProcAddress("glGetStringi")
		if window.context.GetStringi == 0 {
			return fmt.Errorf("Entry point retrieval is broken for glGetStringi")
		}
	}
	if window.context.client == OpenGLAPI {
		// Read back context flags (OpenGL 3.0 and above)
		if window.context.major >= 3 {
			var flags int
			getIntegerv(window, _GL_CONTEXT_FLAGS, &flags)
			if flags&_GL_CONTEXT_FLAG_FORWARD_COMPATIBLE_BIT != 0 {
				window.context.forward = true
			}
			if (flags & _GL_CONTEXT_FLAG_DEBUG_BIT) != 0 {
				window.context.debug = true
			} else if extensionSupported("GL_ARB_debug_output") && ctxconfig.debug {
				// HACK: This is a workaround for older drivers (pre KHR_debug)
				//       not setting the debug bit in the context flags for
				//       debug contexts
				window.context.debug = true
			}
			if flags&_GL_CONTEXT_FLAG_NO_ERROR_BIT_KHR != 0 {
				window.context.noerror = true
			}
		}
		// Read back OpenGL context profile (OpenGL 3.2 and above)
		if window.context.major >= 4 || window.context.major == 3 && window.context.minor >= 2 {
			var mask int
			getIntegerv(window, _GL_CONTEXT_PROFILE_MASK, &mask)
			if (mask & _GL_CONTEXT_COMPATIBILITY_PROFILE_BIT) != 0 {
				window.context.profile = OpenGLCompatProfile
			} else if (mask & _GL_CONTEXT_CORE_PROFILE_BIT) != 0 {
				window.context.profile = OpenGLCoreProfile
			} else if extensionSupported("GL_ARB_compatibility") {
				// HACK: This is a workaround for the compatibility profile bit
				//       not being set in the context flags if an OpenGL 3.2+
				//       context was created without having requested a specific
				//       version
				window.context.profile = OpenGLCompatProfile
			}
		}
		// Read back robustness strategy
		if extensionSupported("GL_ARB_robustness") {
			// NOTE: We avoid using the context flags for detection, as they are
			//       only present from 3.0 while the extension applies from 1.1
			var strategy int
			getIntegerv(window, _GL_RESET_NOTIFICATION_STRATEGY_ARB, &strategy)
			if strategy == _GL_LOSE_CONTEXT_ON_RESET_ARB {
				window.context.robustness = LoseContextOnReset
			} else if strategy == _GL_NO_RESET_NOTIFICATION_ARB {
				window.context.robustness = NoResetNotification
			}
		}
	} else {
		// Read back robustness strategy
		if extensionSupported("GL_EXT_robustness") {
			// NOTE: The values of these constants match those of the OpenGL ARB one, so we can reuse them here
			var strategy int
			getIntegerv(window, _GL_RESET_NOTIFICATION_STRATEGY_ARB, &strategy)
			if strategy == _GL_LOSE_CONTEXT_ON_RESET_ARB {
				window.context.robustness = LoseContextOnReset
			} else if strategy == _GL_NO_RESET_NOTIFICATION_ARB {
				window.context.robustness = NoResetNotification
			}
		}
	}

	if extensionSupported("GL_KHR_context_flush_control") {
		var behavior int
		getIntegerv(window, _GL_CONTEXT_RELEASE_BEHAVIOR, &behavior)
		if behavior == 0 {
			window.context.release = ReleaseBehaviorNone
		} else if behavior == _GL_CONTEXT_RELEASE_BEHAVIOR_FLUSH {
			window.context.release = ReleaseBehaviorFlush
		}
	}
	// Clearing the front buffer to black to avoid garbage pixels left over from
	// previous uses of our bit of VRAM
	glClear := window.context.getProcAddress("glClear")
	syscall.SyscallN(glClear, _GL_COLOR_BUFFER_BIT)
	if window.doublebuffer {
		window.context.swapBuffers(window)
	}
	glfwMakeContextCurrent(previous)
	return nil
}

func getIntegerv(window *Window, name int, value *int) {
	_, _, err := syscall.SyscallN(window.context.GetIntegerv, uintptr(name), uintptr(unsafe.Pointer(value)))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("getIntegerv failed, " + err.Error())
	}
}

func glfwMakeContextCurrent(window *_GLFWwindow) error {
	previous := (*Window)(unsafe.Pointer(glfwPlatformGetTls(&_glfw.contextSlot)))
	if window != nil && window.context.client == NoAPI {
		return fmt.Errorf("glfwMakeContextCurrent failed: Cannot make current with a window that has no OpenGL or OpenGL ES context")
	}
	if previous != nil && (window == nil || window.context.source != previous.context.source) {
		previous.context.makeCurrent(nil)
	}
	if window != nil {
		window.context.makeCurrent(window)
	}
	return nil
}

func glfwGetCurrentContext() *Window {
	p := glfwPlatformGetTls(&_glfw.contextSlot)
	return (*Window)(unsafe.Pointer(p))
}

func glfwSwapBuffers(window *_GLFWwindow) {
	if window == nil {
		panic("glfwSwapBuffers: window == nil")
	}
	window.context.swapBuffers(window)
}

func extensionSupported(extension string) bool {
	p := glfwPlatformGetTls(&_glfw.contextSlot)
	if p == 0 {
		return false
	}
	if extension == "" {
		return false
	}
	window := (*_GLFWwindow)(unsafe.Pointer(p))
	if window.context.major >= 3 {
		// Check if extension is in the modern OpenGL extensions string list
		// count := window.context.getIntegerv(_GL_NUM_EXTENSIONS)
		r, _, _ := syscall.SyscallN(window.context.GetIntegerv, uintptr(_GL_NUM_EXTENSIONS))
		count := int(r)
		for i := 0; i < count; i++ {
			// en := window.context.GetStringi(_GL_EXTENSIONS, i)
			r, _, _ := syscall.SyscallN(window.context.GetStringi, uintptr(_GL_EXTENSIONS), uintptr(i))
			en := GoStr((*uint8)(unsafe.Pointer(r)))
			if en == extension {
				return true
			}
		}
	} else {
		// Check if extension is in the old style OpenGL extensions string
		// extensions := window.context.GetString(_GL_EXTENSIONS)
		r, _, _ := syscall.SyscallN(window.context.GetStringi, uintptr(_GL_EXTENSIONS))
		extensions := GoStr((*uint8)(unsafe.Pointer(r)))
		if strings.Contains(extensions, extension) {
			return true
		}
	}
	// Check if extension is in the platform-specific string
	return window.context.extensionSupported(extension)
}

func glfwGetProcAddress(procname string) *Window {
	p := glfwPlatformGetTls(&_glfw.contextSlot)
	window := (*Window)(unsafe.Pointer(p))
	if window == nil {
		panic("glfwGetProcAddress: window == nil")
	}
	p = window.context.getProcAddress(procname)
	return (*Window)(unsafe.Pointer(p))
}
