package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/digisan/gotk/strs"
	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
)

func rmOneLineField(js, field string, onValues ...string) string {
	js = jt.FmtStr(js, "  ")                              // formatted
	field = fmt.Sprintf(`"%s"`, strings.Trim(field, `"`)) // wrapped with "field"
	idx := 0
	rt, err := strs.StrLineScan(js, func(line string) (bool, string) {
		idx++
		if strings.HasPrefix(strings.TrimSpace(line), field) {
			for _, v := range onValues {
				if v != "[]" {
					v = fmt.Sprintf(`"%v"`, v)
				}
				if strings.HasSuffix(line, v) || strings.HasSuffix(line, v+",") {
					return false, ""
				}
			}
		}
		return true, line
	})
	lk.FailOnErr("%v", err)
	return jt.FmtStr(rt, "  ")
}

// @asn-json
// remove "asn_conceptTerm" with value ["[]", "SCIENCE_TEACHER_BACKGROUND_INFORMATION"]

func main() {

	const (
		inputDir    = "../data-out/asn-json/"
		toBeRemoved = "asn_conceptTerm"
	)

	var (
		onValues = []string{"[]", "SCIENCE_TEACHER_BACKGROUND_INFORMATION"}
	)

	de, err := os.ReadDir(inputDir)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	for _, f := range de {
		fName := f.Name()
		if strs.HasAnySuffix(fName, ".json") {

			fPath := filepath.Join(inputDir, fName)
			data, err := os.ReadFile(fPath)
			lk.FailOnErr("%v", err)
			rt := rmOneLineField(string(data), toBeRemoved, onValues...)
			err = os.WriteFile(fPath, []byte(rt), os.ModePerm)
			lk.FailOnErr("%v", err)

		}
	}
}
