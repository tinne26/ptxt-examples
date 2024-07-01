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

const SampleText = "SIDEWAYS"
var TextColor = color.RGBA{255, 255, 255, 255}
var BackColor = color.RGBA{  0,   0,   0, 255}

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
	renderer.SetAlign(ptxt.Baseline | ptxt.Right)
	renderer.SetScale(4)
	renderer.SetColor(TextColor)

	// create canvas
	px := int(renderer.GetScale())
	w, _ := renderer.Measure(SampleText)
	h := int(renderer.Strand().Font().Metrics().UppercaseAscent())*px
	pad := h/2
	side := w + h + pad*2 + px
	canvas := image.NewRGBA(image.Rect(0, 0, side, side))
	fill(canvas, BackColor)

	// actual drawing
	renderer.SetDirection(ptxt.Horizontal)
	renderer.Draw(canvas, SampleText, side - pad, pad + h)
	renderer.SetDirection(ptxt.SidewaysRight)
	renderer.Draw(canvas, SampleText, side - pad - h, side - pad)
	renderer.SetDirection(ptxt.Horizontal)
	renderer.Draw(canvas, SampleText, side - pad - h - px, side - pad)
	renderer.SetDirection(ptxt.Sideways)
	renderer.Draw(canvas, SampleText, pad + h, pad)

	// export result as png
	filename, err := filepath.Abs("ptxt_examples_cpu_sideways.png")
	if err != nil { log.Fatal(err) }
	fmt.Printf("Output image: %s\n", filename)
	file, err := os.Create(filename)
	if err != nil { log.Fatal(err) }
	err = png.Encode(file, canvas)
	if err != nil { log.Fatal(err) }
	err = file.Close()
	if err != nil { log.Fatal(err) }
	fmt.Print("Program exited successfully.\n")
}

func fill(canvas *image.RGBA, rgba color.RGBA) {
	for i := 0; i < len(canvas.Pix); i += 4 {
		canvas.Pix[i + 0] = rgba.R
		canvas.Pix[i + 1] = rgba.G
		canvas.Pix[i + 2] = rgba.B
		canvas.Pix[i + 3] = rgba.A
	}
}
