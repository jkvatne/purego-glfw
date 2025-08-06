package glfw

import (
	"errors"
	"golang.org/x/sys/windows"
	"sync"
	"syscall"
	"unicode"
	"unicode/utf16"
	"unsafe"
)

type GLFWvidmode struct {
	Width       int32
	Height      int32
	RedBits     int32
	GreenBits   int32
	BlueBits    int32
	RefreshRate int32
}

type (
	_GLFWmakecontextcurrentfun = func(w *Window) error
	_GLFWswapbuffersfun        = func(w *Window)
	_GLFWswapintervalfun       = func(n int)
	_GLFWextensionsupportedfun = func(s string) bool
	_GLFWgetprocaddressfun     = func(s string) uintptr
	_GLFWdestroycontextfun     = func(w *Window)
)

// Context structure
type _GLFWcontext struct {
	client                  int32
	source                  int32
	major, minor, revision  int32
	forward, debug, noerror bool
	profile                 int32
	robustness              int32
	release                 int32
	GetStringi              uintptr
	GetIntegerv             uintptr
	GetString               uintptr
	makeCurrent             _GLFWmakecontextcurrentfun
	swapBuffers             _GLFWswapbuffersfun
	swapInterval            _GLFWswapintervalfun
	extensionSupported      _GLFWextensionsupportedfun
	getProcAddress          _GLFWgetprocaddressfun
	destroy                 _GLFWdestroycontextfun
	wgl                     struct {
		dc       HDC
		handle   HANDLE
		hMonitor HANDLE
		interval int
	}
}

type _GLFWwindow struct {
	next *_GLFWwindow
	// Window settings and state
	resizable              bool
	decorated              bool
	autoIconify            bool
	floating               bool
	focusOnShow            bool
	shouldClose            bool
	userPointer            unsafe.Pointer
	doublebuffer           bool
	videoMode              GLFWvidmode
	monitor                *Monitor
	cursor                 *Cursor
	minwidth               int
	minheight              int
	maxwidth               int
	maxheight              int
	numer                  int
	denom                  int
	stickyKeys             bool
	stickyMouseButtons     bool
	lockKeyMods            bool
	cursorMode             int
	mouseButtons           [MouseButtonLast + 1]byte
	keys                   [KeyLast + 1]byte
	virtualCursorPosX      float64 // Virtual cursor position when cursor is disabled
	virtualCursorPosY      float64 // Virtual cursor position when cursor is disabled
	context                *_GLFWcontext
	lastCursorPosX         float64 // The last received cursor position, regardless of source
	lastCursorPosY         float64 // The last received cursor position, regardless of source
	charCallback           CharCallback
	focusCallback          FocusCallback
	keyCallback            KeyCallback
	mouseButtonCallback    MouseButtonCallback
	cursorPosCallback      CursorPosCallback
	scrollCallback         ScrollCallback
	refreshCallback        RefreshCallback
	sizeCallback           SizeCallback
	dropCallback           DropCallback
	contentScaleCallback   ContentScaleCallback
	windowCloseCallback    func(w *_GLFWwindow)
	fFramebufferSizeHolder func(w *_GLFWwindow, width int, height int)
	fCloseHolder           func(w *_GLFWwindow)
	fMaximizeHolder        func(w *_GLFWwindow, maximized bool)
	fIconifyHolder         func(w *_GLFWwindow, iconified bool)
	fCursorEnterHolder     func(w *_GLFWwindow, entered bool)
	fCharModsHolder        func(w *_GLFWwindow, char rune, mods ModifierKey)
	Win32                  _GLFWwindowWin32
}

type _GLFWwindowWin32 = struct {
	handle         syscall.Handle
	bigIcon        syscall.Handle
	smallIcon      syscall.Handle
	cursorTracked  bool
	frameAction    bool
	iconified      bool
	maximized      bool
	transparent    bool // Whether to enable framebuffer transparency on DWM
	scaleToMonitor bool
	width          int    // Cached size used to filter out duplicate events
	height         int    // Cached size used to filter out duplicate events
	highSurrogate  uint16 // The last recevied high surrogate when decoding pairs of UTF-16 messages
}

type _GLFWinitconfig = struct {
	hatButtons bool
	ns         struct {
		menubar bool
		chdir   bool
	}
	wl struct {
		libdecorMode int
	}
}
type _GLFWwndconfig = struct {
	xpos             int32
	ypos             int32
	width            int32
	height           int32
	title            string
	resizable        bool
	visible          bool
	decorated        bool
	focused          bool
	autoIconify      bool
	floating         bool
	maximized        bool
	centerCursor     bool
	focusOnShow      bool
	mousePassthrough bool
	scaleToMonitor   bool
	scaleFramebuffer bool
	ns               struct {
		retina    bool
		frameName string
	}
}

type _GLFWctxconfig = struct {
	client     int32
	source     int32
	major      int32
	minor      int32
	forward    bool
	debug      bool
	noerror    bool
	profile    int32
	robustness int32
	release    int32
	share      *_GLFWwindow
	nsgl       struct {
		offline bool
	}
}

type hints = struct {
	init        _GLFWinitconfig
	framebuffer _GLFWfbconfig
	window      _GLFWwndconfig
	context     _GLFWctxconfig
	refreshRate int32
}

type _GLFWfbconfig = struct {
	redBits        int32
	greenBits      int32
	blueBits       int32
	alphaBits      int32
	depthBits      int32
	stencilBits    int32
	accumRedBits   int32
	accumGreenBits int32
	accumBlueBits  int32
	accumAlphaBits int32
	auxBuffers     int32
	samples        int32
	sRGB           bool
	doublebuffer   bool
	transparent    bool
	handle         uintptr
}

