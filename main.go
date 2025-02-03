// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Mandelbrot emits a PNG image of the Mandelbrot fractal.
package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math/cmplx"
	"net/http"
)

func main() {
	log.Print("Starting server...")

	http.HandleFunc("/", handler)

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		log.Print("Returning favicon...")
		http.ServeFile(w, r, "static/favicon.ico")
	})

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, _ *http.Request) {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 1024, 1024
	)

	log.Print("Handling request...")

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			// Image point (px, py) represents complex value z.
			img.Set(px, py, mandelbrot(z))
		}
	}

	if err := png.Encode(w, img); err != nil {
		log.Fatal(err)
	}

	log.Print("Request complete.")
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.RGBA{R: (255 - contrast*n*3) % 255, G: 32 + (contrast * n * 2 % 200), B: (200 - contrast*n) % 255, A: 255}
		}
	}
	return color.Black
}
