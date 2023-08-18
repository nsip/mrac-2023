package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/digisan/gotk/strs"
	lk "github.com/digisan/logkit"
	u "github.com/nsip/mrac-2023/util"
)

// remove "dc:text", "dc:title"
func rmOneLineField(js string, fields ...string) string {

	// js = jt.FmtStr(js, "  ") // formatted
	js, err := u.FmtJSON(js)
	lk.FailOnErr("%v", err)

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

	// return jt.FmtStr(rt, "  ")

	rt, err = u.FmtJSON(rt)
	lk.FailOnErr("%v", err)

	return rt
}

// "dc:description" => { "@language", "@value" }
func changeOneLineToStruct(js, field string, f1, f2 string) string {

	// js = jt.FmtStr(js, "  ")                              // formatted
	js, err := u.FmtJSON(js)
	lk.FailOnErr("%v", err)

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

		if strings.HasPrefix(trimed, field+": \"") {
			v := strings.Trim(strings.TrimSpace(strings.TrimPrefix(trimed, field+":")), "\"")
			stru := fmt.Sprintf(`%s: { "%s": "en-au", "%s": "%s" }`, field, f1, f2, v)
			if appendComma {
				stru += ","
			}
			return true, stru
		}
		return true, line
	})
	lk.FailOnErr("%v", err)

	// return jt.FmtStr(rt, "  ")
	rt, err = u.FmtJSON(rt)
	lk.FailOnErr("%v", err)

	return rt
}

// @ asn-json-ld
// remove "dc:text", "dc:title"
// "dc:description" => { "@language", "@value" }

func main() {

	const (
		inputDir = "../data-out/asn-json-ld/"
	)

	var (
		toBeRemoved = []string{"dc:text", "dc:title"}

		toBeChanged1   = "dc:description"
		toBeChanged1F1 = "@language"
		toBeChanged1F2 = "@value"

		// toBeChanged2   = "dc:title"
		// toBeChanged2F1 = "@language"
		// toBeChanged2F2 = "@literal"
	)

	de, err := os.ReadDir(inputDir)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	for _, f := range de {
		fName := f.Name()
		if strs.HasAnySuffix(fName, ".json", ".jsonld") {

			fPath := filepath.Join(inputDir, fName)

			fmt.Printf("processing... %s\n", fPath)

			data, err := os.ReadFile(fPath)
			lk.FailOnErr("%v", err)

			rt := rmOneLineField(string(data), toBeRemoved...)

			rt = changeOneLineToStruct(rt, toBeChanged1, toBeChanged1F1, toBeChanged1F2)
			// rt = changeOneLineToStruct(rt, toBeChanged2, toBeChanged2F1, toBeChanged2F2)

			err = os.WriteFile(fPath, []byte(rt), os.ModePerm)
			lk.FailOnErr("%v", err)

		}
	}
}
