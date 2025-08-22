# GLFW for Windows in pure Go 

## Introduction

This is a translation of the official glfw from C to Go.
It implements the same interface as https://github.com/go-gl/v3.3
The C code comes from https://github.com/glfw/glfw

Only Windows is currently because the other platforms can easily use the
orignal glfw. On Windows there is often problems using the CGO functions.
This code can be used without any C compiler. It is pure Go.

Some of the original tests are also translated, and seems to run fine.

The software is mostly complete. Some functions may be missing.
Please report any errors found.

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

- Only standard OpenGl is supported. No Vulkan or OpenGL ES.
- Only Windows 10 or later is supported (Can perhaps work on Windows 8).
- Monitor connect/disconnect is not supported while the app is running
- Joystick is not supported
