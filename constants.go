package glfw

import (
	"syscall"
)

// WGL Constants
const (
	wgl_NUMBER_PIXEL_FORMATS_ARB                = 0x2000
	wgl_DRAW_TO_WINDOW_ARB                      = 0x2001
	wgl_DRAW_TO_BITMAP_ARB                      = 0x2002
	wgl_ACCELERATION_ARB                        = 0x2003
	wgl_NEED_PALETTE_ARB                        = 0x2004
	wgl_NEED_SYSTEM_PALETTE_ARB                 = 0x2005
	wgl_SWAP_LAYER_BUFFERS_ARB                  = 0x2006
	wgl_SWAP_METHOD_ARB                         = 0x2007
	wgl_NUMBER_OVERLAYS_ARB                     = 0x2008
	wgl_NUMBER_UNDERLAYS_ARB                    = 0x2009
	wgl_TRANSPARENT_ARB                         = 0x200A
	wgl_TRANSPARENT_RED_VALUE_ARB               = 0x2037
	wgl_TRANSPARENT_GREEN_VALUE_ARB             = 0x2038
	wgl_TRANSPARENT_BLUE_VALUE_ARB              = 0x2039
	wgl_TRANSPARENT_ALPHA_VALUE_ARB             = 0x203A
	wgl_TRANSPARENT_INDEX_VALUE_ARB             = 0x203B
	wgl_SHARE_DEPTH_ARB                         = 0x200C
	wgl_SHARE_STENCIL_ARB                       = 0x200D
	wgl_SHARE_ACCUM_ARB                         = 0x200E
	wgl_SUPPORT_GDI_ARB                         = 0x200F
	wgl_SUPPORT_OPENGL_ARB                      = 0x2010
	wgl_DOUBLE_BUFFER_ARB                       = 0x2011
	wgl_STEREO_ARB                              = 0x2012
	wgl_PIXEL_TYPE_ARB                          = 0x2013
	wgl_COLOR_BITS_ARB                          = 0x2014
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
	wgl_GENERIC_ACCELERATION_ARB                = 0x2026
	wgl_FULL_ACCELERATION_ARB                   = 0x2027
	wgl_SWAP_EXCHANGE_ARB                       = 0x2028
	wgl_SWAP_COPY_ARB                           = 0x2029
	wgl_SWAP_UNDEFINED_ARB                      = 0x202A
	wgl_TYPE_RGBA_ARB                           = 0x202B
	wgl_TYPE_COLORINDEX_ARB                     = 0x202C
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
	PFD_TYPE_RGBA      = 0x00
	PFD_DRAW_TO_WINDOW = 0x04
	PFD_SUPPORT_OPENGL = 0x20
	PFD_DOUBLEBUFFER   = 0x01
	PFD_STEREO         = 0x02
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
