package comp

import (
	"github.com/rivo/tview"
)

type Image struct {
	Base[*tview.Image]
}

func CreateImage() Component {
	image := &Image{
		Base: Base[*tview.Image]{
			name:      "image",
			tviewInst: tview.NewImage(),
		},
	}

	image.Base.outerInst = image

	return image
}
