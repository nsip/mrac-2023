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

func insertDescriptionIfHasTitle(js string) string {

	js, err := jt.FmtJS(js)
	lk.FailOnErr("%v", err)

	if rt, err := strs.StrLineScanEx(js, 3, 3, "JUNK", func(line string, cache []string) (bool, string) {

		if ln := strings.TrimSpace(line); strings.HasPrefix(ln, `"dc:title":`) {

			found := false
			for i, cl := range cache {
				if i == 3 {
					continue
				}
				if ln := strings.TrimSpace(cl); strings.HasPrefix(ln, `"dc:description":`) {
					found = true
				}
			}

			if !found {
				// fmt.Println(line)

				nSpace := 0
				for i, c := range line {
					if c != ' ' {
						fmt.Println(i)
						nSpace = i
						break
					}
				}

				descContent := strings.TrimSpace(strings.TrimPrefix(ln, `"dc:title":`)) // content contains double quotes
				descLine := fmt.Sprintf(`%s"dc:description": %s`, strings.Repeat(" ", nSpace), descContent)
				// fmt.Println(descLine)

				insert := strings.Join([]string{line, descLine}, "\n")
				return true, insert

			} //

		}
		return true, line

	}); err == nil {
		return rt
	}

	return ""
}

func main() {

	const (
		// inputDir = "../release/asn-json-ld/"
		inputDir = "../release/json-ld-ccp/"
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

			fmt.Println("---", len(data))

			rt := insertDescriptionIfHasTitle(string(data))

			fmt.Printf("processed... %s\n", fPath)

			err = os.WriteFile(fPath, []byte(rt), os.ModePerm)
			lk.FailOnErr("%v", err)			
		}
	}
}
