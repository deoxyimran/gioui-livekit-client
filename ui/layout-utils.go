package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// Layout utils
func borderLayout(gtx C, inner layout.Widget, width float32, color color.NRGBA) D {
	// Layout inner first
	dims := inner(gtx)

	// Layout border
	defer clip.Stroke{
		Width: width,
		Path:  clip.UniformRRect(image.Rectangle{Max: dims.Size}, 10).Path(gtx.Ops),
	}.Op().Push(gtx.Ops).Pop()
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return dims
}
