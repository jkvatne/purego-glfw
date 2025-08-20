package glfw

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	dwmapi                   = windows.NewLazySystemDLL("dwmapi.dll")
	_DwmIsCompositionEnabled = dwmapi.NewProc("DwmIsCompositionEnabled")
)

var (
	gdi32                = windows.NewLazySystemDLL("gdi32.dll")
	_GetDeviceCaps       = gdi32.NewProc("GetDeviceCaps")
	_CreateDC            = gdi32.NewProc("CreateDCW")
	_DeleteDC            = gdi32.NewProc("DeleteDC")
	_SwapBuffers         = gdi32.NewProc("SwapBuffers")
	_SetPixelFormat      = gdi32.NewProc("SetPixelFormat")
	_ChoosePixelFormat   = gdi32.NewProc("ChoosePixelFormat")
	_DescribePixelFormat = gdi32.NewProc("DescribePixelFormat")
	_CreateDIBSection    = gdi32.NewProc("CreateDIBSection")
	_CreateBitmap        = gdi32.NewProc("CreateBitmap")
	_DeleteObject        = gdi32.NewProc("DeleteObject")
)

var (
	ntdll                 = windows.NewLazySystemDLL("ntdll.dll")
	_RtlVerifyVersionInfo = ntdll.NewProc("RtlVerifyVersionInfo")
)

var (
	kernel32                 = windows.NewLazySystemDLL("kernel32.dll")
	_GetModuleHandleW        = kernel32.NewProc("GetModuleHandleW")
	_SetThreadExecutionState = kernel32.NewProc("SetThreadExecutionState")
	_TlsAlloc                = kernel32.NewProc("TlsAlloc")
	_TlsGetValue             = kernel32.NewProc("TlsGetValue")
	_TlsSetValue             = kernel32.NewProc("TlsSetValue")
	_TlsFree                 = kernel32.NewProc("TlsFree")
	_GetCurrentThreadId      = kernel32.NewProc("GetCurrentThreadId")
)

