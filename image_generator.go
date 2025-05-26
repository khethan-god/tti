// image_generator.go
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/font"
)

func generateStaticImage(text string, config Config) error {
	bgGen := getBackgroundGenerator(config.Background)
	img := bgGen(config.Width, config.Height)

	// Calculate optimal font size and get wrapped lines
	primaryFace, _, lines, err := calculateOptimalFontSize(
		text, config.FontStyle, config.Width, config.Height, config.FontSize)
	if err != nil {
		return err
	}
	defer primaryFace.Close()

	renderer := NewTextRenderer(primaryFace)

	// Calculate positioning for multi-line text
	metrics := primaryFace.Metrics()
	lineHeight := int(float64((metrics.Ascent + metrics.Descent).Ceil()) * 1.2)
	totalHeight := len(lines) * lineHeight

	startY := max((config.Height-totalHeight)/2+metrics.Ascent.Ceil(), metrics.Ascent.Ceil())

	// Center each line horizontally
	if config.RevealBg {
		renderer.RenderRevealBackground(img, lines, -1, startY, lineHeight, false)
	} else {
		// Center each line individually
		for i, line := range lines {
			width, _ := measureText(line, primaryFace)
			lineX := max((config.Width-width)/2, 0)
			y := startY + i*lineHeight

			renderer.renderWithOutline(img, line, lineX, y, color.Black)
			renderer.renderText(img, line, lineX, y, color.White)
		}
	}

	return saveImage(img, text, config.OutputDir)
}

// Fixed generateAnimatedGIF function
func generateAnimatedGIF(text string, config Config) error {
	bgGen := getBackgroundGenerator(config.Background)

	// calculate optimal font size and get wrapped lines
	primaryFace, _, lines, err := calculateOptimalFontSize(
		text, config.FontStyle, config.Width, config.Height, config.FontSize)
	if err != nil {
		return err
	}
	defer primaryFace.Close()

	renderer := NewTextRenderer(primaryFace)

	// Calculate positioning for multi-line text (same as static version)
	metrics := primaryFace.Metrics()
	lineHeight := int(float64((metrics.Ascent + metrics.Descent).Ceil()) * 1.2)
	totalHeight := len(lines) * lineHeight

	startY := max((config.Height-totalHeight)/2+metrics.Ascent.Ceil(), metrics.Ascent.Ceil())

	// Generate four frames with different effects
	frames := []*image.RGBA{
		createFrame(bgGen, config, renderer, lines, primaryFace, startY, lineHeight, "outline"),
		createFrame(bgGen, config, renderer, lines, primaryFace, startY, lineHeight, "reveal"),
		createFrame(bgGen, config, renderer, lines, primaryFace, startY, lineHeight, "plain"),
		createFrame(bgGen, config, renderer, lines, primaryFace, startY, lineHeight, "reveal-outline"),
	}

	// Create and save GIF
	return saveAnimatedGIF(frames, text, config.OutputDir)
}

// Fixed createFrame function
func createFrame(bgGen BackgroundGenFunc, config Config, renderer *TextRenderer, lines []string, primaryFace font.Face, startY, lineHeight int, effect string) *image.RGBA {
	img := bgGen(config.Width, config.Height)

	switch effect {
	case "outline":
		// Center each line individually (like in static version)
		for i, line := range lines {
			width, _ := measureText(line, primaryFace)
			lineX := max((config.Width-width)/2, 0)
			y := startY + i*lineHeight
			renderer.renderWithOutline(img, line, lineX, y, color.Black)
			renderer.renderText(img, line, lineX, y, color.White)
		}
	case "reveal":
		// For reveal effect, center the text properly
		renderer.RenderRevealBackgroundCentered(img, lines, primaryFace, config.Width, startY, lineHeight, false)
	case "plain":
		// Center each line individually
		for i, line := range lines {
			width, _ := measureText(line, primaryFace)
			lineX := max((config.Width-width)/2, 0)
			y := startY + i*lineHeight
			renderer.renderText(img, line, lineX, y, color.White)
		}
	case "reveal-outline":
		// For reveal effect with outline, center the text properly
		renderer.RenderRevealBackgroundCentered(img, lines, primaryFace, config.Width, startY, lineHeight, true)
	}

	return img
}

func saveImage(img *image.RGBA, text, outputDir string) error {
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}

	fileName := filepath.Join(outputDir, sanitizeFilename(text)+".png")
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return err
	}

	fmt.Printf("✅ Image successfully created: %s\n", fileName)
	return nil
}

func saveAnimatedGIF(frames []*image.RGBA, text, outputDir string) error {
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}

	outGif := &gif.GIF{}

	// Add frames to GIF in cycles
	totalFrames := len(frames) * gifNumCylces
	for i := range totalFrames {
		frameIndex := i % len(frames)
		imgRGBA := frames[frameIndex]

		// Convert to paletted image
		palettedImage := image.NewPaletted(imgRGBA.Bounds(), predefinedPalette)
		draw.Draw(palettedImage, palettedImage.Rect, imgRGBA, image.Point{}, draw.Src)

		outGif.Image = append(outGif.Image, palettedImage)
		outGif.Delay = append(outGif.Delay, gifFrameDelay)
	}

	fileName := filepath.Join(outputDir, sanitizeFilename(text)+".gif")
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := gif.EncodeAll(f, outGif); err != nil {
		return err
	}

	fmt.Printf("✅ GIF animation successfully created: %s\n", fileName)
	return nil
}
