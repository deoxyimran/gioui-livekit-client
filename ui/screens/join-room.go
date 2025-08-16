package screens

import (
	"image"
	"image/color"
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deoxyimran/gioui-livekit-client/ui/components"
	"github.com/deoxyimran/gioui-livekit-client/ui/media/video"
	"github.com/deoxyimran/gioui-livekit-client/ui/theme"
	"github.com/deoxyimran/gioui-livekit-client/ui/utils"
)

type JoinRoom struct {
	// Widgets
	th               *material.Theme
	userNameEdit     widget.Editor
	joinRoomClickble widget.Clickable
	// States/control vars
	stateManager *StateManager
	vidCanvas    components.VideoCanvas
}

func NewJoinRoomScreen(stateManager *StateManager) *JoinRoom {
	th := material.NewTheme()
	userNameEdit := widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	vs := video.NewWebcamSource("")
	j := &JoinRoom{
		stateManager: stateManager,
		th:           th,
		userNameEdit: userNameEdit,
		vidCanvas:    components.NewVideoCanvas(&vs),
	}

	// Initialize webcam and start video capture in a separate goroutine
	return j
}

func (j *JoinRoom) StopVideoCapture() {

}

func (j *JoinRoom) Layout(gtx C, screenPointer *Screen) D {
	return layout.Background{}.Layout(gtx,
		// Fullscreen background
		func(gtx C) D {
			defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
			color := theme.BackgroundColor()
			paint.ColorOp{Color: color}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: gtx.Constraints.Max}
		},
		// Main content
		func(gtx C) D {
			gtx.Constraints.Min = image.Pt(0, 0) // Reset Constraints Min
			return layout.UniformInset(10).Layout(gtx,
				func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceAround}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return layout.UniformInset(unit.Dp(10)).Layout(gtx,
										func(gtx C) D {
											return layout.Flex{
												Axis:      layout.Vertical,
												Alignment: layout.Middle,
											}.Layout(gtx,
												// Video canvas
												layout.Rigid(
													func(gtx C) D {
														const w, h = 320, 240
														if vs := j.vidCanvas.GetVideoSource(); vs.IsVideoOn() {
															// Scale the image to fit 320x240 px
															defer clip.Rect(image.Rectangle{Max: image.Pt(w, h)}).Push(gtx.Ops).Pop()

															if vs.GetVideoOutFrame() == nil {
																paint.ColorOp{Color: color.NRGBA{R: 120, G: 120, B: 120, A: 255}}.Add(gtx.Ops)
															} else {
																f := vs.GetVideoOutFrame()
																scale := f32.Affine2D{}.Scale(f32.Point{}, f32.Point{
																	X: float32(w) / float32(f.Bounds().Dx()),
																	Y: float32(h) / float32(f.Bounds().Dy()),
																})
																op.Affine(scale).Add(gtx.Ops)
																paint.NewImageOp(f).Add(gtx.Ops)
															}
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
												// Username editor
												layout.Rigid(
													func(gtx C) D {
														c := gtx.Constraints
														c.Max.X, c.Min.X = 300, 300
														gtx.Constraints = c
														edit := material.Editor(j.th, &j.userNameEdit, "Enter a name")
														edit.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
														edit.HintColor = color.NRGBA{R: 135, G: 135, B: 135, A: 220}
														edit.TextSize = unit.Sp(14)
														return layout.UniformInset(unit.Dp(10)).Layout(gtx,
															func(gtx C) D {
																return utils.BorderLayout(gtx, edit.Layout, 1, 8, color.NRGBA{R: 140, G: 140, B: 140, A: 255})
															},
														)
													},
												),
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
