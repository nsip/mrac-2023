package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	. "github.com/digisan/go-generics"
	"github.com/digisan/gotk/track"
	jt "github.com/digisan/json-tool"
	jts "github.com/digisan/json-tool/scan"
)

var (
	mIdUrl          = loadUrl("../data/id-url.txt")
	js              = ""
	wholePaths      = []string{}
	la              = ""
	mPLUri          = make(map[string][]string)
	mData           = make(map[string]any)
	prevDocTypePath = ""
	profLvl         = ""
	eduLvl          = ""
	mAsnCT          = LoadIdPrefLbl("../data/id-preflabel.txt")
)

func main() {

	mInputLa := map[string]string{
		"la-Languages.json":                      "Languages",
		"la-English.json":                        "English",
		"la-Humanities and Social Sciences.json": "Humanities and Social Sciences",
		"la-Health and Physical Education.json":  "Health and Physical Education",
		"la-Mathematics.json":                    "Mathematics",
		"la-Science.json":                        "Science",
		"la-Technologies.json":                   "Technologies",
		"la-The Arts.json":                       "The Arts",
		"ccp-Cross-curriculum Priorities.json":   "CCP",
		"gc-Critical and Creative Thinking.json": "GC-CCT",
		"gc-Digital Literacy.json":               "GC-DL",
		"gc-Ethical Understanding.json":          "GC-EU",
		"gc-Intercultural Understanding.json":    "GC-IU",
		"gc-Literacy.json":                       "GC-L",
		"gc-Numeracy.json":                       "GC-N",
		"gc-Personal and Social capability.json": "GC-PSC",
	}

	/////////////////////////////////////////////////////////////////////

	fPath := "../data-out/restructure/la-English.json"
	fOut := filepath.Join("./", filepath.Base(fPath))

	La, ok := mInputLa[filepath.Base(fPath)]
	if !ok {
		log.Fatalln("missing la from mInputLa")
	}
	la = La

	/////////////////////////////////////////////////////////////////////

	data, err := os.ReadFile(fPath)
	if err != nil {
		log.Fatalln(err)
	}
	js = string(data)

	mData, err = jt.Flatten(data)
	if err != nil {
		log.Fatalln(err)
	}

	/////////////////////////////////////////////////////////////////////

	progLvlABC := "" // indicate Level 1a, 1b or 1c
	switch {
	case strings.Contains(js, `"Level 1c"`):
		progLvlABC = "1c"
	case strings.Contains(js, `"Level 1b"`):
		progLvlABC = "1b"
	case strings.Contains(js, `"Level 1a"`):
		progLvlABC = "1a"
	}

	switch progLvlABC {
	case "1c":
		mPLUri = MapMerge(mProglvlUri, mProglvlABCUri)
	case "1b":
		mPLUri = MapMerge(mProglvlUri, mProglvlABUri)
	case "1a":
		mPLUri = MapMerge(mProglvlUri, mProglvlAUri)
	default:
		mPLUri = MapMerge(mProglvlUri)
	}

	/////////////////////////////////////////////////////////////////////

	defer track.TrackTime(time.Now())

	paths, err := jts.ScanJsonLine(fPath, fOut, jts.OptLineProc{})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("processing... %s with %d paths\n", fPath, len(paths))
	wholePaths = paths

	/////////////////////////////////////////////////////////////////////

	opt := jts.OptLineProc{
		Fn_KV:          nil,              // proc_kv
		Fn_KV_Str:      proc_kv_str,      // nil
		Fn_KV_Obj_Open: proc_kv_obj_open, // nil
		Fn_KV_Arr_Open: nil,              // proc_kv_arr_open
		Fn_Obj:         proc_obj,         // nil
		Fn_Arr:         nil,              // proc_arr
		Fn_Elem:        nil,              // proc_elem
		Fn_Elem_Str:    nil,              // proc_elem_str
	}

	_, err = jts.ScanJsonLine(fPath, fOut, opt)
	if err != nil {
		log.Fatalln(err)
	}
}

/////////////////////////////////////////////////////////////////////

func proc_kv(I int, path, k string, v any) (bool, string) {
	return true, fmt.Sprintf(`"%v": %v`, k, v)
}

