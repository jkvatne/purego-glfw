package main

import (
	"flag"
)

var testname string

func main() {
	flag.StringVar(&testname, "t", "", "Test can be: title, msaa, reopen, timeout, monitor")
	flag.Parse()
	switch testname {
	case "title":
		title()
	case "reopen":
		reopen()
	case "msaa":
		msaa()
	case "timeout":
		timeout()
	case "monitor":
		monitor()
	}
}
