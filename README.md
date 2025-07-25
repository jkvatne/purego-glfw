# GLFW 3.3 for pure Go 

## Introduction

This is a translation of the official glfw for Go, from C to Go

Only Windows is currently because the other platforms can easily use the
orignal glfw. On Windows there is often problems using the CGO functions.
This code can be used without any C compiler. It is pure Go.

NB: The software is not production ready, and may contain serious errors.
Some functions may be missing. Please report any errors found.

## Installation

```
go get -u github.com/jkvatne/purego-glfw 
```

## Dependencies
```
golang.design/x/clipboard
```
This is needed for clipboard support on windows.
- clipbord.Init()
- clipboard.Read()
- clibport.Write()

```
golang.org/x/sys/windows
```
This is used to import the types LazyDLL and LazyProc.

## Limitations

- Only standard OpenGl revision 3 is supported. No Vulcan or OpenGL ES.
- Only Windows 10 or later is supported (Can perhaps work on Windows 7 and 8).
- Only the default video mode is supported. Full-screen apps will use the system configurated resolution. (Usualy the native monitor resolution)

