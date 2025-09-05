package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

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
	fmt.Printf("ThreadId=%d\n", glfw.GetCurrentThreadId())
	for running.Load() {
		v := float32(math.Abs(math.Sin(glfw.GetTime() * 2)))
		m.Lock()
		gl.ClearColor(self.r*v, self.g*v, self.b*v, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
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
	for i := 0; i < count; i++ {
		m.Lock()
		fmt.Printf("%7d", threadDefs[i].id)
		m.Unlock()
		threadDefs[i].window.Destroy()
	}
}
