package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	. "github.com/digisan/go-generics/v2"
	jt "github.com/digisan/json-tool"
	"github.com/nsip/mrac-2023/asn-json/tool"
	"github.com/tidwall/gjson"
)

var (
	fSf        = fmt.Sprintf
	sJoin      = strings.Join
	sTrim      = strings.Trim
	sLastIndex = strings.LastIndex
	sHasPrefix = strings.HasPrefix
	sHasSuffix = strings.HasSuffix
	sSplit     = strings.Split
)

var (
	mRES = map[string]string{
		"text":               `"text":\s*"[^"]+",?`,
		"uuid":               `"uuid":\s*"[\d\w]{8}-[\d\w]{4}-[\d\w]{4}-[\d\w]{4}-[\d\w]{12}",?`,
		"type":               `"type":\s*"\w+",?`,
		"created_at":         `"created_at":\s*"[^"]+",?`,
		"title":              `"title":\s*"[^"]+",?`,
		"position":           `"position":\s*"[^"]+",?`,
		"doc.typeName":       `"doc":\s*\{[^{}]+\},?`,
		"code":               `"code":\s*"[^"]+",?`,
		"tag":                `"tags":\s*\{[^{}]+\},?`,
		"connections.Levels": `"Levels":\s*\[[^\[\]]+\],?`,
		"connections.OI":     `"Organising Ideas":\s*\[[^\[\]]+\],?`,
		"connections.ASC":    `"Achievement Standard Components":\s*\[[^\[\]]+\],?`,
		"connections.IG":     `"Indicator Groups":\s*\[[^\[\]]+\],?`,
		"connections.CD":     `"Content Descriptions":\s*\[[^\[\]]+\],?`,
	}

	mAsnCT = tool.GetAsnConceptTerm("../data/ACv9_ScOT_BC_20220422.txt", "../data/scot.jsonld")

	reMerged = func() (*regexp.Regexp, map[string]*regexp.Regexp) {
		mRE4Each := map[string]*regexp.Regexp{}
		restr := ""
		for k, v := range mRES {
			restr += fmt.Sprintf("(%s)|", v)    // merged restr for whole regexp
			mRE4Each[k] = regexp.MustCompile(v) // init each regexp
		}
		// remove last '|' and compile to regexp
		return regexp.MustCompile(restr[:len(restr)-1]), mRE4Each
	}
)

func fnGetPathByProp(prop string, paths []string, info string) func() string {
	var idx int64 = -1
	paths = jt.GetLeafPathsOrderly(prop, paths)
	return func() string {
		idx++
		// fmt.Println(prop, idx, info)
		return paths[idx]
	}
}

func kvstrJoin(kvStrGrp ...string) string {
	nonEmptyStrGrp := []string{}
	for _, kvstr := range kvStrGrp {
		if sTrim(kvstr, " \t\n") != "" {
			nonEmptyStrGrp = append(nonEmptyStrGrp, kvstr)
		}
	}
	return sJoin(nonEmptyStrGrp, ",")
}

