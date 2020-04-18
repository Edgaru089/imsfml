package imsfml

import "C"
import (
	"runtime"
	"time"
	"unicode/utf8"

	"github.com/go-gl/gl/v2.1/gl"

	sf "github.com/Edgaru089/gosfml2"
	"github.com/inkyblackness/imgui-go"
)

var (
	hasFocus     bool
	mouseMoved   bool
	mousePressed [3]bool
	mouseHidden  bool

	fontTexture *sf.Texture
)

func encodeRuneUTF8(c rune) string {
	arr := make([]byte, 4)
	return string(arr[:utf8.EncodeRune(arr, c)])
}

// sfClipboard is the wrapper for sf.Clipboard*, satisfying imgui.Clipboard
type sfClipboard struct{}

func (*sfClipboard) Text() (string, error) {
	return sf.ClipboardGetString(), nil
}

func (*sfClipboard) SetText(t string) {
	sf.ClipboardSetString(t)
}

// InitRenderWindow calls Init(win.GetSize(), win.HasFocus(), fontAtlas).
// createDefaultFont controls whether or not to load the default font.
func InitRenderWindow(win *sf.RenderWindow, createDefaultFont bool) error {
	size := win.GetSize()
	return Init(imgui.Vec2{X: float32(size.X), Y: float32(size.Y)}, win.HasFocus(), createDefaultFont)
}

// Init resets internal state, calling imgui.CreateContext.
// createDefaultFont controls whether or not to load the default font.
func Init(displaySize imgui.Vec2, winHasFocus bool, createDefaultFont bool) (err error) {
	runtime.LockOSThread()

	imgui.CreateContext(nil)

	mIO := imgui.CurrentIO()

	// It seems there are no flags to be set?

	// Init keymaps
	mIO.KeyMap(imgui.KeyTab, sf.KeyTab)
	mIO.KeyMap(imgui.KeyLeftArrow, sf.KeyLeft)
	mIO.KeyMap(imgui.KeyRightArrow, sf.KeyRight)
	mIO.KeyMap(imgui.KeyUpArrow, sf.KeyUp)
	mIO.KeyMap(imgui.KeyDownArrow, sf.KeyDown)
	mIO.KeyMap(imgui.KeyPageUp, sf.KeyPageUp)
	mIO.KeyMap(imgui.KeyPageDown, sf.KeyPageDown)
	mIO.KeyMap(imgui.KeyHome, sf.KeyHome)
	mIO.KeyMap(imgui.KeyEnd, sf.KeyEnd)
	mIO.KeyMap(imgui.KeyInsert, sf.KeyInsert)
	mIO.KeyMap(imgui.KeyDelete, sf.KeyDelete)
	mIO.KeyMap(imgui.KeyBackspace, sf.KeyBack)
	mIO.KeyMap(imgui.KeySpace, sf.KeySpace)
	mIO.KeyMap(imgui.KeyEnter, sf.KeyReturn)
	mIO.KeyMap(imgui.KeyEscape, sf.KeyEscape)
	mIO.KeyMap(imgui.KeyA, sf.KeyA)
	mIO.KeyMap(imgui.KeyC, sf.KeyC)
	mIO.KeyMap(imgui.KeyV, sf.KeyV)
	mIO.KeyMap(imgui.KeyX, sf.KeyX)
	mIO.KeyMap(imgui.KeyY, sf.KeyY)
	mIO.KeyMap(imgui.KeyZ, sf.KeyZ)

	// Set display size
	mIO.SetDisplaySize(displaySize)

	// Add clipboard handler
	mIO.SetClipboard(&sfClipboard{})

	hasFocus = winHasFocus

	if createDefaultFont {
		// This is done automatically
		//mIO.Fonts().AddFontDefault()
		UpdateFontTexture()
	}

	return gl.Init()
}

// ProcessEvent is to be called on every SFML Event received
func ProcessEvent(e sf.Event) {
	if hasFocus {
		mIO := imgui.CurrentIO()

		switch e.Type() {
		case sf.EventTypeMouseMoved:
			mouseMoved = true
		case sf.EventTypeMouseButtonReleased, sf.EventTypeMouseButtonPressed:
			var b sf.MouseButton
			switch e.(type) {
			case sf.EventMouseButtonPressed:
				b = e.(sf.EventMouseButtonPressed).Button
			case sf.EventMouseButtonReleased:
				b = e.(sf.EventMouseButtonReleased).Button
			}
			if b >= 0 && b <= 3 {
				mousePressed[b] = (e.Type() == sf.EventTypeMouseButtonPressed)
			}

		case sf.EventTypeMouseWheelMoved:
			d := e.(sf.EventMouseWheelMoved).Delta
			mIO.AddMouseWheelDelta(0, float32(d))

		case sf.EventTypeKeyPressed:
			mIO.KeyPress(int(e.(sf.EventKeyPressed).Code))
		case sf.EventTypeKeyReleased:
			mIO.KeyRelease(int(e.(sf.EventKeyReleased).Code))

		case sf.EventTypeTextEntered:
			code := e.(sf.EventTextEntered).Char
			if !(code < ' ' || code == 127) {
				mIO.AddInputCharacters(encodeRuneUTF8(code))
			}
		}

	}

	switch e.Type() {
	case sf.EventTypeLostFocus:
		hasFocus = false
	case sf.EventTypeGainedFocus:
		hasFocus = true
	}
}

