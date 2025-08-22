package glfw

// constants.go contains the local, non-exported constants for win32

import (
	"syscall"
)

type MSG struct {
	hwnd    syscall.Handle
	message uint16
	wParam  uint16
	lParam  uint32
	time    uint32
	pt      POINT
}

const SM_CXICON = 11
const SM_CXSMICON = 49

type CIEXYZTRIPLE struct {
	ciexyzX int32
	ciexyzY int32
	ciexyzZ int32
}
type BITMAPV5HEADER struct {
	bV5Size          uint32
	bV5Width         int32
	bV5Height        int32
	bV5Planes        uint16
	bV5BitCount      uint16
	bV5Compression   uint32
	bV5SizeImage     uint32
	bV5XPelsPerMeter int32
	bV5YPelsPerMeter int32
	bV5ClrUsed       uint32
	bV5ClrImportant  uint32
	bV5RedMask       uint32
	bV5GreenMask     uint32
	bV5BlueMask      uint32
	bV5AlphaMask     uint32
	bV5CSType        uint32
	bV5Endpoints     CIEXYZTRIPLE
	bV5GammaRed      uint32
	bV5GammaGreen    uint32
	bV5GammaBlue     uint32
	bV5Intent        uint32
	bV5ProfileData   uint32
	bV5ProfileSize   uint32
	bV5Reserved      uint32
}

const (
	BI_BITFIELDS   = 3
	DIB_RGB_COLORS = 0
	SM_CYICON      = 12
	SM_CYSMICON    = 50
	GCLP_HICON     = -14
	GCLP_HICONSM   = -34
	_WM_SETICON    = 0x0080
	ICON_BIG       = 1
	ICON_SMALL     = 0
)

type BITMAPINFO struct {
	biSize          uint32
	biWidth         uint32
	biHeight        uint32
	biPlanes        uint16
	biBitCount      uint16
	biCompression   uint32
	biSizeImage     uint32
	biXPelsPerMeter int32
	biYPelsPerMeter int32
	biClrUsed       uint32
	biClrImportant  uint32
	bmiColors       []uint32
}
type ICONINFO struct {
	fIcon    bool
	xHotspot int32
	yHotspot int32
	hbmMask  syscall.Handle
	hbmColor syscall.Handle
}

type HDC syscall.Handle
type HMONITOR syscall.Handle
type HANDLE syscall.Handle

type MONITORINFO struct {
	CbSize    uint32
	RcMonitor RECT
	RcWork    RECT
	DwFlags   uint32
}
type RECT struct {
	Left, Top, Right, Bottom int32
}

type Msg struct {
	Hwnd     syscall.Handle
	Message  uint16
	WParam   uintptr
	LParam   uintptr
	Time     uint32
	Pt       Point
	LPrivate uint32
}

type Point struct {
	X, Y int32
}

type DISPLAY_DEVICEW struct {
	cb           uint32
	DeviceName   [32]uint16
	DeviceString [128]uint16
	StateFlags   uint32
	DeviceID     [128]uint16
	DeviceKey    [128]uint16
}

type WndClassEx struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CnClsExtra    int32
	CbWndExtra    int32
	HInstance     syscall.Handle
	HIcon         syscall.Handle
	HCursor       syscall.Handle
	HbrBackground syscall.Handle
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       syscall.Handle
}

type DEVMODEW = struct {
	mDeviceName          [32]uint16
	dmSpecVersion        uint16
	dmDriverVersion      uint16
	dmSize               uint16
	dmDriverExtra        uint16
	dmFields             uint32
	dmPosition           POINTL
	dmDisplayOrientation uint32
	dmDisplayFixedOutput uint32
	dmColor              uint16
	dmDuplex             uint16
	dmYResolution        uint16
	dmTTOption           uint16
	dmCollate            uint16
	dmFormName           [32]uint16
	dmLogPixels          uint16
	dmBitsPerPel         int32
	dmPelsWidth          int32
	dmPelsHeight         int32
	dmDisplayFlags       uint32
	dmDisplayFrequency   int32
	dmICMMethod          uint32
	dmICMIntent          uint32
	dmMediaType          uint32
	dmDitherType         uint32
	dmReserved1          uint32
	dmReserved2          uint32
	dmPanningWidth       uint32
	dmPanningHeight      uint32
}

