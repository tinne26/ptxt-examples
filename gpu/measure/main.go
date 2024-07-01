package main

import "os"
import "fmt"
import "unicode"
import "image"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"
import "github.com/hajimehoshi/ebiten/v2/inpututil"

import "github.com/tinne26/ptxt"

const CanvasWidth, CanvasHeight = 160, 90

var BackgroundColor color.RGBA = color.RGBA{142, 166,   4, 255}
var HighlightColor  color.RGBA = color.RGBA{ 38, 131,  17, 255}
var TextColor       color.RGBA = color.RGBA{222, 235,  76, 255}

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
	renderer.SetAlign(ptxt.Top | ptxt.Left)
	renderer.SetColor(TextColor)
	renderer.Advanced().SetParBreakEnabled(true)
	// renderer.Advanced().SetBoundingMode(ptxt.MaskBounding) // uncomment if you want to test this

	// apply rewrite rules (optional)
	err = strand.Mapping().AutoInitRewriteRules()
	if err != nil { panic(err) }

	// heuristic for lowercase support
	var initContent []rune
	var allowLowercase bool
	if renderer.Advanced().AllGlyphsAvailable("abcdefghijklmnopqrstuvwxyz") {
		allowLowercase = true
		initContent = []rune("Type something!")
	} else {
		initContent = []rune("TYPE SOMETHING!")
	}

	// run game
	ebiten.SetWindowTitle("ptxt-examples/gpu/measure")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	err = ebiten.RunGame(&Game{
		text: renderer,
		canvas: ebiten.NewImage(CanvasWidth, CanvasHeight),
		content: initContent,
		allowLowercase: allowLowercase,
	})
	if err != nil { panic(err) }
}

type Game struct {
	text *ptxt.Renderer
	canvas *ebiten.Image
	content []rune // not very efficient, but AppendInputChars uses runes
	underscoreTicker uint8
	allowLowercase bool
}

func (*Game) Layout(_, _ int) (int, int) { panic("F") }
func (*Game) LayoutF(logicWinWidth, logicWinHeight float64) (float64, float64) {
	scale := ebiten.DeviceScaleFactor()
	return logicWinWidth*scale, logicWinHeight*scale
}

func (self *Game) Update() error {
	// helper function
	var keyRepeat = func(key ebiten.Key) bool {
		ticks := inpututil.KeyPressDuration(key)
		return ticks == 1 || (ticks > 14 && (ticks - 14) % 5 == 0)
	}

	// detect enter for newline, backspace for removing text,
	// and otherwise append any new text input we get
	if keyRepeat(ebiten.KeyBackspace) && len(self.content) >= 1 {
		self.content = self.content[0 : len(self.content) - 1]
	} else if keyRepeat(ebiten.KeyEnter) {
		self.content = append(self.content, '\n')
	} else {
		preLen := len(self.content)
		self.content = ebiten.AppendInputChars(self.content)
		if !self.allowLowercase {
			for i := preLen; i < len(self.content); i++ {
				self.content[i] = unicode.ToUpper(self.content[i])
			}
		}
	}

	// update underscore ticker
	self.underscoreTicker += 1
	if self.underscoreTicker > 64 {
		self.underscoreTicker = 0
	}
	
	return nil
}

func (self *Game) Draw(hiResCanvas *ebiten.Image) {
	// background color
	self.canvas.Fill(BackgroundColor)

	// draw highlight rect
	content := string(self.content)
	w, h := self.text.Measure(content)
	ox, _ := self.text.Advanced().LastBoundsOffset() // (only relevant for MaskBounding mode)
	rect := image.Rect(4 + ox, 4, 4 + ox + w, 4 + h)
	textArea := self.canvas.SubImage(rect).(*ebiten.Image)
	textArea.Fill(HighlightColor)

	// draw text
	if self.underscoreTicker < 24 { content += "_" }
	self.text.Draw(self.canvas, content, 4, 4)

	// project logical canvas to main (optional ptxt utility)
	ptxt.Proportional.Project(self.canvas, hiResCanvas)
}
