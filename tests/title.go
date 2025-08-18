package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/jkvatne/jkvgui/gl"
	glfw "github.com/jkvatne/purego-glfw"
)

// UTF-8 window title test

func title() {
	fmt.Printf("Windows title should show non-latin utf characters\n")
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	window, err := glfw.CreateWindow(800, 400, "Temp", nil, nil)
	if err != nil {
		glfw.Terminate()
		os.Exit(1)
	}
	window.MakeContextCurrent()
	gl.Init()
	glfw.SwapInterval(1)
	glfw.SetTime(0)
	window.SetTitle("English 日本語 русский язык 官話")
	for !window.ShouldClose() && glfw.GetTime() < 2.0 {
		gl.ClearColor(100, 100, 0, 256)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		window.SwapBuffers()
		glfw.PollEvents()
	}
	window.Destroy()
	glfw.Terminate()
	fmt.Printf("Window title test finished\n")
}
