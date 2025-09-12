package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	glfw "github.com/jkvatne/purego-glfw"
	"github.com/neclepsio/gl/all-core/gl"
)

var running atomic.Bool

type Thread struct {
	window  *glfw.Window
	title   string
	r, g, b float32
	done    bool
	x, y    int
	id      uint32
}

// Multi-threading test
// This test is intended to verify whether the OpenGL context part of
// the GLFW API is able to be used from multiple threads

func key_callback5(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Press {
		return
	}
	switch key {
	case glfw.KeyEscape:
		window.SetShouldClose(true)
	}
}

var m sync.Mutex

func thread_main(self *Thread) {
	runtime.LockOSThread()
	self.window.MakeContextCurrent()
	glfw.SwapInterval(20)

	// Setup vertex shader
	var vertexShaderSource = `
		#version 110
		attribute vec2 vPos;
		void main()
		{
			gl_Position = vec4(vPos, 0.0, 1.0);	
		}
	` + "\x00"
	vertex_shader := gl.CreateShader(gl.VERTEX_SHADER)
	source, free := gl.Strs(vertexShaderSource)
	gl.ShaderSource(vertex_shader, 1, source, nil)
	free()
	gl.CompileShader(vertex_shader)

	// Setup fragment shader
	var fragmentShaderSource = `
		#version 110
		void main()
		{
			gl_FragColor = vec4(0.8,0.8,0.0,1.0);
		}
	` + "\x00"
	fragment_shader := gl.CreateShader(gl.FRAGMENT_SHADER)
	source, free = gl.Strs(fragmentShaderSource)
	gl.ShaderSource(fragment_shader, 1, source, nil)
	free()
	gl.CompileShader(fragment_shader)

	// Create program
	program := gl.CreateProgram()
	gl.AttachShader(program, vertex_shader)
	gl.AttachShader(program, fragment_shader)
	gl.LinkProgram(program)

	// Create and bind buffers
	var vertex_buffer uint32
	gl.GenBuffers(1, &vertex_buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertex_buffer)

	// Create and setup attribute (vPos)
	vposLocation := gl.GetAttribLocation(program, gl.Str("vPos\x00"))
	gl.EnableVertexAttribArray(uint32(vposLocation))
	gl.VertexAttribPointer(uint32(vposLocation), 2, gl.FLOAT, false, 4, nil)

	for running.Load() {
		// The Thread struct is shared and must be protected by a mutex.
		m.Lock()
		// Clear screen to a color unique to this window
		a := float32(math.Abs(math.Sin(glfw.GetTime() * 2)))
		gl.ClearColor(self.r*a, self.g*a, self.b*a, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		// Draw a triangle
		var v = [3][2]float32{{-a, -a}, {a, a}, {a, -a}}
		gl.UseProgram(program)
		gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(v)), unsafe.Pointer(&v), gl.STREAM_DRAW)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		gl.UseProgram(0)
		// Do end of frame housekeeping
		self.window.SwapBuffers()
		self.id = glfw.GetCurrentThreadId()
		m.Unlock()
		time.Sleep(time.Millisecond * 10)
	}
}

