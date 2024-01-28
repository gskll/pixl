package pxcanvas

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"github.com/gskll/pixl/apptype"
)

type PxCanvasMouseState struct {
	previousCoord *fyne.PointEvent
}

type PxCanvas struct {
	widget.BaseWidget
	apptype.PxCanvasConfig
	renderer    *PxCanvasRenderer
	PixelData   image.Image
	mouseState  PxCanvasMouseState
	appState    *apptype.State
	reloadImage bool
}

func (pc *PxCanvas) Bounds() image.Rectangle {
	x0 := int(pc.CanvasOffset.X)
	y0 := int(pc.CanvasOffset.Y)
	x1 := int(pc.PxCols*pc.PxSize + int(pc.CanvasOffset.X))
	y1 := int(pc.PxRows*pc.PxSize + int(pc.CanvasOffset.Y))

	return image.Rect(x0, y0, x1, y1)
}

func InBounds(pos fyne.Position, bounds image.Rectangle) bool {
	return pos.X >= float32(bounds.Min.X) &&
		pos.X < float32(bounds.Max.X) &&
		pos.Y >= float32(bounds.Min.Y) &&
		pos.Y < float32(bounds.Max.Y)
}

func NewBlankImage(cols, rows int, c color.Color) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, cols, rows))
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

func NewPxCanvas(state *apptype.State, config apptype.PxCanvasConfig) *PxCanvas {
	pxCanvas := &PxCanvas{
		PxCanvasConfig: config,
		appState:       state,
	}
	pxCanvas.PixelData = NewBlankImage(config.PxCols, config.PxRows, color.NRGBA{128, 128, 128, 128})
	pxCanvas.ExtendBaseWidget(pxCanvas)
	return pxCanvas
}

func (pc *PxCanvas) CreateRenderer() fyne.WidgetRenderer {
	canvasImage := canvas.NewImageFromImage(pc.PixelData)
	canvasImage.ScaleMode = canvas.ImageScalePixels
	canvasImage.FillMode = canvas.ImageFillContain

	canvasBorder := make([]canvas.Line, 4)
	for i := 0; i < len(canvasBorder); i++ {
		canvasBorder[i].StrokeColor = color.NRGBA{100, 100, 100, 255}
		canvasBorder[i].StrokeWidth = 2
	}

	renderer := &PxCanvasRenderer{
		pxCanvas:     pc,
		canvasImage:  canvasImage,
		canvasBorder: canvasBorder,
	}
	pc.renderer = renderer
	return renderer
}

func (pc *PxCanvas) TryPan(previousCoord *fyne.PointEvent, ev *desktop.MouseEvent) {
	if previousCoord != nil && ev.Button == desktop.MouseButtonTertiary {
		pc.Pan(*previousCoord, ev.PointEvent)
	}
}

// Brushable interface
func (pc *PxCanvas) SetColor(c color.Color, x, y int) {
	if nrgba, ok := pc.PixelData.(*image.NRGBA); ok {
		nrgba.Set(x, y, c)
	}
	if rgba, ok := pc.PixelData.(*image.RGBA); ok {
		rgba.Set(x, y, c)
	}
	pc.Refresh()
}

func (pc *PxCanvas) MouseToCanvasXY(ev *desktop.MouseEvent) (*int, *int) {
	bounds := pc.Bounds()
	if !InBounds(ev.Position, bounds) {
		return nil, nil
	}

	pxSize := float32(pc.PxSize)
	xOffset := pc.CanvasOffset.X
	yOffset := pc.CanvasOffset.Y

	x := int((ev.Position.X - xOffset) / pxSize)
	y := int((ev.Position.Y - yOffset) / pxSize)

	return &x, &y
}

func (pc *PxCanvas) LoadImage(img image.Image) {
	dimensions := img.Bounds()

	pc.PxCanvasConfig.PxCols = dimensions.Dx()
	pc.PxCanvasConfig.PxRows = dimensions.Dy()

	pc.PixelData = img
	pc.reloadImage = true
	pc.Refresh()
}

func (pc *PxCanvas) NewDrawing(cols, rows int) {
	pc.appState.SetFilePath("")
	pc.PxCols = cols
	pc.PxRows = rows
	pixelData := NewBlankImage(cols, rows, color.NRGBA{128, 128, 128, 255})
	pc.LoadImage(pixelData)
}
