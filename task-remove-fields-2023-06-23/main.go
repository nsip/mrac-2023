package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/digisan/gotk/strs"
	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
)

// remove "dc:text", "dc:title"
func rmOneLineField(js string, fields ...string) string {
	js = jt.FmtStr(js, "  ") // formatted
	for i, field := range fields {
		fields[i] = fmt.Sprintf(`"%s"`, strings.Trim(field, `"`)) // wrapped with "field"
	}
	idx := 0
	rt, err := strs.StrLineScan(js, func(line string) (bool, string) {
		idx++
		// remove fields
		for _, field := range fields {
			if strings.HasPrefix(strings.TrimSpace(line), field+":") {
				// fmt.Println(idx, line)
				return false, ""
			}
		}
		return true, line
	})
	lk.FailOnErr("%v", err)
	return jt.FmtStr(rt, "  ")
}

// "dc:description" => { "@language", "@value" }
func changeOneLineToStruct(js, field string) string {
	js = jt.FmtStr(js, "  ")                              // formatted
	field = fmt.Sprintf(`"%s"`, strings.Trim(field, `"`)) // wrapped with "field"
	idx := 0
	rt, err := strs.StrLineScan(js, func(line string) (bool, string) {
		idx++
		trimed := strings.TrimSpace(line)

		appendComma := false
		if strings.HasSuffix(trimed, ",") {
			appendComma = true
			trimed = strings.TrimSuffix(trimed, ",")
		}

		if strings.HasPrefix(trimed, field+":") {
			v := strings.Trim(strings.TrimSpace(strings.TrimPrefix(trimed, field+":")), "\"")
			stru := fmt.Sprintf(`%s: { "@language": "en-au", "@value": "%s" }`, field, v)
			if appendComma {
				stru += ","
			}
			return true, stru
		}
		return true, line
	})
	lk.FailOnErr("%v", err)
	return jt.FmtStr(rt, "  ")
}

// @ asn-json-ld
// remove "dc:text", "dc:title"
// "dc:description" => { "@language", "@value" }

func main() {

	fPath := "../data-out/asn-json-ld/la-Science.json"
	data, err := os.ReadFile(fPath)
	lk.FailOnErr("%v", err)

	toBeRemoved := []string{"dc:text", "dc:title"}
	rt := rmOneLineField(string(data), toBeRemoved...)

	/////////////////////////////

	toBeChanged := "dc:description"
	rt = changeOneLineToStruct(rt, toBeChanged)

	err = os.WriteFile(fPath, []byte(rt), os.ModePerm)
	lk.FailOnErr("%v", err)
}
