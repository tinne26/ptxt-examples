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

// Usage:
// > go run -tags cputext main.go myfont.ggfnt

var BackColor = color.RGBA{203, 243, 240, 255} // light mint green
var TextColor = color.RGBA{ 46, 196, 182, 255} // sea green

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

	// make sure notdef exists
	notdef := strand.Font().Glyphs().FindIndexByName("notdef")
	if notdef == ggfnt.GlyphMissing {
		fmt.Printf("Font doesn't have a 'notdef' glyph.\n")
		os.Exit(0)
	}

	// create text renderer, set the main properties
	const Scale = 4
	renderer := ptxt.NewRenderer()
	renderer.SetStrand(strand)
	renderer.SetAlign(ptxt.Center)
	renderer.SetScale(Scale)
	renderer.SetColor(TextColor)

	// create canvas
	notdefMask := renderer.Advanced().LoadMask(notdef)
	canvasBounds := notdefMask.Bounds()
	canvasBounds.Min.X = canvasBounds.Min.X*Scale - Scale
	canvasBounds.Max.X = canvasBounds.Max.X*Scale + Scale
	canvasBounds.Min.Y = canvasBounds.Min.Y*Scale - Scale
	canvasBounds.Max.Y = canvasBounds.Max.Y*Scale + Scale
	canvas := image.NewRGBA(canvasBounds)
	fill(canvas, BackColor)

	// actual drawing
	var params ptxt.MaskDrawParameters
	params.X = 0
	params.Y = 0
	params.Scale = Scale
	params.RGBA = [4]float32{
		float32(TextColor.R)/255.0, float32(TextColor.G)/255.0,
		float32(TextColor.B)/255.0, float32(TextColor.A)/255.0,
	}
	renderer.Advanced().DrawMask(canvas, notdefMask, strand, params)

	// export result as png
	filename, err := filepath.Abs("ptxt_examples_cpu_notdef.png")
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
