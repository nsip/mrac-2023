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

				nSpace := 0
				for _, c := range line {
					if c == ' ' {
						// fmt.Println(i)
						// break
						nSpace++
					}
				}

				fmt.Println(nSpace)
				// fmt.Println(len(line) - len(ln))
				fmt.Println(ln)
				fmt.Println(line)

			} //

		}
		return true, ""

	}); err == nil {
		return rt
	}

	return ""
}

func main() {

	const (
		inputDir = "../release/asn-json-ld/"
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

			insertDescriptionIfHasTitle(string(data))

			fmt.Printf("processing... %s\n", fPath)

			break

			// err = os.WriteFile(fPath, []byte(rt), os.ModePerm)
			// lk.FailOnErr("%v", err)
		}
	}
}
