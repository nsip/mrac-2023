package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	. "github.com/digisan/go-generics"
	"github.com/digisan/gotk/strs"
	"github.com/digisan/gotk/track"
	jt "github.com/digisan/json-tool"
	jts "github.com/digisan/json-tool/scan"
	"github.com/nsip/mrac-2023/node2"
	"github.com/tidwall/gjson"
)

var (
	mIdUrl           = loadUrl("../data/id-url.txt")
	js               = ""
	wholePaths       = []string{}
	la               = ""
	mPLUri           = make(map[string][]string)
	mData            = make(map[string]any)
	prevDocTypePath  = ""
	profLvl          = ""
	eduLvl           = ""
	mAsnCT           = LoadIdPrefLbl("../data/id-preflabel.txt")
	fileNode         = "../data/Sofia-API-Node-Data.json" // only for get mCodeChildParent
	mCodeChildParent map[string]string
	fileNodeMeta     = "../data/node-meta.json" // here, it is updated fileNode
	mNodeMeta        map[string]any
)

func main() {

	dataNodeMeta, err := os.ReadFile(fileNodeMeta)
	if err != nil {
		log.Fatalln(err)
	}
	// mUidTitle := scanNodeIdTitle(dataNodeMeta) // title should be node title
	mNodeMeta, err = jt.Flatten(dataNodeMeta)
	if err != nil {
		log.Fatalln(err)
	}

	dataFileNode, err := os.ReadFile(fileNode)
	if err != nil {
		log.Fatalf("%v", err)
	}
	mIdBlock := node2.GenNodeIdBlockMap(dataFileNode)
	_, mCodeChildParent = node2.GenChildParentMap(dataFileNode, mIdBlock)

	/////////////////////////////////////////////////////////////////////

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

	for file, LA := range mInputLa {

		fPath := filepath.Join("../data-out/restructure/", file) // "../data-out/restructure/la-Languages.json"
		// fOut := filepath.Join("./", filepath.Base(fPath))
		fOut := filepath.Join("../data-out/asn-json/", filepath.Base(fPath))
		la = LA

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

		_, paths, _, err := jts.AnalyzeJson(fPath)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("processing... %s with %d paths\n", fPath, len(paths))
		wholePaths = paths

		/////////////////////////////////////////////////////////////////////

		opt := jts.OptLineProc{
			Fn_KV:          proc_kv,          // nil
			Fn_KV_Str:      proc_kv_str,      // nil
			Fn_KV_Obj_Open: proc_kv_obj_open, // nil
			Fn_KV_Arr_Open: proc_kv_arr_open, // nil
			Fn_Obj:         proc_obj,         // nil
			Fn_Arr:         proc_arr,         // nil
			Fn_Elem:        proc_elem,        // nil
			Fn_Elem_Str:    proc_elem_str,    // nil
		}

		if err := jts.ScanJsonLine(fPath, fOut, opt); err != nil {
			log.Fatalln(err)
		}
	}
}

/////////////////////////////////////////////////////////////////////

func proc_kv(I int, k string, v any, lines []string, paths []string) (bool, string, bool) {

	path := paths[I]

	// remove 'tags' object's lines
	if strings.Contains(path, ".tags.") {
		return false, "", true
	}

	return true, fmt.Sprintf(`"%v": %v`, k, v), false
}

func proc_kv_str(I int, k string, v string, lines []string, paths []string) (bool, string, bool) {

	path := paths[I]

	if k == "uuid" {
		return true, fmt.Sprintf(`"id": "%s%s"`, mIdUrl[v], v), false // mIdUrl[value] already append with '/'
	}

	if k == "type" {
		return false, "", true
	}

	if k == "created_at" {
		return true, fmt.Sprintf(`"dcterms_modified": { "literal": "%v" }`, v), false
	}

	if k == "title" {

		// with 'text' sibling
		sibling := jt.NewSibling(path, "text")
		if In(sibling, wholePaths...) {
			return true, fmt.Sprintf(`"dcterms_title": { "language": "en-au", "literal": "%s" }`, v), false
		}

		// without 'text' sibling
		return true, fmt.Sprintf(`"dcterms_title": { "language": "en-au", "literal": "%s" }, "text": "%s"`, v, v), false
	}

	if k == "text" {
		return true, fmt.Sprintf(`"dcterms_description": { "language": "en-au", "literal": "%s" }, "text": "%s"`, v, v), false
	}

	if k == "position" {
		return true, fmt.Sprintf(`"asn_listID": "%v"`, v), false
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

		return true, kvStrJoin(staLbl, profLvl, eduLvl), false
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
		return true, strings.Join(rets, ","), false
	}

	// remove 'tags' object's lines
	if strings.Contains(path, ".tags.") {
		return false, "", true
	}

	return true, fmt.Sprintf(`"%v": "%v"`, k, v), false
}