type POINTL = struct {
	X, Y int32
}

type RAWINPUTDEVICE struct {
	usUsagePage uint16
	usUsage     uint16
	dwFlags     uint32
	hwndTarget  syscall.Handle
}

type WINDOWPLACEMENT struct {
	length           uint32
	flags            uint32
	showCmd          uint32
	ptMinPosition    POINT
	ptMaxPosition    POINT
	rcNormalPosition RECT
	rcDevice         RECT
}

type _OSVERSIONINFOEXW struct {
	dwOSVersionInfoSize uint32
	dwMajorVersion      uint32
	dwMinorVersion      uint32
	dwBuildNumber       uint32
	dwPlatformId        uint32
	szCSDVersion        [128]uint16
	wServicePackMajor   uint16
	wServicePackMinor   uint16
	wSuiteMask          uint16
	wProductType        uint8
	wReserved           uint8
}

const (
	VER_MAJORVERSION     = 0x0000002
	VER_MINORVERSION     = 0x0000001
	VER_BUILDNUMBER      = 0x0000004
	VER_SERVICEPACKMAJOR = 0x00000020
	WIN32_WINNT_WINBLUE  = 0x0603
)

const (
	DPI_AWARENESS_CONTEXT_UNAWARE              = 0xFFFFFFFFFFFFFFFF
	DPI_AWARENESS_CONTEXT_SYSTEM_AWARE         = 0xFFFFFFFFFFFFFFFE
	DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE    = 0xFFFFFFFFFFFFFFFD
	DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2 = 0xFFFFFFFFFFFFFFFC
	DPI_AWARENESS_CONTEXT_UNAWARE_GDISCALED    = 0xFFFFFFFFFFFFFFFB
	PROCESS_DPI_UNAWARE                        = 0
	PROCESS_SYSTEM_DPI_AWARE                   = 1
	PROCESS_PER_MONITOR_DPI_AWARE              = 2
)

const (
	ws_CLIPCHILDREN     = 0x02000000
	ws_CLIPSIBLINGS     = 0x04000000
	ws_MAXIMIZE         = 0x01000000
	ws_ICONIC           = 0x20000000
	ws_VISIBLE          = 0x10000000
	ws_OVERLAPPED       = 0x00000000
	ws_CAPTION          = 0x00C00000
	ws_SYSMENU          = 0x00080000
	ws_THICKFRAME       = 0x00040000
	ws_MINIMIZEBOX      = 0x00020000
	ws_MAXIMIZEBOX      = 0x00010000
	ws_POPUP            = 0x80000000
	ws_OVERLAPPEDWINDOW = ws_OVERLAPPED | ws_CAPTION | ws_SYSMENU | ws_THICKFRAME | ws_MINIMIZEBOX | ws_MAXIMIZEBOX
	ws_EX_APPWINDOW     = 0x40000
	ws_EX_TOPMOST       = 0x00000008
	ws_EX_LAYERED       = 0x00080000
	ws_EX_TRANSPARENT   = 0x00000020
)

