package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	fd "github.com/digisan/gotk/file-dir"
	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
	"github.com/tidwall/gjson"
)

func getAcScotMap(acscotPath string) map[string][]string {
	m := make(map[string][]string)
	_, err := fd.FileLineScan(acscotPath, func(line string) (bool, string) {
		if strings.HasPrefix(line, "AC") {
			k, v := "", ""
			for _, s := range strings.Split(line, "\t") {
				s = strings.Trim(s, " \t\r\n")
				switch {
				case strings.HasPrefix(s, "AC"):
					k = s
				case strings.HasPrefix(s, "http"):
					v = s
				}
			}
			lk.FailOnErrWhen(v == "", "%v @ %s", fmt.Errorf("error in scan http"), line)
			m[k] = append(m[k], v)
		}
		return false, ""
	}, "")
	lk.FailOnErr("%v", err)
	return m
}

// scotJsonLd: http://vocabulary.curriculum.edu.au/scot/export/scot.jsonld
func scanScotJsonLd(scotJsonLdPath string) map[string][]string {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	file, err := os.Open(scotJsonLdPath)
	lk.FailOnErr("%v", err)

	cOut, cErr, err := jt.ScanObjectInArray(ctx, file, true)
	lk.FailOnErr("%v", err)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	m := make(map[string][]string)

	go func(m map[string][]string) {
		for out := range cOut {

			// according to Nick guidance
			// "http://www.w3.org/2004/02/skos/core#prefLabel" + "@language": "en", => "@value"

			idVal := gjson.Get(out, "@id").String()
			idVal = strings.Trim(idVal, " \t\r\n")

			mObj := make(map[string]any)
			lk.FailOnErr("%v", json.Unmarshal([]byte(out), &mObj))

			for k, v := range mObj {
				switch {
				case k == "@id":
					lk.FailOnErrWhen(idVal != v.(string), "%v", fmt.Errorf("error in fetching prefLabel"))

				case strings.HasSuffix(k, "#prefLabel"):
					for _, mPrefLabel := range v.([]any) {
						mpl := mPrefLabel.(map[string]any)
						if langVal, ok := mpl["@language"]; ok && langVal == "en" {
							m[idVal] = append(m[idVal], mpl["@value"].(string)) // fetch @value under @language = "en"
						}
					}
				}
			}
		}
		wg.Done()
	}(m)

	go func() {
		for err := range cErr {
			lk.FailOnErr("%v", err)
		}
		wg.Done()
	}()

	wg.Wait()

	return m
}

func GetAsnConceptTerm(acscotPath, scotJsonLdPath string) map[string]string {

	m := make(map[string]string)

	m1 := getAcScotMap(acscotPath)
	m2 := scanScotJsonLd(scotJsonLdPath)

	for code, scotUris := range m1 {
		// fmt.Println(code)
		asnConceptTerm := "["
		for _, scotUri := range scotUris {
			// fmt.Println(scotUri)
			gotM2 := false
			for _, v := range m2[scotUri] {
				// fmt.Println(v)
				asnConceptTerm += fmt.Sprintf(`{ "uri": "%s", "prefLabel": "%s"	}`, scotUri, v)
				gotM2 = true
			}
			if gotM2 {
				asnConceptTerm += ","
			}
		}
		asnConceptTerm = strings.TrimSuffix(asnConceptTerm, ",")
		asnConceptTerm += "]"

		m[code] = asnConceptTerm

		// fmt.Println(asnConceptTerm)
		// break
	}

	return m
}