func proc_kv_obj_open(I int, k string, v string, lines []string, paths []string) (bool, string, bool) {

	path := paths[I]

	// remove 'tags' object's lines
	if strings.Contains(path, ".tags.") {
		return false, "", true
	}

	// unwrap 'doc' object
	if k == "doc" {
		return false, "", true
	}

	// replace whole 'tags' object to below
	if k == "tags" {
		return true, fmt.Sprintf(`"asn_conceptKeyword": "%s"`, "SCIENCE_TEACHER_BACKGROUND_INFORMATION"), false
	}

	// unwrap 'connections' object
	if k == "connections" {
		return false, "", true
	}

	return true, fmt.Sprintf(`"%v": %v`, k, v), false
}

func proc_kv_arr_open(I int, k string, v string, lines []string, paths []string) (bool, string, bool) {

	path := paths[I]

	// remove 'tags' object's lines
	if strings.Contains(path, ".tags.") {
		return false, "", true
	}

	// remove 'connections' each tag line
	if strings.Contains(path, ".connections.") {
		if In(k,
			"Levels",
			"Organising Ideas",
			"Achievement Standard Components",
			"Indicator Groups",
			"Content Descriptions") {

			if r := gjson.Get(js, path); r.IsArray() {
				for i, rElem := range r.Array() {
					if i == 0 {
						v := rElem.Str
						id := v[strings.LastIndex(v, "/")+1:]
						code := jt.GetStrVal(mNodeMeta[id+"."+"code"])
						nodeType := GetCodeType(code, mCodeChildParent)
						if len(nodeType) == 0 {
							switch {
							case strings.HasPrefix(code, "AS"):
								nodeType = "AS"
							case strings.HasPrefix(code, "LA"):
								nodeType = "LA"
							}
						}
						switch nodeType {
						case "GC":
							return true, `"asn_skillEmbodied": [`, false
						case "LA":
							return true, `"dc_relation": [`, false
						case "AS":
							return true, `"asn_hasLevel": [`, false
						case "CCP":
							return true, `"asn_crossSubjectReference": [`, false
						default:
							log.Fatalf("nodeType '%v' is none of [GC CCP LA AS], code is '%v'", nodeType, code)
						}
					}
				}
			} else {
				log.Fatalln("connections.xxx value should be array")
			}

			return false, "", true
		}
	}

	return true, fmt.Sprintf(`"%v": %v`, k, v), false
}

func proc_obj(I int, v string, lines []string, paths []string) (bool, string, bool) {

	path := paths[I]

	// remove doc '}' and add comma if necessary
	if strings.HasSuffix(path, ".doc}") {
		return true, " ", false // non-empty space, means let outer makes comma if needed
	}

	// remove 'tags' object's lines
	if strings.Contains(path, ".tags.") {
		return false, "", true
	}
	if strings.HasSuffix(path, ".tags}") {
		return true, " ", false // non-empty space, means let outer makes comma if needed
	}

	// remove connections '}' and add comma if necessary
	if strings.HasSuffix(path, ".connections}") {
		return true, " ", false // non-empty space, means let outer makes comma if needed
	}

	return true, v, false
}

func proc_arr(I int, v string, lines []string, paths []string) (bool, string, bool) {

	path := paths[I]

	// remove 'tags' object's lines
	if strings.Contains(path, ".tags.") {
		return false, "", true
	}

	// keep 'connections.xxx' end ']'

	return true, v, false
}

func proc_elem(I int, v any, lines []string, paths []string) (bool, string, bool) {

	path := paths[I]

	// remove 'tags' object's lines
	if strings.Contains(path, ".tags.") {
		return false, "", true
	}

	return true, v.(string), false
}

func proc_elem_str(I int, v string, lines []string, paths []string) (bool, string, bool) {

	path := paths[I]

	// remove 'tags' object's lines
	if strings.Contains(path, ".tags.") {
		return false, "", true
	}

	// process 'connections.xxx' each element
	if strs.ContainsAny(path,
		"connections.Levels.",
		"connections.Organising Ideas.",
		"connections.Achievement Standard Components.",
		"connections.Indicator Groups.",
		"connections.Content Descriptions.") {

		id := v[strings.LastIndex(v, "/")+1:]
		code := jt.GetStrVal(mNodeMeta[id+"."+"code"])
		title := jt.GetStrVal(mNodeMeta[id+"."+"title"])
		nodeType := GetCodeType(code, mCodeChildParent)
		if len(nodeType) == 0 {
			switch {
			case strings.HasPrefix(code, "AS"):
				nodeType = "AS"
			case strings.HasPrefix(code, "LA"):
				nodeType = "LA"
			}
		}

		switch nodeType {
		case "GC", "CCP":
			return true, fmt.Sprintf(`{ "uri": "%s", "prefLabel": "%s" }`, v, code), false
		default:
			return true, fmt.Sprintf(`{ "uri": "%s", "prefLabel": "%s" }`, v, title), false
		}
	}

	return true, fmt.Sprintf(`"%v"`, v), false
}
