package glfw

import "syscall"

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
	pm_REMOVE   = 0x0001
	pm_NOREMOVE = 0x0000

	dm_PELSWIDTH            = 0x00080000
	dm_PELSHEIGHT           = 0x00100000
	dm_BITSPERPEL           = 0x00040000
	dm_DISPLAYFREQUENCY     = 0x00400000
	cds_FULLSCREEN          = 0x00000004
	disp_CHANGE_SUCCESSFUL  = 0
	disp_CHANGE_RESTART     = 1
	disp_CHANGE_FAILED      = -1
	disp_CHANGE_BADMODE     = -2
	disp_CHANGE_NOTUPDATED  = -3
	disp_CHANGE_BADFLAGS    = -4
	disp_CHANGE_BADPARAM    = -5
	disp_CHANGE_BADDUALVIEW = -6

	mdt_EFFECTIVE_DPI = 0
	mdt_ANGULAR_DPI   = 1
	mdt_RAW_DPI       = 2
	mdt_DEFAULT       = 3

	_GL_VERSION     = 0x1F02
	_GWL_WNDPROC    = -4
	_GWL_HINSTANCE  = -6
	_GWL_HWNDPARENT = -8
	_GWL_STYLE      = -16
	_GWL_EXSTYLE    = -20
	_GWL_USERDATA   = -21
	_GWL_ID         = -12

	_IMAGE_BITMAP      = 0
	_IMAGE_ICON        = 1
	_IMAGE_CURSOR      = 2
	_IMAGE_ENHMETAFILE = 3
	_UNICODE_NOCHAR    = 65535
	cw_USEDEFAULT      = -2147483648

	lr_CREATEDIBSECTION = 0x00002000
	lr_DEFAULTCOLOR     = 0x00000000
	lr_DEFAULTSIZE      = 0x00000040
	lr_LOADFROMFILE     = 0x00000010
	lr_LOADMAP3DCOLORS  = 0x00001000
	lr_LOADTRANSPARENT  = 0x00000020
	lr_MONOCHROME       = 0x00000001
	lr_SHARED           = 0x00008000
	lr_VGACOLOR         = 0x00000080

	cs_HREDRAW                = 0x0002
	cs_INSERTCHAR             = 0x2000
	cs_NOMOVECARET            = 0x4000
	cs_VREDRAW                = 0x0001
	cs_OWNDC                  = 0x0020
	kf_EXTENDED               = 0x100
	_MONITOR_DEFAULTTONULL    = 0x00000000
	_MONITOR_DEFAULTTOPRIMARY = 0x00000001
	_MONITOR_DEFAULTTONEAREST = 0x00000002

	_USER_DEFAULT_SCREEN_DPI = 96
	_LOGPIXELSX              = 88
	_LOGPIXELSY              = 90

	_SIZE_RESTORED  = 0
	_SIZE_MINIMIZED = 1
	_SIZE_MAXIMIZED = 2
	_SIZE_MAXSHOW   = 3
	_SIZE_MAXHIDE   = 4

	hwnd_TOPMOST   = 0xFFFFFFFFFFFFFFFF
	hwnd_NOTOPMOST = 0xFFFFFFFFFFFFFFFE

	spi_SETMOUSETRAILS = 0x005D
	spi_GETMOUSETRAILS = 0x005E

	es_CONTINUOUS       = 0x80000000
	es_DISPLAY_REQUIRED = 0x00000002

	_DISPLAY_DEVICE_ACTIVE         = 0x00000001
	_DISPLAY_DEVICE_ATTACHED       = 0x00000002
	_DISPLAY_DEVICE_PRIMARY_DEVICE = 0x00000004
	qs_KEY                         = 0x1
	qs_MOUSEMOVE                   = 0x2
	qs_MOUSEBUTTON                 = 0x4
	qs_MOUSE                       = (qs_MOUSEMOVE | qs_MOUSEBUTTON)
	qs_INPUT                       = (qs_MOUSE | qs_KEY)
	qs_POSTMESSAGE                 = 0x8
	qs_TIMER                       = 0x10
	qs_PAINT                       = 0x20
	qs_SENDMESSAGE                 = 0x40
	qs_HOTKEY                      = 0x80
	qs_REFRESH                     = qs_HOTKEY | qs_KEY | qs_MOUSEBUTTON | qs_PAINT
	qs_ALLEVENTS                   = qs_INPUT | qs_POSTMESSAGE | qs_TIMER | qs_PAINT | qs_HOTKEY
	qs_ALLINPUT                    = qs_SENDMESSAGE | qs_PAINT | qs_TIMER | qs_POSTMESSAGE | qs_MOUSEBUTTON | qs_MOUSEMOVE | qs_HOTKEY | qs_KEY
	qs_ALLPOSTMESSAGE              = 0x100
	qs_RAWINPUT                    = 0x400
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
