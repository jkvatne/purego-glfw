package main

// Test opacity of window

import (
	"fmt"
	"os"
	"runtime"
	"time"

	glfw "github.com/jkvatne/purego-glfw"
	"github.com/neclepsio/gl/all-core/gl"
)

func opacity() {
	fmt.Printf("\nThis windows should gradualy be transparent\n")
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	window, err := glfw.CreateWindow(800, 400, "Transparent", nil, nil)
	if err != nil {
		glfw.Terminate()
		os.Exit(1)
	}
	window.MakeContextCurrent()
	err = gl.Init()
	if err != nil {
		glfw.Terminate()
		fmt.Printf("Could not init gl: %v\n", err)
		os.Exit(2)
	}
	glfw.SwapInterval(1)
	glfw.SetTime(0)
	FadeTime := 2.0
	for !window.ShouldClose() && glfw.GetTime() < FadeTime {
		gl.ClearColor(100, 100, 0, 256)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		window.SwapBuffers()
		glfw.PollEvents()
		window.SetOpacity(max(0.0, 1.0-glfw.GetTime()/FadeTime))
		time.Sleep(time.Millisecond * 100)
	}
	window.Destroy()
	glfw.Terminate()
}
