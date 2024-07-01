package main

import "os"
import "fmt"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"
import "github.com/hajimehoshi/ebiten/v2/inpututil"

import "github.com/tinne26/ptxt"
import "github.com/tinne26/ggfnt"

const CanvasWidth, CanvasHeight = 160, 90

var BackColor = color.RGBA{24, 24, 24, 255}
var TextColor = color.RGBA{250, 250, 250, 255}

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

	// run game
	ebiten.SetWindowTitle("ptxt-examples/gpu/glyphs")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	game := Game{
		text: renderer,
		canvas: ebiten.NewImage(CanvasWidth, CanvasHeight),
	}
	game.Init()
	err = ebiten.RunGame(&game)
	if err != nil { panic(err) }
}

type Game struct {
	text *ptxt.Renderer
	canvas *ebiten.Image
	startIndex int
	glyphCount int
	boxSize int
	topOffset int

	maxGlyphsPerRow int
	maxGlyphsPerCol int
	leftMargin int
	topMargin int
}

func (self *Game) Init() {
	font := self.text.Strand().Font()
	self.glyphCount = int(font.Glyphs().Count())
	self.boxSize = int(font.Metrics().LineHeight()) + 2
	self.topOffset = int(font.Metrics().Ascent()) + 1

	// compute max glyphs per line and so on
	self.maxGlyphsPerRow = CanvasWidth/self.boxSize
	self.leftMargin = (CanvasWidth - (self.maxGlyphsPerRow*self.boxSize))/2
	self.maxGlyphsPerCol = CanvasHeight/self.boxSize
	self.topMargin = (CanvasHeight - (self.maxGlyphsPerCol*self.boxSize))/2
}

func (self *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		if self.startIndex >= self.maxGlyphsPerRow {
			self.startIndex -= self.maxGlyphsPerRow
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		if self.startIndex + self.maxGlyphsPerRow*(max(self.maxGlyphsPerCol/2, 1)) < self.glyphCount {
			self.startIndex += self.maxGlyphsPerRow
		}
	}
	return nil
}

func (*Game) Layout(_, _ int) (int, int) { panic("F") }
func (*Game) LayoutF(logicWinWidth, logicWinHeight float64) (float64, float64) {
	scale := ebiten.DeviceScaleFactor()
	return logicWinWidth*scale, logicWinHeight*scale
}

func (self *Game) Draw(hiResCanvas *ebiten.Image) {
	// background color
	self.canvas.Fill(BackColor)

	// initialize mask drawing params and set the color
	var params ptxt.MaskDrawParameters
	params.Scale = 1
	params.RGBA = [4]float32{
		float32(TextColor.R)/255.0, float32(TextColor.G)/255.0,
		float32(TextColor.B)/255.0, float32(TextColor.A)/255.0,
	}
	
	// draw glyphs
	glyphIndex := self.startIndex
	strand := self.text.Strand()
loop:
	for y := 0; y < self.maxGlyphsPerCol; y++ {
		for x := 0; x < self.maxGlyphsPerRow; x++ {
			mask := self.text.Advanced().LoadMask(ggfnt.GlyphIndex(glyphIndex))
			if mask != nil {
				params.X  = self.leftMargin + x*self.boxSize + (self.boxSize >> 1)
				params.X -= mask.Bounds().Min.X + mask.Bounds().Dx()/2
				params.Y  = self.topMargin + self.topOffset + y*self.boxSize
				self.text.Advanced().DrawMask(self.canvas, mask, strand, params)
			}
			glyphIndex += 1
			if glyphIndex >= self.glyphCount { break loop }
		}
	}

	// project logical canvas to main (optional ptxt utility)
	ptxt.Proportional.Project(self.canvas, hiResCanvas)
}
