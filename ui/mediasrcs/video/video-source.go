package video

import "image"

type VideoSource interface {
	SetDeviceName(string)
	GetDeviceName() string
	GetVideoOutFrame() image.Image
	IsVideoOn() bool
	StartVideo() error
	StopVideo()
}
