package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	glfw "github.com/jkvatne/purego-glfw"
	"github.com/neclepsio/gl/all-core/gl"
)

func key_callback_window(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Press {
		return
	}
	switch key {
	case glfw.KeyI:
		printInfo(window)
	}
}

func printInfo(window *glfw.Window) {
	x, y := window.GetPos()
	fmt.Printf("Window position x=%d, y=%d\n", x, y)
	w, h = window.GetFramebufferSize()
	fmt.Printf("Window Framebuffer Size:%dx%d\n", w, h)
	sx, sy := window.GetContentScale()
	fmt.Printf("ContentScale: %0.3fx%0.3f\n", sx, sy)
	l, t, r, b := window.GetFrameSize()
	fmt.Printf("Frame size: l=%d t=%d r=%d b=%d\n", l, t, r, b)
	o := window.GetOpacity()
	fmt.Printf("Opacity: %v\n", o)
	decorated := window.GetAttrib(glfw.Decorated)
	fmt.Printf("Decorated: %v\n", decorated)
	resizeable := window.GetAttrib(glfw.Resizable)
	fmt.Printf("Resizable: %v\n", resizeable)
	floating := window.GetAttrib(glfw.Floating)
	fmt.Printf("Floating: %v\n", floating)
	mousePassthrough := window.GetAttrib(glfw.MousePassthrough)
	fmt.Printf("MousePassthrough: %v\n", mousePassthrough)
	autoIconify := window.GetAttrib(glfw.AutoIconify)
	fmt.Printf("AutoIconify: %v\n", autoIconify)
	focused := window.GetAttrib(glfw.Focused)
	fmt.Printf("Focused: %v\n", focused)
	hovered := window.GetAttrib(glfw.Hovered)
	fmt.Printf("Hovered: %v\n", hovered)
	visible := window.GetAttrib(glfw.Visible)
	fmt.Printf("Visible: %v\n", visible)
	iconified := window.GetAttrib(glfw.Iconified)
	fmt.Printf("Iconified: %v\n", iconified)
	maximized := window.GetAttrib(glfw.Maximized)
	fmt.Printf("Maximized: %v\n", maximized)
	transparent := window.GetAttrib(glfw.TransparentFramebuffer)
	fmt.Printf("Transparent: %v\n", transparent)
}

func windowinfo() {
	fmt.Printf("\nPrint some window info\n")
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	window, err := glfw.CreateWindow(800, 400, "Window info", nil, nil)
	if err != nil {
		glfw.Terminate()
		os.Exit(1)
	}
	window.SetKeyCallback(key_callback_window)
	window.MakeContextCurrent()
	err = gl.Init()
	if err != nil {
		glfw.Terminate()
		fmt.Printf("Could not init gl: %v\n", err)
		os.Exit(2)
	}
	glfw.SwapInterval(1)
	glfw.SetTime(0)
	for !window.ShouldClose() && glfw.GetTime() < 3.0 {
		gl.ClearColor(0, 1, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		window.SwapBuffers()
		glfw.PollEvents()
		if glfw.GetTime() > 1.0 && glfw.GetTime() < 1.2 {
			printInfo(window)
			window.RequestAttention()
			time.Sleep(time.Millisecond * 200)
		}
	}
	window.Destroy()
	glfw.Terminate()
}
