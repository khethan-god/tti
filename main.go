// CLI parsing and program orchestration

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	config := parseFlags()

	if flag.NArg() < 1 {
		printUsage()
		os.Exit(1)
	}

	text := flag.Arg(0)

	// validate configuration
	if err := validateConfig(config); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		os.Exit(1)
	}

	// generate output
	if config.Animate {
		if err := generateAnimatedGIF(text, config); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate GIF: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := generateStaticImage(text, config); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate image: %v\n", err)
			os.Exit(1)
		}
	}
}

func parseFlags() Config {
	config := Config{
		Width:      600, //default values
		Height:     300,
		FontSize:   48,
		OutputDir:  "images",
		Background: "default",
		FontStyle:  "roboto_bold",
	}

	flag.IntVar(&config.Width, "width", config.Width, "Image width in pixels")
	flag.IntVar(&config.Height, "height", config.Height, "Image height in pixels")
	flag.Float64Var(&config.FontSize, "font-size", config.FontSize, "Font size in points")
	flag.StringVar(&config.OutputDir, "output", config.OutputDir, "Output directory")
	flag.StringVar(&config.Background, "bg", config.Background, "Background pattern: "+strings.Join(getBackgroundTypes(), ", "))
	flag.StringVar(&config.FontStyle, "font", config.FontStyle, "Font style: "+strings.Join(getFontStyles(), ", "))
	flag.BoolVar(&config.RevealBg, "reveal-bg", false, "Display background via Text")
	flag.BoolVar(&config.Animate, "animate", false, "Create animated GIF")

	flag.Usage = printUsage
	flag.Parse()

	return config
}

func printUsage() {
	fmt.Println("Text Image generator")
	fmt.Println("Usage: cli_tool [options] \"YourTextHere\"")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println("  cli_tool \"Hello World\"")
	fmt.Println("  cli_tool -width=800 -height=400 -font=roboto_bold \"Custom Text\"")
	fmt.Println("  cli_tool -animate -bg=perlin \"Animated Text\"")
}
