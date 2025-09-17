package main

import (
	"fmt"
	"image"
	"math"
	"runtime"
	"unsafe"

	glfw "github.com/jkvatne/purego-glfw"
	"github.com/neclepsio/gl/all-core/gl"
)

// Cursor & input mode tests
// This test provides an interface to the cursor image and cursor mode
// parts of the API.
// Custom cursor image generation by urraka.

const CURSOR_FRAME_COUNT = 60

var vertex_shader_text3 = `
#version 110
uniform mat4 MVP;
attribute vec2 vPos;
void main()
{
    gl_Position = MVP * vec4(vPos, 0.0, 1.0);	
}
` + "\x00"

var fragment_shader_text3 = `
#version 110
void main()
{
    gl_FragColor = vec4(1.0);
}
` + "\x00"

var (
	cursor_x         float64
	cursor_y         float64
	swap_interval    = 1
	wait_events      = true
	animate_cursor   = false
	track_cursor     = false
	standard_cursors [10]*glfw.Cursor
	tracking_cursor  *glfw.Cursor
	hasKeyPress      bool
)

func star(x float64, y float64, t float64) float64 {
	c := 64.0 / 2.0
	i := 0.25*math.Sin(2*math.Pi*t) + 0.75
	k := 64 * 0.046875 * i
	dist := math.Sqrt((x-c)*(x-c) + (y-c)*(y-c))
	salpha := 1.0 - dist/c
	xalpha := c
	if x != c {
		xalpha = k / math.Abs(x-c)
	}
	yalpha := c
	if y != c {
		yalpha = k / math.Abs(y-c)
	}
	return max(0.0, min(1.0, i*salpha*0.2+salpha*xalpha*yalpha))
}

func create_cursor(t float64) *glfw.Cursor {
	img := image.NewNRGBA(image.Rect(0, 0, 64, 64))
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			img.Pix[(y*img.Bounds().Dx()+x)*4] = 255
			img.Pix[(y*img.Bounds().Dx()+x)*4+1] = 255
			img.Pix[(y*img.Bounds().Dx()+x)*4+2] = 0
			img.Pix[(y*img.Bounds().Dx()+x)*4+3] = uint8(255.0 * star(float64(x), float64(y), t))
		}
	}
	return glfw.CreateCursor(img, 32, 32)
}

func cursor_position_callback(window *glfw.Window, x float64, y float64) {
	fmt.Printf("%0.3f: Cursor position callback: %f %f (%+f %+f)\n",
		glfw.GetTime(),
		x, y, x-cursor_x, y-cursor_y)
	cursor_x = x
	cursor_y = y
}

var usage = `

Testing cursors
---------------
If no key is pressed, the program will loop through all 
available cursors and disable/hide/show the cursor.
Then it will terminate.

A Animated yellow star cursor
N Normal cursor mode
D Disabled (invisible) cursor mode
H Hidden cursor mode
C Captured cursor mode
R Raw mouse mode
T Track cusor with a big cross
Up Move cursor to upper left corner
Dn Move cursor to lower right corner
0-9 Show the different standard cursors
Esc Exit
`
var x, y, w, h int
var index int

