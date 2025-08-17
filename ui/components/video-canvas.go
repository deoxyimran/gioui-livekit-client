package components

import (
	"image"
	"image/color"
	"sync"
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/deoxyimran/gioui-livekit-client/ui/media/video"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type VideoCanvas struct {
	frame          image.Image
	mutex          sync.Mutex
	vs             video.VideoSource
	size           image.Point
	isVideoRunning bool
	bgJobStopSig   chan bool
}

func NewVideoCanvas(s video.VideoSource, size image.Point) VideoCanvas {
	return VideoCanvas{vs: s, size: size}
}

func (v *VideoCanvas) GetVideoSource() video.VideoSource {
	return v.vs
}

func (v *VideoCanvas) Layout(gtx C) D {
	if vs := v.GetVideoSource(); vs.IsVideoOn() {
		// Scale the image to fit 320x240 px
		defer clip.Rect(image.Rectangle{Max: v.size}).Push(gtx.Ops).Pop()

		if vs.GetVideoOutFrame() == nil {
			paint.ColorOp{Color: color.NRGBA{R: 120, G: 120, B: 120, A: 255}}.Add(gtx.Ops)
		} else {
			f := vs.GetVideoOutFrame()
			scale := f32.Affine2D{}.Scale(f32.Point{}, f32.Point{
				X: float32(v.size.X) / float32(f.Bounds().Dx()),
				Y: float32(v.size.Y) / float32(f.Bounds().Dy()),
			})
			op.Affine(scale).Add(gtx.Ops)
			paint.NewImageOp(f).Add(gtx.Ops)
		}
		paint.PaintOp{}.Add(gtx.Ops)
		gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(time.Second / 30)}) // Cap to 30Fps

		return layout.Dimensions{Size: v.size}
	} else {
		defer clip.Rect{Max: v.size}.Push(gtx.Ops).Pop()
		paint.ColorOp{Color: color.NRGBA{R: 120, G: 120, B: 120, A: 255}}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		return layout.Dimensions{Size: v.size}
	}
}
