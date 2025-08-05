package main

import (
	"flag"
	"runtime"
)

var testname string

func main() {
	flag.StringVar(&testname, "t", "", "Test can be: title, msaa, reopen")
	flag.Parse()
	runtime.LockOSThread()
	switch testname {
	case "title":
		title()
	case "reopen":
		reopen()
	case "msaa":
		msaa()
	}
}
