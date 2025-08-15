package main

import (
	"fmt"
	"math"
	"runtime"
	"unsafe"

	"github.com/jkvatne/jkvgui/gl"
	glfw "github.com/jkvatne/purego-glfw"
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
	var buffer [64 * 64 * 4]uint8
	var image = glfw.GLFWimage{64, 64, &buffer[0]}
	for y := int32(0); y < image.Width; y++ {
		for x := int32(0); x < image.Height; x++ {
			buffer[x*image.Width+y] = 255
			buffer[x*image.Width+y] = 255
			buffer[x*image.Width+y] = 255
			buffer[x*image.Width+y] = uint8(255.0 * star(float64(x), float64(y), t))
		}
	}
	return glfw.CreateCursor(&image, image.Width/2, image.Height/2)
}

func create_tracking_cursor() *glfw.Cursor {
	var buffer [32 * 32 * 4]uint8
	var image = glfw.GLFWimage{32, 32, &buffer[0]}
	var i int
	for y := int32(0); y < image.Width; y++ {
		for x := int32(0); x < image.Height; x++ {
			if x == 7 || y == 7 {
				buffer[i] = 255
				i++
				buffer[i] = 0
				i++
				buffer[i] = 0
				i++
				buffer[i] = 255
				i++
			} else {
				buffer[i] = 0
				i++
				buffer[i] = 0
				i++
				buffer[i] = 0
				i++
				buffer[i] = 0
				i++
			}
		}
	}
	return glfw.CreateCursor(&image, 7, 7)
}

func cursor_position_callback(window *glfw.Window, x float64, y float64) {
	fmt.Printf("%0.3f: Cursor position: %f %f (%+f %+f)\n",
		glfw.GetTime(),
		x, y, x-cursor_x, y-cursor_y)
	cursor_x = x
	cursor_y = y
}

func key_callback9(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
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
			window.SetWindowShouldClose(true)
		}
	case glfw.KeyN:
		window.SetInputMode(glfw.CURSOR_MODE, glfw.CURSOR_NORMAL)
		cursor_x, cursor_y = window.GetCursorPos()
		fmt.Printf("( cursor is normal )\n")
	case glfw.KeyD:
		window.SetInputMode(glfw.CURSOR_MODE, glfw.CURSOR_DISABLED)
		fmt.Printf("( cursor is disabled )\n")
	case glfw.KeyH:
		window.SetInputMode(glfw.CURSOR_MODE, glfw.CURSOR_HIDDEN)
		fmt.Printf("( cursor is hidden )\n")
	case glfw.KeyC:
		window.SetInputMode(glfw.CURSOR_MODE, glfw.CURSOR_CAPTURED)
		fmt.Printf("( cursor is captured )\n")
	case glfw.KeyR:
		if !glfw.RawMouseMotionSupported() {
			break
		}
		if window.GetInputMode(glfw.RAW_MOUSE_MOTION) != 0 {
			window.SetInputMode(glfw.RAW_MOUSE_MOTION, 0)
			fmt.Printf("( raw input is disabled )\n")
		} else {
			window.SetInputMode(glfw.RAW_MOUSE_MOTION, 1)
			fmt.Printf("( raw input is enabled )\n")
		}
	case glfw.KeySpace:
		swap_interval = 1 - swap_interval
		fmt.Printf("( swap interval: %i ))\n", swap_interval)
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
	case glfw.Key1:
	case glfw.Key2:
	case glfw.Key3:
	case glfw.Key4:
	case glfw.Key5:
	case glfw.Key6:
	case glfw.Key7:
	case glfw.Key8:
	case glfw.Key9:
		index := int(key - glfw.Key1)
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
		/*
		   if (glfwGetWindowMonitor(window))
		       glfwSetWindowMonitor(window, NULL, x, y, width, height, 0);
		   else
		   {
		       GLFWmonitor* monitor = glfwGetPrimaryMonitor();
		       const GLFWvidmode* mode = glfwGetVideoMode(monitor);
		       glfwGetWindowPos(window, &x, &y);
		       glfwGetWindowSize(window, &width, &height);
		       glfwSetWindowMonitor(window, monitor,
		                            0, 0, mode->width, mode->height,
		                            mode->refreshRate);
		   }

		   glfwGetCursorPos(window, &cursor_x, &cursor_y);
		*/
	}
}

type vec4 [4]float32

func cursor() {
	type mat4x4 [4]vec4
	var (
		star_cursors                                           [CURSOR_FRAME_COUNT]*glfw.Cursor
		current_frame                                          *glfw.Cursor
		vertex_buffer, vertex_shader, fragment_shader, program uint32
		mvp                                                    mat4x4
	)

	fmt.Println("Cursor test")
	runtime.LockOSThread()
	err := glfw.Init()
	if err != nil {
		panic(err.Error())
	}

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
	}

	for i := 0; i < len(shapes); i++ {
		standard_cursors[i] = glfw.CreateStandardCursor(shapes[i])
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)

	window, err := glfw.CreateWindow(640, 480, "Cursor Test", nil, nil)
	if err != nil {
		glfw.Terminate()
		panic(err.Error())
	}

	window.MakeContextCurrent()
	gl.Init()
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

	cursor_x, cursor_y = window.GetCursorPos()
	fmt.Printf("Cursor position: %f %f\n", cursor_x, cursor_y)
	window.SetCursorPosCallback(cursor_position_callback)
	window.SetKeyCallback(key_callback)
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		if track_cursor {
			wnd_width, _ := window.GetSize()
			fb_width, fb_height := window.GetFramebufferSize()
			gl.Viewport(0, 0, fb_width, fb_height)
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
			// mat4x4_ortho(mvp, 0.0, float32(fb_width), 0.0, float32(fb_height), 0.0, 1.0);
			gl.UniformMatrix4fv(mvp_location, 1, false, &mvp[0][0])
			gl.DrawArrays(gl.LINES, 0, 4)
		}

		window.SwapBuffers()
		if animate_cursor {
			i := int((glfw.GetTime() * 30.0)) % CURSOR_FRAME_COUNT
			if current_frame != star_cursors[i] {
				window.SetCursor(star_cursors[i])
				current_frame = star_cursors[i]
			}
		} else {
			current_frame = nil
		}
		if wait_events {
			if animate_cursor {
				glfw.WaitEventsTimeout(1.0 / 30.0)
			} else {
				glfw.WaitEvents()
			}
		} else {
			glfw.PollEvents()
		}
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
