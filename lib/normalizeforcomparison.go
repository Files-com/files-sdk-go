package lib

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"sync"
)

// path_comparison.json.gz is shared/path_comparison.json, gzip-compressed by
// the generator to keep linked binaries small.
//
//go:embed path_comparison.json.gz
var pathComparisonGzip []byte

// comparisonMap decodes the embedded table on the first lookup outside
// printable ASCII, not at startup. If the table is corrupt, that lookup and
// every later one panic.
var comparisonMap = sync.OnceValue(func() map[rune]string {
	reader, err := gzip.NewReader(bytes.NewReader(pathComparisonGzip))
	if err != nil {
		panic(err)
	}
	pathComparisonJSON, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	if err := reader.Close(); err != nil {
		panic(err)
	}
	var data struct {
		Mapping map[string]string `json:"mapping"`
	}
	if err := json.Unmarshal(pathComparisonJSON, &data); err != nil {
		panic(err)
	}
	mapping := make(map[rune]string, len(data.Mapping))
	for hex, replacement := range data.Mapping {
		scalar, err := strconv.ParseUint(hex, 16, 32)
		if err != nil {
			panic(err)
		}
		mapping[rune(scalar)] = replacement
	}
	return mapping
})

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
		} else if replacement, ok := comparisonMap()[scalar]; ok {
			result.WriteString(replacement)
		} else {
			result.WriteRune(scalar)
		}
	}
	return result.String()
}
