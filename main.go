package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	charHeight     = 8
	linesPerChar   = charHeight + 1
	firstPrintable = 32
	lastPrintable  = 126
)

func main() {
	err := run()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		return fmt.Errorf("usage: go run . <string> [standard|shadow|thinkertoy]")
	}

	input := os.Args[1]

	if input == "" {
		return nil
	}

	bannerStyle := ""
	if len(os.Args) == 3 {
		bannerStyle = os.Args[2]
	}

	filename, err := bannerFilename(bannerStyle)
	if err != nil {
		return err
	}

	banner, err := loadBanner(filename)
	if err != nil {
		return err
	}

	input = strings.ReplaceAll(input, `\n`, "\n")

	return render(input, banner)
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

func loadBanner(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open banner file %s: %s", filename, err)
	}

	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	banner := strings.Split(strings.Trim(content, "\n"), "\n")

	return banner, nil
}

func render(input string, banner []string) error {
	if input == "\n" {
		fmt.Println()
		return nil
	}

	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			fmt.Println()
			continue
		}

		err := printAscii(line, banner)

		if err != nil {
			return err
		}
	}

	return nil
}

func printAscii(text string, banner []string) error {
	for _, ch := range text {
		if ch < firstPrintable || ch > lastPrintable {
			return fmt.Errorf("character %q is not a printable ASCII character", ch)
		}
	}

	var builder strings.Builder

	for row := 0; row < charHeight; row++ {
		builder.Reset()
		for _, ch := range text {
			index := (int(ch)-firstPrintable)*linesPerChar + row
			builder.WriteString(banner[index])
		}
		fmt.Println(builder.String())
	}

	return nil
}
