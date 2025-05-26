// text rendering effects and styles
package main

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// TextRenderer handles different text rendering styles
type TextRenderer struct {
	face font.Face
}

func NewTextRenderer(face font.Face) *TextRenderer {
	return &TextRenderer{face: face}
}

func (tr *TextRenderer) applyTextMask(originalImg, outputImg, mask *image.RGBA) {
	bounds := originalImg.Bounds()
	for py := bounds.Min.Y; py < bounds.Max.Y; py++ {
		for px := bounds.Min.X; px < bounds.Max.X; px++ {
			maskColor := mask.At(px, py).(color.RGBA)
			if maskColor.A > 0 {
				outputImg.Set(px, py, originalImg.At(px, py))
			}
		}
	}
}

func (tr *TextRenderer) isEmoji(r rune) bool {
	return (r >= 0x1F600 && r <= 0x1F64F) || // Emoticons
		(r >= 0x1F300 && r <= 0x1F5FF) || // Misc Symbols
		(r >= 0x1F680 && r <= 0x1F6FF) || // Transport
		(r >= 0x1F1E0 && r <= 0x1F1FF) || // Flags
		(r >= 0x2600 && r <= 0x26FF) || // Misc symbols
		(r >= 0x2700 && r <= 0x27BF) || // Dingbats
		(r >= 0xFE00 && r <= 0xFE0F) || // Variation Selectors
		(r >= 0x1F900 && r <= 0x1F9FF) || // Supplemental Symbols
		(r >= 0x1F018 && r <= 0x1F270) // Various symbols
}

func (tr *TextRenderer) renderText(img *image.RGBA, text string, x, y int, textColor color.Color) {

	currentX := x
	for _, r := range text {
		var face font.Face
		if tr.isEmoji(r) {
			// TODO: Need to add support to display emojis
			continue
		} else {
			face = tr.face
		}

		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(textColor),
			Face: face,
			Dot:  fixed.P(currentX, y),
		}

		d.DrawString(string(r))
		_, advance := d.BoundString(string(r))
		currentX += advance.Ceil()
	}
}

func (tr *TextRenderer) renderWithOutline(img *image.RGBA, text string, x, y int, outlineColor color.Color) {
	offsets := []image.Point{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1},
		{-1, -1}, {-1, 1}, {1, -1}, {1, 1},
	}

	for _, off := range offsets {
		tr.renderText(img, text, x+off.X, y+off.Y, outlineColor)
	}
}

// Enhanced multi-line rendering
func (tr *TextRenderer) RenderLines(img *image.RGBA, lines []string, startX, startY int, lineSpacing int, textColor color.Color) {
	for i, line := range lines {
		y := startY + i*lineSpacing
		tr.renderText(img, line, startX, y, textColor)
	}
}

func (tr *TextRenderer) RenderLinesWithOutline(img *image.RGBA, lines []string, startX, startY int, lineSpacing int) {
	// Draw outline
	for i, line := range lines {
		y := startY + i*lineSpacing
		tr.renderWithOutline(img, line, startX, y, color.Black)
	}

	// Draw main text
	for i, line := range lines {
		y := startY + i*lineSpacing
		tr.renderText(img, line, startX, y, color.White)
	}
}

// Use startX as a signal for centering
func (tr *TextRenderer) RenderRevealBackground(img *image.RGBA, lines []string, startX, startY int, lineSpacing int, withOutline bool) {
	// If startX is -1, center the text; otherwise use startX as-is
	shouldCenter := startX == -1
	imgWidth := img.Bounds().Dx()

	mask := image.NewRGBA(img.Bounds())
	draw.Draw(mask, mask.Bounds(), image.NewUniform(color.Transparent), image.Point{}, draw.Src)

	// Draw text on mask
	for i, line := range lines {
		var x int
		if shouldCenter {
			width, _ := measureText(line, tr.face)
			x = max((imgWidth-width)/2, 0)
		} else {
			x = startX
		}
		y := startY + i*lineSpacing
		tr.renderText(mask, line, x, y, color.Black)
	}

	outputImg := image.NewRGBA(img.Bounds())
	draw.Draw(outputImg, outputImg.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)

	if withOutline {
		for i, line := range lines {
			var x int
			if shouldCenter {
				width, _ := measureText(line, tr.face)
				x = max((imgWidth-width)/2, 0)
			} else {
				x = startX
			}
			y := startY + i*lineSpacing
			tr.renderWithOutline(outputImg, line, x, y, color.Black)
		}
	}

	tr.applyTextMask(img, outputImg, mask)
	draw.Draw(img, img.Bounds(), outputImg, image.Point{}, draw.Src)
}

// New helper method for centered reveal background
func (tr *TextRenderer) RenderRevealBackgroundCentered(img *image.RGBA, lines []string, primaryFace font.Face, imgWidth, startY, lineHeight int, withOutline bool) {
	mask := image.NewRGBA(img.Bounds())
	draw.Draw(mask, mask.Bounds(), image.NewUniform(color.Transparent), image.Point{}, draw.Src)

	// Draw text on mask with proper centering
	for i, line := range lines {
		width, _ := measureText(line, primaryFace)
		lineX := max((imgWidth-width)/2, 0)
		y := startY + i*lineHeight
		tr.renderText(mask, line, lineX, y, color.Black)
	}

	outputImg := image.NewRGBA(img.Bounds())
	draw.Draw(outputImg, outputImg.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)

	if withOutline {
		for i, line := range lines {
			width, _ := measureText(line, primaryFace)
			lineX := max((imgWidth-width)/2, 0)
			y := startY + i*lineHeight
			tr.renderWithOutline(outputImg, line, lineX, y, color.Black)
		}
	}

	tr.applyTextMask(img, outputImg, mask)
	draw.Draw(img, img.Bounds(), outputImg, image.Point{}, draw.Src)
}
