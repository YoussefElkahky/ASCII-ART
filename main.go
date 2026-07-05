package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	charHeight = 8
	linesPerChar = charHeight + 1
	firstPrintable = 32
	lastPrintable = 126
	newlineToken = `\n`
)

func printAscii(text string, banner []string) error {
	for _, ch := range text {
		if ch < firstPrintable || ch > lastPrintable {
			return fmt.Errorf("character %q is not a printable ASCII character", ch,)
		}
	}
	for row := 0; row < charHeight; row++ {
		for _, ch := range text {
			index := (int(ch)-firstPrintable)*linesPerChar + row
			fmt.Print(banner[index])
		}
		fmt.Println()
	}
	return nil
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("Usage: go run . <string> [standard|shadow|thinkertoy]")
		return
	}

	input := os.Args[1]
	if input == "" {
		return
	}

	bannerStyle := ""
	if len(os.Args) == 3 {
		bannerStyle = os.Args[2]
	}

	filename, err := bannerFilename(bannerStyle)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error: could not open banner file:", filename)
		return
	}

	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	banner := strings.Split(strings.Trim(content, "\n"), "\n")

	input = strings.ReplaceAll(input, newlineToken, "\n")

	if input == "\n" {
		fmt.Println()
		return
	}

	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			fmt.Println()
			continue
		}
		if err := printAscii(line, banner); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	}
}

func bannerFilename(style string) (string, error) {
	switch style {
	case "", "standard":
		return "standard.txt", nil
	case "shadow":
		return "shadow.txt", nil
	case "thinkertoy":
		return "thinkertoy.txt", nil
	default:
		return "", fmt.Errorf("unknown banner style %q (expected standard, shadow or thinkertoy)", style)
	}
}