var (
	user32                         = windows.NewLazySystemDLL("user32.dll")
	_SetProcessDpiAwarenessContext = user32.NewProc("SetProcessDpiAwarenessContext")
	_EnumDisplayMonitors           = user32.NewProc("EnumDisplayMonitors")
	_EnumDisplayDevices            = user32.NewProc("EnumDisplayDevicesW")
	_EnumDisplaySettingsEx         = user32.NewProc("EnumDisplaySettingsExW")
	_GetMonitorInfo                = user32.NewProc("GetMonitorInfoW")
	_AdjustWindowRectEx            = user32.NewProc("AdjustWindowRectEx")
	_CreateWindowEx                = user32.NewProc("CreateWindowExW")
	_DefWindowProc                 = user32.NewProc("DefWindowProcW")
	_DestroyWindow                 = user32.NewProc("DestroyWindow")
	_DispatchMessage               = user32.NewProc("DispatchMessageW")
	_GetWindowRect                 = user32.NewProc("GetWindowRect")
	_GetClientRect                 = user32.NewProc("GetClientRect")
	_GetDC                         = user32.NewProc("GetDC")
	_GetDpiForWindow               = user32.NewProc("GetDpiForWindow")
	_GetKeyState                   = user32.NewProc("GetKeyState")
	_LoadCursor                    = user32.NewProc("LoadCursorW")
	_LoadImageW                    = user32.NewProc("LoadImageW")
	_MonitorFromWindow             = user32.NewProc("MonitorFromWindow")
	_PeekMessage                   = user32.NewProc("PeekMessageW")
	_WaitMessage                   = user32.NewProc("WaitMessage")
	_RegisterClassExW              = user32.NewProc("RegisterClassExW")
	_ReleaseDC                     = user32.NewProc("ReleaseDC")
	_ScreenToClient                = user32.NewProc("ScreenToClient")
	_ShowWindow                    = user32.NewProc("ShowWindow")
	_SetCursor                     = user32.NewProc("SetCursor")
	_SetForegroundWindow           = user32.NewProc("SetForegroundWindow")
	_SetFocus                      = user32.NewProc("SetFocus")
	_SetProcessDPIAware            = user32.NewProc("SetProcessDPIAware")
	_SetWindowPos                  = user32.NewProc("SetWindowPos")
	_TranslateMessage              = user32.NewProc("TranslateMessage")
	_UnregisterClass               = user32.NewProc("UnregisterClassW")
	_BringWindowToTop              = user32.NewProc("BringWindowToTop")
	_GetCursorPos                  = user32.NewProc("GetCursorPos")
	_SystemParametersInfoW         = user32.NewProc("SystemParametersInfoW")
	_GetWindowLongW                = user32.NewProc("GetWindowLongW")
	_SetWindowLongW                = user32.NewProc("SetWindowLongW")
	_GetActiveWindow               = user32.NewProc("GetActiveWindow")
	_GetPropA                      = user32.NewProc("GetPropA")
	_SetPropA                      = user32.NewProc("SetPropA")
	_MsgWaitForMultipleObjects     = user32.NewProc("MsgWaitForMultipleObjects")
	_GetSystemMetrics              = user32.NewProc("GetSystemMetrics")
	_CreateIcon                    = user32.NewProc("CreateIcon")
	_DestroyIcon                   = user32.NewProc("DestroyIcon")
	_CreateIconIndirect            = user32.NewProc("CreateIconIndirect")
	_GetClassLongPtrW              = user32.NewProc("GetClassLongPtrW")
	_SendMessage                   = user32.NewProc("SendMessageW")
	_RegisterRawInputDevices       = user32.NewProc("RegisterRawInputDevices")
	_ClientToScreen                = user32.NewProc("ClientToScreen")
	_ClipCursor                    = user32.NewProc("ClipCursor")
	_SetCursorPos                  = user32.NewProc("SetCursorPos")
	_SetWindowTextW                = user32.NewProc("SetWindowTextW")
	_MoveWindow                    = user32.NewProc("MoveWindow")
	_GetLayeredWindowAttributes    = user32.NewProc("GetLayeredWindowAttributes")
	_SetLayeredWindowAttributes    = user32.NewProc("SetLayeredWindowAttributes")
	_IsIconic                      = user32.NewProc("IsIconic")
	_IsWindowVisible               = user32.NewProc("IsWindowVisible")
	_IsZoomed                      = user32.NewProc("IsZoomed")
	_PostMessageW                  = user32.NewProc("PostMessageW")
	_ChangeDisplaySettingsEx       = user32.NewProc("ChangeDisplaySettingsExW")
	_ChangeWindowMessageFilterEx   = user32.NewProc("ChangeWindowMessageFilterEx")
	_GetWindowPlacement            = user32.NewProc("GetWindowPlacement")
	_SetWindowPlacement            = user32.NewProc("SetWindowPlacement")
	_SetCapture                    = user32.NewProc("SetCapture")
	_ReleaseCapture                = user32.NewProc("ReleaseCapture")
	_TrackMouseEvent               = user32.NewProc("TrackMouseEvent")
	_FlashWindow                   = user32.NewProc("FlashWindow")
)

var (
	shcore            = windows.NewLazySystemDLL("shcore.dll")
	_GetDpiForMonitor = shcore.NewProc("GetDpiForMonitor")
)

func GetActiveWindow() syscall.Handle {
	r, _, err := _GetActiveWindow.Call()
	if !errors.Is(err, syscall.Errno(0)) {
		panic("GetActiveWindow failed, " + err.Error())
	}
	return syscall.Handle(r)
}

func GetProp(handle syscall.Handle, key string) uintptr {
	// widestr, _ := syscall.UTF16PtrFromString(key)
	cstr, _ := windows.BytePtrFromString(key)
	r, _, err := _GetPropA.Call(uintptr(handle), uintptr(unsafe.Pointer(cstr)))
	if !errors.Is(err, syscall.Errno(0)) {
		// panic("GetProp failed, " + err.Error())
	}
	return r
}

