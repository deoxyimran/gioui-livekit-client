package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type JoinMeetingScreen struct {
	th *material.Theme
}

func NewJoinMeetingScreen() JoinMeetingScreen {
	return JoinMeetingScreen{th: material.NewTheme()}
}

func (j JoinMeetingScreen) Layout(gtx C) D {
	return layout.Background{}.Layout(gtx,
		func(gtx C) D {
			defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
			color := color.NRGBA{R: 17, G: 17, B: 17, A: 255}
			paint.ColorOp{Color: color}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: gtx.Constraints.Max}
		},
		func(gtx C) D {
			return layout.UniformInset(10).Layout(gtx,
				func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							h := material.H4(j.th, "Family Meet")
							h.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
							h.Alignment = text.Middle
							return h.Layout(gtx)
						}),
					)
				},
			)
		},
	)
}

// Layout utils
func (j JoinMeetingScreen) borderBoxLayout(gtx C, inner layout.Widget, width float32, color color.NRGBA) D {
	// Layout inner first
	inner(gtx)

	sz := image.Point{X: 320, Y: 150}
	defer clip.Stroke{
		Width: width,
		Path:  clip.UniformRRect(image.Rectangle{Max: sz}, 10).Path(gtx.Ops),
	}.Op().Push(gtx.Ops).Pop()
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{Size: sz}
}
