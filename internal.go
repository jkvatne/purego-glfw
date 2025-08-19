package glfw

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"syscall"
	"unicode"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
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
	resizable               bool
	decorated               bool
	autoIconify             bool
	floating                bool
	focusOnShow             bool
	shouldClose             bool
	mousePassthrough        bool
	cursorTracked           bool
	userPointer             unsafe.Pointer
	doublebuffer            bool
	videoMode               GLFWvidmode
	monitor                 *Monitor
	cursor                  *Cursor
	minwidth                int32
	minheight               int32
	maxwidth                int32
	maxheight               int32
	numer                   int
	denom                   int
	stickyKeys              bool
	stickyMouseButtons      bool
	lockKeyMods             bool
	disableMouseButtonLimit bool
	cursorMode              int
	rawMouseMotion          int
	mouseButtons            [MouseButtonLast + 1]Action
	keys                    [KeyLast + 1]Action
	virtualCursorPosX       float64 // Virtual cursor position when cursor is disabled
	virtualCursorPosY       float64 // Virtual cursor position when cursor is disabled
	context                 *_GLFWcontext
	lastCursorPosX          float64 // The last received cursor position, regardless of source
	lastCursorPosY          float64 // The last received cursor position, regardless of source
	charCallback            CharCallback
	focusCallback           FocusCallback
	keyCallback             KeyCallback
	mouseButtonCallback     MouseButtonCallback
	cursorPosCallback       CursorPosCallback
	scrollCallback          ScrollCallback
	refreshCallback         RefreshCallback
	sizeCallback            SizeCallback
	dropCallback            DropCallback
	iconifyCallback         IconifyCallback
	framebufferSizeCallback SizeCallback
	contentScaleCallback    ContentScaleCallback
	cursorEnterCallback     CursorEnterCallback
	maximizeCallback        MaximizeCallback
	windowCloseCallback     func(w *_GLFWwindow)
	fFramebufferSizeHolder  func(w *_GLFWwindow, width int, height int)
	fCloseHolder            func(w *_GLFWwindow)
	fMaximizeHolder         func(w *_GLFWwindow, maximized bool)
	fIconifyHolder          func(w *_GLFWwindow, iconified bool)
	fCursorEnterHolder      func(w *_GLFWwindow, entered bool)
	fCharModsHolder         func(w *_GLFWwindow, char rune, mods ModifierKey)
	Win32                   _GLFWwindowWin32
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
	keyMenu        bool
	showDefault    bool
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
	win32            struct {
		keymenu     bool
		showDefault bool
	}
	ns struct {
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
		blankCursor          syscall.Handle
		keycodes             [512]Key
		scancodes            [512]int16
		instance             syscall.Handle
		acquiredMonitorCount int
		mouseTrailSize       uint32
		restoreCursorPosX    float64
		restoreCursorPosY    float64
		disabledCursorWindow *Window
		capturedCursorWindow *Window
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

type MINMAXINFO struct {
	ptReserved     POINT
	ptMaxSize      POINT
	ptMaxPosition  POINT
	ptMinTrackSize POINT
	ptMaxTrackSize POINT
}

func glfwIconifyWindow(w *Window) {
	ShowWindow(w.Win32.handle, windows.SW_MINIMIZE)
}

func glfwInputKey(window *_GLFWwindow, key Key, scancode int, action Action, mods ModifierKey) {
	var repeated bool
	if key >= 0 && key <= KeyLast {
		repeated = false

		if action == Release && window.keys[key] == Release {
			return
		}

		if action == Press && window.keys[key] == Press {
			repeated = true
		}

		if action == Release && window.stickyKeys {
			window.keys[key] = Stick
		} else {
			window.keys[key] = action
		}
		if repeated {
			action = Repeat
		}
	}
	if !window.lockKeyMods {
		mods &= ^(ModCapsLock | ModNumLock)
	}

	if window.keyCallback != nil {
		window.keyCallback(window, key, scancode, Action(action), mods)
	}
}

// Apply disabled cursor mode to a focused window
//
func disableCursor(window *Window) {
	_glfw.win32.disabledCursorWindow = window
	_glfw.win32.restoreCursorPosX, _glfw.win32.restoreCursorPosY = window.GetCursorPos()
	updateCursorImage(window)
	var width, height int32
	glfwGetWindowSize(window, &width, &height)
	window.SetCursorPos(float64(width)/2, float64(height)/2)
	captureCursor(window)
	if window.rawMouseMotion != 0 {
		enableRawMouseMotion(window)
	}
}

// Exit disabled cursor mode for the specified window
//
func enableCursor(window *Window) {
	if window.rawMouseMotion != 0 {
		disableRawMouseMotion(window)
	}
	_glfw.win32.disabledCursorWindow = nil
	releaseCursor()
	window.SetCursorPos(_glfw.win32.restoreCursorPosX, _glfw.win32.restoreCursorPosY)
	updateCursorImage(window)
}

func glfwInputMouseClick(window *_GLFWwindow, button MouseButton, action Action, mods ModifierKey) {
	if button < 0 || !window.disableMouseButtonLimit && button > MouseButtonLast {
		return
	}
	if !window.lockKeyMods {
		mods &= ^(ModCapsLock | ModNumLock)
	}
	if action == Release && window.stickyMouseButtons {
		window.mouseButtons[button] = Stick
	} else {
		window.mouseButtons[button] = action
	}
	if window.mouseButtonCallback != nil {
		window.mouseButtonCallback(window, button, action, mods)
	}
}

func glfwGetKeyScancode(key Key) int {
	if key < KeySpace || key > KeyLast {
		panic("glfwGetKeyScancode: Invalid key " + strconv.Itoa(int(key)))
	}
	return int(_glfw.win32.scancodes[key])
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
		for k := Key(0); k <= KeyLast; k++ {
			if window.keys[k] == Press {
				scancode := glfwGetKeyScancode(k)
				glfwInputKey(window, k, scancode, Release, 0)
			}
		}
		for button := MouseButton(0); button <= MouseButtonLast; button++ {
			if window.mouseButtons[button] == Press {
				glfwInputMouseClick(window, button, Release, 0)
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
	window := (*Window)(unsafe.Pointer(GetProp(hwnd, "GLFW")))
	if window == nil {
		r1, _, _ := _DefWindowProc.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
		return r1
	}

	switch msg {
	case _WM_CLOSE:
		window.shouldClose = true
	case _WM_UNICHAR:
		if wParam == _UNICODE_NOCHAR {
			// Tell the system that we accept _WM_UNICHAR messages.
			return True
		}
		fallthrough
	case _WM_CHAR, _WM_SYSCHAR:
		if r := rune(wParam); unicode.IsPrint(r) {
			if window.charCallback != nil {
				window.charCallback(nil, r)
			}
		}
		return True
	case _WM_DPICHANGED:
		// Let Windows know we're prepared for runtime DPI changes.
		return True
	case _WM_ERASEBKGND:
		// Avoid flickering between GPU content and background color.
		return True
	case _WM_KEYDOWN, _WM_KEYUP, _WM_SYSKEYDOWN, _WM_SYSKEYUP:
		var key Key
		action := Press
		if (lParam>>16)&0x8000 != 0 {
			action = Release
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
				/* TODO
				// NOTE: Alt Gr sends Left Ctrl followed by Right Alt
				// HACK: We only want one event for Alt Gr, so if we detect
				//       this sequence we discard this Left Ctrl message now
				//       and later report Right Alt normally
				MSG next;
				const DWORD time = GetMessageTime();

				if (PeekMessageW(&next, NULL, 0, 0, pm_NOREMOVE)) {
					if (next.message == _WM_KEYDOWN ||
						next.message == _WM_SYSKEYDOWN ||
						next.message == _WM_KEYUP ||
						next.message == _WM_SYSKEYUP)
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

		if action == Release && wParam == VK_SHIFT {
			// HACK: Release both Shift keys on Shift up event, as when both
			//       are pressed the first release does not emit any event
			// NOTE: The other half of this is in _glfwPlatformPollEvents
			glfwInputKey(window, KeyLeftShift, scancode, action, mods)
			glfwInputKey(window, KeyRightShift, scancode, action, mods)
		} else if wParam == VK_SNAPSHOT {
			// HACK: Key down is not reported for the Print Screen key
			glfwInputKey(window, key, scancode, Press, mods)
			glfwInputKey(window, key, scancode, Release, mods)
		} else {
			glfwInputKey(window, key, scancode, action, mods)
		}
		break

	case _WM_LBUTTONDOWN, _WM_LBUTTONUP, _WM_RBUTTONDOWN, _WM_RBUTTONUP, _WM_MBUTTONDOWN, _WM_MBUTTONUP:
		var button MouseButton
		if msg == _WM_LBUTTONDOWN || msg == _WM_LBUTTONUP {
			button = MouseButtonLeft
		} else if msg == _WM_RBUTTONDOWN || msg == _WM_RBUTTONUP {
			button = MouseButtonRight
		} else if msg == _WM_MBUTTONDOWN || msg == _WM_MBUTTONUP {
			button = MouseButtonMiddle
		}
		var action Action
		if msg == _WM_LBUTTONDOWN || msg == _WM_RBUTTONDOWN || msg == _WM_MBUTTONDOWN {
			action = Press
		} else {
			action = Release
		}
		var i MouseButton
		for i = MouseButtonFirst; i <= MouseButtonLast; i++ {
			if window.mouseButtons[i] == Press {
				break
			}
		}
		if i > MouseButtonLast {
			SetCapture(window.Win32.handle)
		}
		glfwInputMouseClick(window, button, action, getKeyMods())
		for i = MouseButtonFirst; i <= MouseButtonLast; i++ {
			if window.mouseButtons[i] == Press {
				break
			}
		}
		if i > MouseButtonLast {
			ReleaseCapture()
		}
		return 0

	case _WM_SETFOCUS:
		glfwInputWindowFocus(window, true)
		// HACK: Do not disable cursor while the user is interacting with a caption button
		if window.Win32.frameAction {
			break
		}
		if window.cursorMode == CursorDisabled {
			disableCursor(window)
		}
		return 0

	case _WM_KILLFOCUS:
		if window.cursorMode == CursorDisabled {
			enableCursor(window)
		}
		if window.monitor != nil && window.autoIconify {
			glfwIconifyWindow(window)
		}
		glfwInputWindowFocus(window, false)
		return 0

	case _WM_MOUSEMOVE:
		x := float64(int(lParam & 0xffff))
		y := float64(int((lParam >> 16) & 0xffff))
		if !window.Win32.cursorTracked {
			var tme TRACKMOUSEEVENT
			tme.dwFlags = _TME_LEAVE
			tme.hwndTrack = window.Win32.handle
			TrackMouseEvent(&tme)
			window.cursorTracked = true
			if window.cursorEnterCallback != nil {
				window.cursorEnterCallback(window, true)
			}
		}

		if window.cursorMode == CursorDisabled {
			dx := float64(x) - window.lastCursorPosX
			dy := float64(y) - window.lastCursorPosY
			if _glfw.win32.disabledCursorWindow != window {
				break
			}
			glfwInputCursorPos(window, window.virtualCursorPosX+dx, window.virtualCursorPosY+dy)
		} else {
			glfwInputCursorPos(window, x, y)
		}
		window.lastCursorPosX = x
		window.lastCursorPosY = y
		return 0

	case _WM_MOUSEWHEEL:
		glfwInputScroll(window, 0.0, float64(int16(wParam>>16))/120.0)
		return 0

	case _WM_MOUSEHWHEEL:
		glfwInputScroll(window, -float64(int16(wParam>>16))/120.0, 0.0)
		return 0

	case _WM_MOUSELEAVE:
		window.Win32.cursorTracked = false
		if window.cursorEnterCallback != nil {
			window.cursorEnterCallback(window, false)
		}
		return 0

	case _WM_PAINT:
		glfwInputWindowDamage(window)

	case _WM_SIZE:
		width := int(lParam & 0xFFFF)
		height := int(lParam >> 16)
		iconified := wParam == _SIZE_MINIMIZED
		maximized := wParam == _SIZE_MAXIMIZED || (window.Win32.maximized && wParam != _SIZE_RESTORED)
		if _glfw.win32.capturedCursorWindow == window {
			captureCursor(window)
		}

		if window.Win32.iconified != iconified {
			if window.iconifyCallback != nil {
				window.iconifyCallback(window, iconified)
			}
		}

		if window.Win32.maximized != maximized {
			if window.maximizeCallback != nil {
				window.maximizeCallback(window, maximized)
			}
		}

		if width != window.Win32.width || height != window.Win32.height {
			window.Win32.width = width
			window.Win32.height = height
			if window.sizeCallback != nil {
				window.sizeCallback(window, width, height)
			}
			if window.framebufferSizeCallback != nil {
				window.framebufferSizeCallback(window, width, height)
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

	case _WM_GETMINMAXINFO:
		if window.monitor != nil {
			break
		}

		var frame RECT
		mmi := (*MINMAXINFO)(unsafe.Pointer(lParam))
		style := getWindowStyle(window)
		exStyle := getWindowExStyle(window)
		if IsWindows10Version1607OrGreater() {
			AdjustWindowRectExForDpi(&frame, style, 0, exStyle,
				GetDpiForWindow(window.Win32.handle))
		} else {
			AdjustWindowRectEx(&frame, style, 0, exStyle)
		}
		if window.minwidth != DontCare && window.minheight != DontCare {
			mmi.ptMinTrackSize.X = window.minwidth + frame.Right - frame.Left
			mmi.ptMinTrackSize.Y = window.minheight + frame.Bottom - frame.Top
		}

		if window.maxwidth != DontCare && window.maxheight != DontCare {
			mmi.ptMaxTrackSize.X = window.maxwidth + frame.Right - frame.Left
			mmi.ptMaxTrackSize.Y = window.maxheight + frame.Bottom - frame.Top
		}
		if !window.decorated {
			mh := monitorFromWindow(window.Win32.handle, _MONITOR_DEFAULTTONEAREST)
			mi := GetMonitorInfo(mh)
			mmi.ptMaxPosition.X = mi.RcWork.Left - mi.RcMonitor.Left
			mmi.ptMaxPosition.Y = mi.RcWork.Top - mi.RcMonitor.Top
			mmi.ptMaxSize.X = mi.RcWork.Right - mi.RcWork.Left
			mmi.ptMaxSize.Y = mi.RcWork.Bottom - mi.RcWork.Top
		}
		return 0

	case _WM_SETCURSOR:
		if lParam&0xFFFF == 1 {
			updateCursorImage(window)
			return 1
		}
	}

	r1, _, _ := _DefWindowProc.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
	return r1
}

func glfwPlatformPollEvents() {
	var msg Msg
	for PeekMessage(&msg, 0, 0, 0, pm_REMOVE) {
		if msg.Message == _WM_QUIT {
			// NOTE: While GLFW does not itself post _WM_QUIT, other processes may post it to this one, for example Task Manager
			// HACK: Treat _WM_QUIT as a close on all windows
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
	// NOTE: The other half of this is in the _WM_*KEY* handler in windowProc
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
		// NOTE: Re-center the cursor only if it has moved since the last call, to avoid breaking glfwWaitEvents with _WM_MOUSEMOVE
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

func setCursor(handle syscall.Handle) {
	_, _, err := _SetCursor.Call(uintptr(handle))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("_SetCursor failed, " + err.Error())
	}
}

// Updates the cursor image according to its cursor mode
func updateCursorImage(window *Window) {
	if window.cursorMode == CursorNormal || window.cursorMode == CursorCaptured {
		if window.cursor != nil {
			setCursor(window.cursor.handle)
		} else {
			setCursor(LoadCursor(IDC_ARROW))
		}
	} else {
		// NOTE: Via Remote Desktop, setting the cursor to NULL does not hide it.
		// HACK: When running locally, it is set to NULL, but when connected via Remote
		//       Desktop, this is a transparent cursor.
		setCursor(_glfw.win32.blankCursor)
	}
}

func glfwSetCursor(window *_GLFWwindow, cursor *Cursor) {
	window.cursor = cursor
	if cursorInContentArea(window) {
		if window.cursorMode == CursorNormal || window.cursorMode == CursorCaptured {
			if window.cursor != nil {
				setCursor(window.cursor.handle)
			} else {
				setCursor(LoadCursor(IDC_ARROW))
			}
		} else {
			// NOTE: Via Remote Desktop, setting the cursor to NULL does not hide it.
			// HACK: When running locally, it is set to NULL, but when connected via Remote
			//       Desktop, this is a transparent cursor.
			setCursor(_glfw.win32.blankCursor)
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
	_, _, err := _BringWindowToTop.Call(uintptr(unsafe.Pointer(window.Win32.handle)))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("BringWindowToTop failed, " + err.Error())
	}
}

func SetForegroundWindow(window *_GLFWwindow) {
	_, _, err := _SetForegroundWindow.Call(uintptr(unsafe.Pointer(window.Win32.handle)))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("SetForegroundWindow failed, " + err.Error())
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
	fmt.Printf("dmSize=%d, pointl=%d\n", dm.dmSize, unsafe.Sizeof(dm.dmPosition))
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

// Notifies shared code of a monitor connection or disconnection
func glfwInputMonitor(monitor *Monitor, action int, placement int) {
	if action == glfw_CONNECTED {
		_glfw.monitorCount++
		if placement == InsertFirst {
			_glfw.monitors = append([]*Monitor{monitor}, _glfw.monitors...)
		} else {
			_glfw.monitors = append(_glfw.monitors, monitor)
		}
	} else if action == glfw_DISCONNECTED {
		for window := _glfw.windowListHead; window != nil; window = window.next {
			if window.monitor == monitor {
				var width, height int32
				glfwGetWindowSize(window, &width, &height)
				window.SetMonitor(nil, 0, 0, int(width), int(height), 0)
				x, y, _, _ := window.GetFrameSize()
				window.SetPos(x, y)
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
	modes = GetVideoModes(monitor)
	if len(modes) == 0 {
		return false
	}
	// slices.SortFunc(modes, glfwCompareVideoModes)
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
	var sizeDiff, leastSizeDiff int32 = _INT_MAX, _INT_MAX
	var rateDiff, leastRateDiff int32 = _INT_MAX, _INT_MAX
	var colorDiff, leastColorDiff int32 = _INT_MAX, _INT_MAX
	var current GLFWvidmode
	var closest *GLFWvidmode

	if !refreshVideoModes(monitor) {
		return nil
	}

	for i := 0; i < len(monitor.modes); i++ {
		current = monitor.modes[i]
		colorDiff = 0
		if desired.RedBits != DontCare {
			colorDiff += abs(current.RedBits) - desired.RedBits
		}
		if desired.GreenBits != DontCare {
			colorDiff += abs(current.GreenBits - desired.GreenBits)
		}
		if desired.BlueBits != DontCare {
			colorDiff += abs(current.BlueBits - desired.BlueBits)
		}
		sizeDiff = abs((current.Width-desired.Width)*(current.Width-desired.Width) + (current.Height-desired.Height)*(current.Height-desired.Height))
		if desired.RefreshRate != DontCare {
			rateDiff = abs(current.RefreshRate - desired.RefreshRate)
		} else {
			rateDiff = _INT_MAX - current.RefreshRate
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
func glfwSetVideoMode(monitor *Monitor, desired *GLFWvidmode) error {
	best := glfwChooseVideoMode(monitor, desired)
	if best == nil {
		fmt.Printf("NIL")
	}
	current := GetVideoMode(monitor)
	if glfwCompareVideoModes(&current, best) == 0 {
		return nil
	}
	var dm DEVMODEW
	dm.dmSize = uint16(unsafe.Sizeof(dm))
	dm.dmFields = dm_PELSWIDTH | dm_PELSHEIGHT | dm_BITSPERPEL | dm_DISPLAYFREQUENCY
	dm.dmPelsWidth = best.Width
	dm.dmPelsHeight = best.Height
	dm.dmBitsPerPel = best.RedBits + best.GreenBits + best.BlueBits
	dm.dmDisplayFrequency = best.RefreshRate

	if dm.dmBitsPerPel < 15 || dm.dmBitsPerPel >= 24 {
		dm.dmBitsPerPel = 32
	}
	result := ChangeDisplaySettingsEx(&monitor.adapterName[0], &dm, 0, cds_FULLSCREEN, 0)
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
		glfwSetVideoMode(window.monitor, &window.videoMode)
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

func glfwSetPos(w *Window, xPos, yPos int32) {
	rect := RECT{Left: xPos, Top: yPos, Right: xPos, Bottom: yPos}
	adjustWindowRect(&rect, getWindowStyle(w), 0, getWindowExStyle(w), GetDpiForWindow(w.Win32.handle), "glfwSetWindowPos")
	SetWindowPos(w.Win32.handle, 0, rect.Left, rect.Top, 0, 0, SWP_NOACTIVATE|SWP_NOZORDER|SWP_NOSIZE)
}

// Returns the image whose area most closely matches the desired one
//
func chooseImage(count int, images []*GLFWimage, width int32, height int32) *GLFWimage {
	var leastDiff = int32(_INT_MAX)
	var closest int

	for i := 0; i < count; i++ {
		currDiff := abs(images[i].Width*images[i].Height - width*height)
		if currDiff < leastDiff {
			closest = i
			leastDiff = currDiff
		}
	}
	return images[closest]
}

// Creates an RGBA icon or cursor
//
func createIcon(image *GLFWimage, xhot, yhot int32, icon bool) syscall.Handle {
	var handle syscall.Handle
	var bi BITMAPV5HEADER
	var ii ICONINFO
	var target *uint8
	source := (*[16384]uint8)(unsafe.Pointer(image.Pixels))

	bi.bV5Size = uint32(unsafe.Sizeof(bi))
	bi.bV5Width = image.Width
	bi.bV5Height = -image.Height
	bi.bV5Planes = 1
	bi.bV5BitCount = 32
	bi.bV5Compression = BI_BITFIELDS
	bi.bV5RedMask = 0x00ff0000
	bi.bV5GreenMask = 0x0000ff00
	bi.bV5BlueMask = 0x000000ff
	bi.bV5AlphaMask = 0xff000000

	dc := getDC(0)
	color := CreateDIBSection(dc, &bi, DIB_RGB_COLORS, &target, 0, 0)
	releaseDC(0, dc)

	if color == 0 {
		panic("Failed to create RGBA bitmap")
	}

	mask := CreateBitmap(image.Width, image.Height, 1, 1, nil)
	if mask == 0 {
		panic("Failed to create mask bitmap")
	}
	targetArr := (*[16384]uint8)(unsafe.Pointer(target))
	for i := int32(0); i < image.Width*image.Height; i++ {
		(*targetArr)[i*4+0] = (*source)[i*4+2]
		(*targetArr)[i*4+1] = (*source)[i*4+1]
		(*targetArr)[i*4+2] = (*source)[i*4+0]
		(*targetArr)[i*4+3] = (*source)[i*4+3]
	}

	ii.fIcon = icon
	ii.xHotspot = xhot
	ii.yHotspot = yhot
	ii.hbmMask = mask
	ii.hbmColor = color

	handle = CreateIconIndirect(&ii)

	DeleteObject(color)
	DeleteObject(mask)

	return handle
}

func glfwSetWindowIcon(window *Window, count int, images []*GLFWimage) {
	var bigIcon, smallIcon syscall.Handle
	if count > 0 {
		bigImage := chooseImage(count, images, GetSystemMetrics(SM_CXICON), GetSystemMetrics(SM_CYICON))
		smallImage := chooseImage(count, images, GetSystemMetrics(SM_CXSMICON), GetSystemMetrics(SM_CYSMICON))
		bigIcon = createIcon(bigImage, 0, 0, true)
		smallIcon = createIcon(smallImage, 0, 0, true)
	} else {
		bigIcon = GetClassLongPtrW(window.Win32.handle, GCLP_HICON)
		smallIcon = GetClassLongPtrW(window.Win32.handle, GCLP_HICONSM)
	}

	_ = SendMessage(window.Win32.handle, _WM_SETICON, ICON_BIG, uint32(bigIcon))
	_ = SendMessage(window.Win32.handle, _WM_SETICON, ICON_SMALL, uint32(smallIcon))

	if window.Win32.bigIcon != 0 {
		DestroyIcon(window.Win32.bigIcon)
	}

	if window.Win32.smallIcon != 0 {
		DestroyIcon(window.Win32.smallIcon)
	}
	if count > 0 {
		window.Win32.bigIcon = bigIcon
		window.Win32.smallIcon = smallIcon
	}
}

// enableRawMouseMotion enables _WM_INPUT messages for the mouse for the specified window
func enableRawMouseMotion(window *Window) {
	rid := RAWINPUTDEVICE{0x01, 0x02, 0, window.Win32.handle}
	if !RegisterRawInputDevices(&rid, 1, uint32(unsafe.Sizeof(rid))) {
		panic("Win32: Failed to register raw input device")
	}
}

// Disables _WM_INPUT messages for the mouse
//
func disableRawMouseMotion(window *Window) {
	rid := RAWINPUTDEVICE{0x01, 0x02, _RIDEV_REMOVE, 0}
	if !RegisterRawInputDevices(&rid, 1, uint32(unsafe.Sizeof(rid))) {
		panic("Failed to remove raw input device")
	}
}

// Sets the cursor clip rect to the window content area
//
func captureCursor(window *Window) {
	clipRect := GetClientRect(window.Win32.handle)
	p1 := POINT{clipRect.Left, clipRect.Top}
	p2 := POINT{clipRect.Right, clipRect.Bottom}
	p1 = ClientToScreen(window.Win32.handle, p1)
	p2 = ClientToScreen(window.Win32.handle, p2)
	r := RECT{p1.X, p1.Y, p2.X, p2.Y}
	ClipCursor(&r)
	_glfw.win32.capturedCursorWindow = window
}

// Disabled clip cursor
//
func releaseCursor() {
	ClipCursor(nil)
	_glfw.win32.capturedCursorWindow = nil
}

func SetRawMouseMotion(window *Window, enabled bool) {
	if _glfw.win32.disabledCursorWindow != window {
		return
	}
	if enabled {
		enableRawMouseMotion(window)
	} else {
		disableRawMouseMotion(window)
	}
}

func destroyCursor(cursor *Cursor) {
	if cursor.handle != 0 {
		DestroyIcon(cursor.handle)
	}
}

func glfwSetWindowAspectRatio(window *Window, numer, denom int32) {
	if numer == DontCare || denom == DontCare {
		return
	}
	area := GetWindowRect(window.Win32.handle)
	var frame RECT
	ratio := float32(window.numer) / float32(window.denom)
	style := getWindowStyle(window)
	exStyle := getWindowExStyle(window)

	if IsWindows10Version1607OrGreater() {
		AdjustWindowRectExForDpi(&frame, style, 0, exStyle, GetDpiForWindow(window.Win32.handle))
	} else {
		AdjustWindowRectEx(&frame, style, 0, exStyle)
	}
	area.Bottom = area.Top + (frame.Bottom - frame.Top) + int32(float32((area.Right-area.Left)-(frame.Right-frame.Left))/ratio)
	MoveWindow(window.Win32.handle, area.Left, area.Top, area.Right-area.Left, area.Bottom-area.Top, true)
}

func glfwGetWindowOpacity(window *Window) float32 {
	var alpha uint8
	var flags uint32
	exStyle := GetWindowLongW(window.Win32.handle, _GWL_EXSTYLE)
	if (exStyle&ws_EX_LAYERED) != 0 && GetLayeredWindowAttributes(window.Win32.handle, nil, &alpha, &flags) {
		if (flags & _LWA_ALPHA) != 0 {
			return float32(alpha) / 255.0
		}
	}
	return 1.0
}

func glfwSetWindowOpacity(window *Window, opacity float64) {
	exStyle := GetWindowLongW(window.Win32.handle, _GWL_EXSTYLE)
	if opacity < 1.0 || (exStyle&ws_EX_TRANSPARENT) != 0 {
		alpha := uint8(255 * opacity)
		exStyle |= ws_EX_LAYERED
		SetWindowLongW(window.Win32.handle, _GWL_EXSTYLE, exStyle)
		SetLayeredWindowAttributes(window.Win32.handle, 0, alpha, _LWA_ALPHA)
	} else if (exStyle & ws_EX_TRANSPARENT) != 0 {
		SetLayeredWindowAttributes(window.Win32.handle, 0, 0, 0)
	} else {
		exStyle &= ^uint32(ws_EX_LAYERED)
		SetWindowLongW(window.Win32.handle, _GWL_EXSTYLE, exStyle)
	}
}

func glfwRequestWindowAttention(window *Window) {
	// TODO FlashWIndow(window.Win32.handle, 1)
}

func glfwHideWindow(window *Window) {
	ShowWindow(window.Win32.handle, windows.SW_HIDE)
}

func glfwRestoreWindow(window *Window) {
	ShowWindow(window.Win32.handle, windows.SW_RESTORE)
}

func glfwGetWindowAttrib(window *Window, attrib Hint) int32 {
	switch attrib {
	case Focused:
		return int32(toInt(window.Win32.handle == GetActiveWindow()))
	case Iconified:
		return IsIconic(window.Win32.handle)
	case Visible:
		return IsWindowVisible(window.Win32.handle)
	case Maximized:
		return IsZoomed(window.Win32.handle)
	case Hovered:
		return int32(toInt(cursorInContentArea(window)))
	case FocusOnShow:
		return int32(toInt(window.focusOnShow))
	case MousePassthrough:
		return int32(toInt(window.mousePassthrough))
	case TransparentFramebuffer:
		if !window.Win32.transparent || !DwmIsCompositionEnabled() {
			return 0
		}
		return 1
	case Resizable:
		return int32(toInt(window.resizable))
	case Decorated:
		return int32(toInt(window.decorated))
	case Floating:
		return int32(toInt(window.floating))
	case AutoIconify:
		return int32(toInt(window.autoIconify))
	case DoubleBuffer:
		return int32(toInt(window.doublebuffer))
	case ClientAPI:
		return window.context.client
	case ContextCreationAPI:
		return window.context.source
	case ContextVersionMajor:
		return window.context.major
	case ContextVersionMinor:
		return window.context.minor
	case ContextRevision:
		return window.context.revision
	case ContextRobustness:
		return window.context.robustness
	case OpenGLForwardCompatible:
		return int32(toInt(window.context.forward))
	case OpenGLDebugContext:
		return int32(toInt(window.context.debug))
	case OpenGLProfile:
		return window.context.profile
	case ContextReleaseBehavior:
		return window.context.release
	case ContextNoError:
		return int32(toInt(window.context.noerror))
	default:
		return 0
	}
}

func toInt(x bool) int {
	if x {
		return 1
	}
	return 0
}

func toBool(x int32) bool {
	return x != 0
}

// Update native window styles to match attributes
func updateWindowStyles(window *Window) {
	style := GetWindowLongW(window.Win32.handle, _GWL_STYLE)
	style &= ^uint32(ws_OVERLAPPEDWINDOW | ws_POPUP)
	style |= getWindowStyle(window)
	rect := GetClientRect(window.Win32.handle)
	if IsWindows10Version1607OrGreater() {
		AdjustWindowRectExForDpi(&rect, style, 0, getWindowExStyle(window), GetDpiForWindow(window.Win32.handle))
	} else {
		AdjustWindowRectEx(&rect, style, 0, getWindowExStyle(window))
	}
	ClientToScreen(window.Win32.handle, POINT{rect.Left, rect.Top})
	ClientToScreen(window.Win32.handle, POINT{rect.Right, rect.Bottom})
	SetWindowLongW(window.Win32.handle, _GWL_STYLE, style)
	SetWindowPos(window.Win32.handle, hwnd_TOPMOST, rect.Left, rect.Top, rect.Right-rect.Left, rect.Bottom-rect.Top,
		SWP_FRAMECHANGED|SWP_NOACTIVATE|SWP_NOZORDER)
}

func glfwSetWindowAttrib(window *Window, attrib Hint, value int32) {
	switch attrib {
	case AutoIconify:
		window.autoIconify = toBool(value)
	case Resizable:
		window.resizable = toBool(value)
		if window.monitor == nil {
			updateWindowStyles(window)
		}
	case Decorated:
		window.decorated = toBool(value)
		if window.monitor == nil {
			updateWindowStyles(window)
		}
	case Floating:
		window.floating = toBool(value)
		if window.monitor == nil {
			updateWindowStyles(window)
		}
	case FocusOnShow:
		window.focusOnShow = toBool(value)
	case MousePassthrough:
		window.mousePassthrough = toBool(value)
		glfwSetWindowMousePassthrough(window, value != 0)
	default:
		panic("Unknown attribute")
	}
}

func glfwSetWindowMousePassthrough(window *Window, enabled bool) {
	var flags uint32
	var alpha uint8
	var key uint32
	exStyle := GetWindowLongW(window.Win32.handle, _GWL_EXSTYLE)
	if exStyle&ws_EX_LAYERED != 0 {
		GetLayeredWindowAttributes(window.Win32.handle, &key, &alpha, &flags)
	}
	if enabled {
		exStyle |= ws_EX_TRANSPARENT | ws_EX_LAYERED
	} else {
		exStyle &= ^uint32(ws_EX_TRANSPARENT)
		// NOTE: Window opacity also needs the layered window style so do not
		//       remove it if the window is alpha blended
		if exStyle&ws_EX_LAYERED != 0 {
			if (flags & _LWA_ALPHA) == 0 {
				exStyle &= ^uint32(ws_EX_LAYERED)
			}
		}
	}
	SetWindowLongW(window.Win32.handle, _GWL_EXSTYLE, exStyle)
	if enabled {
		SetLayeredWindowAttributes(window.Win32.handle, key, alpha, flags)
	}
}

func glfwSetTitle(w *Window, title string) {
	SetWindowTextW(w.Win32.handle, title)
}

func glfwGetPos(w *Window) (x, y int32) {
	p := ClientToScreen(w.Win32.handle, POINT{0, 0})
	return p.X, p.Y
}

func glfwGetFramebufferSize(w *Window) (width int, height int) {
	var area RECT
	_, _, err := _GetClientRect.Call(uintptr(unsafe.Pointer(w.Win32.handle)), uintptr(unsafe.Pointer(&area)))
	if !errors.Is(err, syscall.Errno(0)) {
		panic(err)
	}
	width = int(area.Right)
	height = int(area.Bottom)
	return width, height
}

func glfwDetachCurrentContext() {
	makeContextCurrentWGL(nil)
}

func glfwPostEmptyEvent(w *Window) {
	PostMessageW(w.Win32.handle, 0, 0, 0)
}
