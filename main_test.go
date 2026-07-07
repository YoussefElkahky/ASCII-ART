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

func assertSliceEqual(t *testing.T, got, expected []string) {
	t.Helper()

	if len(got) != len(expected) {
		t.Errorf("got %d elements, want %d", len(got), len(expected))
		return
	}

	for i := range expected {
		if got[i] != expected[i] {
			t.Errorf("index %d: got %q, want %q", i, got[i], expected[i])
		}
	}
}

func stubBanner() []string {
	var banner []string

	for code := 32; code <= 126; code++ {
		for row := 0; row < 8; row++ {
			banner = append(banner, "X")
		}
		if code < 126 {
			banner = append(banner, "")
		}
	}

	return banner
}

func loadStandardBanner(t *testing.T) []string {
	t.Helper()

	data, err := os.ReadFile("standard.txt")
	if err != nil {
		t.Skip("standard.txt not found")
	}

	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	content = strings.Trim(content, "\n")

	return strings.Split(content, "\n")
}

// ============================================================
// bannerFilename Tests
// ============================================================

func TestBannerFilenameDefault(t *testing.T) {
	got, err := bannerFilename("")
	want := "standard.txt"

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBannerFilenameShadow(t *testing.T) {
	got, err := bannerFilename("shadow")
	want := "shadow.txt"

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBannerFilenameThinkertoy(t *testing.T) {
	got, err := bannerFilename("thinkertoy")
	want := "thinkertoy.txt"

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBannerFilenameInvalid(t *testing.T) {
	_, err := bannerFilename("invalid")

	if err == nil {
		t.Error("expected error")
	}
}

// ============================================================
// printAscii Tests
// ============================================================

func TestPrintAsciiSingleChar(t *testing.T) {
	banner := stubBanner()

	var err error

	got := captureOutput(func() {
		err = printAscii("A", banner)
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(got, "\n"), "\n")

	if len(lines) != 8 {
		t.Errorf("got %d lines, want 8", len(lines))
	}
}

func TestPrintAsciiMultipleChars(t *testing.T) {
	banner := stubBanner()

	var err error

	got := captureOutput(func() {
		err = printAscii("ABC", banner)
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(got, "\n"), "\n")

	if len(lines) != 8 {
		t.Errorf("got %d lines, want 8", len(lines))
	}
}

func TestPrintAsciiSpace(t *testing.T) {
	banner := stubBanner()

	var err error

	got := captureOutput(func() {
		err = printAscii(" ", banner)
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got == "" {
		t.Error("expected output")
	}
}

func TestPrintAsciiInvalidCharacter(t *testing.T) {
	banner := stubBanner()

	err := printAscii("ñ", banner)

	if err == nil {
		t.Error("expected error")
	}
}

func TestPrintAsciiAllPrintableASCII(t *testing.T) {
	banner := stubBanner()

	var all strings.Builder

	for i := 32; i <= 126; i++ {
		all.WriteRune(rune(i))
	}

	var err error

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic: %v", r)
		}
	}()

	captureOutput(func() {
		err = printAscii(all.String(), banner)
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPrintAsciiRealBannerHello(t *testing.T) {
	banner := loadStandardBanner(t)

	var err error

	got := captureOutput(func() {
		err = printAscii("Hello", banner)
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(got, "\n"), "\n")

	if len(lines) != 8 {
		t.Errorf("got %d lines, want 8", len(lines))
	}
}

func TestPrintAsciiRealBannerNumbers(t *testing.T) {
	banner := loadStandardBanner(t)

	var err error

	captureOutput(func() {
		err = printAscii("1234567890", banner)
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPrintAsciiRealBannerSymbols(t *testing.T) {
	banner := loadStandardBanner(t)

	var err error

	captureOutput(func() {
		err = printAscii("!@#$%^&*()", banner)
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ============================================================
// Input Parsing Tests
// ============================================================

func TestParseInputNoNewline(t *testing.T) {
	input := strings.ReplaceAll("hello", "\\n", "\n")

	got := strings.Split(input, "\n")
	want := []string{"hello"}

	assertSliceEqual(t, got, want)
}

func TestParseInputSingleNewline(t *testing.T) {
	input := strings.ReplaceAll("hello\\nworld", "\\n", "\n")

	got := strings.Split(input, "\n")
	want := []string{"hello", "world"}

	assertSliceEqual(t, got, want)
}

func TestParseInputDoubleNewline(t *testing.T) {
	input := strings.ReplaceAll("hello\\n\\nworld", "\\n", "\n")

	got := strings.Split(input, "\n")
	want := []string{"hello", "", "world"}

	assertSliceEqual(t, got, want)
}

func TestParseInputOnlyNewline(t *testing.T) {
	got := strings.ReplaceAll("\\n", "\\n", "\n")
	want := "\n"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestParseInputNewlineAtEnd(t *testing.T) {
	input := strings.ReplaceAll("hello\\n", "\\n", "\n")

	got := strings.Split(input, "\n")
	want := []string{"hello", ""}

	assertSliceEqual(t, got, want)
}

// ============================================================
// Banner Loading Tests
// ============================================================

func TestBannerLoadingFileNotFound(t *testing.T) {
	_, err := os.ReadFile("notfound.txt")

	if err == nil {
		t.Error("expected error")
	}
}

func TestBannerLoadingLineCount(t *testing.T) {
	banner := loadStandardBanner(t)

	if len(banner) < 854 {
		t.Errorf("got %d lines, expected at least 854", len(banner))
	}
}

func TestBannerLoadingNoCarriageReturn(t *testing.T) {
	banner := loadStandardBanner(t)

	for _, line := range banner {
		if strings.Contains(line, "\r") {
			t.Error("banner still contains carriage return")
		}
	}
}

// ============================================================
// Index Tests
// ============================================================

func TestSpaceIndex(t *testing.T) {
	got := (32 - 32) * linesPerChar
	want := 0

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestLastCharacterIndex(t *testing.T) {
	banner := stubBanner()

	got := (126-firstPrintable)*linesPerChar + (charHeight - 1)

	if got >= len(banner) {
		t.Errorf("index %d out of range", got)
	}
}