func key_callback_cursor(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	fmt.Printf("Key %v, action %v, mods %v\n", key, action, mods)
	hasKeyPress = true
	if action != glfw.Press {
		return
	}
	switch key {
	case glfw.KeyA:
		animate_cursor = !animate_cursor
		if !animate_cursor {
			window.SetCursor(nil)
		}
	case glfw.KeyEscape:
		mode := window.GetInputMode(glfw.CursorMode)
		if mode != glfw.CursorDisabled && mode != glfw.CursorCaptured {
			window.SetShouldClose(true)
		}
	case glfw.KeyN:
		window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		cursor_x, cursor_y = window.GetCursorPos()
		fmt.Printf("( cursor is normal )\n")
	case glfw.KeyD:
		window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
		fmt.Printf("( cursor is disabled )\n")
	case glfw.KeyH:
		window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
		fmt.Printf("( cursor is hidden )\n")
	case glfw.KeyC:
		window.SetInputMode(glfw.CursorMode, glfw.CursorCaptured)
		fmt.Printf("( cursor is captured )\n")
	case glfw.KeyR:
		if !glfw.RawMouseMotionSupported() {
			break
		}
		if window.GetInputMode(glfw.RawMouseMotion) != 0 {
			window.SetInputMode(glfw.RawMouseMotion, 0)
			fmt.Printf("( raw input is disabled )\n")
		} else {
			window.SetInputMode(glfw.RawMouseMotion, 1)
			fmt.Printf("( raw input is enabled )\n")
		}
	case glfw.KeySpace:
		swap_interval = 1 - swap_interval
		fmt.Printf("( swap interval: %d ))\n", swap_interval)
		glfw.SwapInterval(swap_interval)
	case glfw.KeyW:
		wait_events = !wait_events
		if wait_events {
			fmt.Printf("( waiting for events )\n")
		} else {
			fmt.Printf("( polling for events )\n")
		}
	case glfw.KeyT:
		track_cursor = !track_cursor
		if track_cursor {
			window.SetCursor(tracking_cursor)
		} else {
			window.SetCursor(nil)
		}
	case glfw.KeyP:
		x, y := window.GetCursorPos()
		fmt.Printf("Query before set: %f %f (%+f %+f)\n",
			x, y, x-cursor_x, y-cursor_y)
		cursor_x = x
		cursor_y = y
		window.SetCursorPos(cursor_x, cursor_y)
		x, y = window.GetCursorPos()
		fmt.Printf("Query after set: %f %f (%+f %+f)\n",
			x, y, x-cursor_x, y-cursor_y)
		cursor_x = x
		cursor_y = y
	case glfw.KeyUp:
		window.SetCursorPos(0, 0)
		cursor_x, cursor_y = window.GetCursorPos()
	case glfw.KeyDown:
		width, height := window.GetSize()
		window.SetCursorPos(float64(width-1), float64(height-1))
		cursor_x, cursor_y = window.GetCursorPos()
	case glfw.Key0, glfw.Key1, glfw.Key2, glfw.Key3, glfw.Key4, glfw.Key5, glfw.Key6, glfw.Key7, glfw.Key8, glfw.Key9:
		index = int(key - glfw.Key0)
		if mods&glfw.ModShift != 0 {
			index += 9
		}
		if index < len(standard_cursors) {
			window.SetCursor(standard_cursors[index])
		}
	case glfw.KeyEnter:
		if mods != glfw.ModAlt {
			return
		}
		if window.GetMonitor() != nil {
			// Set windowed
			window.SetMonitor(nil, x, y, w, h, 0)
		} else {
			// Set windowed
			m := glfw.GetPrimaryMonitor()
			mode := m.GetVideoMode()
			x, y = window.GetPos()
			w, h = window.GetSize()
			window.SetMonitor(m, 0, 0, int(mode.Width), int(mode.Height), int(mode.RefreshRate))
		}
		cursor_x, cursor_y = window.GetCursorPos()
	}
}

