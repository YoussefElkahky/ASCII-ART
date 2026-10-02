package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	charHeight     = 8
	linesPerChar   = charHeight + 1 // 8 lines per character + separator
	firstPrintable = 32
	lastPrintable  = 126
)

func main() {

	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "Error: usage: go run . <string> [standard|shadow|thinkertoy]")
		os.Exit(1)
	}

	input := os.Args[1]

	if input == "" {
		return
	}

	bannerStyle := ""

	if len(os.Args) == 3 {
		bannerStyle = os.Args[2]
	}

	err := run(input, bannerStyle)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run(input string, bannerStyle string) error {

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

	// (126 - 32 + 1) printable ASCII characters * 9 lines per character.
	expectedLines := (lastPrintable - firstPrintable + 1) * linesPerChar

	banner := strings.Split(strings.TrimPrefix(content, "\n"), "\n")

	if len(banner) != expectedLines {
		return nil, fmt.Errorf(
			"invalid banner file %s: expected %d lines, got %d",
			filename,
			expectedLines,
			len(banner),
		)
	}

	return banner, nil
}

func render(input string, banner []string) error {

	for _, ch := range input {
		if ch == '\n' {
			continue
		}
		if ch < firstPrintable || ch > lastPrintable {
			return fmt.Errorf("unsupported character: %q", ch)
		}
	}

	trimmed := input
	if strings.Trim(input, "\n") == "" {
		trimmed = strings.TrimSuffix(input, "\n")
	}

	lines := strings.Split(trimmed, "\n")

	for _, line := range lines {
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