func proc_kv_str(I int, path, k, v string) (bool, string) {

	if k == "uuid" {
		return true, fmt.Sprintf(`"id": "%s%s"`, mIdUrl[v], v) // mIdUrl[value] already append with '/'
	}

	if k == "type" {
		return false, ""
	}

	if k == "created_at" {
		return true, fmt.Sprintf(`"dcterms_modified": { "literal": "%v" }`, v)
	}

	if k == "title" {

		// with 'text' sibling
		sibling := jt.NewSibling(path, "text")
		if In(sibling, wholePaths...) {
			return true, fmt.Sprintf(`"dcterms_title": { "language": "en-au", "literal": "%s" }`, v)
		}

		// without 'text' sibling
		return true, fmt.Sprintf(`"dcterms_title": { "language": "en-au", "literal": "%s" }, "text": "%s"`, v, v)
	}

	if k == "text" {
		return true, fmt.Sprintf(`"dcterms_description": { "language": "en-au", "literal": "%s" }, "text": "%s"`, v, v)
	}

	if k == "position" {
		return true, fmt.Sprintf(`"asn_listID": "%v"`, v)
	}

	if k == "typeName" {

		// "asn_statementLabel"
		staLbl := fmt.Sprintf(`"asn_statementLabel": { "language": "%s", "literal": "%s" }`, "en-au", v)

		// “asn_proficiencyLevel”
		if strings.HasPrefix(la, "GC-") {
			if v == "Level" {
				lvl := getProLevel(mData, path)
				uri := mPLUri[lvl][0] // only one uri
				profLvl = fmt.Sprintf(`"asn_proficiencyLevel": { "uri": "%s", "prefLabel": "%s" }`, uri, lvl)
				prevDocTypePath = path
			}
			// only children path can keep retEL
			if strings.Count(path, ".") < strings.Count(prevDocTypePath, ".") {
				profLvl = ""
			}
		} else {
			profLvl = ""
		}

		// "dcterms_educationLevel"
		if NotIn(la, "CCP", "GC-L", "GC-N") {
			if v == "Level" { // see doc.typeName: 'Level', update global retEL
				outArrs := []string{}
				for _, y := range getYears(mData, path) {
					outArrs = append(outArrs, fmt.Sprintf(`{ "uri": "%s", "prefLabel": "%s" }`, mYrlvlUri[y], y))
				}
				if len(outArrs) > 0 {
					eduLvl = strings.Join(outArrs, ",")
				}
				eduLvl = fmt.Sprintf(`"dcterms_educationLevel": [%s]`, eduLvl)
				prevDocTypePath = path
			}
			// only children path can keep retEL
			if strings.Count(path, ".") < strings.Count(prevDocTypePath, ".") {
				eduLvl = ""
			}
		} else {
			eduLvl = ""
		}

		return true, kvStrJoin(staLbl, profLvl, eduLvl)
	}

	if k == "code" {

		retSN := fmt.Sprintf(`"asn_statementNotation": { "language": "%s", "literal": "%s" }`, "en-au", v)

		retAS := fmt.Sprintf(`"asn_authorityStatus": { "uri": "%s" }`, `http://purl.org/ASN/scheme/ASNAuthorityStatus/Original`)

		retIS := fmt.Sprintf(`"asn_indexingStatus": { "uri": "%s" }`, `http://purl.org/ASN/scheme/ASNIndexingStatus/No`)

		retCT := ""
		if conceptTerm, ok := mAsnCT[v]; ok {
			retCT = fmt.Sprintf(`"asn_conceptTerm": %s`, conceptTerm)
		}

		// retTxt := ""
		// if !gjson.Get(js, jt.NewSibling(path, "text")).Exists() {
		// 	retTxt = `"text": null`
		// }

		retSub := ``
		if In(v, "ENG", "HAS", "HPE", "LAN", "MAT", "SCI", "TEC", "ART") {
			retS := []string{}
			if subUri, okSubUri := mLaUri[la]; okSubUri {
				retS = append(retS, fmt.Sprintf(`"dcterms_subject": { "prefLabel": "%s", "uri": "%s" }`, la, subUri))
			}
			retSub = strings.Join(retS, ",")
		}

		retRT, retRTH := ``, ``
		if In(v, "root", "LA") {
			retRT = fmt.Sprintf(`"dcterms_rights": { "language": "%s", "literal": "%s" }`, "en-au", `©Copyright Australian Curriculum, Assessment and Reporting Authority`)
			retRTH = fmt.Sprintf(`"dcterms_rightsHolder": { "language": "%s", "literal": "%s" }`, "en-au", `Australian Curriculum, Assessment and Reporting Authority`)
		}

		retCLS, retLEAF := ``, ``

		sibling := jt.NewSibling(path, "children")
		if In(sibling, wholePaths...) {
			retCLS = `"cls": "folder"`
		} else {
			retLEAF = `"leaf": "true"`
		}

		rets := []string{}
		for _, r := range []string{retSN, retAS, retCT, retIS, retSub, retRT, retRTH, retCLS, retLEAF} {
			if r != "" {
				rets = append(rets, r)
			}
		}
		return true, strings.Join(rets, ",")
	}

	if k == "tag" {
		return true, fmt.Sprintf(`"asn_conceptKeyword": "%s"`, "SCIENCE_TEACHER_BACKGROUND_INFORMATION")
	}

	// if In(k, "Levels", "OI", "ASC", "IG", "CD") {
	// }

	return true, fmt.Sprintf(`"%v": "%v"`, k, v)
}

func proc_kv_obj_open(I int, path, k, v string) (bool, string) {

	// unwrap 'doc' object
	if k == "doc" {
		fmt.Println("--->", path, k, v)
		return false, ""
	}

	return true, fmt.Sprintf(`"%v": %v`, k, v)
}

func proc_kv_arr_open(I int, path, k, v string) (bool, string) {
	return true, fmt.Sprintf(`"%v": %v`, k, v)
}

func proc_obj(I int, path, v string) (bool, string) {

	// remove doc '}' and add comma if necessary
	if strings.HasSuffix(path, ".doc}") {
		fmt.Println("===>", path, v)
		return true, " " // non-empty space, means let outer makes comma if needed
	}

	return true, v
}

func proc_arr(I int, path, v string) (bool, string) {
	return true, v
}

func proc_elem(I int, path string, v any) (bool, string) {
	return true, v.(string)
}

func proc_elem_str(I int, path, v string) (bool, string) {
	return true, fmt.Sprintf(`"%v"`, v)
}