func SetProp(handle syscall.Handle, key string, data uintptr) {
	// widestr, _ := syscall.UTF16PtrFromString(key)
	cstr, _ := windows.BytePtrFromString(key)
	_, _, err := _SetPropA.Call(uintptr(handle), uintptr(unsafe.Pointer(cstr)), data)
	if !errors.Is(err, syscall.Errno(0)) {
		panic("SetProp failed, " + err.Error())
	}
}

func GetKeyState(nVirtKey int) uint16 {
	c, _, _ := _GetKeyState.Call(uintptr(nVirtKey))
	return uint16(c)
}

func GetModuleHandle() (syscall.Handle, error) {
	h, _, err := _GetModuleHandleW.Call(uintptr(0))
	if h == 0 {
		return 0, fmt.Errorf("GetModuleHandleW failed: %v", err)
	}
	return syscall.Handle(h), nil
}

func RegisterClassEx(cls *WndClassEx) (uint16, error) {
	a, _, err := _RegisterClassExW.Call(uintptr(unsafe.Pointer(cls)))
	if a == 0 {
		return 0, fmt.Errorf("RegisterClassExW failed: %v", err)
	}
	return uint16(a), nil
}

func LoadImage(hInst syscall.Handle, res uintptr, typ uint32, cx, cy int, fuload uint32) (syscall.Handle, error) {
	h, _, err := _LoadImageW.Call(uintptr(hInst), res, uintptr(typ), uintptr(cx), uintptr(cy), uintptr(fuload))
	if h == 0 {
		return 0, fmt.Errorf("LoadImage failed: %v", err)
	}
	return syscall.Handle(h), nil
}

func PeekMessage(m *Msg, hwnd syscall.Handle, wMsgFilterMin, wMsgFilterMax, wRemoveMsg uint32) bool {
	r, _, _ := _PeekMessage.Call(uintptr(unsafe.Pointer(m)), uintptr(hwnd), uintptr(wMsgFilterMin), uintptr(wMsgFilterMax), uintptr(wRemoveMsg))
	return r != 0
}

func TranslateMessage(m *Msg) {
	_TranslateMessage.Call(uintptr(unsafe.Pointer(m)))
}

func DispatchMessage(m *Msg) {
	_DispatchMessage.Call(uintptr(unsafe.Pointer(m)))
}

func WaitMessage() {
	_, _, err := _WaitMessage.Call()
	if !errors.Is(err, syscall.Errno(0)) {
		// panic("WaitMessage failed, " + err.Error())
		fmt.Printf("WaitMessage failed: %v", err)
	}
}

func CreateWindowEx(dwExStyle uint32, lpClassName uint16, lpWindowName string, dwStyle uint32, x, y, w, h int32, hWndParent, hMenu, hInstance syscall.Handle, lpParam uintptr) (syscall.Handle, error) {
	wname, _ := syscall.UTF16PtrFromString(lpWindowName)
	hwnd, _, err := _CreateWindowEx.Call(
		uintptr(dwExStyle),
		uintptr(lpClassName),
		uintptr(unsafe.Pointer(wname)),
		uintptr(dwStyle),
		uintptr(x), uintptr(y),
		uintptr(w), uintptr(h),
		uintptr(hWndParent),
		uintptr(hMenu),
		uintptr(hInstance),
		uintptr(lpParam))
	if hwnd == 0 {
		return 0, fmt.Errorf("CreateWindowEx failed: %v", err)
	}
	return syscall.Handle(hwnd), nil
}

func SetWindowTextW(window syscall.Handle, title string) {
	wname, _ := syscall.UTF16PtrFromString(title)
	_SetWindowTextW.Call(uintptr(window), uintptr(unsafe.Pointer(wname)))
}

func DestroyWindow(h syscall.Handle) {
	_, _, err := _DestroyWindow.Call(uintptr(h))
	if !errors.Is(err, syscall.Errno(0)) {
		// An error 'invalid window handle' can occur without any specific reasons (#2551).
		if !errors.Is(err, syscall.Errno(1400)) {
			panic("DestroyWindow failed, " + err.Error())
		}
	}
}

