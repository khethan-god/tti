// background pattern generator algorithms
package main

import (
	// "fmt"
	"image"
	"image/color"
	"math"
)

// defines the signature for background generation functions
type BackgroundGenFunc func(w, h int) *image.RGBA

// Enhanced Studio Ghibli color palette
var studioPalette = []color.RGBA{
	// Warm forest tones
	{34, 89, 34, 255},   // Deep Forest Green
	{85, 139, 47, 255},  // Olive Green
	{154, 205, 50, 255}, // Yellow Green
	{173, 255, 47, 255}, // Green Yellow

	// Sky and water tones
	{70, 130, 180, 255},  // Steel Blue
	{135, 206, 235, 255}, // Sky Blue
	{176, 224, 230, 255}, // Powder Blue
	{240, 248, 255, 255}, // Alice Blue

	// Earth and warmth
	{160, 82, 45, 255},   // Saddle Brown
	{205, 133, 63, 255},  // Peru
	{222, 184, 135, 255}, // Burlywood
	{245, 222, 179, 255}, // Wheat

	// Magical accents
	{255, 182, 193, 255}, // Light Pink
	{221, 160, 221, 255}, // Plum
	{230, 230, 250, 255}, // Lavender
	{255, 228, 181, 255}, // Moccasin
}

// var studioPalette_ = []color.RGBA{
// 	{34, 139, 34, 255},   // Forest Green
// 	{70, 130, 180, 255},  // Steel Blue
// 	{205, 133, 63, 255},  // Peru Brown
// 	{255, 218, 185, 255}, // Peach Puff
// 	{176, 196, 222, 255}, // Light Steel Blue
// 	{255, 182, 193, 255}, // Light Pink
// 	{152, 251, 152, 255}, // Pale Green
// 	{255, 239, 213, 255}, // Papaya Whip
// 	{175, 238, 238, 255}, // Pale Turquoise
// 	{255, 228, 196, 255}, // Bisque
// 	{230, 230, 250, 255}, // Lavender
// 	{240, 248, 255, 255}, // Alice Blue
// 	{144, 238, 144, 255}, // Light Green
// 	{255, 228, 181, 255}, // Moccasin
// 	{221, 160, 221, 255}, // Plum
// 	{173, 216, 230, 255}, // Light Blue
// }

var predefinedPalette color.Palette

func init() {
	initializePalette()
}

func initializePalette() {
	// Add black and white first
	predefinedPalette = append(predefinedPalette, color.Black, color.White)

	// Add Studio Ghibli colors
	for _, c := range studioPalette {
		predefinedPalette = append(predefinedPalette, c)
		if len(predefinedPalette) >= maxPaletteSize {
			break
		}
	}

	// fmt.Printf("Initialized GIF palette with %d colors.", len(predefinedPalette))
}

// create XOR pattern
func generatePatternBackground(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			v := uint8(x ^ y + (x+y)/2)
			img.Set(x, y, color.RGBA{v, 255 - v, (v * 3) % 255, 255})
		}
	}
	return img
}

// Add this helper function before generatePerlinLikeBackground
func interpolateColor(c1, c2 color.RGBA, t float64) color.RGBA {
	return color.RGBA{
		uint8(float64(c1.R)*(1-t) + float64(c2.R)*t),
		uint8(float64(c1.G)*(1-t) + float64(c2.G)*t),
		uint8(float64(c1.B)*(1-t) + float64(c2.B)*t),
		255,
	}
}

// create pseudo-perlin noise background
func generatePerlinLikeBackground(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := range h {
		for x := range w {
			val1 := math.Sin(float64(x)/50.0) + math.Cos(float64(y)/40.0)
			val2 := math.Sin(float64(x)/25.0) + math.Cos(float64(y)/20.0)*0.5
			val3 := math.Sin(float64(x)/12.5) + math.Cos(float64(y)/10.0)*0.25
			val4 := math.Sin(float64(x)/80.0) + math.Cos(float64(y)/60.0)*1.5

			combinedNoise := val1 + val2 + val3 + val4
			normalizedNoise := math.Max(0, math.Min(1, (combinedNoise+4.0)/8.0))

			paletteIndex := normalizedNoise * float64(len(studioPalette)-1)
			index := int(paletteIndex)
			fraction := paletteIndex - float64(index)

			var finalColor color.RGBA
			if index >= len(studioPalette)-1 {
				finalColor = studioPalette[len(studioPalette)-1]
			} else {
				finalColor = interpolateColor(studioPalette[index], studioPalette[index+1], fraction)
			}

			img.Set(x, y, finalColor)
		}
	}
	return img
}

// Significantly smoother Perlin-like noise using higher precision
func generatePerlinSmootherBackground(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := range h {
		for x := range w {
			fx, fy := float64(x)/float64(w), float64(y)/float64(h)

			// Multi-octave noise with smoother interpolation
			noise := 0.0
			amplitude := 1.0
			frequency := 1.0

			// 6 octaves for ultra-smooth result
			for range 6 {
				n1 := math.Sin(fx*frequency*math.Pi*4) * math.Cos(fy*frequency*math.Pi*3)
				n2 := math.Cos(fx*frequency*math.Pi*3) * math.Sin(fy*frequency*math.Pi*4)
				n3 := math.Sin((fx + fy) * frequency * math.Pi * 2)
				n4 := math.Cos((fx - fy) * frequency * math.Pi * 2.5)

				octaveNoise := (n1 + n2 + n3 + n4) / 4.0
				noise += amplitude * octaveNoise

				amplitude *= 0.6 // Gentler amplitude decrease
				frequency *= 1.8 // Gentler frequency increase
			}

			// Smooth normalization
			normalized := (math.Tanh(noise) + 1.0) / 2.0 // Tanh for smoother distribution

			// Map to enhanced palette
			colorIndex := normalized * float64(len(studioPalette)-1)
			idx := int(colorIndex)
			frac := colorIndex - float64(idx)

			if idx >= len(studioPalette)-1 {
				img.Set(x, y, studioPalette[len(studioPalette)-1])
			} else {
				// Smooth color interpolation
				interpolated := interpolateColor(studioPalette[idx], studioPalette[idx+1], frac)
				img.Set(x, y, interpolated)
			}
		}
	}
	return img
}

// creates a radial pattern
func generateRadialPatternBackground(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	centerX, centerY := float64(w)/2.0, float64(h)/2.0

	for y := range h {
		for x := range w {
			dx := float64(x) - centerX
			dy := float64(y) - centerY

			dist := math.Sqrt(dx*dx + dy*dy)
			angle := math.Atan2(dy, dx)

			r := uint8(math.Abs(math.Sin(dist/20.0+angle*5.0) * 255.0))
			g := uint8(math.Abs(math.Cos(dist/30.0-angle*3.0) * 255.0))
			b := uint8(math.Abs(math.Sin(dist/40.0+angle*7.0) * 255.0))

			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	return img
}

// creates diagonal grid pattern
func generateDiagonalGridBackground(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	gridSize := 50

	for y := range h {
		for x := range w {
			patternVal := (x + y) % gridSize

			colorIdx := ((x + y) / gridSize) % len(studioPalette)
			baseColor := studioPalette[colorIdx]

			if patternVal < gridSize/2 {
				img.Set(x, y, baseColor)
			} else {
				// Slightly darker version
				darkColor := color.RGBA{
					uint8(float64(baseColor.R) * 0.7),
					uint8(float64(baseColor.G) * 0.7),
					uint8(float64(baseColor.B) * 0.7),
					255,
				}
				img.Set(x, y, darkColor)
			}
		}
	}
	return img
}