const (
	_WM_CANCELMODE           = 0x001F
	_WM_CHAR                 = 0x0102
	_WM_SYSCHAR              = 0x0106
	_WM_CLOSE                = 0x0010
	_WM_CREATE               = 0x0001
	_WM_DPICHANGED           = 0x02E0
	_WM_DESTROY              = 0x0002
	_WM_ERASEBKGND           = 0x0014
	_WM_GETMINMAXINFO        = 0x0024
	_WM_IME_COMPOSITION      = 0x010F
	_WM_IME_ENDCOMPOSITION   = 0x010E
	_WM_IME_STARTCOMPOSITION = 0x010D
	_WM_KEYDOWN              = 0x0100
	_WM_KEYUP                = 0x0101
	_WM_KILLFOCUS            = 0x0008
	_WM_LBUTTONDOWN          = 0x0201
	_WM_LBUTTONUP            = 0x0202
	_WM_MBUTTONDOWN          = 0x0207
	_WM_MBUTTONUP            = 0x0208
	_WM_MOUSEMOVE            = 0x0200
	_WM_MOUSEWHEEL           = 0x020A
	_WM_MOUSEHWHEEL          = 0x020E
	_WM_MOUSELEAVE           = 0x02A3
	_WM_MOUSEHOVER           = 0x02A1
	_WM_NCACTIVATE           = 0x0086
	_WM_NCHITTEST            = 0x0084
	_WM_NCCALCSIZE           = 0x0083
	_WM_PAINT                = 0x000F
	_WM_QUIT                 = 0x0012
	_WM_SETCURSOR            = 0x0020
	_WM_SETFOCUS             = 0x0007
	_WM_SHOWWINDOW           = 0x0018
	_WM_SIZE                 = 0x0005
	_WM_STYLECHANGED         = 0x007D
	_WM_SYSKEYDOWN           = 0x0104
	_WM_SYSKEYUP             = 0x0105
	_WM_RBUTTONDOWN          = 0x0204
	_WM_RBUTTONUP            = 0x0205
	_WM_TIMER                = 0x0113
	_WM_UNICHAR              = 0x0109
	_WM_USER                 = 0x0400
	_WM_WINDOWPOSCHANGED     = 0x0047
	_WM_DROPFILES            = 0x0233
	_WM_COPYDATA             = 0x004A
	_WM_COPYGLOBALDATA       = 0x0049
	_MSGFLT_ALLOW            = 1
)

