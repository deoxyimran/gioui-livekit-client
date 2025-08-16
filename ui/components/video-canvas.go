package components

import (
	"image"
	"sync"

	"gioui.org/layout"
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
	isVideoRunning bool
	bgJobStopSig   chan bool
}

func NewVideoCanvas(s video.VideoSource) VideoCanvas {
	return VideoCanvas{vs: s}
}

func (v *VideoCanvas) GetVideoSource() video.VideoSource {
	return v.vs
}

func (v *VideoCanvas) Layout(gtx C) D {
	return layout.Dimensions{}
}
