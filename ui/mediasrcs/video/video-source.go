package video

import "image"

type VideoSource interface {
	SetDevice(string)
	GetVideoOutFrame() image.Image
	IsVideoOn() bool
	StartVideo() error
	StopVideo()
}
