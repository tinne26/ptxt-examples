package main

import "os"
import "fmt"
import "image"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"
import "github.com/hajimehoshi/ebiten/v2/inpututil"

import "github.com/tinne26/ptxt"

const CanvasWidth, CanvasHeight = 160, 90

var BackgroundColor color.RGBA = color.RGBA{ 59,  82,  73, 255}
var WrapLineColor   color.RGBA = color.RGBA{ 79, 102,  93, 255}
var TextColor       color.RGBA = color.RGBA{  6, 167, 125, 255}

func main() {
	// usage check
	if len(os.Args) != 2 {
		fmt.Print("Usage: go run main.go font.ggfnt\n")
		os.Exit(1)
	}

	// open font file
	fontName := os.Args[1]
	fontFile, err := os.Open(fontName)
	if err != nil { panic(err) }

	// parse font and create strand
	strand, err := ptxt.NewStrand(fontFile)
	if err != nil { panic(err) }
	fmt.Printf("Font loaded: %s\n", strand.Font().Header().Name())

	// create strand and renderer
	renderer := ptxt.NewRenderer()
	renderer.SetStrand(strand)
	renderer.SetScale(1)
	renderer.SetAlign(ptxt.CapLine | ptxt.Left)
	renderer.SetColor(TextColor)
	renderer.Advanced().SetParBreakEnabled(true)
	
	// set up a soft shadow
	// strand.Shadow().SetStrand(strand)
	// strand.Shadow().SetOffsets(1, 1)
	// strand.Shadow().SetColor(color.RGBA{52, 21, 21, 52})

	// run game
	ebiten.SetWindowTitle("ptxt-examples/gpu/wrap")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	err = ebiten.RunGame(&Game{
		text: renderer,
		canvas: ebiten.NewImage(CanvasWidth, CanvasHeight),
		wrapX: CanvasWidth - (3*CanvasHeight/20),
	})
	if err != nil { panic(err) }
}

type Game struct {
	text *ptxt.Renderer
	canvas *ebiten.Image
	lastWidth, lastHeight float64
	wrapX int
}

func (*Game) Layout(_, _ int) (int, int) { panic("F") }
func (self *Game) LayoutF(logicWinWidth, logicWinHeight float64) (float64, float64) {
	scale := ebiten.DeviceScaleFactor()
	self.lastWidth  = logicWinWidth*scale
	self.lastHeight = logicWinHeight*scale
	return self.lastWidth, self.lastHeight
}

func (self *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		hiWidth, hiHeight := int(self.lastWidth), int(self.lastHeight)
		x, y := ebiten.CursorPosition()
		x,  _ = ptxt.Proportional.Remap(x, y, hiWidth, hiHeight, CanvasWidth, CanvasHeight)
		self.wrapX = x
	}
	return nil
}

func (self *Game) Draw(hiResCanvas *ebiten.Image) {
	// background color
	self.canvas.Fill(BackgroundColor)

	// draw lines
	pad    := CanvasHeight/10
	pad1p5 := (pad*3)/2
	self.canvas.SubImage(image.Rect(pad, 0, pad + 1, CanvasHeight)).(*ebiten.Image).Fill(WrapLineColor)
	self.canvas.SubImage(image.Rect(CanvasWidth - pad, 0, CanvasWidth - pad + 1, CanvasHeight)).(*ebiten.Image).Fill(WrapLineColor)
	self.canvas.SubImage(image.Rect(0, pad, CanvasWidth, pad + 1)).(*ebiten.Image).Fill(WrapLineColor)
	self.canvas.SubImage(image.Rect(0, CanvasHeight - pad, CanvasWidth, CanvasHeight - pad + 1)).(*ebiten.Image).Fill(WrapLineColor)

	// draw occupied rect
	wrapXStart := min(CanvasWidth - pad, max(pad, self.wrapX))
	self.canvas.SubImage(image.Rect(wrapXStart, pad, CanvasWidth - pad + 1, CanvasHeight - pad)).(*ebiten.Image).Fill(WrapLineColor)

	// draw text
	text := "You may click at any point within this rectangle in order to adjust the wrapping point, which is displayed as a lighter box on the right."
	self.text.DrawWithWrap(self.canvas, text, pad1p5, pad1p5, wrapXStart - pad1p5)

	// project logical canvas to main (optional ptxt utility)
	ptxt.Proportional.Project(self.canvas, hiResCanvas)
}
