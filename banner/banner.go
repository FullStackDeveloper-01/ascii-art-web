package banner

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
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

// Known-good SHA-256 of each banner file (computed once with
// `sha256sum banner/*.txt` while the files were still correct). If a file
// gets edited later — even a single character, same line count and all —
// its hash won't match anymore and Load() rejects it.
const (
	standardChecksum   = "e194f1033442617ab8a78e1ca63a2061f5cc07a3f05ac226ed32eb9dfd22a6bf"
	shadowChecksum     = "26b94d0b134b77e9fd23e0360bfd81740f80fb7f6541d1d8c5d85e73ee550f73"
	thinkertoyChecksum = "64285e4960d199f4819323c4dc6319ba34f1f0dd9da14d07111345f5d76c3fa3"
)

func Load(name string) (map[rune][]string, error) {
	var data, checksum string
	switch name {
	case "", "standard":
		data, checksum = standardData, standardChecksum
	case "shadow":
		data, checksum = shadowData, shadowChecksum
	case "thinkertoy":
		data, checksum = thinkertoyData, thinkertoyChecksum
	default:
		return nil, fmt.Errorf("unknown banner: %q", name)
	}

	sum := sha256.Sum256([]byte(data))
	if hex.EncodeToString(sum[:]) != checksum {
		return nil, fmt.Errorf("banner file %q has been modified", name)
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