func proc(

	js, s, name, value string,
	mLvlSiblings map[int][]string,
	mData map[string]interface{},
	la, uri4id string,
	mCodeParent map[string]string,
	mNodeData map[string]interface{},
	fnPathWithTitle func() string,
	fnPathWithDocType func() string,
	fnPathWithCode func() string,
	// static for filling
	pPrevDocTypePath *string,
	pRetEL *string,
	pRetPL *string,
	// mProglvlUri Extra with 1a, 1b, 1c
	mPLUri map[string][]string,

) (bool, string) {

	// if name == "doc.typeName" {
	// 	fmt.Println(name, s)
	// }

	switch name {
	case "uuid":
		return true, fSf(`"id": "%s/%s"`, uri4id, value)

	case "type":
		return true, ""

	case "created_at":
		return true, fSf(`"dcterms_modified": { "literal": "%s" }`, value)

	case "title":
		path := fnPathWithTitle()
		// with 'text' sibling
		if gjson.Get(js, jt.NewSibling(path, "text")).Exists() {
			return true, fSf(`"dcterms_title": { "language": "en-au", "literal": "%s" }`, value)
		}
		// no 'text' sibling
		return true, fSf(`"dcterms_title": { "language": "en-au", "literal": "%s" }, "text": "%s"`, value, value)

	case "text":
		return true, fSf(`"dcterms_description": { "language": "en-au", "literal": "%s" }, "text": "%s"`, value, value)

	case "position":
		return true, fSf(`"asn_listID": "%s"`, value)

	case "doc.typeName":

		path := fnPathWithDocType()

		// "asn_statementLabel"
		retSL := fSf(`"asn_statementLabel": { "language": "%s", "literal": "%s" }`, "en-au", value)

		// “asn_proficiencyLevel”
		if sHasPrefix(la, "GC-") {
			if value == "Level" {
				lvl := getProLevel(mData, path)
				uri := mPLUri[lvl][0] // only one uri
				*pRetPL = fSf(`"asn_proficiencyLevel": { "uri": "%s", "prefLabel": "%s" }`, uri, lvl)
				*pPrevDocTypePath = path
			}
			// only children path can keep retEL
			if strings.Count(path, ".") < strings.Count(*pPrevDocTypePath, ".") {
				*pRetPL = ""
			}
		} else {
			*pRetPL = ""
		}

		// "dcterms_educationLevel"
		if NotIn(la, "CCP", "GC-L", "GC-N") {
			if value == "Level" { // see doc.typeName: 'Level', update global retEL
				outArrs := []string{}
				for _, y := range getYears(mData, path) {
					outArrs = append(outArrs, fSf(`{ "uri": "%s", "prefLabel": "%s" }`, mYrlvlUri[y], y))
				}
				if len(outArrs) > 0 {
					*pRetEL = sJoin(outArrs, ",")
				}
				*pRetEL = fSf(`"dcterms_educationLevel": [%s]`, *pRetEL)
				*pPrevDocTypePath = path
			}
			// only children path can keep retEL
			if strings.Count(path, ".") < strings.Count(*pPrevDocTypePath, ".") {
				*pRetEL = ""
			}
		} else {
			*pRetEL = ""
		}

		return true, kvstrJoin(retSL, *pRetPL, *pRetEL)

	case "code":

		path := fnPathWithCode()

		retSN := fSf(`"asn_statementNotation": { "language": "%s", "literal": "%s" }`, "en-au", value)

		retAS := fSf(`"asn_authorityStatus": { "uri": "%s" }`, `http://purl.org/ASN/scheme/ASNAuthorityStatus/Original`)

		retIS := fSf(`"asn_indexingStatus": { "uri": "%s" }`, `http://purl.org/ASN/scheme/ASNIndexingStatus/No`)

		retCT := ""
		if conceptTerm, ok := mAsnCT[value]; ok {
			retCT = fSf(`"asn_conceptTerm": %s`, conceptTerm)
		}

		// retTxt := ""
		// if !gjson.Get(js, jt.NewSibling(path, "text")).Exists() {
		// 	retTxt = `"text": null`
		// }

		retSub := ``
		if In(value, "ENG", "HAS", "HPE", "LAN", "MAT", "SCI", "TEC", "ART") {
			retS := []string{}
			if subUri, okSubUri := mLaUri[la]; okSubUri {
				retS = append(retS, fSf(`"dcterms_subject": { "prefLabel": "%s", "uri": "%s" }`, la, subUri))
			}
			retSub = sJoin(retS, ",")
		}

		retRT, retRTH := ``, ``
		if In(value, "root", "LA") {
			retRT = fSf(`"dcterms_rights": { "language": "%s", "literal": "%s" }`, "en-au", `©Copyright Australian Curriculum, Assessment and Reporting Authority`)
			retRTH = fSf(`"dcterms_rightsHolder": { "language": "%s", "literal": "%s" }`, "en-au", `Australian Curriculum, Assessment and Reporting Authority`)
		}

		retCLS, retLEAF := ``, ``
		if jt.HasSiblings(path, mLvlSiblings, "children") {
			retCLS = fSf(`"cls": "folder"`)
		} else {
			retLEAF = fSf(`"leaf": "true"`)
		}

		rets := []string{}
		for _, r := range []string{retSN, retAS, retCT, retIS, retSub, retRT, retRTH, retCLS, retLEAF} {
			if r != "" {
				rets = append(rets, r)
			}
		}
		return true, sJoin(rets, ",")

	case "tag":
		return true, fSf(`"asn_conceptTerm": "%s"`, "SCIENCE_TEACHER_BACKGROUND_INFORMATION")

	case "connections.Levels",
		"connections.OI",
		"connections.ASC",
		"connections.IG",
		"connections.CD":

		items := sSplit(value, "|")
		// fmt.Println(items)

		code := ""
		nodeType := ""
		outArrs := []string{}
		for _, item := range items {
			id := item[sLastIndex(item, "/")+1:]
			code = jt.GetStrVal(mNodeData[id+"."+"code"])
			title := jt.GetStrVal(mNodeData[id+"."+"title"])
			nodeType = tool.GetCodeAncestor(mCodeParent, code, 0)
			switch nodeType {
			case "GC", "CCP":
				outArrs = append(outArrs, fSf(`{ "uri": "%s", "prefLabel": "%s" }`, item, code))
			default:
				outArrs = append(outArrs, fSf(`{ "uri": "%s", "prefLabel": "%s" }`, item, title))
			}
		}

		outArrStr := sJoin(outArrs, ",")
		ret := ""

		switch nodeType {
		case "GC":
			ret = fSf(`"asn_skillEmbodied": [%s]`, outArrStr)
		case "LA":
			ret = fSf(`"dc_relation": [%s]`, outArrStr)
		case "AS":
			ret = fSf(`"asn_hasLevel": [%s]`, outArrStr)
		case "CCP":
			ret = fSf(`"asn_crossSubjectReference": [%s]`, outArrStr)
		default:
			log.Fatalf("nodeType '%v' is none of [GC CCP LA AS], code is '%v'", nodeType, code)
		}

		return true, ret

	default:
		return false, ""
	}
}

