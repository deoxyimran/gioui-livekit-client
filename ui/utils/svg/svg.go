package svg

import (
	"fmt"
	"image"
	"io"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

func LoadSvg(r io.Reader, targetSize image.Point) (image.Image, error) {
	svg, err := oksvg.ReadIconStream(r)
	if err != nil {
		return nil, fmt.Errorf("error decoding SVG file: %v", err)
	}
	var w, h int
	if targetSize != (image.Point{}) {
		w, h = targetSize.X, targetSize.Y
		svg.SetTarget(0, 0, float64(w), float64(h))
	} else {
		w, h = int(svg.ViewBox.W), int(svg.ViewBox.H)
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	scanner := rasterx.NewScannerGV(w, h, img, img.Bounds())
	raster := rasterx.NewDasher(w, h, scanner)
	svg.Draw(raster, 1.0)

	return img, nil
}
