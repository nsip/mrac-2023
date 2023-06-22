package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/digisan/gotk/strs"
	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
)

func rmOneLineField(js, field string, onValues ...string) string {
	js = jt.FmtStr(js, "  ")                              // formatted
	field = fmt.Sprintf(`"%s"`, strings.Trim(field, `"`)) // warpped with "field"
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
			// fmt.Println(1, idx, line)
		}
		return true, line
	})
	lk.FailOnErr("%v", err)
	return jt.FmtStr(rt, "  ")
}

func main() {

	toBeRemoved := "asn_conceptTerm"
	onValues := []string{"[]", "SCIENCE_TEACHER_BACKGROUND_INFORMATION"}

	fPath := "../data-out/asn-json/la-English.json"
	data, err := os.ReadFile(fPath)
	lk.FailOnErr("%v", err)

	rt := rmOneLineField(string(data), toBeRemoved, onValues...)

	err = os.WriteFile(fPath, []byte(rt), os.ModePerm)
	lk.FailOnErr("%v", err)
}
