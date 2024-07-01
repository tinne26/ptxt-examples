package main

import "os"
import "fmt"

import "github.com/tinne26/ggfnt"

func main() {
	// usage check
	if len(os.Args) != 2 {
		fmt.Print("Usage: go run main.go font.ggfnt\n")
		os.Exit(1)
	}

	// parse font
	file, err := os.Open(os.Args[1])
	if err != nil { panic(err) }
	defer file.Close()
	font, err := ggfnt.Parse(file)
	if err != nil { panic(err) }
	
	// print font info
	fmt.Print("Header\n")
	fmt.Printf("  Font name : %s\n", font.Header().Name())
	fmt.Printf("  Author    : %s\n", font.Header().Author())
	fmt.Print("Metrics\n")
	fmt.Printf("  Ascent  : %d (+%d)\n", font.Metrics().Ascent(), font.Metrics().ExtraAscent())
	fmt.Printf("  CapLine : %d\n", font.Metrics().UppercaseAscent())
	fmt.Printf("  Midline : %d\n", font.Metrics().MidlineAscent())
	fmt.Printf("  Descent : %d (+%d)\n", font.Metrics().Descent(), font.Metrics().ExtraDescent())
}