// UpdateRenderWindow calls Update(sf.Mouse.GetPosition(win), win.GetSize(), deltaT)
// Call it before every frame.
func UpdateRenderWindow(win *sf.RenderWindow, deltaT time.Duration) {
	pos := sf.MouseGetPosition(win)
	size := win.GetSize()
	Update(imgui.Vec2{X: float32(pos.X), Y: float32(pos.Y)}, imgui.Vec2{X: float32(size.X), Y: float32(size.Y)}, deltaT)
}

// converts bool to int
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Update updates the initial state, calling imgui.NewFrame.
// Call it before every frame.
func Update(mousePos, displaySize imgui.Vec2, deltaT time.Duration) {

	io := imgui.CurrentIO()

	io.SetDisplaySize(displaySize)

	io.SetMousePosition(mousePos)
	for i := 0; i < 3; i++ {
		io.SetMouseButtonDown(i, mousePressed[i])
	}

	// Update the state of Ctrl, Alt, Shift and Super
	io.KeyCtrl(sf.KeyLControl, sf.KeyRControl)
	io.KeyAlt(sf.KeyLAlt, sf.KeyRAlt)
	io.KeyShift(sf.KeyLShift, sf.KeyRShift)
	io.KeySuper(sf.KeyLSystem, sf.KeyRSystem)

	io.SetDeltaTime(float32(deltaT.Seconds()))

	imgui.NewFrame()

}

// Render calls imgui.Render, and then draws the frame.
func Render(win *sf.RenderWindow) {
	imgui.Render()

	win.SetActive(true)
	size := win.GetSize()
	renderDrawlist(size, sf.Vector2f{X: float32(size.X), Y: float32(size.Y)}, imgui.RenderedDrawData())
}

// UpdateFontTexture uploads the font texture.
func UpdateFontTexture() {
	io := imgui.CurrentIO()
	im := io.Fonts().TextureDataRGBA32()

	var err error
	fontTexture, err = sf.NewTexture(uint(im.Width), uint(im.Height))
	if err != nil {
		panic("imsfml.UpdateFontTexture(): sf.NewTexture Error: " + err.Error())
	}

	fontTexture.UpdateFromPixelsUnsafe(im.Pixels, uint(im.Width), uint(im.Height), 0, 0)

	io.Fonts().SetTextureID(imgui.TextureID(fontTexture.GetNativeHandle()))
}

// FontTexture returns the font texture
func FontTexture() *sf.Texture {
	return fontTexture
}

// ImageTextureV calls imgui.ImageV, with the parameters filled in as expected
func ImageTextureV(texture *sf.Texture, size sf.Vector2f, textureRect sf.IntRect, tintColor, borderColor sf.Color) {
	tSize := texture.GetSize()

	uv0 := imgui.Vec2{
		X: float32(textureRect.Left) / float32(tSize.X),
		Y: float32(textureRect.Top) / float32(tSize.Y),
	}
	uv1 := imgui.Vec2{
		X: float32(textureRect.Left+textureRect.Width) / float32(tSize.X),
		Y: float32(textureRect.Top+textureRect.Height) / float32(tSize.Y),
	}

	imgui.ImageV(
		imgui.TextureID(texture.GetNativeHandle()),
		imgui.Vec2{X: size.X, Y: size.Y},
		uv0, uv1,
		ColorToVec4(tintColor),
		ColorToVec4(borderColor),
	)
}

// Image calls ImageTextureV(sprite.Texture, sprite.GetSize(), sprite.GetTextureRect(), tintColor, borderColor)
func Image(sprite *sf.Sprite, tintColor, borderColor sf.Color) {
	t := sprite.GetTexture()
	s := sprite.GetScale()
	ts := t.GetSize()
	ImageTextureV(t, sf.Vector2f{X: float32(ts.X) * s.X, Y: float32(ts.Y) * s.Y}, sprite.GetTextureRect(), tintColor, borderColor)
}

// ImageButtonTextureV calls imgui.ImageButtonV, with the parameters filled in as expected
func ImageButtonTextureV(texture *sf.Texture, size sf.Vector2f, textureRect sf.IntRect, framePadding int, bgColor, tintColor sf.Color) bool {
	tSize := texture.GetSize()

	uv0 := imgui.Vec2{
		X: float32(textureRect.Left) / float32(tSize.X),
		Y: float32(textureRect.Top) / float32(tSize.Y),
	}
	uv1 := imgui.Vec2{
		X: float32(textureRect.Left+textureRect.Width) / float32(tSize.X),
		Y: float32(textureRect.Top+textureRect.Height) / float32(tSize.Y),
	}

	return imgui.ImageButtonV(
		imgui.TextureID(texture.GetNativeHandle()),
		imgui.Vec2{X: size.X, Y: size.Y},
		uv0, uv1,
		framePadding,
		ColorToVec4(bgColor),
		ColorToVec4(tintColor),
	)
}

// ImageButton calls ImageButtonTextureV(sprite.Texture, sprite.GetSize(), sprite.GetTextureRect(), framePadding, bgColor, tintColor)
func ImageButton(sprite *sf.Sprite, framePadding int, bgColor, tintColor sf.Color) {
	t := sprite.GetTexture()
	s := sprite.GetScale()
	ts := t.GetSize()
	ImageButtonTextureV(t, sf.Vector2f{X: float32(ts.X) * s.X, Y: float32(ts.Y) * s.Y}, sprite.GetTextureRect(), framePadding, bgColor, tintColor)
}
