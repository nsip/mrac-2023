package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/digisan/gotk/strs"
	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
)

func addSquareBrackets(js, field string) string {
	js = jt.FmtStr(js, "  ")
	field = fmt.Sprintf(`"%s"`, strings.Trim(field, `"`))
	idx := 0
	strs.StrLineScan(js, func(line string) (bool, string) {
		idx++
		if strings.HasPrefix(strings.TrimSpace(line), field+": {") {
			fmt.Println(1, idx, line)
		}
		if strings.HasPrefix(strings.TrimSpace(line), field+": [") {
			fmt.Println(2, idx, line)
		}
		return true, ""
	})
	return field
}

func main() {

	data, err := os.ReadFile("../data-out/asn-json/la-English.json")
	lk.FailOnErr("%v", err)

	r := addSquareBrackets(string(data), "asn_hasLevel")
	fmt.Println(r)
}
