package banner

import (
	_ "embed"
	"fmt"
	"strings"
)

const characters = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"

//go:embed standard.txt
var standardData string

//go:embed shadow.txt
var shadowData string

//go:embed thinkertoy.txt
var thinkertoyData string

func Load(name string) (map[rune][]string, error) {
	var data string
	switch name {
	case "", "standard":
		data = standardData
	case "shadow":
		data = shadowData
	case "thinkertoy":
		data = thinkertoyData
	default:
		return nil, fmt.Errorf("unknown banner: %q", name)
	}

	return parse(data)
}

func parse(data string) (map[rune][]string, error) {
	data = strings.ReplaceAll(data, "\r\n", "\n")
	lines := strings.Split(strings.TrimRight(data, "\n"), "\n")

	m := make(map[rune][]string)
	cursor := 0
	for _, char := range characters {
		cursor++
		if cursor+8 > len(lines) {
			return nil, fmt.Errorf("banner file is incomplete")
		}
		m[char] = lines[cursor : cursor+8]
		cursor += 8
	}

	if cursor != len(lines) {
		return nil, fmt.Errorf("banner file has unexpected length")
	}

	return m, nil
}
