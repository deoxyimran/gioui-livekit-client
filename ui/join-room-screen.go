package ui

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"sync"
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
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
	videoStop     chan bool
	isVideoOn     bool
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
		videoStop:     make(chan bool),
	}

	// Initialize webcam and start video capture in a separate goroutine
	j.startVideoCapture()
	return j
}

func (j *JoinRoomScreen) StopVideoCapture() {
	if j.isVideoOn {
		j.isVideoOn = false
		j.videoStop <- true // Stop video capture
	}
}

func (j *JoinRoomScreen) startVideoCapture() {
	j.isVideoOn = true
	go func() {
		cap, err := gocv.VideoCaptureDevice(0)
		if err != nil {
			log.Printf("Error opening video capture device: %v", err)
			return
		}
		defer cap.Close()
		mat := gocv.NewMat()
		defer mat.Close()

		loop := true
		for loop {
			select {
			case b := <-j.videoStop:
				if b {
					loop = !loop
				}
			default:
				fmt.Println("Capturing video frame...")
				if ok := cap.Read(&mat); !ok {
					log.Println("Error reading from video capture device")
					continue
				}

				j.mutex.Lock()
				j.frame, err = mat.ToImage()
				if err != nil {
					log.Printf("Error converting frame to image: %v\n", err)
					j.mutex.Unlock()
					continue
				}
				j.mutex.Unlock()
			}
		}
	}()
}

func (j *JoinRoomScreen) Layout(gtx C) D {
	return layout.Background{}.Layout(gtx, // Fullscreen background
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
										// Layout the content
										func(gtx C) D {
											return layout.UniformInset(unit.Dp(10)).Layout(gtx,
												func(gtx C) D {
													return layout.Flex{
														Axis:      layout.Vertical,
														Alignment: layout.Middle,
													}.Layout(gtx,
														layout.Rigid(
															func(gtx C) D {
																const w, h = 320, 240
																if j.isVideoOn {
																	// Scale the image to fit 320x240 px
																	defer clip.Rect(image.Rectangle{Max: image.Pt(w, h)}).Push(gtx.Ops).Pop()

																	j.mutex.Lock()
																	if j.frame == nil {
																		paint.ColorOp{Color: color.NRGBA{R: 120, G: 120, B: 120, A: 255}}.Add(gtx.Ops)
																	} else {
																		scale := f32.Affine2D{}.Scale(f32.Point{}, f32.Point{
																			X: float32(w) / float32(j.frame.Bounds().Dx()),
																			Y: float32(h) / float32(j.frame.Bounds().Dy()),
																		})
																		op.Affine(scale).Add(gtx.Ops)
																		paint.NewImageOp(j.frame).Add(gtx.Ops)
																	}
																	j.mutex.Unlock()
																	paint.PaintOp{}.Add(gtx.Ops)
																	gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(time.Second / 30)}) // Cap to 30Fps

																	return layout.Dimensions{Size: image.Pt(w, h)}
																} else {
																	defer clip.Rect{Max: image.Pt(w, h)}.Push(gtx.Ops).Pop()
																	paint.ColorOp{Color: color.NRGBA{R: 120, G: 120, B: 120, A: 255}}.Add(gtx.Ops)
																	paint.PaintOp{}.Add(gtx.Ops)
																	return layout.Dimensions{Size: image.Pt(w, h)}
																}
															},
														),
														layout.Rigid(
															func(gtx C) D {
																c := gtx.Constraints
																c.Max.X, c.Min.X = 200, 200
																gtx.Constraints = c
																edit := material.Editor(j.th, &j.userNameEdit, "Enter a name")
																edit.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
																edit.TextSize = unit.Sp(14)
																return edit.Layout(gtx)
															},
														),
													)
												},
											)
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
