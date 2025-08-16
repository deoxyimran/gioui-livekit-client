package video

import "image"

type VideoSource interface {
	SetDevPath(string)
	GetVideoOutFrame() image.Image
	IsVideoOn() bool
	StartVideo() error
	StopVideo()
}
