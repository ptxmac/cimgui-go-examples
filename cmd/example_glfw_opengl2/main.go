//go:build glfw
// +build glfw

package main

import (
	"fmt"
	"os"

	imgui "github.com/AllenDang/cimgui-go"

	"github.com/ptxmac/cimgui-go-examples/internal/example"
	"github.com/ptxmac/cimgui-go-examples/internal/platforms"
	"github.com/ptxmac/cimgui-go-examples/internal/renderers"
)

func main() {
	context := imgui.CreateContext(0)
	defer context.Destroy()
	io := imgui.GetIO()

	platform, err := platforms.NewGLFW(io, platforms.GLFWClientAPIOpenGL2)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	defer platform.Dispose()

	renderer, err := renderers.NewOpenGL2(io)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	defer renderer.Dispose()

	example.Run(platform, renderer)
}
