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

type ProcPos struct {
	flag  bool
	idx   int
	space string
}

func (pp *ProcPos) Init() {
	pp.flag = false
	pp.idx = 0
	pp.space = ""
}

func (pp *ProcPos) Start(i int, prefixSpace string) {
	pp.flag = true
	pp.idx = i
	pp.space = prefixSpace
}

func (pp *ProcPos) End(i int) {
	pp.flag = false
	pp.idx = i
	pp.space = ""
}

func addSquareBrackets(js, field string) string {
	js = jt.FmtStr(js, "  ")                              // formatted
	field = fmt.Sprintf(`"%s"`, strings.Trim(field, `"`)) // wrapped with "field"
	idx := 0

	pp := ProcPos{}
	pp.Init()

	rt, err := strs.StrLineScan(js, func(line string) (bool, string) {

		idx++

		if strings.HasPrefix(strings.TrimSpace(line), field+": {") {
			// fmt.Println(1, idx, line)
			space := strs.TrimTailFromFirst(line, "\"")
			pp.Start(idx, space)
			return true, strings.Replace(line, "{", "[{", 1)
		}

		// if strings.HasPrefix(strings.TrimSpace(line), field+": [") {
		// 	fmt.Println(2, idx, line)
		//  return true, ""
		// }

		if pp.flag {
			if line == pp.space+"}" || line == pp.space+"}," {
				pp.End(idx)
				return true, strings.Replace(line, "}", "}]", 1)
			}
		}

		return true, line
	})

	lk.FailOnErr("%v", err)
	return jt.FmtStr(rt, "  ")
}

// @asn-json
// make all the "asn_hasLevel" value be array type

func main() {

	const (
		mustBeArray = "asn_hasLevel"
		inputDir    = "../data-out/asn-json/"
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
			rt := addSquareBrackets(string(data), mustBeArray)
			err = os.WriteFile(fPath, []byte(rt), os.ModePerm)
			lk.FailOnErr("%v", err)

		}
	}
}
