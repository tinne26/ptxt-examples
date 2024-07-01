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
	renderer.SetScale(4)
	renderer.SetColor(color.RGBA{255, 116, 119, 255}) // light red

	// create canvas
	const CanvasWidth, CanvasHeight = 360, 180
	canvas := image.NewRGBA(image.Rect(0, 0, CanvasWidth, CanvasHeight))
	fill(canvas, color.RGBA{23, 18, 25, 255}) // licorice

	// actual drawing
	renderer.Draw(canvas, "GETTING STARTED", CanvasWidth/2, CanvasHeight/2)

	// export result as png
	filename, err := filepath.Abs("ptxt_examples_cpu_getstarted.png")
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
