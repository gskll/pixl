package ui

import (
	fyne "fyne.io/fyne/v2"

	"github.com/gskll/pixl/apptype"
	"github.com/gskll/pixl/swatch"
)

type AppInit struct {
	PixlWindow fyne.Window
	State      *apptype.State
	Swatches   []*swatch.Swatch
}
