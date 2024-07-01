package main

import "os"
import "fmt"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/ptxt"

const CanvasWidth, CanvasHeight = 160, 90

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
	renderer.SetAlign(ptxt.Baseline | ptxt.Right)
	renderer.SetColor(color.RGBA{239, 91, 91, 255})

	// run game
	ebiten.SetWindowTitle("ptxt-examples/gpu/sideways")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	err = ebiten.RunGame(&Game{
		text: renderer,
		canvas: ebiten.NewImage(CanvasWidth, CanvasHeight),
	})
	if err != nil { panic(err) }
}

type Game struct {
	text *ptxt.Renderer
	canvas *ebiten.Image
}

func (*Game) Update() error { return nil }
func (*Game) Layout(_, _ int) (int, int) { panic("F") }
func (*Game) LayoutF(logicWinWidth, logicWinHeight float64) (float64, float64) {
	scale := ebiten.DeviceScaleFactor()
	return logicWinWidth*scale, logicWinHeight*scale
}

func (self *Game) Draw(hiResCanvas *ebiten.Image) {
	// background color
	self.canvas.Fill(color.RGBA{229, 255, 222, 255})

	// draw text
	const SampleText = "SIDEWAYS"
	px := int(self.text.GetScale())
	w, _ := self.text.Measure(SampleText)
	h := int(self.text.Strand().Font().Metrics().UppercaseAscent())*px
	side := w + h + px

	// actual drawing
	cx, cy := CanvasWidth/2, CanvasHeight/2
	cx -= side/2
	cy -= side/2
	self.text.SetDirection(ptxt.Horizontal)
	self.text.Draw(self.canvas, SampleText, cx + side, cy + h)
	self.text.SetDirection(ptxt.SidewaysRight)
	self.text.Draw(self.canvas, SampleText, cx + side - h, cy + side)
	self.text.SetDirection(ptxt.Horizontal)
	self.text.Draw(self.canvas, SampleText, cx + side - h - px, cy + side)
	self.text.SetDirection(ptxt.Sideways)
	self.text.Draw(self.canvas, SampleText, cx + h, cy)

	// project logical canvas to main (optional ptxt utility)
	ptxt.Proportional.Project(self.canvas, hiResCanvas)
}
