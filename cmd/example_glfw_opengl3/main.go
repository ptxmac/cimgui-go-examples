//go:build glfw
// +build glfw

package main

import (
	"fmt"
	"os"

	"github.com/AllenDang/cimgui-go"
	"github.com/ptxmac/cimgui-go-examples/internal/example"
	"github.com/ptxmac/cimgui-go-examples/internal/platforms"
	"github.com/ptxmac/cimgui-go-examples/internal/renderers"
)

func main() {
	context := cimgui.CreateContext(0)

	defer context.Destroy()
	io := cimgui.GetIO()

	platform, err := platforms.NewGLFW(io, platforms.GLFWClientAPIOpenGL3)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	defer platform.Dispose()

	renderer, err := renderers.NewOpenGL3(io)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	defer renderer.Dispose()

	example.Run(platform, renderer)
}