func UnregisterClass(class uint16, instance syscall.Handle) {
	_, _, _ = _UnregisterClass.Call(uintptr(class), uintptr(instance))
}

func IsWindows10Version1607OrGreater() bool {
	var osvi _OSVERSIONINFOEXW
	osvi.dwOSVersionInfoSize = uint32(unsafe.Sizeof(osvi))
	osvi.dwMajorVersion = 10
	osvi.dwMinorVersion = 0
	osvi.dwBuildNumber = 14393
	var mask uint32 = VER_MAJORVERSION | VER_MINORVERSION | VER_BUILDNUMBER
	r, _, err := _RtlVerifyVersionInfo.Call(uintptr(unsafe.Pointer(&osvi)), uintptr(mask), uintptr(0x80000000000000db))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("SetProcessDpiAwarenessContext failed, " + err.Error())
	}
	return r == 0
}

func IsWindows10Version1703OrGreater() bool {
	var osvi _OSVERSIONINFOEXW
	osvi.dwOSVersionInfoSize = uint32(unsafe.Sizeof(osvi))
	osvi.dwMajorVersion = 10
	osvi.dwMinorVersion = 0
	osvi.dwBuildNumber = 15063
	var mask uint32 = VER_MAJORVERSION | VER_MINORVERSION | VER_BUILDNUMBER
	r, _, err := _RtlVerifyVersionInfo.Call(uintptr(unsafe.Pointer(&osvi)), uintptr(mask), uintptr(0x80000000000000db))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("isWindows10Version1703OrGreater failed, " + err.Error())
	}
	return r == 0
}

func IsWindows8Point1OrGreater() bool {
	var osvi _OSVERSIONINFOEXW
	osvi.dwOSVersionInfoSize = uint32(unsafe.Sizeof(osvi))
	osvi.dwMajorVersion = uint32(WIN32_WINNT_WINBLUE >> 8)
	osvi.dwMinorVersion = uint32(WIN32_WINNT_WINBLUE & 0xFF)
	osvi.wServicePackMajor = 0
	var mask uint32 = VER_MAJORVERSION | VER_MINORVERSION | VER_SERVICEPACKMAJOR
	r, _, err := _RtlVerifyVersionInfo.Call(uintptr(unsafe.Pointer(&osvi)), uintptr(mask), uintptr(0x800000000001801b))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("SetProcessDpiAwarenessContext failed, " + err.Error())
	}
	return r == 0
}

func SetProcessDpiAwareness() {
	if IsWindows10Version1703OrGreater() {
		_, _, err := _SetProcessDpiAwarenessContext.Call(uintptr(DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2))
		if !errors.Is(err, syscall.Errno(0)) {
			panic("SetProcessDpiAwarenessContext failed, " + err.Error())
		}
	} else if IsWindows8Point1OrGreater() {
		_, _, err := _SetProcessDpiAwarenessContext.Call(uintptr(PROCESS_PER_MONITOR_DPI_AWARE))
		if !errors.Is(err, syscall.Errno(0)) {
			panic("SetProcessDpiAwarenessContext failed, " + err.Error())
		}
	} else if isWindowsVistaOrGreater() {
		_, _, _ = _SetProcessDPIAware.Call()
	}
}

func SetWindowPos(hwnd syscall.Handle, after syscall.Handle, x, y, w, h int32, flags uint32) {
	_, _, err := _SetWindowPos.Call(uintptr(hwnd), uintptr(after), uintptr(x), uintptr(y), uintptr(w), uintptr(h), uintptr(flags))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("SetWindowPos failed, " + err.Error())
	}
}

func GetWindowLongW(hWnd syscall.Handle, index int32) uint32 {
	r1, _, err := _GetWindowLongW.Call(uintptr(hWnd), uintptr(index))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("GetWindowLongW failed, " + err.Error())
	}
	return uint32(r1)
}

