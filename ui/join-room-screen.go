package ui

import (
	"image"
	"image/color"
	"log"
	"sync"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deoxyimran/gioui-livekit-client/ui/theme"
	"gocv.io/x/gocv"
)

type JoinRoomScreen struct {
	screenPointer *Screen
	th            *material.Theme
	userNameEdit  widget.Editor
	joinRoomBtn   widget.Clickable
	videoEnabled  chan bool
	Destroyed     bool
	frame         image.Image
	mutex         sync.Mutex
}

func NewJoinRoomScreen(screenPointer *Screen) *JoinRoomScreen {
	th := material.NewTheme()
	userNameEdit := widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	j := &JoinRoomScreen{
		th:            th,
		userNameEdit:  userNameEdit,
		screenPointer: screenPointer,
		videoEnabled:  make(chan bool),
	}

	// Initialize webcam and start video capture in a separate goroutine

	j.videoEnabled <- true // Start video capture
	return j
}

func (j *JoinRoomScreen) Destroy() {
	j.videoEnabled <- false // Stop video capture
	j.Destroyed = true
}

func (j *JoinRoomScreen) startVideoCapture() {
	cap, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		log.Printf("Error opening video capture device: %v", err)
		return
	}
	defer cap.Close()
	img := gocv.NewMat()
	defer img.Close()

	for enabled := range j.videoEnabled {
		if !enabled {
			return
		}

		if ok := cap.Read(&img); !ok {
			log.Println("Error reading from video capture device")
			continue
		}

		j.mutex.Lock()
		j.frame, err = img.ToImage()
		if err != nil {
			log.Printf("Error converting frame to image: %v\n", err)
			j.mutex.Unlock()
			continue
		}
		j.mutex.Unlock()
	}
}

func (j *JoinRoomScreen) Layout(gtx C) D {
	return layout.Background{}.Layout(gtx,
		func(gtx C) D {
			defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
			color := theme.BackgroundColor()
			paint.ColorOp{Color: color}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: gtx.Constraints.Max}
		},
		func(gtx C) D {
			gtx.Constraints.Min = image.Pt(0, 0) // Reset Constraints Min
			return layout.UniformInset(10).Layout(gtx,
				func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceAround}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
								// layout Username editor
								layout.Rigid(func(gtx C) D {
									return layout.Background{}.Layout(gtx,
										// Set a background
										func(gtx C) D {
											sz := gtx.Constraints.Min
											defer clip.UniformRRect(image.Rectangle{Max: sz}, 5).Push(gtx.Ops).Pop()
											paint.ColorOp{Color: color.NRGBA{255, 255, 255, 20}}.Add(gtx.Ops)
											paint.PaintOp{}.Add(gtx.Ops)
											return layout.Dimensions{Size: sz}
										},
										// Layout the webcam video canvas
										// func(gtx C) D {
										// 	return layout
										// }
										// Layout the editor
										func(gtx C) D {
											c := gtx.Constraints
											c.Max.X, c.Min.X = 200, 200
											gtx.Constraints = c
											edit := material.Editor(j.th, &j.userNameEdit, "Enter a name")
											edit.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
											edit.TextSize = unit.Sp(14)
											return layout.UniformInset(unit.Dp(10)).Layout(gtx, edit.Layout)
										},
									)
								}),
							)
						}),
					)
				},
			)
		},
	)
}
