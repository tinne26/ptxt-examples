package main

import "os"
import "fmt"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"
import "github.com/hajimehoshi/ebiten/v2/inpututil"

import "github.com/tinne26/ptxt"

const CanvasWidth, CanvasHeight = 160, 90
var BackColor = color.RGBA{25, 100, 126, 255}
var TextColor = color.RGBA{40, 175, 176, 255}
var ShadowColor = color.RGBA{55, 57, 46, 255}

// TODO: yeah, this example doesn't seem to set or use the ascent
//       offset appropriately. Check ptxt and addd more tests...

func main() {
	// usage check
	if len(os.Args) != 2 {
		fmt.Print("Usage: go run main.go font.ggfnt\n")
		os.Exit(1)
	}

	// parse font and create strand
	strand, err := ptxt.NewStrand(os.Args[1])
	if err != nil { panic(err) }
	fmt.Printf("Font loaded: %s\n", strand.Font().Header().Name())

	// create strand and renderer
	renderer := ptxt.NewRenderer()
	renderer.SetStrand(strand)
	renderer.SetScale(1)
	renderer.SetAlign(ptxt.Center)
	renderer.SetColor(TextColor)

	// shadow
	// strand.Shadow().SetStrand(strand)
	// strand.Shadow().SetOffsets(-1, 0)
	// strand.Shadow().SetColor(ShadowColor)

	// run game
	ebiten.SetWindowTitle("ptxt-examples/gpu/bounding")
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

func (*Game) Layout(_, _ int) (int, int) { panic("F") }
func (self *Game) LayoutF(logicWinWidth, logicWinHeight float64) (float64, float64) {
	scale := ebiten.DeviceScaleFactor()
	return logicWinWidth*scale, logicWinHeight*scale
}

func (self *Game) Update() error {
	// detect bounding mode changes
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		switch self.text.Advanced().GetBoundingMode() {
		case ptxt.LogicalBounding       : self.text.Advanced().SetBoundingMode(ptxt.NoDescLogicalBounding)
		case ptxt.MaskBounding          : self.text.Advanced().SetBoundingMode(ptxt.NoDescMaskBounding)
		case ptxt.NoDescLogicalBounding : self.text.Advanced().SetBoundingMode(ptxt.MaskBounding)
		case ptxt.NoDescMaskBounding    : self.text.Advanced().SetBoundingMode(ptxt.LogicalBounding)
		default:
			panic("unexpected bounding mode")
		}
	}

	return nil
}

func (self *Game) Draw(hiResCanvas *ebiten.Image) {
	// background color
	self.canvas.Fill(BackColor)

	mode := self.text.Advanced().GetBoundingMode().String()
	info := "" + 
		"TESTING BOUNDING MODE:\n" +
		FmtCamelASCII(mode) + "\n" + 
		"...PRESS [M] TO CHANGE...\n"
	self.text.Draw(self.canvas, info, CanvasWidth/2, CanvasHeight/2)

	// project logical canvas to main (optional ptxt utility)
	ptxt.Proportional.Project(self.canvas, hiResCanvas)
}

// format camel case ascii into hyphen separated uppercase
func FmtCamelASCII(str string) string {
	var bytes []byte = make([]byte, 0, len(str))
	var fromLower bool = false
	for _, codePoint := range str {
		if codePoint > 128 { panic("expected only ascii") }
		if codePoint >= 'a' && codePoint <= 'z' {
			fromLower = true
			bytes = append(bytes, byte(codePoint - 'a' + 'A'))
		} else if codePoint >= 'A' && codePoint <= 'Z' {
			if fromLower { bytes = append(bytes, '-') }
			fromLower = false
			bytes = append(bytes, byte(codePoint))
		} else if codePoint == '|' {
			fromLower = false
			bytes[len(bytes) - 1] = ','
		} else {
			fromLower = false
			bytes = append(bytes, byte(codePoint))
		}
	}
	return string(bytes)
}
