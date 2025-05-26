// output configuration, validation and mappings

package main

import (
	"embed"
	"errors"
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// go build does not include external files by default so we need to
// embed .ttf files if we want to use them without the need of sharing them

//go:embed assets/text/*.ttf
var fontAssets embed.FS

// config holds all configuration paramets
type Config struct {
	Width      int
	Height     int
	FontSize   float64
	OutputDir  string
	Background string
	FontStyle  string
	RevealBg   bool
	Animate    bool
}

// Font mapping - maps user-friendly names to font files
var fontMap = map[string]string{
	"roboto_black":       "assets/text/Roboto-Black.ttf",
	"roboto_c_black":     "assets/text/Roboto_Condensed-Black.ttf",
	"roboto_sc_ebold":    "assets/text/Roboto_SemiCondensed-ExtraBold.ttf",
	"roboto_bold":        "assets/text/Roboto-Bold.ttf",
	"roboto_c_bold":      "assets/text/Roboto_Condensed-Bold.ttf",
	"roboto_sc_italic":   "assets/text/Roboto_SemiCondensed-Italic.ttf",
	"roboto_ebold":       "assets/text/Roboto-ExtraBold.ttf",
	"roboto_c_ebitalic":  "assets/text/Roboto_Condensed-ExtraBoldItalic.ttf",
	"roboto_sc_light":    "assets/text/Roboto_SemiCondensed-Light.ttf",
	"roboto_elight":      "assets/text/Roboto-ExtraLight.ttf",
	"roboto_c_elight":    "assets/text/Roboto_Condensed-ExtraLight.ttf",
	"roboto_sc_litalic":  "assets/text/Roboto_SemiCondensed-LightItalic.ttf",
	"roboto_italic":      "assets/text/Roboto-Italic.ttf",
	"roboto_c_elitalic":  "assets/text/Roboto_Condensed-ExtraLightItalic.ttf",
	"roboto_sc_medium":   "assets/text/Roboto_SemiCondensed-Medium.ttf",
	"roboto_light":       "assets/text/Roboto-Light.ttf",
	"roboto_c_regular":   "assets/text/Roboto_Condensed-Regular.ttf",
	"roboto_sc_mitalic":  "assets/text/Roboto_SemiCondensed-MediumItalic.ttf",
	"roboto_medium":      "assets/text/Roboto-Medium.ttf",
	"roboto_c_titalic":   "assets/text/Roboto_Condensed-ThinItalic.ttf",
	"roboto_sc_sbold":    "assets/text/Roboto_SemiCondensed-SemiBold.ttf",
	"roboto_regular":     "assets/text/Roboto-Regular.ttf",
	"roboto_sc_blitalic": "assets/text/Roboto_SemiCondensed-BlackItalic.ttf",
	"roboto_sc_sbitalic": "assets/text/Roboto_SemiCondensed-SemiBoldItalic.ttf",
	"roboto_sbold":       "assets/text/Roboto-SemiBold.ttf",
	"roboto_sc_bitalic":  "assets/text/Roboto_SemiCondensed-BoldItalic.ttf",
	"roboto_sc_thin":     "assets/text/Roboto_SemiCondensed-Thin.ttf",
}

// Backgroung generator mapping
var backgroundMap = map[string]BackgroundGenFunc{
	"default":  generatePatternBackground,
	"perlin":   generatePerlinLikeBackground,
	"perlin-s": generatePerlinSmootherBackground,
	"radial":   generateRadialPatternBackground,
	"diagonal": generateDiagonalGridBackground,
}

const (
	backgroundDPI  = 72
	gifFrameDelay  = 15
	gifNumCylces   = 2
	maxPaletteSize = 256
)

func validateConfig(config Config) error {
	if config.Width <= 0 || config.Height <= 0 {
		return errors.New("width and height must be positive")
	}
	if config.FontSize <= 0 {
		return errors.New("font size must be positive")
	}
	if _, exists := fontMap[config.FontStyle]; !exists {
		return errors.New("invalid font style: " + config.FontStyle)
	}
	if _, exists := backgroundMap[config.Background]; !exists {
		return errors.New("invalid background type: " + config.Background)
	}

	return nil
}

// generic function to get sorted keys from a map
func getSortedKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func getFontStyles() []string {
	return getSortedKeys(fontMap)
}

func getBackgroundTypes() []string {
	return getSortedKeys(backgroundMap)
}

// Generic function to get a value from a map with fallback
func getValueOrDefault[T any](key string, m map[string]T, defaultValue T) T {
	val, exists := m[key]
	if !exists {
		fmt.Printf("%s does not exist, reverting to default\n", key)
		return defaultValue // default fallback
	}
	return val
}

func getFontPath(fontStyle string) string {
	return getValueOrDefault(fontStyle, fontMap, fontMap["roboto_bold"])
}

func getBackgroundGenerator(bgType string) BackgroundGenFunc {
	return getValueOrDefault(bgType, backgroundMap, backgroundMap["default"])
}

func sanitizeFilename(text string) string {
	// remove or replace invalid characters (keep emojis, numbers, letters)
	var result strings.Builder // to create strings efficiently
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsNumber(r) ||
			unicode.IsSymbol(r) || unicode.IsPrint(r) {
			if unicode.IsSpace(r) {
				result.WriteRune('_')
			} else if r != '/' && r != '\\' && r != ':' && r != '*' &&
				r != '?' && r != '"' && r != '<' && r != '>' && r != '|' {
				result.WriteRune(r)
			}
		}
		// stop at 25 characters
		if result.Len() > 20 {
			break
		}
	}

	filename := result.String()

	if filename == "" {
		filename = "output"
	}

	return filename
}
