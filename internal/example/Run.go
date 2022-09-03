package example

import (
	"fmt"
	"time"

	"github.com/AllenDang/cimgui-go"
	"github.com/ptxmac/cimgui-go-examples/internal/demo"
)

// Platform covers mouse/keyboard/gamepad inputs, cursor shape, timing, windowing.
type Platform interface {
	// ShouldStop is regularly called as the abort condition for the program loop.
	ShouldStop() bool
	// ProcessEvents is called once per render loop to dispatch any pending events.
	ProcessEvents()
	// DisplaySize returns the dimension of the display.
	DisplaySize() [2]float32
	// FramebufferSize returns the dimension of the framebuffer.
	FramebufferSize() [2]float32
	// NewFrame marks the begin of a render pass. It must update the cimgui IO state according to user input (mouse, keyboard, ...)
	NewFrame()
	// PostRender marks the completion of one render pass. Typically this causes the display buffer to be swapped.
	PostRender()
	// ClipboardText returns the current text of the clipboard, if available.
	ClipboardText() (string, error)
	// SetClipboardText sets the text as the current text of the clipboard.
	SetClipboardText(text string)
}

type clipboard struct {
	platform Platform
}

func (board clipboard) Text() (string, error) {
	return board.platform.ClipboardText()
}

func (board clipboard) SetText(text string) {
	board.platform.SetClipboardText(text)
}

// Renderer covers rendering cimgui draw data.
type Renderer interface {
	// PreRender causes the display buffer to be prepared for new output.
	PreRender(clearColor [3]*float32)
	// Render draws the provided cimgui draw data.
	Render(displaySize [2]float32, framebufferSize [2]float32, drawData cimgui.ImDrawData)
}

const (
	millisPerSecond = 1000
	sleepDuration   = time.Millisecond * 25
)

// Run implements the main program loop of the demo. It returns when the platform signals to stop.
// This demo application shows some basic features of ImGui, as well as exposing the standard demo window.
func Run(p Platform, r Renderer) {

	// TODO:
	//cimgui.CurrentIO().SetClipboard(clipboard{platform: p})

	showDemoWindow := false
	showGoDemoWindow := false
	var c1, c2, c3 float32 = 0.0, 0.0, 0.0
	clearColor := [3]*float32{&c1, &c2, &c3}
	f := float32(0)
	counter := 0
	showAnotherWindow := false

	for !p.ShouldStop() {
		p.ProcessEvents()

		// Signal start of a new frame
		p.NewFrame()
		cimgui.NewFrame()

		// 1. Show a simple window.
		// Tip: if we don't call cimgui.Begin()/cimgui.End() the widgets automatically appears in a window called "Debug".
		{
			cimgui.Text("ภาษาไทย测试조선말")                                                     // To display these, you'll need to register a compatible font
			cimgui.Text("Hello, world!")                                                    // Display some text
			cimgui.SliderFloat("float", &f, 0.0, 1.0, "%.3f", cimgui.ImGuiSliderFlags_None) // Edit 1 float using a slider from 0.0f to 1.0f

			cimgui.ColorEdit3("clear color", clearColor, cimgui.ImGuiColorEditFlags_None) // Edit 3 floats representing a color

			cimgui.Checkbox("Demo Window", &showDemoWindow) // Edit bools storing our window open/close state
			cimgui.Checkbox("Go Demo Window", &showGoDemoWindow)
			cimgui.Checkbox("Another Window", &showAnotherWindow)

			if cimgui.Button("Button", cimgui.ImVec2{}) { // Buttons return true when clicked (most widgets return true when edited/activated)
				counter++
			}
			cimgui.SameLine(0, -1)
			cimgui.Text(fmt.Sprintf("counter = %d", counter))

			cimgui.Text(fmt.Sprintf("Application average %.3f ms/frame (%.1f FPS)",
				millisPerSecond/cimgui.GetIO().GetFramerate(), cimgui.GetIO().GetFramerate()))
		}

		// 2. Show another simple window. In most cases you will use an explicit Begin/End pair to name your windows.
		if showAnotherWindow {
			// Pass a pointer to our bool variable (the window will have a closing button that will clear the bool when clicked)
			cimgui.Begin("Another window", &showAnotherWindow, 0)
			cimgui.Text("Hello from another window!")
			if cimgui.Button("Close Me", cimgui.ImVec2{}) {
				showAnotherWindow = false
			}
			cimgui.End()
		}

		// 3. Show the ImGui demo window. Most of the sample code is in cimgui.ShowDemoWindow().
		// Read its code to learn more about Dear ImGui!
		if showDemoWindow {
			// Normally user code doesn't need/want to call this because positions are saved in .ini file anyway.
			// Here we just want to make the demo initial state a bit more friendly!
			const demoX = 650
			const demoY = 20
			cimgui.SetNextWindowPos(cimgui.ImVec2{X: demoX, Y: demoY}, cimgui.ImGuiCond_FirstUseEver, cimgui.ImVec2{})

			cimgui.ShowDemoWindow(&showDemoWindow)
		}
		if showGoDemoWindow {
			demo.Show(&showGoDemoWindow)
		}

		// Rendering
		cimgui.Render() // This call only creates the draw data list. Actual rendering to framebuffer is done below.

		r.PreRender(clearColor)
		// A this point, the application could perform its own rendering...
		// app.RenderScene()

		r.Render(p.DisplaySize(), p.FramebufferSize(), cimgui.GetDrawData())
		p.PostRender()

		// sleep to avoid 100% CPU usage for this demo
		<-time.After(sleepDuration)
	}
}
