package main

import "image/color"
import "github.com/hajimehoshi/ebiten/v2"
import "github.com/tinne26/ptxt"
import "github.com/tinne26/ggfnt-fonts/jumpy"

const CanvasWidth, CanvasHeight = 160, 90

type Game struct {
	canvas *ebiten.Image // logical canvas
	text *ptxt.Renderer
}

func (*Game) Layout(int, int) (int, int) { panic("F") }
func (*Game) LayoutF(logicWinWidth, logicWinHeight float64) (float64, float64) {
	scale := ebiten.Monitor().DeviceScaleFactor()
	return logicWinWidth*scale, logicWinHeight*scale
}

func (self *Game) Update() error { return nil }
func (self *Game) Draw(hiResCanvas *ebiten.Image) {
	// fill background
	self.canvas.Fill(color.RGBA{246, 242, 240, 255})

	// draw text
	self.text.Draw(self.canvas, "LITTLE LETTERS DANCING", CanvasWidth/2, CanvasHeight/2)

	// project logical canvas to main (optional ptxt utility)
	ptxt.PixelPerfect.Project(self.canvas, hiResCanvas)
}

// ---- main function ----

func main() {
	// initialize font strand
	strand, err := ptxt.NewStrand(jumpy.Font())
	if err != nil { panic(err) }
	var picker jumpy.GoldenPicker // see also PulsePicker
	strand.GlyphPickers().Add(&picker)
	
	// create text renderer, set the main properties
	renderer := ptxt.NewRenderer()
	renderer.SetStrand(strand)
	renderer.SetAlign(ptxt.Center)
	renderer.SetColor(color.RGBA{242, 143, 59, 255})

	// set up Ebitengine and start the game
	ebiten.SetWindowTitle("ptxt-examples/gpu/animate")
	canvas := ebiten.NewImage(CanvasWidth, CanvasHeight)
	err = ebiten.RunGame(&Game{ text: renderer, canvas: canvas })
	if err != nil { panic(err) }
}
