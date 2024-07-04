package main

import "image/color"
import "github.com/hajimehoshi/ebiten/v2"
import "github.com/hajimehoshi/ebiten/v2/inpututil"
import "github.com/tinne26/ptxt"
import "github.com/tinne26/ggfnt-fonts/jammy"

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

func (self *Game) Update() error {
	shift := ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeyShiftLeft)
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		strand := self.text.Strand()
		numOpts := strand.Font().Settings().GetNumOptions(jammy.NumericStyleSettingKey)
		value   := strand.GetSetting(jammy.NumericStyleSettingKey)
		if shift { // previous numeric style
			strand.SetSetting(jammy.NumericStyleSettingKey, (value + (numOpts - 1)) % numOpts)
		} else { // next numeric style
			strand.SetSetting(jammy.NumericStyleSettingKey, (value + 1) % numOpts)
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		// switch zero disambiguation style
		strand := self.text.Strand()
		numOpts := strand.Font().Settings().GetNumOptions(jammy.ZeroDisambiguationMarkSettingKey)
		value   := strand.GetSetting(jammy.ZeroDisambiguationMarkSettingKey)
		strand.SetSetting(jammy.ZeroDisambiguationMarkSettingKey, (value + 1) % numOpts)
	}
	return nil
}
func (self *Game) Draw(hiResCanvas *ebiten.Image) {
	// fill background
	self.canvas.Fill(color.RGBA{246, 242, 240, 255})

	// compose text
	strand := self.text.Strand()
	font   := strand.Font()
	info := "0123456789\n\n"
	settingValue := strand.GetSetting(jammy.ZeroDisambiguationMarkSettingKey)
	settingState := font.Settings().GetOptionName(jammy.ZeroDisambiguationMarkSettingKey, settingValue)
	info += jammy.ZeroDisambiguationMarkSettingName + ": " + settingState + "\n"
	settingValue = strand.GetSetting(jammy.NumericStyleSettingKey)
	settingState = font.Settings().GetOptionName(jammy.NumericStyleSettingKey, settingValue)
	info += jammy.NumericStyleSettingName + ": " + settingState

	// draw text
	self.text.SetAlign(ptxt.LastBaseline | ptxt.Left)
	self.text.Draw(self.canvas, "[Z] Switch zero mark\n[N] Switch numeric style", 3, CanvasHeight - 3)
	self.text.SetAlign(ptxt.Center)
	self.text.Draw(self.canvas, info, CanvasWidth/2, CanvasHeight/2)

	// project logical canvas to main (optional ptxt utility)
	ptxt.PixelPerfect.Project(self.canvas, hiResCanvas)
}

// ---- main function ----

func main() {
	// initialize font strand
	strand, err := ptxt.NewStrand(jammy.Font())
	if err != nil { panic(err) }
	err = strand.Mapping().AutoInitRewriteRules()
	if err != nil { panic(err) }
	
	// create text renderer, set the main properties
	renderer := ptxt.NewRenderer()
	renderer.SetStrand(strand)
	renderer.SetAlign(ptxt.Center)
	renderer.SetColor(color.RGBA{242, 143, 59, 255})

	// set up Ebitengine and start the game
	ebiten.SetWindowTitle("ptxt-examples/gpu/settingmap")
	canvas := ebiten.NewImage(CanvasWidth, CanvasHeight)
	err = ebiten.RunGame(&Game{ text: renderer, canvas: canvas })
	if err != nil { panic(err) }
}
