package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"time"

	"github.com/jkvatne/jkvgui/gl"
	glfw "github.com/jkvatne/purego-glfw"
)

var running bool

type Thread struct {
	window  *glfw.Window
	title   string
	r, g, b float32
	done    bool
	x, y    int32
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
		window.SetWindowShouldClose(true)
	}
}

func thread_main(self *Thread) {
	runtime.LockOSThread()
	self.window.MakeContextCurrent()
	glfw.SwapInterval(20)
	fmt.Printf("ThreadId=%d\n", glfw.GetCurrentThreadId())
	for running {
		v := float32(math.Abs(math.Sin(glfw.GetTime() * 2)))
		gl.ClearColor(self.r*v, self.g*v, self.b*v, 0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		self.window.SwapBuffers()
		self.id = glfw.GetCurrentThreadId()
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
	fmt.Printf("Window count %d\n", count)
	fmt.Printf("CPU count=%d\n", runtime.NumCPU())
	fmt.Printf("ProcCount=%d\n", runtime.GOMAXPROCS(0))
	// glfwSetErrorCallback(error_callback);
	err := glfw.Init()
	if err != nil {
		panic(err.Error())
	}
	for i := 0; i < count; i++ {
		glfw.WindowHint(glfw.POSITION_X, threadDefs[i].x)
		glfw.WindowHint(glfw.POSITION_Y, threadDefs[i].y)
		threadDefs[i].window, err = glfw.CreateWindow(400, 120, threadDefs[i].title, nil, nil)
		if threadDefs[i].done {
			glfw.Terminate()
			os.Exit(1)
		}
		threadDefs[i].window.SetKeyCallback(key_callback)
	}

	threadDefs[0].window.MakeContextCurrent()
	gl.Init()
	glfw.DetachCurrentContext()
	running = true
	for i := 0; i < count; i++ {
		go thread_main(&threadDefs[i])
	}

	for running {
		glfw.PollEvents()
		for i := 0; i < count; i++ {
			if threadDefs[i].window.ShouldClose() {
				running = false
			}
			fmt.Printf("%7d", threadDefs[i].id)

		}
		fmt.Printf("\n")
		time.Sleep(time.Millisecond * 1000)
	}
	for i := 0; i < count; i++ {
		threadDefs[i].window.Destroy()
	}
	os.Exit(0)

}
