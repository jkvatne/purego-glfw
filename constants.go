package glfw

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
	_TRUE                   = 1
	pm_REMOVE               = 0x0001
	pm_NOREMOVE             = 0x0000
	wm_CANCELMODE           = 0x001F
	wm_CHAR                 = 0x0102
	wm_SYSCHAR              = 0x0106
	wm_CLOSE                = 0x0010
	wm_CREATE               = 0x0001
	wm_DPICHANGED           = 0x02E0
	wm_DESTROY              = 0x0002
	wm_ERASEBKGND           = 0x0014
	wm_GETMINMAXINFO        = 0x0024
	wm_IME_COMPOSITION      = 0x010F
	wm_IME_ENDCOMPOSITION   = 0x010E
	wm_IME_STARTCOMPOSITION = 0x010D
	wm_KEYDOWN              = 0x0100
	wm_KEYUP                = 0x0101
	wm_KILLFOCUS            = 0x0008
	wm_LBUTTONDOWN          = 0x0201
	wm_LBUTTONUP            = 0x0202
	wm_MBUTTONDOWN          = 0x0207
	wm_MBUTTONUP            = 0x0208
	wm_MOUSEMOVE            = 0x0200
	wm_MOUSEWHEEL           = 0x020A
	wm_MOUSEHWHEEL          = 0x020E
	wm_NCACTIVATE           = 0x0086
	wm_NCHITTEST            = 0x0084
	wm_NCCALCSIZE           = 0x0083
	wm_PAINT                = 0x000F
	wm_QUIT                 = 0x0012
	wm_SETCURSOR            = 0x0020
	wm_SETFOCUS             = 0x0007
	wm_SHOWWINDOW           = 0x0018
	wm_SIZE                 = 0x0005
	wm_STYLECHANGED         = 0x007D
	wm_SYSKEYDOWN           = 0x0104
	wm_SYSKEYUP             = 0x0105
	wm_RBUTTONDOWN          = 0x0204
	wm_RBUTTONUP            = 0x0205
	wm_TIMER                = 0x0113
	wm_UNICHAR              = 0x0109
	wm_USER                 = 0x0400
	wm_WINDOWPOSCHANGED     = 0x0047
	_UNICODE_NOCHAR         = 65535
	cw_USEDEFAULT           = -2147483648

	ws_CLIPCHILDREN         = 0x02000000
	ws_CLIPSIBLINGS         = 0x04000000
	ws_MAXIMIZE             = 0x01000000
	ws_ICONIC               = 0x20000000
	ws_VISIBLE              = 0x10000000
	ws_OVERLAPPED           = 0x00000000
	ws_CAPTION              = 0x00C00000
	ws_SYSMENU              = 0x00080000
	ws_THICKFRAME           = 0x00040000
	ws_MINIMIZEBOX          = 0x00020000
	ws_MAXIMIZEBOX          = 0x00010000
	ws_POPUP                = 0x80000000
	ws_OVERLAPPEDWINDOW     = ws_OVERLAPPED | ws_CAPTION | ws_SYSMENU | ws_THICKFRAME | ws_MINIMIZEBOX | ws_MAXIMIZEBOX
	ws_EX_APPWINDOW         = 0x40000
	ws_EX_TOPMOST           = 0x00000008
	ws_EX_LAYERED           = 0x00080000
	ws_EX_TRANSPARENT       = 0x00000020
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
	glfw_DONT_CARE = -1
	INT_MAX        = 0x7FFFFFFF

	glfw_RELEASE            = 0
	glfw_PRESS              = 1
	glfw_REPEAT             = 2
	glfw_OPENGL_API         = 0x00030001
	glfw_NATIVE_CONTEXT_API = 0x00036001
	glfw_OPENGL_ES_API      = 0x00030002
	glfw_EGL_CONTEXT_API    = 0x00036002
	glfw_OSMESA_CONTEXT_API = 0x00036003
	glfw_NO_API             = 0
	glfw_OPENGL_ANY_PROFILE = 0
	glfw_ANY_POSISTION      = -2147483648
	glfw_MOD_CAPS_LOCK      = 0x0010
	glfw_MOD_NUM_LOCK       = 0x0020
	glfw_CONNECTED          = 0x00040001
	glfw_DISCONNECTED       = 0x00040002
	glfw_STICK              = 3
	glfw_INSERT_FIRST       = 0
	glfw_INSERT_LAST        = 1

	glfw_NO_RESET_NOTIFICATION  = 0x00031001
	glfw_LOSE_CONTEXT_ON_RESET  = 0x00031002
	glfw_RELEASE_BEHAVIOR_FLUSH = 0x00035001
	glfw_RELEASE_BEHAVIOR_NONE  = 0x00035002

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

// Hints
const (
	glfw_OPENGL_CORE_PROFILE   = 0x00032001
	glfw_OPENGL_COMPAT_PROFILE = 0x00032002
)
