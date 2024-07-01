package main

import "os"
import "fmt"
import "image"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"
import "github.com/hajimehoshi/ebiten/v2/inpututil"

import "github.com/tinne26/ptxt"

const CanvasWidth, CanvasHeight = 640, 360
const UpperInstructions = "CLICK AROUND TO SET DRAW COORDINATES\nUSE ARROWS TO CHANGE ALIGNS\nUSE D TO CHANGE TEXT DIRECTION"
const LowerInstructions = "Click around to set draw coordinates\nUse ARROWS to change aligns\nUse D to change text direction"

var VertAligns []ptxt.Align = []ptxt.Align{
	ptxt.Bottom, ptxt.LastBaseline, ptxt.Baseline,
	ptxt.VertCenter, ptxt.Midline, ptxt.CapLine, ptxt.Top,
}
var HorzAligns []ptxt.Align = []ptxt.Align{
	ptxt.Right, ptxt.HorzCenter, ptxt.Left,
}
var Directions []ptxt.Direction = []ptxt.Direction{
	ptxt.Horizontal, ptxt.Sideways, ptxt.SidewaysRight,
}

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
	renderer.SetScale(4)
	renderer.SetAlign(VertAligns[0] | HorzAligns[0])
	renderer.SetDirection(Directions[0])
	renderer.SetColor(color.RGBA{255, 214, 175, 255})

	// create helper text renderer
	infoRenderer := ptxt.NewRenderer()
	infoRenderer.SetStrand(strand)
	infoRenderer.SetScale(2)

	// run game
	ebiten.SetWindowTitle("ptxt-examples/gpu/aligns")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	err = ebiten.RunGame(&Game{
		text: renderer,
		info: infoRenderer,
		canvas: ebiten.NewImage(CanvasWidth, CanvasHeight),
		cx: CanvasWidth/2.0,
		cy: CanvasHeight/2.0,
		uppercaseOnly: !renderer.Advanced().AllGlyphsAvailable("abcdefghijklmnopqrstuvwzyx"),
	})
	if err != nil { panic(err) }
}

type Game struct {
	text *ptxt.Renderer
	info *ptxt.Renderer
	canvas *ebiten.Image
	cx, cy int
	horzAlignIndex, vertAlignIndex int
	dirIndex int
	hiResWidth, hiResHeight float64
	uppercaseOnly bool
}

func (*Game) Layout(_, _ int) (int, int) { panic("F") }
func (self *Game) LayoutF(logicWinWidth, logicWinHeight float64) (float64, float64) {
	scale := ebiten.DeviceScaleFactor()
	self.hiResWidth  = logicWinWidth*scale
	self.hiResHeight = logicWinHeight*scale
	return self.hiResWidth, self.hiResHeight
}

func (self *Game) Update() error {
	// detect horz align changes
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		self.horzAlignIndex -= 1
		if self.horzAlignIndex < 0 {
			self.horzAlignIndex = len(HorzAligns) - 1
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		self.horzAlignIndex += 1
		if self.horzAlignIndex >= len(HorzAligns) {
			self.horzAlignIndex = 0
		}
	}

	// detect vert align changes
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		self.vertAlignIndex -= 1
		if self.vertAlignIndex < 0 {
			self.vertAlignIndex = len(VertAligns) - 1
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		self.vertAlignIndex += 1
		if self.vertAlignIndex >= len(VertAligns) {
			self.vertAlignIndex = 0
		}
	}

	// update align
	align := VertAligns[self.vertAlignIndex] | HorzAligns[self.horzAlignIndex]
	self.text.SetAlign(align)

	// detect text dir changes
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
			self.dirIndex -= 1
			if self.dirIndex < 0 {
				self.dirIndex = len(Directions) - 1
			}
		} else {
			self.dirIndex += 1
			if self.dirIndex >= len(Directions) {
				self.dirIndex = 0
			}
		}
		self.text.SetDirection(Directions[self.dirIndex])
	}

	// detect cursor position
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		fromWidth, fromHeight := int(self.hiResWidth), int(self.hiResHeight)
		toWidth, toHeight     := CanvasWidth, CanvasHeight
		self.cx, self.cy = ptxt.Proportional.Remap(x, y, fromWidth, fromHeight, toWidth, toHeight)
	}

	return nil
}

func (self *Game) Draw(hiResCanvas *ebiten.Image) {
	// background color
	self.canvas.Fill(color.RGBA{131, 151, 136, 255})

	// crossing lines
	horz := self.canvas.SubImage(image.Rect(0, self.cy, CanvasWidth, self.cy + 1)).(*ebiten.Image)
	horz.Fill(color.RGBA{161, 181, 166, 255})
	vert := self.canvas.SubImage(image.Rect(self.cx, 0, self.cx + 1, CanvasHeight)).(*ebiten.Image)
	vert.Fill(color.RGBA{161, 181, 166, 255})

	// draw instructions
	self.info.SetColor(color.RGBA{171, 191, 176, 255})
	self.info.SetAlign(ptxt.Top | ptxt.Left)
	if self.uppercaseOnly {
		self.info.Draw(self.canvas, UpperInstructions, 6, 6)
	} else {
		self.info.Draw(self.canvas, LowerInstructions, 6, 6)
	}

	// draw text
	alignStr := self.text.GetAlign().String()
	if self.uppercaseOnly { alignStr = fmtAlignString(alignStr) }
	self.text.Draw(self.canvas, alignStr, self.cx, self.cy)

	// aux info warnings
	self.info.SetColor(color.RGBA{140, 70, 40, 255})
	self.info.SetAlign(ptxt.Baseline | ptxt.Left)
	vertAlign := self.text.GetAlign().Vert()
	switch vertAlign {
	case ptxt.Midline:
		if self.text.Strand().Font().Metrics().MidlineAscent() == 0 {
			self.info.Draw(self.canvas, "WARNING: FONT MIDLINE ASCENT IS ZERO", 6, CanvasHeight - 6)
		}
	case ptxt.CapLine:
		if self.text.Strand().Font().Metrics().UppercaseAscent() == 0 {
			self.info.Draw(self.canvas, "WARNING: FONT UPPERCASE ASCENT IS ZERO", 6, CanvasHeight - 6)
		}
	}

	// project logical canvas to main (optional ptxt utility)
	ptxt.Proportional.Project(self.canvas, hiResCanvas)
}

// Helper method to make the align strings more compatible with
// pixel art fonts that often only include uppercase letters.
func fmtAlignString(alignString string) string {
	var bytes []byte = make([]byte, 0, len(alignString))
	var fromLower bool = false
	for _, codePoint := range alignString {
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
