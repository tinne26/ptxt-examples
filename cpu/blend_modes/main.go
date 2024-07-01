package main

import "os"
import "fmt"
import "log"
import "image"
import "image/png"
import "image/color"
import "path/filepath"

import "github.com/tinne26/ptxt"

// Usage:
// > go run -tags cputext main.go myfont.ggfnt

const Alpha = 255 // can be changed (e.g. 144) if you want to see how
                  // color modes work with semi-transparency too

func main() {
	// get font path
	if len(os.Args) != 2 {
		msg := "Usage: expects one argument with the path to the font to be used\n"
		fmt.Fprint(os.Stderr, msg)
		os.Exit(1)
	}

	// parse font and create strand
	strand, err := ptxt.NewStrand(os.Args[1])
	if err != nil { log.Fatal(err) }
	fmt.Printf("Font loaded: %s\n", strand.Font().Header().Name())

	// create text renderer, set the main properties
	renderer := ptxt.NewRenderer()
	renderer.SetStrand(strand)
	renderer.SetAlign(ptxt.Center)
	renderer.SetScale(5)

	// determine canvas size, create canvas, fill with three colors
	barHeight := (strand.Font().Metrics().LineHeight() + 4)*int(renderer.GetScale())
	width, height := barHeight*12, barHeight*4
	wpad := 16
	target := image.NewRGBA(image.Rect(0, 0, width + wpad*2, height))
	fillRows(target, 0*height/3, 1*height/3, color.RGBA{0, 255, 255, 255})
	fillRows(target, 1*height/3, 2*height/3, color.RGBA{255, 0, 255, 255})
	fillRows(target, 2*height/3, 3*height/3, color.RGBA{255, 255, 0, 255})

	// actual drawing
	// draw first row of blend modes
	renderer.SetColor(color.RGBA{0, 0, 0, Alpha})
	renderer.SetBlendMode(ptxt.BlendOver)
	renderer.Draw(target, "OVER", wpad + 1*width/8, 1*height/6)
	renderer.SetBlendMode(ptxt.BlendCut)
	renderer.Draw(target, "CUT", wpad + 3*width/8, 1*height/6)
	renderer.SetBlendMode(ptxt.BlendHue)
	renderer.Draw(target, "HUE", wpad + 5*width/8, 1*height/6)
	renderer.SetBlendMode(ptxt.BlendReplace)
	renderer.Draw(target, "REPLACE", wpad + 7*width/8, 1*height/6)
	
	// draw second row of blend modes
	renderer.SetColor(color.RGBA{0, Alpha, Alpha, Alpha})
	renderer.SetBlendMode(ptxt.BlendSub)
	renderer.Draw(target, "SUBTRACT", wpad + 1*width/8, 3*height/6)
	renderer.SetBlendMode(ptxt.BlendAdd)
	renderer.Draw(target, "ADD", wpad + 3*width/8, 3*height/6)
	renderer.SetBlendMode(ptxt.BlendOver)
	renderer.Draw(target, "OVER", wpad + 5*width/8, 3*height/6)
	renderer.SetBlendMode(ptxt.BlendMultiply)
	renderer.Draw(target, "MULTIPLY", wpad + 7*width/8, 3*height/6)

	// draw third row of blend modes
	renderer.SetColor(color.RGBA{Alpha, 0, 0, Alpha})
	renderer.SetBlendMode(ptxt.BlendOver)
	renderer.Draw(target, "OVER", wpad + 1*width/8, 5*height/6)
	renderer.SetBlendMode(ptxt.BlendMultiply)
	renderer.Draw(target, "MULTIPLY", wpad + 3*width/8, 5*height/6)
	renderer.SetBlendMode(ptxt.BlendHue)
	renderer.Draw(target, "HUE", wpad + 5*width/8, 5*height/6)
	renderer.SetBlendMode(ptxt.BlendSub)
	renderer.Draw(target, "SUBTRACT", wpad + 7*width/8, 5*height/6)

	// export result as png
	filename, err := filepath.Abs("ptxt_examples_cpu_blend_modes.png")
	if err != nil { log.Fatal(err) }
	fmt.Printf("Output image: %s\n", filename)
	file, err := os.Create(filename)
	if err != nil { log.Fatal(err) }
	err = png.Encode(file, target)
	if err != nil { log.Fatal(err) }
	err = file.Close()
	if err != nil { log.Fatal(err) }
	fmt.Print("Program exited successfully.\n")
}

func fillRows(canvas *image.RGBA, fromRow, toRow int, rgba color.RGBA) {
	rowLength := canvas.Bounds().Dx()
	for i := fromRow*rowLength*4; i < toRow*rowLength*4; i += 4 {
		canvas.Pix[i + 0] = rgba.R
		canvas.Pix[i + 1] = rgba.G
		canvas.Pix[i + 2] = rgba.B
		canvas.Pix[i + 3] = rgba.A
	}
}