func threads() {
	var threadDefs = []Thread{
		{title: "Red", r: 1, g: 0, b: 0, x: 50, y: 50},
		{title: "Green", r: 0, g: 1, b: 0, x: 50, y: 250},
		{title: "Blue", r: 0, g: 0, b: 1, x: 50, y: 450},
		{title: "White", r: 1, g: 1, b: 1, x: 50, y: 650},
		{title: "Yellow", r: 1, g: 1, b: 0, x: 50, y: 850},
		{title: "Cyan", r: 0, g: 1, b: 1, x: 50, y: 1050},
		{title: "Magenta", r: 1, g: 0, b: 1, x: 50, y: 1250},
		{title: "Red", r: 1, g: 0, b: 0, x: 550, y: 50},
		{title: "Green", r: 0, g: 1, b: 0, x: 550, y: 250},
		{title: "Blue", r: 0, g: 0, b: 1, x: 550, y: 450},
		{title: "White", r: 1, g: 1, b: 1, x: 550, y: 650},
		{title: "Yellow", r: 1, g: 1, b: 0, x: 550, y: 850},
		{title: "Cyan", r: 0, g: 1, b: 1, x: 550, y: 1050},
		{title: "Magenta", r: 1, g: 0, b: 1, x: 550, y: 1250},
		{title: "Red", r: 1, g: 0, b: 0, x: 1050, y: 50},
		{title: "Green", r: 0, g: 1, b: 0, x: 1050, y: 250},
		{title: "Blue", r: 0, g: 0, b: 1, x: 1050, y: 450},
		{title: "White", r: 1, g: 1, b: 1, x: 1050, y: 650},
		{title: "Yellow", r: 1, g: 1, b: 0, x: 1050, y: 850},
		{title: "Cyan", r: 0, g: 1, b: 1, x: 1050, y: 1050},
		{title: "Magenta", r: 1, g: 0, b: 1, x: 1050, y: 1250},
		{title: "Red", r: 1, g: 0, b: 0, x: 1550, y: 50},
		{title: "Green", r: 0, g: 1, b: 0, x: 1550, y: 250},
		{title: "Blue", r: 0, g: 0, b: 1, x: 1550, y: 450},
		{title: "White", r: 1, g: 1, b: 1, x: 1550, y: 650},
		{title: "Yellow", r: 1, g: 1, b: 0, x: 1550, y: 850},
		{title: "Cyan", r: 0, g: 1, b: 1, x: 1550, y: 1050},
		{title: "Magenta", r: 1, g: 0, b: 1, x: 1550, y: 1250},
		{title: "Red", r: 1, g: 0, b: 0, x: 2050, y: 50},
		{title: "Green", r: 0, g: 1, b: 0, x: 2050, y: 250},
		{title: "Blue", r: 0, g: 0, b: 1, x: 2050, y: 450},
		{title: "White", r: 1, g: 1, b: 1, x: 2050, y: 650},
		{title: "Yellow", r: 1, g: 1, b: 0, x: 2050, y: 850},
		{title: "Cyan", r: 0, g: 1, b: 1, x: 2050, y: 1050},
		{title: "Magenta", r: 1, g: 0, b: 1, x: 2050, y: 1250},
	}
	count := len(threadDefs)
	runtime.LockOSThread()
	runtime.GOMAXPROCS(1)
	fmt.Printf("\nTesting drawing windows in different goroutines\n")
	fmt.Printf("Window count %d\n", count)
	fmt.Printf("CPU count=%d\n", runtime.NumCPU())
	fmt.Printf("ProcCount=%d\n", runtime.GOMAXPROCS(0))
	glfw.SetErrorCallback(error_callback)
	err := glfw.Init()
	if err != nil {
		panic(err.Error())
	}
	for i := 0; i < count; i++ {
		_ = glfw.WindowHint(glfw.PositionX, threadDefs[i].x)
		_ = glfw.WindowHint(glfw.PositionY, threadDefs[i].y)
		threadDefs[i].window, err = glfw.CreateWindow(400, 120, threadDefs[i].title, nil, nil)
		if threadDefs[i].done {
			glfw.Terminate()
			os.Exit(1)
		}
		threadDefs[i].window.SetKeyCallback(key_callback5)
	}

	threadDefs[0].window.MakeContextCurrent()
	err = gl.Init()
	if err != nil {
		glfw.Terminate()
		fmt.Printf("gl Init error, " + err.Error())
	}
	glfw.DetachCurrentContext()
	running.Store(true)
	for i := 0; i < count; i++ {
		go thread_main(&threadDefs[i])
	}

	glfw.SetTime(0)
	for running.Load() && glfw.GetTime() < 4.0 {
		glfw.PollEvents()
		for i := 0; i < count; i++ {
			if threadDefs[i].window.ShouldClose() || glfw.GetTime() > 4.0 {
				break
			}
		}
		time.Sleep(time.Millisecond * 1000)
	}
	running.Store(false)
	time.Sleep(time.Millisecond * 100)
	fmt.Printf("Thread IDs\n")
	for i := 0; i < count; i++ {
		m.Lock()
		fmt.Printf("%7d", threadDefs[i].id)
		m.Unlock()
		threadDefs[i].window.Destroy()
	}
	fmt.Printf("\n")
}
