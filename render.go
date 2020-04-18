package imsfml

import (
	"unsafe"

	sf "github.com/Edgaru089/gosfml2"
	"github.com/go-gl/gl/v2.1/gl" // Adjust the GL version, if you need to
	"github.com/inkyblackness/imgui-go"
)

func renderDrawlist(displaySize sf.Vector2u, framebufferSize sf.Vector2f, drawData imgui.DrawData) {
	displayWidth, displayHeight := float32(displaySize.X), float32(displaySize.Y)
	fbWidth, fbHeight := framebufferSize.X, framebufferSize.Y
	if (fbWidth <= 0) || (fbHeight <= 0) {
		return
	}
	drawData.ScaleClipRects(imgui.Vec2{
		X: fbWidth / displayWidth,
		Y: fbHeight / displayHeight,
	})

	gl.PushAttrib(gl.ENABLE_BIT | gl.COLOR_BUFFER_BIT | gl.TRANSFORM_BIT)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.LIGHTING)
	gl.Enable(gl.SCISSOR_TEST)
	gl.Enable(gl.TEXTURE_2D)
	gl.Disable(gl.LIGHTING)
	gl.EnableClientState(gl.VERTEX_ARRAY)
	gl.EnableClientState(gl.COLOR_ARRAY)
	gl.EnableClientState(gl.TEXTURE_COORD_ARRAY)

	gl.Viewport(0, 0, int32(fbWidth), int32(fbHeight))

	gl.MatrixMode(gl.TEXTURE)
	gl.LoadIdentity()

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()

	gl.Ortho(0, float64(displayWidth), float64(displayHeight), 0, -1, 1)

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	vertexSize, vertexOffsetPos, vertexOffsetUv, vertexOffsetCol := imgui.VertexBufferLayout()
	indexSize := imgui.IndexBufferLayout()

	drawType := gl.UNSIGNED_SHORT
	if indexSize == 4 {
		drawType = gl.UNSIGNED_INT
	}

	// Render command lists
	for _, commandList := range drawData.CommandLists() {
		vertexBuffer, _ := commandList.VertexBuffer()
		indexBuffer, _ := commandList.IndexBuffer()
		indexBufferOffset := uintptr(indexBuffer)

		gl.VertexPointer(2, gl.FLOAT, int32(vertexSize), unsafe.Pointer(uintptr(vertexBuffer)+uintptr(vertexOffsetPos)))
		gl.TexCoordPointer(2, gl.FLOAT, int32(vertexSize), unsafe.Pointer(uintptr(vertexBuffer)+uintptr(vertexOffsetUv)))
		gl.ColorPointer(4, gl.UNSIGNED_BYTE, int32(vertexSize), unsafe.Pointer(uintptr(vertexBuffer)+uintptr(vertexOffsetCol)))

		for _, command := range commandList.Commands() {
			if command.HasUserCallback() {
				command.CallUserCallback(commandList)
			} else {
				clipRect := command.ClipRect()
				gl.BindTexture(gl.TEXTURE_2D, uint32(command.TextureID()))
				gl.Scissor(int32(clipRect.X), int32(fbHeight)-int32(clipRect.W), int32(clipRect.Z-clipRect.X), int32(clipRect.W-clipRect.Y))
				gl.DrawElements(gl.TRIANGLES, int32(command.ElementCount()), uint32(drawType), unsafe.Pointer(indexBufferOffset))
			}

			indexBufferOffset += uintptr(command.ElementCount() * indexSize)
		}
	}

	gl.PopAttrib()

}