type _GLFWerror struct {
	next        *_GLFWerror
	code        int
	description string
}

type _GLFWtls = struct {
	allocated bool
	index     int
}

// Library global Data
var _glfw struct {
	hints
	class           uint16
	available       bool
	instance        syscall.Handle
	initialized     bool
	errorListHead   *_GLFWerror
	cursorListHead  *Cursor
	windowListHead  *_GLFWwindow
	monitors        []*Monitor
	monitorCallback func(w *Monitor, action int)
	errorCallback   ErrorCallbackFunc
	monitorCount    int
	errorSlot       _GLFWtls
	contextSlot     _GLFWtls
	errorLock       sync.Mutex
	win32           struct {
		helperWindowHandle   syscall.Handle
		helperWindowClass    uint16
		mainWindowClass      uint16
		blankCursor          HANDLE
		keycodes             [512]Key
		scancodes            [512]int16
		instance             syscall.Handle
		acquiredMonitorCount int
		mouseTrailSize       uint32
	}
	wgl struct {
		dc       HDC
		handle   syscall.Handle
		interval int
		// _GLFWlibraryWGL
		instance                       *windows.LazyDLL
		wglCreateContext               *windows.LazyProc
		wglDeleteContext               *windows.LazyProc
		wglGetProcAddress              *windows.LazyProc
		wglGetCurrentDC                *windows.LazyProc
		wglGetCurrentContext           *windows.LazyProc
		wglMakeCurrent                 *windows.LazyProc
		wglShareLists                  *windows.LazyProc
		SwapIntervalEXT                uintptr
		GetPixelFormatAttribivARB      uintptr
		GetExtensionsStringEXT         uintptr
		GetExtensionsStringARB         uintptr
		wglCreateContextAttribsARB     uintptr
		EXT_swap_control               bool
		EXT_colorspace                 bool
		ARB_multisample                bool
		ARB_framebuffer_sRGB           bool
		EXT_framebuffer_sRGB           bool
		ARB_pixel_format               bool
		ARB_create_context             bool
		ARB_create_context_profile     bool
		EXT_create_context_es2_profile bool
		ARB_create_context_robustness  bool
		ARB_create_context_no_error    bool
		ARB_context_flush_control      bool
	}
}

func glfwInputKey(window *_GLFWwindow, key Key, scancode int, action int, mods ModifierKey) {
	var repeated bool
	if key >= 0 && key <= KeyLast {
		repeated = false

		if action == glfw_RELEASE && window.keys[key] == glfw_RELEASE {
			return
		}

		if action == glfw_PRESS && window.keys[key] == glfw_PRESS {
			repeated = true
		}

		if action == glfw_RELEASE && window.stickyKeys {
			window.keys[key] = glfw_STICK
		} else {
			window.keys[key] = uint8(action)
		}
		if repeated {
			action = glfw_REPEAT
		}
	}
	if !window.lockKeyMods {
		mods &= ^(glfw_MOD_CAPS_LOCK | glfw_MOD_NUM_LOCK)
	}

	if window.keyCallback != nil {
		window.keyCallback(window, key, scancode, Action(action), mods)
	}
}

func glfwInputMouseClick(window *_GLFWwindow, button MouseButton, action Action, mods ModifierKey) {
	// TODO if (!window.lockKeyMods)	mods &= ~(glfw_MOD_CAPS_LOCK | glfw_MOD_NUM_LOCK);
	// TODO if (action == glfw_RELEASE && window.stickyMouseButtons) window.mouseButtons[button] = glfw_STICK; else window.mouseButtons[button] = (char) action;
	if window.mouseButtonCallback != nil {
		window.mouseButtonCallback(window, button, action, mods)
	}
}

// Notifies shared code that a window has lost or received input focus
func glfwInputWindowFocus(window *_GLFWwindow, focused bool) {
	if window == nil {
		return
	}
	if window.focusCallback != nil {
		window.focusCallback(window, focused)
	}
	if !focused {
		// Force release of buttons
		/* TODO
		for k := Key(0);  k <= KeyLast;  k++ {
			if (window.keys[k] == glfw_PRESS) {
				scancode := glfwPlatformGetKeyScancode(k);
				glfwInputKey(window, k, scancode, glfw_RELEASE, 0);
			}
		}*/
		for button := MouseButton(0); button <= MouseButtonLast; button++ {
			if window.mouseButtons[button] == glfw_PRESS {
				glfwInputMouseClick(window, button, glfw_RELEASE, 0)
			}
		}
	}
}

func glfwInputCursorPos(window *_GLFWwindow, xpos, ypos float64) {
	if window.virtualCursorPosX == xpos && window.virtualCursorPosY == ypos {
		return
	}
	window.virtualCursorPosX = xpos
	window.virtualCursorPosY = ypos
	if window.cursorPosCallback != nil {
		window.cursorPosCallback(window, xpos, ypos)
	}
}

func glfwInputScroll(window *_GLFWwindow, xoffset, yoffset float64) {
	if window.scrollCallback != nil {
		window.scrollCallback(window, xoffset, yoffset)
	}
}

func glfwInputWindowDamage(window *_GLFWwindow) {
	if window.refreshCallback != nil {
		window.refreshCallback(window)
	}
}

func glfwInputWindowCloseRequest(window *_GLFWwindow) {
	// slog.Error("Got CloseRequest")
}

