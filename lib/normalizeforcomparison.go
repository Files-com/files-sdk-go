package lib

import (
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/width"
)

const TRANSLIT_MAP = "ÀA,ÁA,ÂA,ÃA,ÄA,ÅA,ÆAE,ÇC,ÈE,ÉE,ÊE,ËE,ÌI,ÍI,ÎI,ÏI,ÐD,ÑN,ÒO,ÓO,ÔO,ÕO,ÖO,ØO,ÙU,ÚU,ÛU,ÜU,ÝY,ßss,àa,áa,âa,ãa,äa,åa,æae,çc,èe,ée,êe,ëe,ìi,íi,îi,ïi,ðd,ñn,òo,óo,ôo,õo,öo,øo,ùu,úu,ûu,üu,ýy,ÿy,ĀA,āa,ĂA,ăa,ĄA,ąa,ĆC,ćc,ĈC,ĉc,ĊC,ċc,ČC,čc,ĎD,ďd,ĐD,đd,ĒE,ēe,ĔE,ĕe,ĖE,ėe,ĘE,ęe,ĚE,ěe,ĜG,ĝg,ĞG,ğg,ĠG,ġg,ĢG,ģg,ĤH,ĥh,ĦH,ħh,ĨI,ĩi,ĪI,īi,ĬI,ĭi,ĮI,įi,İI,ĲIJ,ĳij,ĴJ,ĵj,ĶK,ķk,ĹL,ĺl,ĻL,ļl,ĽL,ľl,ŁL,łl,ŃN,ńn,ŅN,ņn,ŇN,ňn,ŉ'n,ŌO,ōo,ŎO,ŏo,ŐO,őo,ŒOE,œoe,ŔR,ŕr,ŖR,ŗr,ŘR,řr,ŚS,śs,ŜS,ŝs,ŞS,şs,ŠS,šs,ŢT,ţt,ŤT,ťt,ŨU,ũu,ŪU,ūu,ŬU,ŭu,ŮU,ůu,ŰU,űu,ŲU,ųu,ŴW,ŵw,ŶY,ŷy,ŸY,ŹZ,źz,ŻZ,żz,ŽZ,žz"

func NormalizeForComparison(apiPath string) string {
	// Normalize the path
	normalizedPath := strings.Replace(apiPath, "\x00", "", -1)
	normalizedPath = strings.ReplaceAll(normalizedPath, "..", "")
	normalizedPath = path.Clean(normalizedPath)
	normalizedPath = filepath.Clean(normalizedPath)
	normalizedPath = filepath.ToSlash(normalizedPath)
	normalizedPath = strings.ReplaceAll(normalizedPath, "\\", "/")
	normalizedPath = strings.ReplaceAll(normalizedPath, "//", "/")
	normalizedPath = strings.TrimPrefix(normalizedPath, "/")
	normalizedPath = trimTrailingWhitespace(normalizedPath)
	normalizedPath = strings.TrimRight(normalizedPath, " ")

	// Normalize Unicode characters
	normalizedPath = unicodeNormalizeAndTransliterate(strings.ToLower(normalizedPath))
	return normalizedPath
}

func createTransliterationMap() map[rune]string {
	mapping := make(map[rune]string)
	pairs := strings.Split(TRANSLIT_MAP, ",")

	for _, pair := range pairs {
		runes := []rune(pair)
		mapping[runes[0]] = string(runes[1:])
	}

	return mapping
}

func unicodeNormalizeAndTransliterate(s string) string {
	var result strings.Builder
	t := transform.Chain(norm.NFKC, width.Fold)
	normalized, _, _ := transform.String(t, s)

	for _, r := range normalized {
		if v, ok := createTransliterationMap()[r]; ok {
			result.WriteString(v)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func trimTrailingWhitespace(s string) string {
	return strings.TrimRightFunc(s, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == '\r'
	})
}
