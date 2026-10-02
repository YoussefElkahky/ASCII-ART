package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// ============================================================
// Helpers
// ============================================================

func captureOutput(f func()) string {
	r, w, _ := os.Pipe()

	old := os.Stdout
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

func stubBanner() []string {
	var banner []string

	for code := firstPrintable; code <= lastPrintable; code++ {
		for row := 0; row < charHeight; row++ {
			banner = append(banner, "X")
		}

		if code != lastPrintable {
			banner = append(banner, "")
		}
	}

	return banner
}

func loadTestBanner(t *testing.T) []string {
	t.Helper()

	banner, err := loadBanner("standard.txt")
	if err != nil {
		t.Skip("standard.txt not found")
	}

	return banner
}

// ============================================================
// bannerFilename Tests
// ============================================================

func TestBannerFilename(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
		err   bool
	}{
		{"default", "", "standard.txt", false},
		{"standard", "standard", "standard.txt", false},
		{"shadow", "shadow", "shadow.txt", false},
		{"thinkertoy", "thinkertoy", "thinkertoy.txt", false},
		{"invalid", "abc", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bannerFilename(tt.input)

			if tt.err {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("got %q want %q", got, tt.want)
			}
		})
	}
}

// ============================================================
// loadBanner Tests
// ============================================================

func TestLoadBannerNotFound(t *testing.T) {
	_, err := loadBanner("missing.txt")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadBannerStandard(t *testing.T) {
	banner := loadTestBanner(t)

	if len(banner) == 0 {
		t.Fatal("banner should not be empty")
	}

	for _, line := range banner {
		if strings.Contains(line, "\r") {
			t.Error("banner contains carriage return")
		}
	}
}

// ============================================================
// printAscii Tests
// ============================================================

func TestPrintAsciiSingleCharacter(t *testing.T) {
	var err error

	output := captureOutput(func() {
		err = printAscii("A", stubBanner())
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")

	if len(lines) != charHeight {
		t.Errorf("got %d lines want %d", len(lines), charHeight)
	}
}

func TestPrintAsciiMultipleCharacters(t *testing.T) {
	var err error

	output := captureOutput(func() {
		err = printAscii("ABC", stubBanner())
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")

	if len(lines) != charHeight {
		t.Errorf("got %d lines want %d", len(lines), charHeight)
	}
}

func TestPrintAsciiRealBanner(t *testing.T) {
	banner := loadTestBanner(t)

	err := printAscii("Hello", banner)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ============================================================
// render Tests
// ============================================================

func TestRenderEmptyLine(t *testing.T) {
	output := captureOutput(func() {
		err := render("", stubBanner())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if output != "\n" {
		t.Errorf("got %q", output)
	}
}

func TestRenderNewLine(t *testing.T) {
	output := captureOutput(func() {
		err := render("\n", stubBanner())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if output != "\n" {
		t.Errorf("got %q", output)
	}
}

func TestRenderMultipleLines(t *testing.T) {
	output := captureOutput(func() {
		err := render("A\nB", stubBanner())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")

	if len(lines) != charHeight*2 {
		t.Errorf("got %d lines", len(lines))
	}
}

// ============================================================
// Index Tests
// ============================================================

func TestCharacterIndexCalculation(t *testing.T) {
	tests := []struct {
		char rune
		want int
	}{
		{' ', 0},
		{'!', 9},
		{'A', (65 - firstPrintable) * linesPerChar},
		{'~', (126 - firstPrintable) * linesPerChar},
	}

	for _, tt := range tests {
		got := (int(tt.char) - firstPrintable) * linesPerChar

		if got != tt.want {
			t.Errorf("char %q: got %d want %d", tt.char, got, tt.want)
		}
	}
}