func getKeyMods() ModifierKey {
	var mods ModifierKey
	if GetKeyState(VK_SHIFT)&0x8000 != 0 {
		mods |= ModShift
	}
	if GetKeyState(VK_CONTROL)&0x8000 != 0 {
		mods |= ModControl
	}
	if GetKeyState(VK_MENU)&0x8000 != 0 {
		mods |= ModAlt
	}
	if (GetKeyState(VK_LWIN)|GetKeyState(VK_RWIN))&0x8000 != 0 {
		mods |= ModSuper
	}
	if (GetKeyState(VK_CAPITAL) & 1) != 0 {
		mods |= ModCapsLock
	}
	if (GetKeyState(VK_NUMLOCK) & 1) != 0 {
		mods |= ModNumLock
	}
	return mods
}

func windowProc(hwnd syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	window := (*Window)(unsafe.Pointer(GetProp(HANDLE(hwnd), "GLFW")))
	if window == nil {
		r1, _, _ := _DefWindowProc.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
		return r1
	}

	switch msg {
	case wm_CLOSE:
		window.shouldClose = true
	case wm_UNICHAR:
		if wParam == _UNICODE_NOCHAR {
			// Tell the system that we accept wm_UNICHAR messages.
			return _TRUE
		}
		fallthrough
	case wm_CHAR, wm_SYSCHAR:
		if r := rune(wParam); unicode.IsPrint(r) {
			window.charCallback(nil, r)
		}
		return _TRUE
	case wm_DPICHANGED:
		// Let Windows know we're prepared for runtime DPI changes.
		return _TRUE
	case wm_ERASEBKGND:
		// Avoid flickering between GPU content and background color.
		return _TRUE
	case wm_KEYDOWN, wm_KEYUP, wm_SYSKEYDOWN, wm_SYSKEYUP:
		var key Key
		action := glfw_PRESS
		if (lParam>>16)&0x8000 != 0 {
			action = glfw_RELEASE
		}
		mods := getKeyMods()
		scancode := int((lParam >> 16) & 0x1ff)
		switch scancode {
		case 0: // scancode = MapVirtualKeyW((UINT) wParam, MAPVK_VK_TO_VSC);
		case 0x54:
			scancode = 0x137 // Alt+PrtSc
		case 0x146:
			scancode = 0x45 // Ctrl+Pause
		case 0x136:
			scancode = 0x36 // CJK IME sets the extended bit for right Shift
		}

		key = _glfw.win32.keycodes[scancode]
		if wParam == VK_CONTROL {
			if lParam>>16&kf_EXTENDED != 0 {
				// Right side keys have the extended key bit set
				key = KeyRightControl
			} else {
				/*
					// NOTE: Alt Gr sends Left Ctrl followed by Right Alt
					// HACK: We only want one event for Alt Gr, so if we detect
					//       this sequence we discard this Left Ctrl message now
					//       and later report Right Alt normally
					MSG next;
					const DWORD time = GetMessageTime();

					if (PeekMessageW(&next, NULL, 0, 0, pm_NOREMOVE)) {
						if (next.message == wm_KEYDOWN ||
							next.message == wm_SYSKEYDOWN ||
							next.message == wm_KEYUP ||
							next.message == wm_SYSKEYUP)
						{
							if (next.wParam == VK_MENU &&
								(HIWORD(next.lParam) & KF_EXTENDED) &&
								next.time == time)
							{
								// Next message is Right Alt down so discard this
								break;
							}
						}
					}
				*/
				// This is a regular Left Ctrl message
				key = KeyLeftControl
			}
		}

		if action == glfw_RELEASE && wParam == VK_SHIFT {
			// HACK: Release both Shift keys on Shift up event, as when both
			//       are pressed the first release does not emit any event
			// NOTE: The other half of this is in _glfwPlatformPollEvents
			glfwInputKey(window, KeyLeftShift, scancode, action, mods)
			glfwInputKey(window, KeyRightShift, scancode, action, mods)
		} else if wParam == VK_SNAPSHOT {
			// HACK: Key down is not reported for the Print Screen key
			glfwInputKey(window, key, scancode, glfw_PRESS, mods)
			glfwInputKey(window, key, scancode, glfw_RELEASE, mods)
		} else {
			glfwInputKey(window, key, scancode, action, mods)
		}
		break

	case wm_LBUTTONDOWN, wm_LBUTTONUP, wm_RBUTTONDOWN, wm_RBUTTONUP, wm_MBUTTONDOWN, wm_MBUTTONUP:
		var button MouseButton
		if msg == wm_LBUTTONDOWN || msg == wm_LBUTTONUP {
			button = MouseButtonLeft
		} else if msg == wm_RBUTTONDOWN || msg == wm_RBUTTONUP {
			button = MouseButtonRight
		} else if msg == wm_MBUTTONDOWN || msg == wm_MBUTTONUP {
			button = MouseButtonMiddle
		}
		var action Action
		if msg == wm_LBUTTONDOWN || msg == wm_RBUTTONDOWN || msg == wm_MBUTTONDOWN {
			action = glfw_PRESS
		} else {
			action = glfw_RELEASE
		}
		var i MouseButton
		for i = MouseButtonFirst; i <= MouseButtonLast; i++ {
			if window.mouseButtons[i] == glfw_PRESS {
				break
			}
		}
		// if i > MouseButtonLast {
		// TODO SetCapture(hWnd);
		// }

		glfwInputMouseClick(window, button, action, getKeyMods())
		for i = MouseButtonFirst; i <= MouseButtonLast; i++ {
			if window.mouseButtons[i] == glfw_PRESS {
				break
			}
		}
		// if (i > MouseButtonLast)
		// TODO ReleaseCapture();
		// }

		return 0

	// TODO case wm_CANCELMODE:

	case wm_SETFOCUS:
		glfwInputWindowFocus(window, true)
		// HACK: Do not disable cursor while the user is interacting with a caption button
		// TODO if (window.Win32.frameAction) break;
		// TODO if (window.cursorMode == glfw_CURSOR_DISABLED)	disableCursor(window);
		return 0
	case wm_KILLFOCUS:
		// TODO if (window.cursorMode == glfw_CURSOR_DISABLED) enableCursor(window);
		// TODO if (window.monitor && window.autoIconify) _glfwPlatformIconifyWindow(window);
		glfwInputWindowFocus(window, false)
		return 0

	case wm_MOUSEMOVE:
		x := float64(int(lParam & 0xffff))
		y := float64(int((lParam >> 16) & 0xffff))
		if !window.Win32.cursorTracked {
			// tme.dwFlags = TME_LEAVE;
			// tme.hwndTrack = window.hMonitor;
			// TrackMouseEvent(&tme);
			// window.cursorTracked = true;
			// glfwInputCursorEnter(window, glfw_TRUE);
		}

		if window.cursorMode == glfw_CURSOR_DISABLED {
			dx := float64(x) - window.lastCursorPosX
			dy := float64(y) - window.lastCursorPosY
			// TODO if _glfw.Win32.disabledCursorWindow != window {			break			}
			glfwInputCursorPos(window, window.virtualCursorPosX+dx, window.virtualCursorPosY+dy)
		} else {
			glfwInputCursorPos(window, x, y)
		}
		window.lastCursorPosX = x
		window.lastCursorPosY = y
		return 0

	case wm_MOUSEWHEEL:
		glfwInputScroll(window, 0.0, float64(int16(wParam>>16))/120.0)
		return 0

	case wm_MOUSEHWHEEL:
		glfwInputScroll(window, -float64(int16(wParam>>16))/120.0, 0.0)
		return 0

	case wm_PAINT:
		glfwInputWindowDamage(window)

	case wm_SIZE:
		width := int(lParam & 0xFFFF)
		height := int(lParam >> 16)
		iconified := wParam == _SIZE_MINIMIZED
		maximized := wParam == _SIZE_MAXIMIZED || (window.Win32.maximized && wParam != _SIZE_RESTORED)
		// if (_glfw.win32.capturedCursorWindow == window) {
		//	captureCursor(window)
		// }

		if window.Win32.iconified != iconified {
			// TODO _glfwInputWindowIconify(window, iconified)
		}

		if window.Win32.maximized != maximized {
			// TODO _glfwInputWindowMaximize(window, maximized);
		}

		if width != window.Win32.width || height != window.Win32.height {
			window.Win32.width = width
			window.Win32.height = height
			// TODO _glfwInputFramebufferSize(window, width, height);
			if window.sizeCallback != nil {
				window.sizeCallback(window, width, height)
			}
		}
		if window.monitor != nil && window.Win32.iconified != iconified {
			if iconified {
				releaseMonitor(window)
			} else {
				acquireMonitor(window)
				fitToMonitor(window)
			}
		}
		window.Win32.iconified = iconified
		window.Win32.maximized = maximized
		return 0

	case wm_GETMINMAXINFO:
		// TODO

	case wm_SETCURSOR:
		// TODO
	}

	r1, _, _ := _DefWindowProc.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
	return r1
}

