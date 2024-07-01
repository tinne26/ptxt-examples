package main

import "os"
import "fmt"
import "log"
import "image"
import "image/png"
import "image/color"
import "path/filepath"

import "github.com/tinne26/ptxt"
import "github.com/tinne26/ggfnt"
import "github.com/tinne26/ggfnt-fonts/jammy"

// Usage:
// > go run -tags cputext main.go

const Text = "WE <3 PIXELS"
const Scale = 4

func main() {
	// parse font and create strand
	font := jammy.Font() // we use this specific font because it contains a rewrite rule for <3 to ❤
	strand, err := ptxt.NewStrand(font)
	if err != nil { log.Fatal(err) }
	fmt.Printf("Font loaded: %s\n", font.Header().Name())

	// create text renderer, set the main properties
	renderer := ptxt.NewRenderer()
	renderer.SetStrand(strand)
	renderer.SetAlign(ptxt.Center)
	renderer.SetScale(Scale)
	renderer.SetColor(color.RGBA{255, 116, 119, 255}) // light red

	// set up rewrite rules
	err = strand.Mapping().AutoInitRewriteRules()
	if err != nil { panicDebugRule(err) }
	
	// create canvas
	renderer.Advanced().SetBoundingMode(ptxt.MaskBounding)
	w, h := renderer.Measure(Text)
	canvas := image.NewRGBA(image.Rect(0, 0, w + Scale*2, h + Scale*2))
	fill(canvas, color.RGBA{23, 18, 25, 255}) // licorice

	// actual drawing
	renderer.Advanced().DrawFromBuffer(canvas, canvas.Bounds().Dx()/2, canvas.Bounds().Dy()/2)

	// export result as png
	filename, err := filepath.Abs("ptxt_examples_cpu_rewrite.png")
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

func panicDebugRule(err error) {
	errWithRule, hasRule := err.(interface { Rule() ggfnt.Utf8RewriteRule })
	if !hasRule { panic(err) }
	rule := errWithRule.Rule()
	panic(err.Error() + "\n" + rule.String())
}
