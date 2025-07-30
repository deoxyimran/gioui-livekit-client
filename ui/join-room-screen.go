package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deoxyimran/gioui-livekit-client/ui/theme"
)

type JoinRoomScreen struct {
	screenPointer *Screen
	th            *material.Theme
	userNameEdit  widget.Editor
	joinRoomBtn   widget.Clickable
}

func NewJoinRoomScreen(screenPointer *Screen) JoinRoomScreen {
	th := material.NewTheme()
	userNameEdit := widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	return JoinRoomScreen{
		th:            th,
		userNameEdit:  userNameEdit,
		screenPointer: screenPointer,
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
			gtx.Constraints.Min = image.Pt(0, 0) // Reset Min
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
										// Layout the editor
										func(gtx C) D {
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
