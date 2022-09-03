//go:build glfw
// +build glfw

package platforms

import (
	"fmt"
	"math"
	"runtime"

	"github.com/AllenDang/cimgui-go"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// GLFWClientAPI identifies the render system that shall be initialized.
type GLFWClientAPI string

// This is a list of GLFWClientAPI constants.
const (
	GLFWClientAPIOpenGL2 GLFWClientAPI = "OpenGL2"
	GLFWClientAPIOpenGL3 GLFWClientAPI = "OpenGL3"
)

// GLFW implements a platform based on github.com/go-gl/glfw (v3.2).
type GLFW struct {
	imguiIO cimgui.ImGuiIO

	window *glfw.Window

	keyMap map[glfw.Key]cimgui.ImGuiKey

	time             float64
	mouseJustPressed [3]bool
}

// NewGLFW attempts to initialize a GLFW context.
func NewGLFW(io cimgui.ImGuiIO, clientAPI GLFWClientAPI) (*GLFW, error) {
	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize glfw: %w", err)
	}

	switch clientAPI {
	case GLFWClientAPIOpenGL2:
		glfw.WindowHint(glfw.ContextVersionMajor, 2)
		glfw.WindowHint(glfw.ContextVersionMinor, 1)
	case GLFWClientAPIOpenGL3:
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 2)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, 1)
	default:
		glfw.Terminate()
		return nil, ErrUnsupportedClientAPI
	}

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "CImGui-Go GLFW+"+string(clientAPI)+" example", nil, nil)
	if err != nil {
		glfw.Terminate()
		return nil, fmt.Errorf("failed to create window: %w", err)
	}
	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	platform := &GLFW{
		imguiIO: io,
		window:  window,
	}
	platform.setKeyMapping()
	platform.installCallbacks()

	return platform, nil
}

// Dispose cleans up the resources.
func (platform *GLFW) Dispose() {
	platform.window.Destroy()
	glfw.Terminate()
}

// ShouldStop returns true if the window is to be closed.
func (platform *GLFW) ShouldStop() bool {
	return platform.window.ShouldClose()
}

// ProcessEvents handles all pending window events.
func (platform *GLFW) ProcessEvents() {
	glfw.PollEvents()
}

// DisplaySize returns the dimension of the display.
func (platform *GLFW) DisplaySize() [2]float32 {
	w, h := platform.window.GetSize()
	return [2]float32{float32(w), float32(h)}
}

// FramebufferSize returns the dimension of the framebuffer.
func (platform *GLFW) FramebufferSize() [2]float32 {
	w, h := platform.window.GetFramebufferSize()
	return [2]float32{float32(w), float32(h)}
}

// NewFrame marks the begin of a render pass. It forwards all current state to imgui IO.
func (platform *GLFW) NewFrame() {
	// Setup display size (every frame to accommodate for window resizing)
	displaySize := platform.DisplaySize()
	platform.imguiIO.SetDisplaySize(cimgui.ImVec2{X: displaySize[0], Y: displaySize[1]})

	// Setup time step
	currentTime := glfw.GetTime()
	if platform.time > 0 {
		platform.imguiIO.SetDeltaTime(float32(currentTime - platform.time))
	}
	platform.time = currentTime

	// Setup inputs
	if platform.window.GetAttrib(glfw.Focused) != 0 {
		x, y := platform.window.GetCursorPos()
		platform.imguiIO.SetMousePos(cimgui.ImVec2{X: float32(x), Y: float32(y)})
	} else {
		platform.imguiIO.SetMousePos(cimgui.ImVec2{X: -math.MaxFloat32, Y: -math.MaxFloat32})
	}

	for i := 0; i < len(platform.mouseJustPressed); i++ {
		down := platform.mouseJustPressed[i] || (platform.window.GetMouseButton(glfwButtonIDByIndex[i]) == glfw.Press)
		platform.imguiIO.SetMouseButtonDown(i, down)
		platform.mouseJustPressed[i] = false
	}
}

// PostRender performs a buffer swap.
func (platform *GLFW) PostRender() {
	platform.window.SwapBuffers()
}