func SetWindowLongW(hWnd syscall.Handle, index int, newValue uint32) {
	_, _, err := _SetWindowLongW.Call(uintptr(hWnd), uintptr(index), uintptr(newValue))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("GetWindowLongW failed, " + err.Error())
	}
}

func EnumDisplayDevices(device uintptr, no int, adapter *DISPLAY_DEVICEW, flags uint32) bool {
	ret, _, err := _EnumDisplayDevices.Call(device, uintptr(no), uintptr(unsafe.Pointer(adapter)), uintptr(flags))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("EnumDisplayDevices failed")
	}
	return ret == 1
}

func GetMonitorInfo(hMonitor HMONITOR) *MONITORINFO {
	lmpi := MONITORINFO{}
	lmpi.CbSize = uint32(unsafe.Sizeof(lmpi))
	_, _, err := _GetMonitorInfo.Call(uintptr(hMonitor), uintptr(unsafe.Pointer(&lmpi)))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("GetMonitorInfo failed, " + err.Error())
	}
	return &lmpi
}

func GetDeviceCaps(dc HDC, flags int) int {
	r1, _, err := _GetDeviceCaps.Call(uintptr(unsafe.Pointer(dc)), uintptr(flags))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("GetDeviceCaps failed, " + err.Error())
	}
	return int(r1)
}

func GetDpiForMonitor(h HMONITOR, kind uint32) (dpiX int, dpiY int) {
	r1, _, err := _GetDpiForMonitor.Call(uintptr(h), uintptr(kind), uintptr(unsafe.Pointer(&dpiX)), uintptr(unsafe.Pointer(&dpiY)))
	if !errors.Is(err, syscall.Errno(0)) || r1 != 0 {
		panic("GetDpiForMonitor failed, " + err.Error())
	}
	return dpiX, dpiY
}

func EnumDisplayMonitors(hdc HDC, clip *RECT, lpfnEnum uintptr, data uintptr) error {
	ret, _, _ := _EnumDisplayMonitors.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(clip)),
		lpfnEnum,
		uintptr(unsafe.Pointer(data)),
	)
	if ret == 0 {
		return fmt.Errorf("w32.EnumDisplayMonitors returned FALSE")
	}
	return nil
}

func isWindowsVistaOrGreater() bool {
	return true
}

func AdjustWindowRectEx(rect *RECT, style uint32, menu int, exStyle uint32) {
	_, _, err := _AdjustWindowRectEx.Call(uintptr(unsafe.Pointer(rect)), uintptr(style), uintptr(menu), uintptr(exStyle))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("AdjustWindowRectEx failed, " + err.Error())
	}
}

func GetDpiForWindow(handle syscall.Handle) int {
	r, _, err := _GetDpiForWindow.Call(uintptr(handle))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("GetDpiForWindow failed, " + err.Error())
	}
	return int(r)
}

func AdjustWindowRectExForDpi(rect *RECT, style uint32, menu int, exStyle uint32, dpi int) {
	_, _, err := _AdjustWindowRectEx.Call(uintptr(unsafe.Pointer(rect)), uintptr(style), uintptr(menu), uintptr(exStyle), uintptr(dpi))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("AdjustWindowRectExForDpi failed, " + err.Error())
	}
}

func EnumDisplaySettingsEx(name *uint16, mode int, dm *DEVMODEW, flags int) int {
	r, _, err := _EnumDisplaySettingsEx.Call(uintptr(unsafe.Pointer(name)), uintptr(mode), uintptr(unsafe.Pointer(dm)), uintptr(flags))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("EnumDisplySettingsEx failed, " + err.Error())
	}
	return int(r)
}
func TlsSetValue(index int, value uintptr) {
	_, _, err := _TlsSetValue.Call(uintptr(index), value)
	if !errors.Is(err, syscall.Errno(0)) {
		panic("_TTlsGetValue failed, " + err.Error())
	}
}

func TlsGetValue(index int) uintptr {
	r, _, err := _TlsGetValue.Call(uintptr(index))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("_TTlsGetValue failed, " + err.Error())
	}
	return r
}

