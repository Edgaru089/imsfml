package main

import (
	"runtime"
	"time"

	sf "github.com/Edgaru089/gosfml2"
	"github.com/Edgaru089/imsfml"
	"github.com/inkyblackness/imgui-go"
)

func main() {
	runtime.LockOSThread()

	win := sf.NewRenderWindow(
		sf.VideoMode{Width: 1600, Height: 900, BitsPerPixel: 32},
		"ImGUI+GoSFML Window",
		sf.StyleDefault,
		sf.ContextSettings{MajorVersion: 2, MinorVersion: 1},
	)

	err := imsfml.InitRenderWindow(win, true)
	if err != nil {
		panic(err)
	}

	win.SetVSyncEnabled(true)
	win.ResetGLStates()

	lastUpdate := time.Now()

	for win.IsOpen() {
		for {
			e := win.PollEvent()
			if e == nil {
				break
			}

			switch e.Type() {
			case sf.EventTypeClosed:
				win.Close()
			}

			imsfml.ProcessEvent(e)
		}

		now := time.Now()
		imsfml.UpdateRenderWindow(win, now.Sub(lastUpdate))
		lastUpdate = now

		// We can now call ImGui widgets
		imgui.ShowDemoWindow(nil)

		win.Clear(sf.ColorBlack())
		imsfml.Render(win)

		win.Display()
	}

	// TODO Cleanup
}