func (platform *GLFW) setKeyMapping() {
	// Keyboard mapping. ImGui will use those indices to peek into the io.KeysDown[] array.

	platform.keyMap = map[glfw.Key]cimgui.ImGuiKey{
		glfw.KeyTab:       cimgui.ImGuiKey_Tab,
		glfw.KeyLeft:      cimgui.ImGuiKey_LeftArrow,
		glfw.KeyRight:     cimgui.ImGuiKey_RightArrow,
		glfw.KeyUp:        cimgui.ImGuiKey_UpArrow,
		glfw.KeyDown:      cimgui.ImGuiKey_DownArrow,
		glfw.KeyPageUp:    cimgui.ImGuiKey_PageUp,
		glfw.KeyPageDown:  cimgui.ImGuiKey_PageDown,
		glfw.KeyHome:      cimgui.ImGuiKey_Home,
		glfw.KeyEnd:       cimgui.ImGuiKey_End,
		glfw.KeyInsert:    cimgui.ImGuiKey_Insert,
		glfw.KeyDelete:    cimgui.ImGuiKey_Delete,
		glfw.KeyBackspace: cimgui.ImGuiKey_Backspace,
		glfw.KeySpace:     cimgui.ImGuiKey_Space,
		glfw.KeyEnter:     cimgui.ImGuiKey_Enter,
		glfw.KeyEscape:    cimgui.ImGuiKey_Escape,
		glfw.KeyA:         cimgui.ImGuiKey_A,
		glfw.KeyC:         cimgui.ImGuiKey_C,
		glfw.KeyV:         cimgui.ImGuiKey_V,
		glfw.KeyX:         cimgui.ImGuiKey_X,
		glfw.KeyY:         cimgui.ImGuiKey_Y,
		glfw.KeyZ:         cimgui.ImGuiKey_Z,

		glfw.KeyLeftControl:  cimgui.ImGuiKey_ModCtrl,
		glfw.KeyRightControl: cimgui.ImGuiKey_ModCtrl,
		glfw.KeyLeftAlt:      cimgui.ImGuiKey_ModAlt,
		glfw.KeyRightAlt:     cimgui.ImGuiKey_ModAlt,
		glfw.KeyLeftSuper:    cimgui.ImGuiKey_ModSuper,
		glfw.KeyRightSuper:   cimgui.ImGuiKey_ModSuper,
	}

}

func (platform *GLFW) installCallbacks() {
	platform.window.SetMouseButtonCallback(platform.mouseButtonChange)
	platform.window.SetScrollCallback(platform.mouseScrollChange)
	platform.window.SetKeyCallback(platform.keyChange)
	platform.window.SetCharCallback(platform.charChange)
}

var glfwButtonIndexByID = map[glfw.MouseButton]int{
	glfw.MouseButton1: mouseButtonPrimary,
	glfw.MouseButton2: mouseButtonSecondary,
	glfw.MouseButton3: mouseButtonTertiary,
}

var glfwButtonIDByIndex = map[int]glfw.MouseButton{
	mouseButtonPrimary:   glfw.MouseButton1,
	mouseButtonSecondary: glfw.MouseButton2,
	mouseButtonTertiary:  glfw.MouseButton3,
}

func (platform *GLFW) mouseButtonChange(window *glfw.Window, rawButton glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	buttonIndex, known := glfwButtonIndexByID[rawButton]

	if known && (action == glfw.Press) {
		platform.mouseJustPressed[buttonIndex] = true
	}
}

func (platform *GLFW) mouseScrollChange(window *glfw.Window, x, y float64) {
	platform.imguiIO.AddMouseWheelDelta(float32(x), float32(y))
}

func (platform *GLFW) keyChange(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	imKey := cimgui.ImGuiKey(key)
	if mapped, ok := platform.keyMap[key]; ok {
		imKey = mapped
	}

	platform.imguiIO.AddKeyEvent(imKey, action == glfw.Press)
}

func (platform *GLFW) charChange(window *glfw.Window, char rune) {
	platform.imguiIO.AddInputCharactersUTF8(string(char))
	//platform.imguiIO.AddInputCharacters(string(char))
}

// ClipboardText returns the current clipboard text, if available.
func (platform *GLFW) ClipboardText() (string, error) {
	return platform.window.GetClipboardString()
}

// SetClipboardText sets the text as the current clipboard text.
func (platform *GLFW) SetClipboardText(text string) {
	platform.window.SetClipboardString(text)
}