func TlsAlloc() int {
	r, _, err := _TlsAlloc.Call()
	if !errors.Is(err, syscall.Errno(0)) {
		panic("_TlsAlloc failed, " + err.Error())
	}
	return int(r)
}

func TlsFree(index int) {
	_, _, err := _TlsFree.Call(uintptr(index))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("TlsFree failed, " + err.Error())
	}
}

func glfwPlatformSetTls(tls *_GLFWtls, value uintptr) {
	if !tls.allocated {
		panic("_glfwPlatformGetTls failed: tls not allocated")
	}
	TlsSetValue(tls.index, value)
}

func glfwPlatformGetTls(tls *_GLFWtls) uintptr {
	if !tls.allocated {
		panic("_glfwPlatformGetTls failed: tls not allocated")
	}
	return TlsGetValue(tls.index)
}

func glfwPlatformDestroyTls(tls *_GLFWtls) {
	if tls.allocated {
		TlsFree(tls.index)
	}
}

func glfwPlatformCreateTls(tls *_GLFWtls) error {
	if tls.allocated {
		return fmt.Errorf("glfwPlatformCreateTls: already allocated")
	}
	tls.index = TlsAlloc()
	if tls.index == 4294967295 { // TLS_OUT_OF_INDEXES
		return fmt.Errorf("glfwPlatformCreateTls: Failed to allocate TLS index")
	}
	tls.allocated = true
	return nil
}

func MsgWaitForMultipleObjects(nCount uint32, pHandles *HANDLE, fWaitAll uint32, dwMilliseconds uint32, dwWakeMask uint32) uint32 {
	r, _, err := _MsgWaitForMultipleObjects.Call(uintptr(nCount), uintptr(unsafe.Pointer(pHandles)), uintptr(fWaitAll),
		uintptr(dwMilliseconds), uintptr(dwWakeMask))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("MsgWaitForMultipleObjects failed, " + err.Error())
	}
	return uint32(r)
}

func GetCurrentThreadId() uint32 {
	r, _, err := _GetCurrentThreadId.Call()
	if !errors.Is(err, syscall.Errno(0)) {
		panic("GetCurrentThreadId failed, " + err.Error())
	}
	return uint32(r)
}

func GetSystemMetrics(index int32) int32 {
	r, _, err := _GetSystemMetrics.Call(uintptr(index))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("GetSystemMetrics failed, " + err.Error())
	}
	return int32(r)
}

func CreateIcon(hInstance, nWidth, nHeight int, cPlanes int, cBitsPixel, AndBits *uint8, XorBits *uint8) syscall.Handle {
	r, _, err := _CreateIcon.Call(uintptr(hInstance), uintptr(nWidth), uintptr(nHeight), uintptr(cPlanes),
		uintptr(unsafe.Pointer(AndBits)), uintptr(unsafe.Pointer(XorBits)))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("CreateIcon failed, " + err.Error())
	}
	return syscall.Handle(r)
}

func CreateDIBSection(hdc HDC, pbmi *BITMAPV5HEADER, usage uint32, ppvBits **uint8, hSection syscall.Handle, offset uint32) syscall.Handle {
	r, _, err := _CreateDIBSection.Call(uintptr(hdc), uintptr(unsafe.Pointer(pbmi)), uintptr(usage), uintptr(unsafe.Pointer(ppvBits)),
		uintptr(hSection), uintptr(offset))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("CreateDIBSection failed, " + err.Error())
	}
	return syscall.Handle(r)
}

func CreateBitmap(nWidth int32, nHeight int32, nPlanes uint32, nBitCount uint32, lpBits *uint8) syscall.Handle {
	r, _, err := _CreateBitmap.Call(uintptr(nWidth), uintptr(nHeight), uintptr(nPlanes), uintptr(nBitCount), uintptr(unsafe.Pointer(lpBits)))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("CreateBitmap failed, " + err.Error())
	}
	return syscall.Handle(r)
}

