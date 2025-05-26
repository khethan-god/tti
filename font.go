// font loading and text positioning calculations
package main

import (
	"errors"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func loadFontFace(fontStyle string, size float64) (font.Face, error) {
	fontPath := getFontPath(fontStyle)

	data, err := fontAssets.ReadFile(fontPath)
	if err != nil {
		return nil, errors.New("unable to read font file: " + err.Error())
	}

	ft, err := opentype.Parse(data)
	if err != nil {
		return nil, errors.New("failed to parse font: " + err.Error())
	}

	face, err := opentype.NewFace(ft, &opentype.FaceOptions{
		Size:    size,
		DPI:     backgroundDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, errors.New("failed to create font face: " + err.Error())
	}

	return face, nil
}

// Enhanced text measurement and wrapping
func measureText(text string, face font.Face) (int, int) {
	d := &font.Drawer{Face: face}
	bounds, _ := d.BoundString(text)
	width := (bounds.Max.X - bounds.Min.X).Ceil()
	height := (bounds.Max.Y - bounds.Min.Y).Ceil()
	return width, height
}

func wrapText(text string, maxWidth int, face font.Face) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{text} // Return original if no spaces
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		testLine := currentLine.String()
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		width, _ := measureText(testLine, face)

		if width <= maxWidth || currentLine.Len() == 0 {
			if currentLine.Len() > 0 {
				currentLine.WriteString(" ")
			}
			currentLine.WriteString(word)
		} else {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentLine.WriteString(word)
		}
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

// Calculate optimal font size that fits both width and height constraints
func calculateOptimalFontSize(text string, fontStyle string, maxWidth, maxHeight int, startingSize float64) (font.Face, float64, []string, error) {
	for fontSize := startingSize; fontSize >= 8.0; fontSize -= 2.0 {
		face, err := loadFontFace(fontStyle, fontSize)
		if err != nil {
			continue
		}

		// Try single line first
		width, height := measureText(text, face)
		if width <= (maxWidth*9)/10 && height <= (maxHeight*9)/10 {
			return face, fontSize, []string{text}, nil
		}

		// Try text wrapping
		lines := wrapText(text, int(float64(maxWidth)*0.9), face)

		metrics := face.Metrics()
		lineHeight := (metrics.Ascent + metrics.Descent).Ceil()
		totalHeight := len(lines) * int(float64(lineHeight)*1.2)

		if totalHeight <= int(float64(maxHeight)*0.9) {
			// Check if all lines fit width-wise
			allLinesFit := true
			for _, line := range lines {
				lineWidth, _ := measureText(line, face)
				if lineWidth > int(float64(maxWidth)*0.9) {
					allLinesFit = false
					break
				}
			}

			if allLinesFit {
				return face, fontSize, lines, nil
			}
		}

		face.Close()
	}

	// Fallback to minimum size
	face, err := loadFontFace(fontStyle, 8.0)
	if err != nil {
		return nil, 0, nil, err
	}

	lines := wrapText(text, int(float64(maxWidth)*0.9), face)
	return face, 8.0, lines, nil
}