func CursorMain() {
	runtime.LockOSThread()
	fmt.Printf(usage)

	// type mat4x4 [4]vec4
	var (
		star_cursors                                           [CURSOR_FRAME_COUNT]*glfw.Cursor
		current_frame                                          *glfw.Cursor
		vertex_buffer, vertex_shader, fragment_shader, program uint32
		mvp                                                    Mat4
	)

	fmt.Println("Cursor test")
	// Initialize the glfw library
	err := glfw.Init()
	if err != nil {
		panic(err.Error())
	}
	defer glfw.Terminate()
	glfw.SetErrorCallback(error_callback)

	// tracking_cursor := create_tracking_cursor()
	for i := 0; i < CURSOR_FRAME_COUNT; i++ {
		star_cursors[i] = create_cursor(float64(i) / CURSOR_FRAME_COUNT)
	}
	shapes := []int{
		glfw.ArrowCursor,
		glfw.IBeamCursor,
		glfw.CrosshairCursor,
		glfw.HResizeCursor,
		glfw.VResizeCursor,
		glfw.HandCursor,
		glfw.ResizeNwseCursor,
		glfw.ResizeNeswCursor,
		glfw.ResizeAllCursor,
		glfw.NotAllowedCursor,
	}

	for i := 0; i < len(shapes); i++ {
		standard_cursors[i] = glfw.CreateStandardCursor(shapes[i])
	}

	_ = glfw.WindowHint(glfw.ContextVersionMajor, 2)
	_ = glfw.WindowHint(glfw.ContextVersionMinor, 0)
	window, err := glfw.CreateWindow(640, 480, "Cursor Test", nil, nil)
	if err != nil {
		glfw.Terminate()
		panic(err.Error())
	}

	window.MakeContextCurrent()
	err = gl.Init()
	if err != nil {
		glfw.Terminate()
		fmt.Printf("gl Init error, " + err.Error())
	}
	gl.GenBuffers(1, &vertex_buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertex_buffer)
	vertex_shader = gl.CreateShader(gl.VERTEX_SHADER)
	csources, free := gl.Strs(vertex_shader_text3)
	gl.ShaderSource(vertex_shader, 1, csources, nil)
	free()
	gl.CompileShader(vertex_shader)

	fragment_shader = gl.CreateShader(gl.FRAGMENT_SHADER)
	csources, free = gl.Strs(fragment_shader_text3)
	gl.ShaderSource(fragment_shader, 1, csources, nil)
	free()
	gl.CompileShader(fragment_shader)

	program = gl.CreateProgram()
	gl.AttachShader(program, vertex_shader)
	gl.AttachShader(program, fragment_shader)
	gl.LinkProgram(program)

	mvp_location := gl.GetUniformLocation(program, gl.Str("MVP\x00"))
	vpos_location := gl.GetAttribLocation(program, gl.Str("vPos\x00"))

	gl.EnableVertexAttribArray(uint32(vpos_location))
	gl.VertexAttribPointer(uint32(vpos_location), 2, gl.FLOAT, false, 8, nil)
	gl.UseProgram(program)

	fmt.Println("Centers cursor in window")
	x, y := window.GetSize()
	window.SetCursorPos(float64(x/2), float64(y/2))

	cursor_x, cursor_y = window.GetCursorPos()
	fmt.Printf("Cursor position: %f %f\n", cursor_x, cursor_y)
	window.SetCursorPosCallback(cursor_position_callback)
	window.SetKeyCallback(key_callback_cursor)
	glfw.SetTime(0)
	for !window.ShouldClose() {
		if glfw.GetTime() > 1.0 {
			glfw.SetTime(0)
			index = index + 1
			if !hasKeyPress {
				if index == len(standard_cursors) {
					fmt.Println("Cursor disabled")
					window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
				} else if index == len(standard_cursors)+1 {
					fmt.Println("Cursor hidden")
					window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
				} else if index == len(standard_cursors)+2 {
					fmt.Println("Cursor normal")
					window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
				} else if index == len(standard_cursors)+3 {
					fmt.Println("Cursor captured")
					window.SetInputMode(glfw.CursorMode, glfw.CursorCaptured)
				} else if index >= len(standard_cursors)+4 {
					fmt.Println("Cursor test finished")
					break
				} else {
					fmt.Printf("Testing cursor type %d\n", index)
				}
			}
			if index >= len(standard_cursors) {
				animate_cursor = true
			} else {
				window.SetCursor(standard_cursors[index])
				glfw.PostEmptyEvent()
			}
		}
		gl.Clear(gl.COLOR_BUFFER_BIT)
		if track_cursor {
			wnd_width, _ := window.GetSize()
			fb_width, fb_height := window.GetFramebufferSize()
			gl.Viewport(0, 0, int32(fb_width), int32(fb_height))
			scale := float32(fb_width) / float32(wnd_width)
			vertices[0][0] = 0.5
			vertices[0][1] = float32(fb_height) - float32(cursor_y)*scale - 1.0 + 0.5
			vertices[1][0] = float32(fb_width) + 0.5
			vertices[1][1] = float32(fb_height) - float32(cursor_y)*scale - 1.0 + 0.5
			vertices[2][0] = (float32(cursor_x) * scale) + 0.5
			vertices[2][1] = 0.5
			vertices[3][0] = (float32(cursor_x) * scale) + 0.5
			vertices[3][1] = float32(fb_height) + 0.5

			gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices)), unsafe.Pointer(&vertices), gl.STREAM_DRAW)
			mvp = mat4x4_ortho(0.0, float32(fb_width), 0.0, float32(fb_height), 0.0, 1.0)
			gl.UniformMatrix4fv(mvp_location, 1, false, &mvp[0])
			gl.DrawArrays(gl.LINES, 0, 4)
		}

		window.SwapBuffers()
		if animate_cursor {
			i := int(glfw.GetTime()*30.0) % CURSOR_FRAME_COUNT
			if current_frame != star_cursors[i] {
				window.SetCursor(star_cursors[i])
				current_frame = star_cursors[i]
			}
		} else {
			current_frame = nil
		}
		glfw.WaitEventsTimeout(1.0 / 30.0)
	}
	window.Destroy()
	for i := 0; i < CURSOR_FRAME_COUNT; i++ {
		glfw.DestroyCursor(star_cursors[i])
	}
	for i := 0; i < len(standard_cursors); i++ {
		glfw.DestroyCursor(standard_cursors[i])
	}
	glfw.Terminate()
}
