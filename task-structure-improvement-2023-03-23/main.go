package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	. "github.com/digisan/go-generics/v2"
	dt "github.com/digisan/gotk/data-type"
	fd "github.com/digisan/gotk/file-dir"
	"github.com/digisan/gotk/strs"
	jt "github.com/digisan/json-tool"
)

func main() {

	var (
		todo = [6]bool{true, true, true, true, true, true}
	)

	inDir := "../task-split-ccp-2022-11-12/json-ld/MRAC/2023/06/GC/CCP"
	outDir := "./out/json-ld-ccp"
	fd.MustCreateDir(outDir)

	files, err := os.ReadDir(inDir)
	if err != nil {
		log.Fatalln(err)
		return
	}

	for _, file := range files {

		fName := file.Name()
		fPath := filepath.Join(inDir, fName)

		data, err := os.ReadFile(fPath)
		if err != nil {
			log.Panicln(err)
			return
		}

		if !dt.IsJSON(data) {
			fmt.Println("1 --->", fPath)
			return
		}

		///////////

		met := false
		head := ""

		fd.FileLineScanEx(
			fPath,
			3, 3, "", func(line string, cache []string) (bool, string) {

				// *** 1. date-time format ***
				//
				if todo[0] {
					if strings.Contains(line, `"@value"`) {
						ln := strings.TrimSpace(line)
						above1, _, _ := strings.TrimSpace(cache[2]), strings.TrimSpace(cache[1]), strings.TrimSpace(cache[0])
						below1, _, _ := strings.TrimSpace(cache[4]), strings.TrimSpace(cache[5]), strings.TrimSpace(cache[6])
						if In(`"@type": "xsd:dateTime"`, above1, below1) {
							val := strings.TrimPrefix(ln, `"@value": `)
							val = strings.TrimSuffix(val, `,`)
							val = strings.Trim(val, `"`)
							dt, ok := TryToDateTime(val)
							if !ok {
								log.Println("cannot be date-time")
								os.Exit(-1)
							}
							return true, fmt.Sprintf(`%s"@value": "%s",`, strs.TrimTailFromLast(line, `"@value"`), dt.Format(`2006-01-02T15:04:05.000Z`))
						}
					}
				}

				// *** 2. prefLabel on curriculum statements as concepts ***
				//
				if todo[1] {
					if strings.Contains(line, `"asn:statementNotation"`) {
						ln := strings.TrimSpace(line)
						// above1, _, _ := strings.TrimSpace(cache[2]), strings.TrimSpace(cache[1]), strings.TrimSpace(cache[0])
						// below1, _, _ := strings.TrimSpace(cache[4]), strings.TrimSpace(cache[5]), strings.TrimSpace(cache[6])

						head := strs.TrimTailFromLast(line, `"asn:statementNotation"`)

						valStr := strings.TrimPrefix(ln, `"asn:statementNotation"`)
						valStr = strings.TrimPrefix(valStr, ":")
						valStr = strings.TrimSpace(valStr)
						valStr = strings.TrimSuffix(valStr, ",")
						valStr = strings.Trim(valStr, `"`)

						newFieldVal := fmt.Sprintf(`"skos:prefLabel": "%s",`, valStr)
						return true, line + "\n" + head + newFieldVal
					}
				}

				// *** 3. skos shadowing of hasChild/isChildOf *** (NOT IMPLEMENT YET)
				//
				if todo[2] {
					if strings.Contains(line, `"gem:hasChild"`) {
						// skos:narrower
					}
					if strings.Contains(line, `"gem:isChildOf"`) {
						// skos:broader
					}
				}

				// *** 4. *** (No issue found as 'Remove the "connects" structure from JSON and JSON-LD output')
				//
				if todo[3] {

				}

				// *** 5. ***
				if todo[4] {
					if strings.Contains(line, `"language"`) {
						return true, strings.Replace(line, `"language"`, `"@language"`, 1)
					}
					if strings.Contains(line, `"literal"`) {
						return true, strings.Replace(line, `"literal"`, `"@value"`, 1)
					}
					if strings.Contains(line, `"position"`) {
						return true, strings.Replace(line, `"position"`, `"@asn:listID"`, 1)
					}
				}
				//

				// *** 6. unify asn:hasLevel style ***
				if todo[5] {

					// above1, _, _ := strings.TrimSpace(cache[2]), strings.TrimSpace(cache[1]), strings.TrimSpace(cache[0])
					below1, _, _ := strings.TrimSpace(cache[4]), strings.TrimSpace(cache[5]), strings.TrimSpace(cache[6])

					if !met && strings.TrimSpace(line) == `"asn:hasLevel": {` {
						head = strs.TrimTailFromLast(line, `"asn:hasLevel"`)
						met = true
						return true, head + `"asn:hasLevel": [{`
					}
					if met && line == head+"}" {
						met = false
						return true, head + `}]`
					}

					if below1 == `"asn:hasLevel": []` {
						return true, strings.TrimSuffix(line, ",")
					}
					if strings.Contains(line, `"asn:hasLevel": []`) {
						return false, ""
					}
				}

				return true, line
			},
			filepath.Join(outDir, fName),
		)
	}

	// *** format output ***

	files, err = os.ReadDir(outDir)
	if err != nil {
		log.Fatalln(err)
		return
	}
	for _, file := range files {
		fName := file.Name()
		fPath := filepath.Join(outDir, fName)
		data, err := os.ReadFile(fPath)
		if err != nil {
			log.Panicln(err)
			return
		}
		if err := os.WriteFile(fPath, jt.Fmt(data, "    "), os.ModePerm); err != nil {
			fmt.Println(fPath)
			log.Panicln(err)
			return
		}
	}

	// *** check valid json ***

	files, err = os.ReadDir(outDir)
	if err != nil {
		log.Fatalln(err)
		return
	}
	for _, file := range files {
		fName := file.Name()
		fPath := filepath.Join(outDir, fName)
		data, err := os.ReadFile(fPath)
		if err != nil {
			log.Panicln(err)
			return
		}
		if !dt.IsJSON(data) {
			fmt.Println("2 --->", fPath)
			return
		}
	}
}
