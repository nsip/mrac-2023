package tool

import (
	"encoding/json"
	"fmt"
	"strings"

	dt "github.com/digisan/gotk/data-type"
	fd "github.com/digisan/gotk/file-dir"
	"github.com/digisan/gotk/strs"
	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
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

//////////////////////////////////////////////////////////////////////////////////////

func scanSCOT(scotPath string) map[string][]string {

	js, err := jt.FmtFileJS(scotPath)
	lk.FailOnErr("%v", err)

	// fmt.Println("js length:", len(js))

	// @ 2space indents
	const hs_aboveOCB = "      "
	const hs_belowCCB = hs_aboveOCB
	const hs_prefLabel = "        " // pre-running scan to check 'strings.Contains(line, "#prefLabel\":")'

	var mBlock = make(map[string]string)

	// adjust below 4096 to search '}'
	const JUNK = "JUNK"
	strs.StrLineScanEx(js, 4, 4096, JUNK, func(line string, cache []string) (bool, string) {

		// if strings.Contains(line, "#prefLabel\":") {
		// 	hs := strs.HeadSpace(line)
		// 	fmt.Println(line, len(hs))
		// 	if hs != hs_prefLabel {
		// 		panic("FAILED")
		// 	}
		// 	// above := cache[:8]
		// 	// below := cache[9:]
		// 	// for i, bl := range below {
		// 	// 	if bl != JUNK && strings.Contains(bl, "@language") && strings.Contains(bl, "en") {
		// 	// 		fmt.Println(below[i+1])
		// 	// 	}
		// 	// }
		// }

		if strings.Contains(line, "@id") {

			above, below := cache[:4], cache[5:]
			getStart, getEnd := false, false
			pEnd := 0

			al1 := above[len(above)-1]
			if al1 != JUNK && al1 == hs_aboveOCB+"{" {
				if strs.HeadSpace(line) == hs_prefLabel {
					// fmt.Println(line)
					getStart = true
					for i, bl := range below {
						if bl != JUNK && (bl == hs_belowCCB+"}" || bl == hs_belowCCB+"},") {
							// fmt.Printf("%02d%s\n", i, bl)
							getEnd = true
							pEnd = i
							break
						}
						if i == len(below)-1 {
							panic("CANNOT FIND")
						}
					}
				}
			}

			if getStart && getEnd {
				blocks := append([]string{al1, line}, below[:pEnd+1]...)
				block := strings.TrimSuffix(strings.Join(blocks, "\n"), ",")
				lk.FailOnErrWhen(!dt.IsJSON([]byte(block)), "%v", fmt.Errorf("JSON ERROR"))
				mBlock[line] = block
			}
		}

		return true, ""
	})

	// fmt.Println("mBlock length:", len(mBlock))

	rt := make(map[string][]string)
	for id, block := range mBlock {

		obj := make(map[string]any)
		lk.FailOnErr("%v", json.Unmarshal([]byte(block), &obj))

		for field, value := range obj {

			if strings.HasSuffix(field, "#prefLabel") {

				lsPrefLabelValue, ok := value.([]any)
				lk.FailOnErrWhen(!ok, "%v", "#prefLabel value ERROR 1")

				for _, m := range lsPrefLabelValue {

					mpl, ok := m.(map[string]any)
					lk.FailOnErrWhen(!ok, "%v", "#prefLabel value ERROR 2")

					if langVal, ok := mpl["@language"]; ok && langVal == "en" {

						id = strings.TrimSpace(id)
						id = strings.TrimPrefix(id, `"@id":`)
						id = strings.TrimSpace(id)
						id = strings.TrimSuffix(id, ",")
						id = strings.Trim(id, `"`)

						rt[id] = append(rt[id], mpl["@value"].(string))
					}
				}
			}
		}
	}
	return rt
}

//////////////////////////////////////////////////////////////////////////////////////

// scotJsonLd: http://vocabulary.curriculum.edu.au/scot/export/scot.jsonld
// func scanScotJsonLd(scotJsonLdPath string) map[string][]string {

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	file, err := os.Open(scotJsonLdPath)
// 	lk.FailOnErr("%v", err)

// 	cOut, cErr, err := jt.ScanObjectInArray(ctx, file, true)
// 	lk.FailOnErr("%v", err)

// 	wg := &sync.WaitGroup{}
// 	wg.Add(2)

// 	m := make(map[string][]string)
// 	go func(m map[string][]string) {
// 		for out := range cOut { // { "@id", "@graph" }

// 			fmt.Println("test")

// 			// according to Nick guidance
// 			// "http://www.w3.org/2004/02/skos/core#prefLabel" + "@language": "en", => "@value"

// 			idVal := gjson.Get(out, "@id").String()
// 			idVal = strings.TrimSpace(idVal) // e.g. "http://vocabulary.curriculum.edu.au/scot/linkeddata/dbpedia/en"

// 			mObj := make(map[string]any)
// 			lk.FailOnErr("%v", json.Unmarshal([]byte(out), &mObj))

// 			for k, v := range mObj {
// 				switch {
// 				case k == "@id":
// 					lk.FailOnErrWhen(idVal != v.(string), "%v", fmt.Errorf("error in fetching prefLabel"))

// 				case strings.HasSuffix(k, "#prefLabel"):
// 					for _, mPrefLabel := range v.([]any) {
// 						mpl := mPrefLabel.(map[string]any)
// 						if langVal, ok := mpl["@language"]; ok && langVal == "en" {
// 							m[idVal] = append(m[idVal], mpl["@value"].(string)) // fetch @value under @language = "en"
// 						}
// 					}
// 				}
// 			}
// 		}
// 		wg.Done()
// 	}(m)

// 	go func() {
// 		for err := range cErr {
// 			lk.FailOnErr("%v", err)
// 		}
// 		wg.Done()
// 	}()

// 	wg.Wait()

// 	return m
// }

func GetAsnConceptTerm(acscotPath, scotJsonLdPath string) map[string]string {

	m := make(map[string]string)

	m1 := getAcScotMap(acscotPath)
	m2 := scanSCOT(scotJsonLdPath) // scanScotJsonLd(scotJsonLdPath)

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

func LoadIdPrefLbl(fpath string) map[string]string {

	lk.FailOnErrWhen(!fd.FileExists(fpath), "%v", fmt.Errorf("id-preflabel.txt (%s) [made from scot.txt & .jsonld] doesn't exist", fpath))

	rt := make(map[string]string)
	fd.FileLineScan(fpath, func(line string) (bool, string) {
		kv := strings.SplitN(line, "\t", 2)
		rt[kv[0]] = kv[1]
		return true, ""
	}, "")
	return rt
}
