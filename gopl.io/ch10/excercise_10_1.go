// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 287.

//!+main

// The jpeg command reads a PNG image from the standard input
// and writes it as a JPEG image to the standard output.
package main

import (
	"flag"
	"fmt"
	"image" // register GIF decoder
	"image/gif"
	"image/jpeg" // register PNG decoder
	"image/png"
	"io"
	"log"
	"os"
)

var converts = make(map[string]func(io.Reader, io.Writer) error)

func main() {
	var imgType string
	flag.StringVar(&imgType, "itype", "jpeg", "image type")
	flag.Parse()

	if decoder, ok := converts[imgType]; ok {
		if err := decoder(os.Stdin, os.Stdout); err != nil {
			log.Fatalf("convert to %v error: %v", imgType, err)
		}
		os.Exit(0)
	} else {
		log.Fatalf("unsupport image type")
	}
}

func decode(in io.Reader) (image.Image, error) {
	img, kind, err := image.Decode(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input format =", kind)
	}
	return img, err
}

func toJPEG(in io.Reader, out io.Writer) error {
	img, err := decode(in)
	if err != nil {
		return err
	}
	return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
}

func toPNG(in io.Reader, out io.Writer) error {
	img, err := decode(in)
	if err != nil {
		return err
	}
	return png.Encode(out, img)
}

func toGIF(in io.Reader, out io.Writer) error {
	img, err := decode(in)
	if err != nil {
		return err
	}
	return gif.Encode(out, img, &gif.Options{NumColors: 100})
}

func init() {
	converts["png"] = toPNG
	converts["jpeg"] = toJPEG
	converts["gif"] = toGIF
}

//!-main

/*
//!+with
$ go build gopl.io/ch3/mandelbrot
$ go build gopl.io/ch10/jpeg
$ ./mandelbrot | ./jpeg >mandelbrot.jpg
Input format = png
//!-with

//!+without
$ go build gopl.io/ch10/jpeg
$ ./mandelbrot | ./jpeg >mandelbrot.jpg
jpeg: image: unknown format
//!-without
*/
