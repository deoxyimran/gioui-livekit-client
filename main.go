package main

import (
	"fmt"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"github.com/deoxyimran/gioui-livekit-client/ui/screens"
	"github.com/deoxyimran/gioui-livekit-client/ui/state"
)

var screenPointer screens.Screen = screens.JOIN_MEETING
var prevScreenPointer screens.Screen = screens.JOIN_MEETING

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
	appState := state.App{}

	joinMeetingScreen := screens.NewJoinMeetingScreen(&appState) // Default screen
	var joinRoomScreen *screens.JoinRoom
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			fmt.Println("Exiting application...")
			if joinRoomScreen != nil {
				joinRoomScreen.StopVideoCapture()
			}
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			// Handle screen switching
			switch screenPointer {
			case screens.JOIN_MEETING:
				temp := screenPointer
				if prevScreenPointer != screenPointer {
					joinMeetingScreen = screens.NewJoinMeetingScreen(&appState)
					joinMeetingScreen.Layout(gtx, &screenPointer)
				} else {
					joinMeetingScreen.Layout(gtx, &screenPointer)
				}
				if temp == screenPointer {
					prevScreenPointer = temp
				} else { // if screenPointer now points to different screen
					joinMeetingScreen = nil
				}
			case screens.JOIN_ROOM:
				temp := screenPointer
				if prevScreenPointer != screenPointer {
					joinRoomScreen = screens.NewJoinRoomScreen(&appState)
					joinRoomScreen.Layout(gtx, &screenPointer)
				} else {
					joinRoomScreen.Layout(gtx, &screenPointer)
				}
				if temp == screenPointer {
					prevScreenPointer = temp
				} else { // if screenPointer now points to different screen
					joinRoomScreen = nil
				}
			}
			e.Frame(gtx.Ops)
		}
	}
}
