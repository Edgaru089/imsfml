package imsfml

import (
	sf "github.com/Edgaru089/gosfml2"
	"github.com/inkyblackness/imgui-go"
)

// ColorToVec4 converts sf.Color to imgui.Vec4
func ColorToVec4(col sf.Color) imgui.Vec4 {
	return imgui.Vec4{
		X: float32(col.R) / 255,
		Y: float32(col.G) / 255,
		Z: float32(col.B) / 255,
		W: float32(col.A) / 255,
	}
}
