package lib

import (
	_ "embed"
	"encoding/json"
	"strconv"
	"strings"
)

//go:embed path_comparison.json
var pathComparisonJSON []byte

var comparisonMap map[rune]string

func init() {
	var data struct {
		Mapping map[string]string `json:"mapping"`
	}
	if err := json.Unmarshal(pathComparisonJSON, &data); err != nil {
		panic(err)
	}
	comparisonMap = make(map[rune]string, len(data.Mapping))
	for hex, replacement := range data.Mapping {
		scalar, err := strconv.ParseUint(hex, 16, 32)
		if err != nil {
			panic(err)
		}
		comparisonMap[rune(scalar)] = replacement
	}
}

func NormalizeForComparison(path string) string {
	path = NormalizeAPIPath(path)
	var result strings.Builder
	result.Grow(len(path))
	for _, scalar := range path {
		if scalar >= ' ' && scalar <= '~' {
			if scalar >= 'A' && scalar <= 'Z' {
				scalar += 'a' - 'A'
			}
			result.WriteByte(byte(scalar))
		} else if replacement, ok := comparisonMap[scalar]; ok {
			result.WriteString(replacement)
		} else {
			result.WriteRune(scalar)
		}
	}
	return result.String()
}
