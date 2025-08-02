package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"github.com/deoxyimran/gioui-livekit-client/ui"
)

var screenPointer ui.Screen = ui.JOIN_MEETING_SCREEN
var prevScreenPointer ui.Screen = ui.JOIN_MEETING_SCREEN

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
	stateManger := ui.NewStateManager()
	joinMeetingScreen := ui.NewJoinMeetingScreen(&stateManger)
	var joinRoomScreen *ui.JoinRoomScreen
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			if joinRoomScreen != nil {
				joinRoomScreen.StopVideoCapture()
			}
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			// Handle screen switching
			switch screenPointer {
			case ui.JOIN_MEETING_SCREEN:
				temp := screenPointer
				if prevScreenPointer != screenPointer {
					joinMeetingScreen = ui.NewJoinMeetingScreen(&stateManger)
					joinMeetingScreen.Layout(gtx, &screenPointer)
				} else {
					joinMeetingScreen.Layout(gtx, &screenPointer)
				}
				if temp == screenPointer {
					prevScreenPointer = temp
				}
			case ui.JOIN_ROOM_SCREEN:
				temp := screenPointer
				if prevScreenPointer != screenPointer {
					joinRoomScreen = ui.NewJoinRoomScreen(&stateManger)
					joinRoomScreen.Layout(gtx, &screenPointer)
				} else {
					joinRoomScreen.Layout(gtx, &screenPointer)
				}
				if temp == screenPointer {
					prevScreenPointer = temp
				}
			}
			e.Frame(gtx.Ops)
		}
	}
}
