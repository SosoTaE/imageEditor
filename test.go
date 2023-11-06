package main

import (
	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"log"
)

func main() {
	// Load an image from file (assuming it's in the same directory and named "input.jpg")
	inputImg, err := imgio.Open("0.jpg")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	// Resize the image to a manageable size (optional)
	resizedImg := transform.Resize(inputImg, 800, 800, transform.Linear)

	// Apply a Gaussian blur to the image
	blurredImg := blur.Gaussian(resizedImg, 3.0)  // The sigma parameter controls the strength of the blur

	// Save the blurred image to file
	if err := imgio.Save("output.jpg", blurredImg, imgio.JPEGEncoder(95)); err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
}