func glfwPlatformPollEvents() {
	var msg Msg
	for PeekMessage(&msg, 0, 0, 0, pm_REMOVE) {
		if msg.Message == wm_QUIT {
			// NOTE: While GLFW does not itself post wm_QUIT, other processes may post it to this one, for example Task Manager
			// HACK: Treat wm_QUIT as a close on all windows
			window := _glfw.windowListHead
			for window != nil {
				glfwInputWindowCloseRequest(window)
				window = window.next
			}
		} else {
			TranslateMessage(&msg)
			DispatchMessage(&msg)
		}
	}

	// HACK: Release modifier keys that the system did not emit KEYUP for
	// NOTE: Shift keys on Windows tend to "stick" when both are pressed as no key up message is generated by the first key release
	// NOTE: Windows key is not reported as released by the Win+V hotkey. Other Win hotkeys are handled implicitly by _glfwInputWindowFocus
	//       because they change the input focus
	// NOTE: The other half of this is in the wm_*KEY* handler in windowProc
	/* TODO
	hMonitor = GetActiveWindow();
	if (hMonitor!=nil) {
		window := 74W(hMonitor, "GLFW");
		if window != nil {
			//const keys[4][2] = {{ VK_LSHIFT, glfw_KEY_LEFT_SHIFT },    { VK_RSHIFT, glfw_KEY_RIGHT_SHIFT },{ VK_LWIN, glfw_KEY_LEFT_SUPER },{ VK_RWIN, glfw_KEY_RIGHT_SUPER }}
			for i := 0; i < 4; i++ {
				vk := keys[i][0];
				key := keys[i][1];
				// scancode := Win32.scancodes[key];
				if GetKeyState(vk) & 0x8000 || window.keys[key] != glfw_PRESS {
					continue;
				}
				_glfwInputKey(window, key, scancode, glfw_RELEASE, getKeyMods());
			}
		}
	}
	window := _glfw.Win32.disabledCursorWindow;
	if window!=nil {
		var width, height int
		glfwPlatformGetWindowSize(window, &width, &height);
		// NOTE: Re-center the cursor only if it has moved since the last call, to avoid breaking glfwWaitEvents with wm_MOUSEMOVE
		if window.lastCursorPosX != width / 2 || window.lastCursorPosY != height / 2 {
			glfwPlatformSetCursorPos(window, width / 2, height / 2);
		}
	}
	*/
}

