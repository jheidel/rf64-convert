package main

import (
	"flag"
	"log"

	"rf64-convert/convert"
)

var (
	input  = flag.String("input", "", "Input path for RIFF WAV")
	output = flag.String("output", "", "Output path for RF64 WAV")
)

// RF64 Specification
// https://tech.ebu.ch/docs/tech/tech3306v1_0.pdf

func main() {
	flag.Parse()
	in, err := convert.OpenInput(*input)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer in.Close()

	out, err := convert.NewOutputFile(*output)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer out.Close()

	if err := out.CopyFrom(in); err != nil {
		log.Fatalf("%v", err)
	}
}
