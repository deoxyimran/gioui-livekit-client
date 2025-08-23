package theme

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/deoxyimran/gioui-livekit-client/ui/res/img"
)

var decodedBg image.Image

func AppBackground() image.Image {
	if decodedBg == nil {
		decodedBg, _ = jpeg.Decode(bytes.NewReader(img.AppBg))
	}
	return decodedBg
}
