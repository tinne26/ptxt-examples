package main

import ( "math" ; "image/color" )
import "github.com/hajimehoshi/ebiten/v2"
import "github.com/tinne26/ptxt"
import "github.com/tinne26/ggfnt-fonts/jammy"

const CanvasWidth, CanvasHeight = 80, 45 // (1/24th of 1920x1080)
const WordsPerSec = 2.71828
var Words = []string {
	"PIXEL", "BLOCK", "RETRO", "WORLD", "CLICKY", "GAME",
	"SHARP", "CONTROL", "SIMPLE", "PLAIN", "COLOR", "PALETTE",
}

// ---- Ebitengine's Game interface implementation ----

type Game struct {
	canvas *ebiten.Image // logical canvas
	text *ptxt.Renderer
	wordIndex float64
}

func (*Game) Layout(int, int) (int, int) { panic("F") }
func (*Game) LayoutF(logicWinWidth, logicWinHeight float64) (float64, float64) {
	scale := ebiten.Monitor().DeviceScaleFactor()
	return logicWinWidth*scale, logicWinHeight*scale
}

func (self *Game) Update() error {
	newIndex := (self.wordIndex + WordsPerSec/60.0)
	self.wordIndex = math.Mod(newIndex, float64(len(Words)))
	return nil
}

func (self *Game) Draw(hiResCanvas *ebiten.Image) {
	// fill background
	self.canvas.Fill(color.RGBA{246, 242, 240, 255})

	// draw text
	word := Words[int(self.wordIndex)]
	self.text.Draw(self.canvas, word, 6, CanvasHeight - 6)

	// project logical canvas to main (optional ptxt utility)
	ptxt.PixelPerfect.Project(self.canvas, hiResCanvas)
}

// ---- main function ----

func main() {
	// initialize font strand
	strand, err := ptxt.NewStrand(jammy.Font())
	if err != nil { panic(err) }
	
	// create text renderer, set the main properties
	renderer := ptxt.NewRenderer()
	renderer.SetStrand(strand)
	renderer.SetAlign(ptxt.Baseline | ptxt.Left)
	renderer.SetColor(color.RGBA{242, 143, 59, 255})

	// set up Ebitengine and start the game
	ebiten.SetWindowTitle("ptxt-examples/gpu/words")
	canvas := ebiten.NewImage(CanvasWidth, CanvasHeight)
	err = ebiten.RunGame(&Game{ text: renderer, canvas: canvas })
	if err != nil { panic(err) }
}