func cursorInContentArea(w *_GLFWwindow) bool {
	var x, y, width, height int32
	glfwGetCursorPos(w, &x, &y)
	glfwGetWindowSize(w, &width, &height)
	return x >= 0 && y >= 0 && x < width && y < height // PtInRect(&area, pos);
}

func setCursorWin32(handle HANDLE) {
	_, _, err := _SetCursor.Call(uintptr(handle))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("_SetCursor failed, " + err.Error())
	}
}

func updateCursorImage(window *_GLFWwindow) {
	// TODO
}

func glfwSetCursor(window *_GLFWwindow, cursor *Cursor) {
	window.cursor = cursor
	if cursorInContentArea(window) {
		if window.cursorMode == glfw_CURSOR_NORMAL || window.cursorMode == glfw_CURSOR_CAPTURED {
			if window.cursor != nil {
				setCursorWin32(window.cursor.handle)
			} else {
				setCursorWin32(LoadCursor(IDC_ARROW))
			}
		} else {
			// NOTE: Via Remote Desktop, setting the cursor to NULL does not hide it.
			// HACK: When running locally, it is set to NULL, but when connected via Remote
			//       Desktop, this is a transparent cursor.
			setCursorWin32(_glfw.win32.blankCursor)
		}
	}
}

func SetFocus(window *_GLFWwindow) {
	r1, _, err := _SetFocus.Call(uintptr(unsafe.Pointer(window.Win32.handle)))
	if r1 == 0 || err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("SetFocus failed, " + err.Error())
	}
	if r1 == 0 {
		panic("SetFocus failed")
	}
}

func BringWindowToTop(window *_GLFWwindow) {
	r1, _, err := _BringWindowToTop.Call(uintptr(unsafe.Pointer(window.Win32.handle)))
	if r1 == 0 || err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("BringWindowToTop failed, " + err.Error())
	}
	if r1 == 0 {
		panic("BringWindowToTop failed")
	}
}

func SetForegroundWindow(window *_GLFWwindow) {
	r1, _, err := _SetForegroundWindow.Call(uintptr(unsafe.Pointer(window.Win32.handle)))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("SetForegroundWindow failed, " + err.Error())
	}
	if r1 == 0 {
		panic("SetForegroundWindow failed, returned 0")
	}
}

func glfwFocusWindow(window *_GLFWwindow) {
	BringWindowToTop(window)
	SetForegroundWindow(window)
	SetFocus(window)
}

const (
	ENUM_CURRENT_SETTINGS      = -1
	HORZSIZE                   = 4
	VERTSIZE                   = 6
	HORZRES                    = 8
	VERTRES                    = 10
	DISPLAY_DEVICE_MODESPRUNED = 0x08000000
	DISPLAY_DEVICE_REMOTE      = 0x04000000
	DISPLAY_DEVICE_DISCONNECT  = 0x02000000
)

func createMonitor(adapter *DISPLAY_DEVICEW, display *DISPLAY_DEVICEW) *Monitor {
	var monitor Monitor
	var widthMM, heightMM int
	var rect RECT
	var dm DEVMODEW

	dm.dmSize = uint16(unsafe.Sizeof(dm))
	EnumDisplaySettingsEx(&adapter.DeviceName[0], ENUM_CURRENT_SETTINGS, &dm, 0)
	pName, _ := syscall.UTF16PtrFromString("DISPLAY")
	ret, _, err := _CreateDC.Call(uintptr(unsafe.Pointer(pName)), uintptr(unsafe.Pointer(&adapter.DeviceName)), 0, 0)
	if !errors.Is(err, syscall.Errno(0)) {
		panic("CreateDC failed, " + err.Error())
	}
	dc := HDC(ret)
	if IsWindows8Point1OrGreater() {
		widthMM = GetDeviceCaps(dc, HORZSIZE)
		heightMM = GetDeviceCaps(dc, VERTSIZE)
	} else {
		widthMM = int(float64(dm.dmPelsWidth) * 25.4 / float64(GetDeviceCaps(dc, _LOGPIXELSX)))
		heightMM = int(float64(dm.dmPelsHeight) * 25.4 / float64(GetDeviceCaps(dc, _LOGPIXELSY)))
	}
	ret, _, err = _DeleteDC.Call(uintptr(dc))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("CreateDC failed, " + err.Error())
	}
	monitor.heightMM = heightMM
	monitor.widthMM = widthMM

	if adapter.StateFlags&DISPLAY_DEVICE_MODESPRUNED != 0 {
		monitor.modesPruned = true
	}
	for i := 0; i < len(adapter.DeviceName); i++ {
		monitor.adapterName[i] = adapter.DeviceName[i]
	}
	// WideCharToMultiByte(CP_UTF8, 0, adapter.DeviceName, -1, monitor.win32.publicAdapterName, sizeof(monitor.win32.publicAdapterName), NULL, NULL)
	if display != nil {
		for i := 0; i < len(adapter.DeviceName); i++ {
			monitor.displayName[i] = display.DeviceName[i]
		}
	}
	//	WideCharToMultiByte(CP_UTF8, 0, display.DeviceName, -1, monitor.win32.publicDisplayName, sizeof(monitor.win32.publicDisplayName), NULL, NULL)
	monitor.publicDisplayName = string(utf16.Decode(display.DeviceName[:]))
	monitor.publicAdapterName = string(utf16.Decode(adapter.DeviceName[:]))
	rect.Left = dm.dmPosition.X
	rect.Top = dm.dmPosition.Y
	rect.Right = dm.dmPosition.X + dm.dmPelsWidth
	rect.Bottom = dm.dmPosition.Y + dm.dmPelsHeight
	_ = EnumDisplayMonitors(0, &rect, NewEnumDisplayMonitorsCallback(enumMonitorCallback), uintptr(unsafe.Pointer(&monitor)))
	return &monitor
}

