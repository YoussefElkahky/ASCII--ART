package AsciiArt

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed standard.txt shadow.txt thinkertoy.txt
var bannerFiles embed.FS

const (
	charHeight     = 8
	linesPerChar   = charHeight + 1
	firstPrintable = 32
	lastPrintable  = 126
)

func BannerFilename(style string) (string, error) { 
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

func LoadBanner(filename string) ([]string, error) {

	data, err := bannerFiles.ReadFile(filename) 
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

func Render(input string, banner []string) (string, error) {

	for _, ch := range input {
		if ch == '\n' {
			continue
		}
		if ch < firstPrintable || ch > lastPrintable {
			return "", fmt.Errorf("unsupported character: %q", ch)
		}
	}

	trimmed := input
	if strings.Trim(input, "\n") == "" {
		trimmed = strings.TrimSuffix(input, "\n")
	}

	lines := strings.Split(trimmed, "\n")

	var output strings.Builder

	for _, line := range lines {
		if line == "" {
			output.WriteString("\n")
			continue
		}

		output.WriteString(printAscii(line, banner))
	}

	return output.String(), nil
}

func printAscii(text string, banner []string) string {

	var builder strings.Builder
	var result strings.Builder

	for row := 0; row < charHeight; row++ {
		builder.Reset()
		for _, ch := range text {
			index := (int(ch)-firstPrintable)*linesPerChar + row
			builder.WriteString(banner[index])
		}
		result.WriteString(builder.String())
		result.WriteString("\n")
	}

	return result.String()
}