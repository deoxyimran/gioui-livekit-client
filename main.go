package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"github.com/deoxyimran/gioui-livekit-client/ui"
)

func main() {
	go func() {
		window := new(app.Window)
		window.Option(
			app.Title("Family Meet"),
		)
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window) error {
	var ops op.Ops
	joinScreen := ui.NewJoinMeetingScreen()
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			joinScreen.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}
