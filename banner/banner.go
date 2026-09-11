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
	lines := strings.Split(strings.TrimRight(data, "\n"), "\n")

	m := make(map[rune][]string)
	cursor := 0
	for _, char := range characters {
		cursor++
		m[char] = lines[cursor : cursor+8]
		cursor += 8
	}

	return m, nil
}