func CreateIconIndirect(piconinfo *ICONINFO) syscall.Handle {
	r, _, err := _CreateIconIndirect.Call(uintptr(unsafe.Pointer(piconinfo)))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("CreateIconIndirect failed, " + err.Error())
	}
	return syscall.Handle(r)
}

func DestroyIcon(h syscall.Handle) bool {
	r, _, _ := _DestroyIcon.Call(uintptr(h))
	return r != 0
}

func DeleteObject(h syscall.Handle) bool {
	r, _, err := _DeleteObject.Call(uintptr(h))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("DeleteObject failed, " + err.Error())
	}
	return r != 0
}

func GetClassLongPtrW(hWnd syscall.Handle, nIndex int32) syscall.Handle {
	r, _, err := _GetClassLongPtrW.Call(uintptr(hWnd), uintptr(nIndex))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("GetClassLongPtrW failed, " + err.Error())
	}
	return syscall.Handle(r)
}

func SendMessage(hWnd syscall.Handle, Msg uint32, wParam uint16, Lparam uint32) uintptr {
	r, _, err := _SendMessage.Call(uintptr(hWnd), uintptr(Msg), uintptr(wParam), uintptr(Lparam))
	if !errors.Is(err, syscall.Errno(0)) {
		fmt.Println("SendMessage failed, " + err.Error())
	}
	return r
}

func RegisterRawInputDevices(pRawInputDevices *RAWINPUTDEVICE, uiNumDevices uint32, cbSize uint32) bool {
	r, _, err := _RegisterRawInputDevices.Call(uintptr(unsafe.Pointer(pRawInputDevices)), uintptr(uiNumDevices), uintptr(cbSize))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("RegisterRawInputDevices failed, " + err.Error())
	}
	return r != 0
}

func GetClientRect(hWnd syscall.Handle) RECT {
	var area RECT
	_, _, err := _GetClientRect.Call(uintptr(unsafe.Pointer(hWnd)), uintptr(unsafe.Pointer(&area)))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("GetClientRect failed, " + err.Error())
	}
	return area
}

func ClientToScreen(hWnd syscall.Handle, p POINT) POINT {
	_, _, err := _ClientToScreen.Call(uintptr(hWnd), uintptr(unsafe.Pointer(&p)))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("ClientToScreen failed, " + err.Error())
	}
	return p
}

func ClipCursor(rect *RECT) {
	_, _, err := _ClipCursor.Call(uintptr(unsafe.Pointer(rect)))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("ClipCursor failed, " + err.Error())
	}
}

func SetCursorPos(x, y int32) {
	_, _, err := _SetCursorPos.Call(uintptr(x), uintptr(y))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("SetCursorPos failed, " + err.Error())
	}
}

func GetWindowRect(handle syscall.Handle) RECT {
	var area RECT
	_, _, err := _GetWindowRect.Call(uintptr(handle), uintptr(unsafe.Pointer(&area)))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("GetWindowRect failed, " + err.Error())
	}
	return area
}

func MoveWindow(hWnd syscall.Handle, x, y, w, h int32, repaint bool) {
	rp := 0
	if repaint {
		rp = 1
	}
	_, _, err := _MoveWindow.Call(uintptr(hWnd), uintptr(x), uintptr(y), uintptr(w), uintptr(h), uintptr(rp))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("MoveWindow failed, " + err.Error())
	}
}

func ShowWindow(hWnd syscall.Handle, mode int32) {
	_, _, err := _ShowWindow.Call(uintptr(hWnd), uintptr(mode))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("ShowWindow failed, " + err.Error())
	}
}

func GetLayeredWindowAttributes(hWnd syscall.Handle, pcrKey *uint32, pbAlpha *uint8, pdwFlags *uint32) bool {
	r, _, err := _GetLayeredWindowAttributes.Call(uintptr(hWnd), uintptr(unsafe.Pointer(pcrKey)), uintptr(unsafe.Pointer(pbAlpha)), uintptr(unsafe.Pointer(pdwFlags)))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("GetLayeredWindowAttributes failed, " + err.Error())
	}
	return r != 0
}

