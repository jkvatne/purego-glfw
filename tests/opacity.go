package main

// Test opacity of window

import (
	"fmt"
	"os"
	"runtime"
	"time"
	"unsafe"

	glfw "github.com/jkvatne/purego-glfw"
	"github.com/jkvatne/purego-glfw/gl"
)

func OpacityMain() {
	fmt.Printf("\nThis windows should gradually fade in and out\n")
	runtime.LockOSThread()
	// Initialize the glfw library
	err := glfw.Init()
	if err != nil {
		panic(err.Error())
	}
	defer glfw.Terminate()
	glfw.SetErrorCallback(error_callback)

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	// Need to make initial window invisible in order to start with opacity=0
	// If not, the window will show briefly with full intensity.
	glfw.WindowHint(glfw.Visible, 0)
	window, err := glfw.CreateWindow(400, 400, "Transparent", nil, nil)
	if err != nil {
		glfw.Terminate()
		os.Exit(1)
	}

	// Initialize Open-gl on current window
	window.MakeContextCurrent()
	err = gl.Init()
	if err != nil {
		glfw.Terminate()
		fmt.Printf("Could not init gl: %v\n", err)
		os.Exit(2)
	}
	glfw.SwapInterval(1)

	// Setup vertex shader
	var vertexShaderSource = `
		#version 110
		attribute vec2 pos;
		void main()
		{
			gl_Position = vec4(pos, 0.0, 1.0);	
		}
	` + "\x00"
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	source, free := gl.Strs(vertexShaderSource)
	gl.ShaderSource(vertexShader, 1, source, nil)
	free()
	gl.CompileShader(vertexShader)

	// Setup fragment shader
	var fragmentShaderSource = `
		#version 110
		void main()
		{
			gl_FragColor = vec4(1.0);
		}
	` + "\x00"
	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	source, free = gl.Strs(fragmentShaderSource)
	gl.ShaderSource(fragmentShader, 1, source, nil)
	free()
	gl.CompileShader(fragmentShader)

	// Create program
	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	// Create and bind buffers
	var vertexBuffer uint32
	gl.GenBuffers(1, &vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)

	// Create and setup attribute (vPos)
	vPosLocation := gl.GetAttribLocation(program, gl.Str("pos\x00"))
	gl.EnableVertexAttribArray(uint32(vPosLocation))
	gl.VertexAttribPointer(uint32(vPosLocation), 2, gl.FLOAT, false, 4, nil)

	glfw.SetTime(0)
	// Set opacity to 0 before showing the window to make it gradually appear
	window.SetOpacity(0)
	window.Show()
	for !window.ShouldClose() && glfw.GetTime() < 5.0 {
		t := glfw.GetTime()
		// Ramp up over 1.5 seconds, then max in two second and ramp down over 1.5 seconds
		o := min(t/1.5, (5.0-t)/1.5, 1.0)
		window.SetOpacity(o)

		// Clear screen to a brown color
		gl.ClearColor(0.5, 0.5, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Draw a triangle
		var v = [3][2]float32{{0.0, 0.7}, {-0.7, -0.7}, {0.7, -0.7}}
		gl.UseProgram(program)
		gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(v)), unsafe.Pointer(&v), gl.STREAM_DRAW)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		gl.UseProgram(0)

		// Do end of frame housekeeping
		window.SwapBuffers()
		glfw.PollEvents()
		time.Sleep(time.Millisecond * 20)
	}
	window.Destroy()
	fmt.Printf("Now the window should have disappeared completely\n")
	glfw.Terminate()
}
