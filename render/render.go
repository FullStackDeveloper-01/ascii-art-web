package render

import (
	"fmt"
	"strings"
)

func Render(text string, chars map[rune][]string) (string, error) {
	if strings.Trim(text, "\n") == "" {
		return text, nil
	}

	lines := strings.Split(text, "\n")

	var out strings.Builder
	for _, line := range lines {
		if line == "" {
			out.WriteString("\n")
			continue
		}
		rendered, err := renderLine(line, chars)
		if err != nil {
			return "", err
		}
		out.WriteString(rendered)
	}
	return out.String(), nil
}

func renderLine(line string, chars map[rune][]string) (string, error) {
	runes := []rune(line)

	blocks := make([][]string, len(runes))
	for i, r := range runes {
		block, ok := chars[r]
		if !ok {
			return "", fmt.Errorf("character not supported: %q", r)
		}
		blocks[i] = block
	}

	var out strings.Builder
	for row := 0; row < 8; row++ {
		for _, block := range blocks {
			out.WriteString(block[row])
		}
		out.WriteString("\n")
	}
	return out.String(), nil
}