func treeProc3(

	data []byte,
	la string,
	mCodeParent map[string]string,
	mNodeData map[string]interface{},
	paths []string,
	// static for filling
	pPrevDocTypePath *string,
	pRetEL *string,
	pRetPL *string,
	// indicator for mProglvlUri Extra 1a, 1b, 1c
	progLvlABC string,

) string {

	mPLUri := make(map[string][]string)
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

	js := string(data)
	mLvlSiblings, _ := jt.FamilyTree(js)

	mData, err := jt.Flatten(data)
	if err != nil {
		log.Fatalln(err)
	}

	re4json, mRE4Each := reMerged()
	// fmt.Println(re4json, len(mRE4Each))

	getPathWithTitle := fnGetPathByProp("title", paths, "")
	getPathWithTypeName := fnGetPathByProp("typeName", paths, "")
	getPathWithCode := fnGetPathByProp("code", paths, "")

	js = re4json.ReplaceAllStringFunc(js, func(s string) string {

		hasComma := false
		if sHasSuffix(s, ",") {
			hasComma = true
		}

		for name, v := range mRE4Each {
			if v.MatchString(s) {

				// if name == "doc.typeName" {
				// 	fmt.Println(name, s)
				// }

				if ok, repl := proc(
					js,
					s,
					name,
					tool.FetchValue(s, "|"),
					mLvlSiblings,
					mData,
					la,
					uri4id,
					mCodeParent,
					mNodeData,
					getPathWithTitle,
					getPathWithTypeName,
					getPathWithCode,
					pPrevDocTypePath,
					pRetEL,
					pRetPL,
					mPLUri,
				); ok {
					if hasComma && repl != "" {
						return repl + ","
					}
					return repl
				}
			}
		}
		return s
	})

	// *** further process ***

	// remove "connections" wrapper
	_, _, mPropLocs, mPropValues := jt.GetProperties(js)
	js = jt.RemoveParent(js, "connections", mPropLocs, mPropValues)

	return js
}