const (
	/* Printable keys */
	glfw_KEY_SPACE         = 32
	glfw_KEY_APOSTROPHE    = 39 /* ' */
	glfw_KEY_COMMA         = 44 /* , */
	glfw_KEY_MINUS         = 45 /* - */
	glfw_KEY_PERIOD        = 46 /* . */
	glfw_KEY_SLASH         = 47 /* / */
	glfw_KEY_0             = 48
	glfw_KEY_1             = 49
	glfw_KEY_2             = 50
	glfw_KEY_3             = 51
	glfw_KEY_4             = 52
	glfw_KEY_5             = 53
	glfw_KEY_6             = 54
	glfw_KEY_7             = 55
	glfw_KEY_8             = 56
	glfw_KEY_9             = 57
	glfw_KEY_SEMICOLON     = 59 /* ; */
	glfw_KEY_EQUAL         = 61 /* = */
	glfw_KEY_A             = 65
	glfw_KEY_B             = 66
	glfw_KEY_C             = 67
	glfw_KEY_D             = 68
	glfw_KEY_E             = 69
	glfw_KEY_F             = 70
	glfw_KEY_G             = 71
	glfw_KEY_H             = 72
	glfw_KEY_I             = 73
	glfw_KEY_J             = 74
	glfw_KEY_K             = 75
	glfw_KEY_L             = 76
	glfw_KEY_M             = 77
	glfw_KEY_N             = 78
	glfw_KEY_O             = 79
	glfw_KEY_P             = 80
	glfw_KEY_Q             = 81
	glfw_KEY_R             = 82
	glfw_KEY_S             = 83
	glfw_KEY_T             = 84
	glfw_KEY_U             = 85
	glfw_KEY_V             = 86
	glfw_KEY_W             = 87
	glfw_KEY_X             = 88
	glfw_KEY_Y             = 89
	glfw_KEY_Z             = 90
	glfw_KEY_LEFT_BRACKET  = 91  /* [ */
	glfw_KEY_BACKSLASH     = 92  /* \ */
	glfw_KEY_RIGHT_BRACKET = 93  /* ] */
	glfw_KEY_GRAVE_ACCENT  = 96  /* ` */
	glfw_KEY_WORLD_1       = 161 /* non-US #1 */
	glfw_KEY_WORLD_2       = 162 /* non-US #2 */

	/* Function keys */
	glfw_KEY_ESCAPE        = 256
	glfw_KEY_ENTER         = 257
	glfw_KEY_TAB           = 258
	glfw_KEY_BACKSPACE     = 259
	glfw_KEY_INSERT        = 260
	glfw_KEY_DELETE        = 261
	glfw_KEY_RIGHT         = 262
	glfw_KEY_LEFT          = 263
	glfw_KEY_DOWN          = 264
	glfw_KEY_UP            = 265
	glfw_KEY_PAGE_UP       = 266
	glfw_KEY_PAGE_DOWN     = 267
	glfw_KEY_HOME          = 268
	glfw_KEY_END           = 269
	glfw_KEY_CAPS_LOCK     = 280
	glfw_KEY_SCROLL_LOCK   = 281
	glfw_KEY_NUM_LOCK      = 282
	glfw_KEY_PRINT_SCREEN  = 283
	glfw_KEY_PAUSE         = 284
	glfw_KEY_F1            = 290
	glfw_KEY_F2            = 291
	glfw_KEY_F3            = 292
	glfw_KEY_F4            = 293
	glfw_KEY_F5            = 294
	glfw_KEY_F6            = 295
	glfw_KEY_F7            = 296
	glfw_KEY_F8            = 297
	glfw_KEY_F9            = 298
	glfw_KEY_F10           = 299
	glfw_KEY_F11           = 300
	glfw_KEY_F12           = 301
	glfw_KEY_KP_0          = 320
	glfw_KEY_KP_1          = 321
	glfw_KEY_KP_2          = 322
	glfw_KEY_KP_3          = 323
	glfw_KEY_KP_4          = 324
	glfw_KEY_KP_5          = 325
	glfw_KEY_KP_6          = 326
	glfw_KEY_KP_7          = 327
	glfw_KEY_KP_8          = 328
	glfw_KEY_KP_9          = 329
	glfw_KEY_KP_DECIMAL    = 330
	glfw_KEY_KP_DIVIDE     = 331
	glfw_KEY_KP_MULTIPLY   = 332
	glfw_KEY_KP_SUBTRACT   = 333
	glfw_KEY_KP_ADD        = 334
	glfw_KEY_KP_ENTER      = 335
	glfw_KEY_KP_EQUAL      = 336
	glfw_KEY_LEFT_SHIFT    = 340
	glfw_KEY_LEFT_CONTROL  = 341
	glfw_KEY_LEFT_ALT      = 342
	glfw_KEY_LEFT_SUPER    = 343
	glfw_KEY_RIGHT_SHIFT   = 344
	glfw_KEY_RIGHT_CONTROL = 345
	glfw_KEY_RIGHT_ALT     = 346
	glfw_KEY_RIGHT_SUPER   = 347
	glfw_KEY_MENU          = 348
)

// func ToUnicode(vk uint32, scancode uint32, state *[512]byte , chars, len, 0) {
// r1,_,err := _ToUnicode.Call(uintptr(vk), uintptr(scancode), uintptr(state), uintptr(chars), size)
// }