// Windows constants
const (
	_PM_REMOVE               = 0x0001
	_PM_NOREMOVE             = 0x0000
	_DM_PELSWIDTH            = 0x00080000
	_DM_PELSHEIGHT           = 0x00100000
	_DM_BITSPERPEL           = 0x00040000
	_DM_DISPLAYFREQUENCY     = 0x00400000
	_CDS_FULLSCREEN          = 0x00000004
	_DISP_CHANGE_SUCCESSFUL  = 0
	_DISP_CHANGE_RESTART     = 1
	_DISP_CHANGE_FAILED      = -1
	_DISP_CHANGE_BADMODE     = -2
	_DISP_CHANGE_NOTUPDATED  = -3
	_DISP_CHANGE_BADFLAGS    = -4
	_DISP_CHANGE_BADPARAM    = -5
	_DISP_CHANGE_BADDUALVIEW = -6
	_MDT_EFFECTIVE_DPI       = 0
	_MDT_ANGULAR_DPI         = 1
	_MDT_RAW_DPI             = 2
	_MDT_DEFAULT             = 3

	_GL_VERSION     = 0x1F02
	_GWL_WNDPROC    = -4
	_GWL_HINSTANCE  = -6
	_GWL_HWNDPARENT = -8
	_GWL_STYLE      = -16
	_GWL_EXSTYLE    = -20
	_GWL_USERDATA   = -21
	_GWL_ID         = -12

	_IMAGE_ICON     = 1
	_IMAGE_CURSOR   = 2
	_UNICODE_NOCHAR = 65535
	_CW_USEDEFAULT  = -2147483648

	_LR_DEFAULTCOLOR = 0x00000000
	_LR_DEFAULTSIZE  = 0x00000040
	_LR_SHARED       = 0x00008000

	_CS_HREDRAW               = 0x0002
	_CS_VREDRAW               = 0x0001
	_CS_OWNDC                 = 0x0020
	_KF_EXTENDED              = 0x100
	_MONITOR_DEFAULTTONULL    = 0x00000000
	_MONITOR_DEFAULTTOPRIMARY = 0x00000001
	_MONITOR_DEFAULTTONEAREST = 0x00000002

	_USER_DEFAULT_SCREEN_DPI = 96
	_LOGPIXELSX              = 88
	_LOGPIXELSY              = 90

	_SIZE_RESTORED  = 0
	_SIZE_MINIMIZED = 1
	_SIZE_MAXIMIZED = 2

	_HWND_TOPMOST   = 0xFFFFFFFFFFFFFFFF
	_HWND_NOTOPMOST = 0xFFFFFFFFFFFFFFFE

	_SPI_SETMOUSETRAILS = 0x005D
	_SPI_GETMOUSETRAILS = 0x005E

	_ES_CONTINUOUS       = 0x80000000
	_ES_DISPLAY_REQUIRED = 0x00000002

	_DISPLAY_DEVICE_ACTIVE         = 0x00000001
	_DISPLAY_DEVICE_ATTACHED       = 0x00000002
	_DISPLAY_DEVICE_PRIMARY_DEVICE = 0x00000004
	_QS_KEY                        = 0x1
	_QS_MOUSEMOVE                  = 0x2
	_QS_MOUSEBUTTON                = 0x4
	_QS_MOUSE                      = (_QS_MOUSEMOVE | _QS_MOUSEBUTTON)
	_QS_INPUT                      = (_QS_MOUSE | _QS_KEY)
	_QS_POSTMESSAGE                = 0x8
	_QS_TIMER                      = 0x10
	_QS_PAINT                      = 0x20
	_QS_SENDMESSAGE                = 0x40
	_QS_HOTKEY                     = 0x80
	_QS_REFRESH                    = _QS_HOTKEY | _QS_KEY | _QS_MOUSEBUTTON | _QS_PAINT
	_QS_ALLEVENTS                  = _QS_INPUT | _QS_POSTMESSAGE | _QS_TIMER | _QS_PAINT | _QS_HOTKEY
	_QS_ALLINPUT                   = _QS_SENDMESSAGE | _QS_PAINT | _QS_TIMER | _QS_POSTMESSAGE | _QS_MOUSEBUTTON | _QS_MOUSEMOVE | _QS_HOTKEY | _QS_KEY
	_QS_ALLPOSTMESSAGE             = 0x100
	_QS_RAWINPUT                   = 0x400
	_RIDEV_REMOVE                  = 1

	_LWA_COLORKEY = 0x00000001
	_LWA_ALPHA    = 0x00000002
)

const (
	_INT_MAX = 0x7FFFFFFF

	glfw_CONNECTED    = 0x00040001
	glfw_DISCONNECTED = 0x00040002

	_GL_NUM_EXTENSIONS                      = 0x821d
	_GL_EXTENSIONS                          = 0x1f03
	_GL_CONTEXT_FLAGS                       = 0x821e
	_GL_CONTEXT_FLAG_FORWARD_COMPATIBLE_BIT = 0x00000001
	_GL_CONTEXT_FLAG_DEBUG_BIT              = 0x00000002
	_GL_CONTEXT_FLAG_NO_ERROR_BIT_KHR       = 0x00000008
	_GL_CONTEXT_PROFILE_MASK                = 0x9126
	_GL_CONTEXT_COMPATIBILITY_PROFILE_BIT   = 0x00000002
	_GL_CONTEXT_CORE_PROFILE_BIT            = 0x00000001
	_GL_RESET_NOTIFICATION_STRATEGY_ARB     = 0x8256
	_GL_LOSE_CONTEXT_ON_RESET_ARB           = 0x8252
	_GL_NO_RESET_NOTIFICATION_ARB           = 0x8261
	_GL_CONTEXT_RELEASE_BEHAVIOR            = 0x82fb
	_GL_CONTEXT_RELEASE_BEHAVIOR_FLUSH      = 0x82fc
	_GL_COLOR_BUFFER_BIT                    = 0x00004000
)