func SetLayeredWindowAttributes(hWnd syscall.Handle, pcrKey uint32, pbAlpha uint8, pdwFlags uint32) bool {
	r, _, err := _SetLayeredWindowAttributes.Call(uintptr(hWnd), uintptr(pcrKey), uintptr(pbAlpha), uintptr(pdwFlags))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("SetLayeredWindowAttributes failed, " + err.Error())
	}
	return r != 0
}

func IsIconic(hWnd syscall.Handle) int32 {
	r, _, err := _IsIconic.Call(uintptr(hWnd))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("IsIconic failed, " + err.Error())
	}
	return int32(r)
}

func IsWindowVisible(hwnd syscall.Handle) int32 {
	r, _, err := _IsWindowVisible.Call(uintptr(hwnd))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("IsWindowVisible failed, " + err.Error())
	}
	return int32(r)
}

func IsZoomed(hwnd syscall.Handle) int32 {
	r, _, err := _IsZoomed.Call(uintptr(hwnd))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("IsZoomed failed, " + err.Error())
	}
	return int32(r)
}

func LoadCursor(cursorID uint16) syscall.Handle {
	h, err := LoadImage(0, uintptr(cursorID), _IMAGE_CURSOR, 0, 0, _LR_DEFAULTSIZE|_LR_SHARED)
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("LoadCursor failed, " + err.Error())
	}
	if h == 0 {
		panic("LoadCursor failed")
	}
	return syscall.Handle(h)
}

func PostMessageW(hWnd syscall.Handle, msg uint32, wParam, lParam uintptr) {
	_, _, _ = _PostMessageW.Call(uintptr(hWnd), uintptr(msg), uintptr(wParam), uintptr(lParam))
}

func ChangeDisplaySettingsEx(name *uint16, mode *DEVMODEW, hWnd syscall.Handle, flags uint32, lParam uintptr) int32 {
	r, _, err := _ChangeDisplaySettingsEx.Call(uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(mode)), uintptr(hWnd), uintptr(flags), lParam)
	if !errors.Is(err, syscall.Errno(0)) {
		panic("IsZoomed failed, " + err.Error())
	}
	return int32(r)
}

func DwmIsCompositionEnabled() bool {
	var flag uint32
	r, _, err := _DwmIsCompositionEnabled.Call(uintptr(unsafe.Pointer(&flag)))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		return false
	}
	return r != 0
}

func ChangeWindowMessageFilterEx(hWnd syscall.Handle, msg uint32, action uint32, filter uintptr) bool {
	r, _, err := _ChangeWindowMessageFilterEx.Call(uintptr(hWnd), uintptr(msg), uintptr(action), filter)
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		return false
	}
	return r != 0
}

func GetWindowPlacement(hWnd syscall.Handle, wp *WINDOWPLACEMENT) bool {
	r, _, err := _GetWindowPlacement.Call(uintptr(hWnd), uintptr(unsafe.Pointer(wp)))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		return false
	}
	return r != 0
}

func SetWindowPlacement(hWnd syscall.Handle, wp *WINDOWPLACEMENT) bool {
	r, _, err := _SetWindowPlacement.Call(uintptr(hWnd), uintptr(unsafe.Pointer(wp)))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		return false
	}
	return r != 0
}

func OffsetRect(r *RECT, dx, dy int32) {
	r.Top += dy
	r.Left += dx
	r.Bottom += dy
	r.Right += dx
}

func SetCapture(hWnd syscall.Handle) {
	_, _, _ = _SetCapture.Call(uintptr(hWnd))
}

func ReleaseCapture() {
	_, _, _ = _ReleaseCapture.Call()
}

func TrackMouseEvent(tme *TRACKMOUSEEVENT) {
	_TrackMouseEvent.Call(uintptr(unsafe.Pointer(tme)))
}

func FlashWindow(hWnd syscall.Handle, invert uint32) {
	_, _, _ = _FlashWindow.Call(uintptr(hWnd), uintptr(invert))
}