// TODO :Updates key names according to the current keyboard layout
func glfwUpdateKeyNamesWin32() {
	for key := glfw_KEY_SPACE; key <= glfw_KEY_MENU; key++ {
		/* TODO: Make readable key names
		scancode := _glfw.win32.scancodes[key]
		var vk uint16
		if scancode == -1 {
			continue
		}
		if key >= glfw_KEY_KP_0 && key <= glfw_KEY_KP_ADD {
			vks := []uint16{VK_NUMPAD0, VK_NUMPAD1, VK_NUMPAD2, VK_NUMPAD3, VK_NUMPAD4, VK_NUMPAD5, VK_NUMPAD6, VK_NUMPAD7, VK_NUMPAD8, VK_NUMPAD9, VK_DECIMAL, VK_DIVIDE, VK_MULTIPLY, VK_SUBTRACT, VK_ADD}
			vk = vks[key-glfw_KEY_KP_0]
		} else {
			r1, _, err := _MapVirtualKeyW.Call(uintptr(scancode), uintptr(MAPVK_VSC_TO_VK))
			if !errors.Is(err, syscall.Errno(0)) {
				panic("MapVirtualKeyW failed, " + err.Error())
			}
			vk = uint16(r1)
		}
		var state [256]uint8
		var vk uint16
		length := ToUnicode(vk, scancode, state, chars, sizeof(chars)/sizeof(WCHAR), 0);
		if length == -1 {
			// This is a dead key, so we need a second simulated key press
			// to make it output its own character (usually a diacritic)
			length = ToUnicode(vk, scancode, state, chars, sizeof(chars)/sizeof(WCHAR), 0);
		}

		if (length < 1) {
			continue;
		}
		WideCharToMultiByte(CP_UTF8, 0, chars, 1, _glfw.win32.keynames[key], sizeof(_glfw.win32.keynames[key]), NULL, NULL);
		*/
	}
}

// Notifies shared code of a monitor connection or disconnection
func glfwInputMonitor(monitor *Monitor, action int, placement int) {
	if action == glfw_CONNECTED {
		_glfw.monitorCount++
		if placement == glfw_INSERT_FIRST {
			_glfw.monitors = append([]*Monitor{monitor}, _glfw.monitors...)
		} else {
			_glfw.monitors = append(_glfw.monitors, monitor)
		}
	} else if action == glfw_DISCONNECTED {
		for window := _glfw.windowListHead; window != nil; window = window.next {
			if window.monitor == monitor {
				// TODO var width, height, xoff, yoff int
				// glfwGetWindowSize(window, &width, &height);
				// _glfw.platform.setWindowMonitor(window, NULL, 0, 0, width, height, 0);
				// _glfw.platform.getWindowFrameSize(window, &xoff, &yoff, NULL, NULL);
				// _glfw.platform.SetWindowPos(window, xoff, yoff);
			}
		}
		for i := 0; i < _glfw.monitorCount; i++ {
			if _glfw.monitors[i] == monitor {
				_glfw.monitors = append(_glfw.monitors[:i], _glfw.monitors[i+1:]...)
				_glfw.monitorCount--
				break
			}
		}
	}

	if _glfw.monitorCallback != nil {
		_glfw.monitorCallback(monitor, action)
	}

}

// glfwInputMonitorWindow Notifies shared code that a full screen window has acquired or released a monitor
func glfwInputMonitorWindow(monitor *Monitor, window *_GLFWwindow) {
	monitor.window = window
}

// glfwInputWindowMonitor Notifies shared code that a full screen window has acquired or released a monitor
func glfwInputWindowMonitor(window *_GLFWwindow, monitor *Monitor) {
	window.monitor = monitor
}

// Retrieves the available modes for the specified monitor
func refreshVideoModes(monitor *Monitor) bool {
	var modes []GLFWvidmode

	if len(monitor.modes) != 0 {
		return true
	}
	// modes = monitor.GetVideoModes()
	if len(monitor.modes) == 0 {
		return false
	}
	// slices.SortFunc(modes, compareVideoModes)
	monitor.modes = modes
	return true
}

func abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

// Chooses the video mode most closely matching the desired one
// const GLFWvidmode* _glfwChooseVideoMode(_GLFWmonitor* monitor,const GLFWvidmode* desired)
func glfwChooseVideoMode(monitor *Monitor, desired *GLFWvidmode) *GLFWvidmode {
	var sizeDiff, leastSizeDiff int32 = INT_MAX, INT_MAX
	var rateDiff, leastRateDiff int32 = INT_MAX, INT_MAX
	var colorDiff, leastColorDiff int32 = INT_MAX, INT_MAX
	var current GLFWvidmode
	var closest *GLFWvidmode

	if !refreshVideoModes(monitor) {
		return nil
	}

	for i := 0; i < len(monitor.modes); i++ {
		current = monitor.modes[i]
		colorDiff = 0
		if desired.RedBits != glfw_DONT_CARE {
			colorDiff += abs(current.RedBits - desired.RedBits)
		}
		if desired.GreenBits != glfw_DONT_CARE {
			colorDiff += abs(current.GreenBits - desired.GreenBits)
		}
		if desired.BlueBits != glfw_DONT_CARE {
			colorDiff += abs(current.BlueBits - desired.BlueBits)
		}
		sizeDiff = abs((current.Width-desired.Width)*(current.Width-desired.Width) + (current.Height-desired.Height)*(current.Height-desired.Height))
		if desired.RefreshRate != glfw_DONT_CARE {
			rateDiff = abs(current.RefreshRate - desired.RefreshRate)
		} else {
			rateDiff = INT_MAX - current.RefreshRate
		}
		if (colorDiff < leastColorDiff) || (colorDiff == leastColorDiff && sizeDiff < leastSizeDiff) ||
			(colorDiff == leastColorDiff && sizeDiff == leastSizeDiff && rateDiff < leastRateDiff) {
			closest = &current
			leastSizeDiff = sizeDiff
			leastRateDiff = rateDiff
			leastColorDiff = colorDiff
		}
	}
	return closest
}