type TRACKMOUSEEVENT struct {
	cbSize      uint32
	dwFlags     uint32
	hwndTrack   syscall.Handle
	dwHoverTime uint32
}

const (
	_TME_LEAVE    = 2
	_TME_CANCEL   = 0
	_TME_HOVER    = 1
	_TME_NOCLIENT = 0x10
	_TME_QUERY    = 0x40000000
)

type PIXELFORMATDESCRIPTOR = struct {
	nSize           uint16
	nVersion        uint16
	dwFlags         uint32
	iPixelType      uint8
	cColorBits      uint8
	cRedBits        uint8
	cRedShift       uint8
	cGreenBits      uint8
	cGreenShift     uint8
	cBlueBits       uint8
	cBlueShift      uint8
	cAlphaBits      uint8
	cAlphaShift     uint8
	cAccumBits      uint8
	cAccumRedBits   uint8
	cAccumGreenBits uint8
	cAccumBlueBits  uint8
	cAccumAlphaBits uint8
	cDepthBits      uint8
	cStencilBits    uint8
	cAuxBuffers     uint8
	iLayerType      uint8
	bReserved       uint8
	dwLayerMask     uint32
	dwVisibleMask   uint32
	dwDamageMask    uint32
}

const (
	PFD_TYPE_RGBA           = 0x00
	PFD_DRAW_TO_WINDOW      = 0x04
	PFD_SUPPORT_OPENGL      = 0x20
	PFD_DOUBLEBUFFER        = 0x01
	PFD_STEREO              = 0x02
	PFD_GENERIC_FORMAT      = 0x00000040
	PFD_GENERIC_ACCELERATED = 0x00001000
)
const (
	SWP_NOSIZE         = 0x0001
	SWP_NOMOVE         = 0x0002
	SWP_NOZORDER       = 0x0004
	SWP_NOREDRAW       = 0x0008
	SWP_NOACTIVATE     = 0x0010
	SWP_FRAMECHANGED   = 0x0020
	SWP_SHOWWINDOW     = 0x0040
	SWP_HIDEWINDOW     = 0x0080
	SWP_NOCOPYBITS     = 0x0100
	SWP_NOOWNERZORDER  = 0x0200
	SWP_NOSENDCHANGING = 0x0400
)

// Internal cursor types
const (
	IDC_ARROW       = 32512 // Standard arrow
	IDC_IBEAM       = 32513 // I-beam
	IDC_WAIT        = 32514 // Hour
	IDC_CROSS       = 32515 // Crosshair
	IDC_UPARROW     = 32516 // Vertical arrow
	IDC_SIZENWSE    = 32642 // Double-pointed arrow pointing northwest and southeast
	IDC_SIZENESW    = 32643 // Double-pointed arrow pointing northeast and southwest
	IDC_SIZEWE      = 32644 // Double-pointed arrow pointing west and east
	IDC_SIZENS      = 32645 // Double-pointed arrow pointing north and south
	IDC_SIZEALL     = 32646 // Four-pointed arrow pointing north, south, east, and west
	IDC_NO          = 32648 // Slashed circle
	IDC_HAND        = 32649 // Hand
	IDC_APPSTARTING = 32650 // Standard arrow and small hourglass
	IDC_HELP        = 32651 // Arrow and question mark
)

type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]uint8
}

type DEV_BROADCAST_DEVICEINTERFACE_W struct {
	dbcc_size       uint32
	dbcc_devicetype uint32
	dbcc_reserved   uint32
	dbcc_classguid  GUID
	dbcc_name       *uint16
}

const (
	DBT_DEVTYP_DEVICEINTERFACE  = 0x00000005
	DEVICE_NOTIFY_WINDOW_HANDLE = 0x00000000
)
