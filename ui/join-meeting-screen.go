package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type JoinMeetingScreen struct {
	screenPointer  *Screen
	th             *material.Theme
	joinMeetingBtn widget.Clickable
}

func NewJoinMeetingScreen(screenPointer *Screen) JoinMeetingScreen {
	return JoinMeetingScreen{th: material.NewTheme(), screenPointer: screenPointer}
}

func (j *JoinMeetingScreen) Layout(gtx C) D {
	return layout.Background{}.Layout(gtx,
		func(gtx C) D {
			defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
			color := color.NRGBA{R: 17, G: 17, B: 17, A: 255}
			paint.ColorOp{Color: color}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: gtx.Constraints.Max}
		},
		func(gtx C) D {
			gtx.Constraints.Min = image.Pt(0, 0) // Reset Min
			return layout.UniformInset(10).Layout(gtx,
				func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceAround}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
								// Heading
								layout.Rigid(func(gtx C) D {
									h := material.H4(j.th, "Family Meet")
									h.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
									h.Alignment = text.Middle
									return layout.UniformInset(unit.Dp(10)).Layout(gtx, h.Layout)
								}),
								// Join Meeting box
								layout.Rigid(func(gtx C) D {
									w := func(gtx C) D {
										x := (50.0 / 100.0) * float32(gtx.Constraints.Max.X)
										gtx.Constraints.Max = image.Pt(int(x), 200)
										gtx.Constraints.Min = image.Pt(int(x), 0)
										return layout.Flex{
											Axis:      layout.Vertical,
											Alignment: layout.Middle,
										}.Layout(gtx,
											layout.Rigid(func(gtx C) D {
												t := material.Label(j.th, unit.Sp(14), "Join meeting by clicking \"Join Meeting\" below")
												t.Alignment = text.Middle
												t.Color = color.NRGBA{255, 255, 255, 255}
												return layout.UniformInset(unit.Dp(5)).Layout(gtx, t.Layout)
											}),
											// Join Meeting Button
											layout.Flexed(1.0, func(gtx C) D {
												gtx.Constraints.Max.Y = 60
												b := material.Button(j.th, &j.joinMeetingBtn, "Join Meeting")
												if j.joinMeetingBtn.Clicked(gtx) {
													*j.screenPointer = JOIN_ROOM_SCREEN
												}
												return layout.UniformInset(unit.Dp(10)).Layout(gtx, b.Layout)
											}),
										)
									}
									return borderLayout(gtx, w, 1, color.NRGBA{150, 150, 150, 255})
								}),
							)
						}),
					)
				},
			)
		},
	)
}
