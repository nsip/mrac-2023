package tool

import (
	"regexp"
	"strings"
)

var (
	sJoin       = strings.Join
	sTrim       = strings.Trim
	sTrimSuffix = strings.TrimSuffix
	sLastIndex  = strings.LastIndex
)

var (
	r4value  = regexp.MustCompile(`:\s*"`)
	r4array  = regexp.MustCompile(`:\s*\[`)
	r4arrstr = regexp.MustCompile(`(("[^"]+"),?)+`)
)

func FetchValue(kvstr, sep string) string {

	// if sContains(kvstr, "Levels") {
	// 	fmt.Println("connections.Levels")
	// }

	// simple value
	if loc := r4value.FindStringIndex(kvstr); loc != nil {
		start := loc[1]
		end := sLastIndex(kvstr, `"`)
		return kvstr[start:end]
	}

	// array value content
	if loc := r4array.FindStringIndex(kvstr); loc != nil {
		start := loc[1] + 1
		end := sLastIndex(kvstr, `]`)
		arrcont := sTrim(kvstr[start:end], " \n\t")
		items := r4arrstr.FindAllString(arrcont, -1)
		for i, item := range items {
			item = sTrimSuffix(item, ",")
			items[i] = sTrim(item, "\"")
		}
		return sJoin(items, sep)
	}

	panic(kvstr)
}