// Change the current video mode
//
/*
func glfwSetVideoModeWin32(monitor *Monitor, desired *GLFWvidmode) error {
	var current GLFWvidmode
	var best *GLFWvidmode
	var dm DEVMODEW
	var result int

	best = _glfwChooseVideoMode(monitor, *desired)
	//TODO _glfwPlatformGetVideoMode(monitor, &current)
	//TODO if _glfwCompareVideoModes(&current, best) == 0 {
	//	return nil
	//}

	dm.dmSize = uint16(unsafe.Sizeof(dm))
	dm.dmFields = DM_PELSWIDTH | DM_PELSHEIGHT | DM_BITSPERPEL | DM_DISPLAYFREQUENCY
	dm.dmPelsWidth = int32(best.Width)
	dm.dmPelsHeight = int32(best.Height)
	dm.dmBitsPerPel = uint32(best.RedBits + best.GreenBits + best.BlueBits)
	dm.dmDisplayFrequency = uint32(best.RefreshRate)

	if dm.dmBitsPerPel < 15 || dm.dmBitsPerPel >= 24 {
		dm.dmBitsPerPel = 32
	}
	result = ChangeDisplaySettingsExW(
		&monitor.adapterName,
		&dm,
		nil,
		cds_FULLSCREEN,
		0)
	if result == disp_CHANGE_SUCCESSFUL {
		monitor.modeChanged = true
	} else {
		description := "Unknown error"
		if result == disp_CHANGE_BADDUALVIEW {
			description = "The system uses DualView"
		} else if result == disp_CHANGE_BADFLAGS {
			description = "Invalid flags"
		} else if result == disp_CHANGE_BADMODE {
			description = "Graphics mode not supported"
		} else if result == disp_CHANGE_BADPARAM {
			description = "Invalid parameter"
		} else if result == disp_CHANGE_FAILED {
			description = "Graphics mode failed"
		} else if result == disp_CHANGE_NOTUPDATED {
			description = "Failed to write to registry"
		} else if result == disp_CHANGE_RESTART {
			description = "Computer restart required"
		}
		return fmt.Errorf("Win32: Failed to set video mode: %s", description)
	}
	return nil
}
*/

func fitToMonitor(window *Window) {
	mi := GetMonitorInfo(window.monitor.hMonitor)
	_, _, err := _SetWindowPos.Call(
		uintptr(window.Win32.handle),
		uintptr(hwnd_TOPMOST),
		uintptr(mi.RcMonitor.Left),
		uintptr(mi.RcMonitor.Top),
		uintptr(mi.RcMonitor.Right-mi.RcMonitor.Left),
		uintptr(mi.RcMonitor.Bottom-mi.RcMonitor.Top),
		uintptr(SWP_NOZORDER|SWP_NOACTIVATE|SWP_NOCOPYBITS))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("fitToMonitor failed, " + err.Error())
	}
}

func systemParametersInfoW(uiAction uint32, uiParam uint32, pvParam *uint32, fWinIni uint32) {
	_SystemParametersInfoW.Call(uintptr(uiAction), uintptr(uiParam), uintptr(unsafe.Pointer(pvParam)), uintptr(fWinIni))
}

// Make the specified window and its video mode active on its monitor
//
func acquireMonitor(window *Window) {
	if _glfw.win32.acquiredMonitorCount > 0 {
		_SetThreadExecutionState.Call(uintptr(es_CONTINUOUS | es_DISPLAY_REQUIRED))
		// HACK: When mouse trails are enabled the cursor becomes invisible when the OpenGL ICD switches to page flipping
		systemParametersInfoW(spi_GETMOUSETRAILS, 0, &_glfw.win32.mouseTrailSize, 0)
		systemParametersInfoW(spi_SETMOUSETRAILS, 0, nil, 0)
	}

	if window.monitor.window == nil {
		_glfw.win32.acquiredMonitorCount++
		// TODO _glfwSetVideoModeWin32(window.monitor, &window.videoMode)
		glfwInputMonitorWindow(window.monitor, window)
	}
}

// Remove the window and restore the original video mode
//
func releaseMonitor(window *Window) {
	if window.monitor.window != window {
		return
	}

	_glfw.win32.acquiredMonitorCount--
	if _glfw.win32.acquiredMonitorCount == 0 {
		_SetThreadExecutionState.Call(uintptr(es_CONTINUOUS))
		// HACK: Restore mouse trail length saved in acquireMonitor
		systemParametersInfoW(spi_SETMOUSETRAILS, _glfw.win32.mouseTrailSize, nil, 0)
	}
	glfwInputMonitorWindow(window.monitor, nil)
	// TODO _glfwRestoreVideoModeWin32(window.monitor)
}

func glfwSetPos(w *Window, xPos, yPos int) {
	rect := RECT{Left: int32(xPos), Top: int32(yPos), Right: int32(xPos), Bottom: int32(yPos)}
	AdjustWindowRect(&rect, getWindowStyle(w), 0, getWindowExStyle(w), GetDpiForWindow(w.Win32.handle), "glfwSetWindowPos")
	SetWindowPos(w.Win32.handle, 0, rect.Left, rect.Top, 0, 0, SWP_NOACTIVATE|SWP_NOZORDER|SWP_NOSIZE)
}
