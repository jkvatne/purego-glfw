// This is the tests for the purego-glfw implementations.
// Type go run ./tests from the main directory.
// To avoid using gcc the "github.com/neclepsio/gl/all-core/gl" is used
// instead of the "github.com/go-gl/gl/all-core/gl" repository.
package main

import (
	"flag"
)

var testName string

func main() {
	flag.StringVar(&testName, "t", "", "string can tearing, window, cursor, threads, msaa, reopen, timeout, monitor or opacity")
	flag.Parse()
	switch testName {
	case "title":
		TitleMain()
	case "reopen":
		ReopenMain()
	case "msaa":
		MsaaMain()
	case "timeout":
		TimeoutMain()
	case "monitor":
		MonitorMain()
	case "threads":
		ThreadsMain()
	case "icon":
		IconMain()
	case "cursor":
		CursorMain()
	case "window":
		WindowInfoMain()
	case "opacity":
		OpacityMain()
	case "tearing":
		TearingMain()
	case "all":
		TitleMain()
		ThreadsMain()
		OpacityMain()
		MsaaMain()
		TimeoutMain()
		ReopenMain()
		IconMain()
		WindowInfoMain()
		MonitorMain()
		CursorMain()
	case "":
		OpacityMain()
	default:
		println("")
	}
}
