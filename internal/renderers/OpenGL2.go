package renderers

import (
	"fmt"
	"unsafe"

	"github.com/AllenDang/cimgui-go"
	"github.com/ptxmac/cimgui-go-examples/internal/renderers/gl/v2.1/gl"
)

// OpenGL2 implements a renderer based on github.com/go-gl/gl (v2.1).
type OpenGL2 struct {
	imguiIO imgui.IO

	fontTexture uint32
}

// NewOpenGL2 attempts to initialize a renderer.
// An OpenGL context has to be established before calling this function.
func NewOpenGL2(io imgui.IO) (*OpenGL2, error) {
	err := gl.Init()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize OpenGL: %w", err)
	}

	renderer := &OpenGL2{
		imguiIO: io,
	}
	renderer.createFontsTexture()
	return renderer, nil
}

// Dispose cleans up the resources.
func (renderer *OpenGL2) Dispose() {
	renderer.destroyFontsTexture()
}

// PreRender clears the framebuffer.
func (renderer *OpenGL2) PreRender(clearColor [3]float32) {
	gl.ClearColor(clearColor[0], clearColor[1], clearColor[2], 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

// Render translates the ImGui draw data to OpenGL3 commands.
func (renderer *OpenGL2) Render(displaySize [2]float32, framebufferSize [2]float32, drawData imgui.DrawData) {
	// Avoid rendering when minimized, scale coordinates for retina displays (screen coordinates != framebuffer coordinates)
	displayWidth, displayHeight := displaySize[0], displaySize[1]
	fbWidth, fbHeight := framebufferSize[0], framebufferSize[1]
	if (fbWidth <= 0) || (fbHeight <= 0) {
		return
	}
	drawData.ScaleClipRects(imgui.Vec2{
		X: fbWidth / displayWidth,
		Y: fbHeight / displayHeight,
	})

	// Setup render state: alpha-blending enabled, no face culling, no depth testing, scissor enabled, vertex/texcoord/color pointers, polygon fill.
	var lastTexture int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	var lastPolygonMode [2]int32
	gl.GetIntegerv(gl.POLYGON_MODE, &lastPolygonMode[0])
	var lastViewport [4]int32
	gl.GetIntegerv(gl.VIEWPORT, &lastViewport[0])
	var lastScissorBox [4]int32
	gl.GetIntegerv(gl.SCISSOR_BOX, &lastScissorBox[0])
	gl.PushAttrib(gl.ENABLE_BIT | gl.COLOR_BUFFER_BIT | gl.TRANSFORM_BIT)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.LIGHTING)
	gl.Disable(gl.COLOR_MATERIAL)
	gl.Enable(gl.SCISSOR_TEST)
	gl.EnableClientState(gl.VERTEX_ARRAY)
	gl.EnableClientState(gl.TEXTURE_COORD_ARRAY)
	gl.EnableClientState(gl.COLOR_ARRAY)
	gl.Enable(gl.TEXTURE_2D)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	// You may want this if using this code in an OpenGL 3+ context where shaders may be bound
	// gl.UseProgram(0)

	// Setup viewport, orthographic projection matrix
	// Our visible cimgui space lies from draw_data->DisplayPos (top left) to draw_data->DisplayPos+data_data->DisplaySize (bottom right).
	// DisplayMin is typically (0,0) for single viewport apps.
	gl.Viewport(0, 0, int32(fbWidth), int32(fbHeight))
	gl.MatrixMode(gl.PROJECTION)
	gl.PushMatrix()
	gl.LoadIdentity()
	gl.Ortho(0, float64(displayWidth), float64(displayHeight), 0, -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	gl.LoadIdentity()

	vertexSize, vertexOffsetPos, vertexOffsetUv, vertexOffsetCol := imgui.VertexBufferLayout()
	indexSize := imgui.IndexBufferLayout()

	drawType := gl.UNSIGNED_SHORT
	const bytesPerUint32 = 4
	if indexSize == bytesPerUint32 {
		drawType = gl.UNSIGNED_INT
	}

	// Render command lists
	for _, commandList := range drawData.CommandLists() {
		vertexBuffer, _ := commandList.GetVertexBuffer()
		indexBuffer, _ := commandList.GetIndexBuffer()
		indexBufferOffset := uintptr(indexBuffer)

		gl.VertexPointer(2, gl.FLOAT, int32(vertexSize), unsafe.Pointer(uintptr(vertexBuffer)+uintptr(vertexOffsetPos)))
		gl.TexCoordPointer(2, gl.FLOAT, int32(vertexSize), unsafe.Pointer(uintptr(vertexBuffer)+uintptr(vertexOffsetUv)))
		gl.ColorPointer(4, gl.UNSIGNED_BYTE, int32(vertexSize), unsafe.Pointer(uintptr(vertexBuffer)+uintptr(vertexOffsetCol)))

		for _, command := range commandList.Commands() {
			if command.HasUserCallback() {
				command.CallUserCallback(commandList)
			} else {
				clipRect := command.ClipRect()
				gl.Scissor(int32(clipRect.X), int32(fbHeight)-int32(clipRect.W), int32(clipRect.Z-clipRect.X), int32(clipRect.W-clipRect.Y))
				gl.BindTexture(gl.TEXTURE_2D, uint32(uintptr(command.TextureId())))
				gl.DrawElementsWithOffset(gl.TRIANGLES, int32(command.ElemCount()), uint32(drawType), indexBufferOffset)
			}

			indexBufferOffset += uintptr(command.ElemCount() * uint32(indexSize))
		}
	}

	// Restore modified state
	gl.DisableClientState(gl.COLOR_ARRAY)
	gl.DisableClientState(gl.TEXTURE_COORD_ARRAY)
	gl.DisableClientState(gl.VERTEX_ARRAY)
	gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))
	gl.MatrixMode(gl.MODELVIEW)
	gl.PopMatrix()
	gl.MatrixMode(gl.PROJECTION)
	gl.PopMatrix()
	gl.PopAttrib()
	gl.PolygonMode(gl.FRONT, uint32(lastPolygonMode[0]))
	gl.PolygonMode(gl.BACK, uint32(lastPolygonMode[1]))
	gl.Viewport(lastViewport[0], lastViewport[1], lastViewport[2], lastViewport[3])
	gl.Scissor(lastScissorBox[0], lastScissorBox[1], lastScissorBox[2], lastScissorBox[3])
}

func (renderer *OpenGL2) createFontsTexture() {
	// Build texture atlas
	pixels, width, height, _ := renderer.imguiIO.Fonts().GetTextureDataAsRGBA32()

	// Upload texture to graphics system
	var lastTexture int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	gl.GenTextures(1, &renderer.fontTexture)
	gl.BindTexture(gl.TEXTURE_2D, renderer.fontTexture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(width), int32(height),
		0, gl.RGBA, gl.UNSIGNED_BYTE, pixels)

	// Store our identifier

	renderer.imguiIO.Fonts().SetTexID(imgui.TextureID(uintptr(renderer.fontTexture)))

	// Restore state
	gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))
}

func (renderer *OpenGL2) destroyFontsTexture() {
	if renderer.fontTexture != 0 {
		gl.DeleteTextures(1, &renderer.fontTexture)
		imgui.CurrentIO().Fonts().SetTexID(nil)
		renderer.fontTexture = 0
	}
}
