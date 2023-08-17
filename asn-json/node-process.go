package main

import (
	"bytes"
	"path/filepath"

	// "context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	jt "github.com/digisan/json-tool"
	"github.com/nsip/mrac-2023/asn-json/tool"
	"github.com/tidwall/gjson"
)

type asnjson struct {
	Id string `json:"id"` // DIRECT

	///////////////

	Dcterms_modified struct {
		Literal string `json:"literal"` // DIRECT
	} `json:"dcterms_modified"`

	Dcterms_subject struct {
		Uri       string `json:"uri"`       // derived
		PrefLabel string `json:"prefLabel"` // derived
	} `json:"dcterms_subject"`

	Dcterms_educationLevel []struct {
		Uri       string `json:"uri"`       // derived
		PrefLabel string `json:"prefLabel"` // derived
	} `json:"dcterms_educationLevel"`

	Dcterms_description struct {
		Uri       string `json:"uri"`
		PrefLabel string `json:"prefLabel"`
	} `json:"dcterms_description"`

	Dcterms_title struct {
		Literal  string `json:"literal"`  // DIRECT
		Language string `json:"language"` // boilerplate
	} `json:"dcterms_title"`

	Dcterms_rights struct {
		Literal  string `json:"literal"`  // boilerplate
		Language string `json:"language"` // boilerplate
	} `json:"dcterms_rights"`

	Dcterms_rightsHolder struct {
		Literal  string `json:"literal"`  // boilerplate
		Language string `json:"language"` // boilerplate
	} `json:"dcterms_rightsHolder"`

	///////////////

	Asn_statementLabel struct {
		Literal  string `json:"literal"`  // DIRECT
		Language string `json:"language"` // boilerplate
	} `json:"asn_statementLabel"`

	Asn_statementNotation struct {
		Literal  string `json:"literal"`  // DIRECT
		Language string `json:"language"` // boilerplate
	} `json:"asn_statementNotation"`

	Asn_skillEmbodied []struct {
		Uri       string `json:"uri"`       // Predicate
		PrefLabel string `json:"prefLabel"` // Predicate
	} `json:"asn_skillEmbodied"`

	Asn_authorityStatus struct {
		Uri string `json:"uri"` // boilerplate
	} `json:"asn_authorityStatus"`

	Asn_indexingStatus struct {
		Uri string `json:"uri"` // boilerplate
	} `json:"asn_indexingStatus"`

	Asn_hasLevel []struct {
		Uri       string `json:"uri"`       // Predicate
		PrefLabel string `json:"prefLabel"` // Predicate
	} `json:"asn_hasLevel"`

	Asn_crossSubjectReference []struct {
		Uri       string `json:"uri"`       // Predicate
		PrefLabel string `json:"prefLabel"` // Predicate
	} `json:"asn_crossSubjectReference"`

	Asn_conceptTerm string `json:"asn_conceptTerm"` // tag key

	///////////////

	Dc_relation []struct {
		Uri       string `json:"uri"`       // Predicate
		PrefLabel string `json:"prefLabel"` // Predicate
	} `json:"dc_relation"`

	///////////////

	Cls string `json:"cls"` // boilerplate

	Leaf string `json:"leaf"` // boilerplate

	Text string `json:"text"` // DIRECT

	Children []string `json:"children"` //DIRECT
}

// "Year 9 and 10"
// "Level 2 (Years 1- 2)"
func yearsSplit(yearStr string) (rt []string) {

	switch {

	case strings.Contains(yearStr, "Foundation"):
		rt = append(rt, "Foundation Year")

	case strings.Contains(yearStr, "(") && strings.Contains(yearStr, ")"): // "Level 2 (Years 1- 2)"
		s := strings.Index(yearStr, "(")
		e := strings.LastIndex(yearStr, ")")
		r := regexp.MustCompile(`\d+(\s*-\s*\d+)?$`)
		r.ReplaceAllStringFunc(yearStr[s+1:e], func(s string) string {
			yn := strings.Split(s, "-")
			for _, y := range yn {
				y = strings.Trim(y, " ")
				rt = append(rt, "Year "+y)
			}
			return s
		})

	default: // "Year 9 and 10"
		r := regexp.MustCompile(`\d+( and \d+)*$`)
		ss := r.FindAllString(yearStr, 1)
		if len(ss) > 0 {
			s := ss[0]
			yn := strings.Split(s, "and")
			for _, y := range yn {
				y = strings.Trim(y, " ")
				rt = append(rt, "Year "+y)
			}
		}
	}

	return
}

func scanNodeIdTitle(data []byte) map[string]string {
	m := make(map[string]string)
	tool.ScanNode(data, func(i int, id, block string) bool {
		uid := gjson.Get(block, "id").String()
		title := gjson.Get(block, "title").String()
		m[uid] = title
		return true
	})
	return m
}

func nodeProc(dataNM []byte, mCodeChildParent map[string]string, outDir, outName string) {

	const pref4children = "http://vocabulary.curriculum.edu.au/"

	e := bytes.LastIndexAny(dataNM, "}")
	dataNM = dataNM[:e+1]

	outDir = filepath.Clean(outDir)
	parts := []string{}
	out := ""

	mUidTitle := scanNodeIdTitle(dataNM)

	tool.ScanNode(dataNM, func(i int, id, block string) bool {

		code := gjson.Get(block, "code").String()
		// fmt.Println(i, id, code)

		////////////////////////////////////////////////////////

		aj := asnjson{}

		rstChildren := gjson.Get(block, "children")

		// Direct
		aj.Id = gjson.Get(block, "id").String()
		aj.Dcterms_modified.Literal = gjson.Get(block, "created_at").String()
		aj.Dcterms_title.Literal = gjson.Get(block, "title").String()
		aj.Dcterms_title.Language = "en-au"
		aj.Asn_statementLabel.Literal = gjson.Get(block, "doc.typeName").String()
		aj.Asn_statementLabel.Language = "en-au"
		aj.Asn_statementNotation.Literal = gjson.Get(block, "code").String()
		aj.Asn_statementNotation.Language = "en-au"
		aj.Text = gjson.Get(block, "text").String()
		for _, c := range rstChildren.Array() {
			aj.Children = append(aj.Children, pref4children+c.String())
		}

		// Derived
		laTitle := tool.GetAncestorTitle(code, "", mCodeChildParent)
		if laTitle == "" {
			// fmt.Println("Learning area missing:", code)
		}
		subUri, okSubUri := mLaUri[laTitle]
		if okSubUri {
			aj.Dcterms_subject.Uri = subUri
			aj.Dcterms_subject.PrefLabel = laTitle
		}
		if tn := gjson.Get(block, "doc.typeName").String(); tn == "Level" {
			yrTitle := gjson.Get(block, "title").String()
			for _, y := range yearsSplit(yrTitle) {
				aj.Dcterms_educationLevel = append(aj.Dcterms_educationLevel, struct {
					Uri       string "json:\"uri\""
					PrefLabel string "json:\"prefLabel\""
				}{
					Uri:       mYrlvlUri[y],
					PrefLabel: y,
				})
			}
		}

		// Boilerplate
		aj.Asn_authorityStatus.Uri = `http://purl.org/ASN/scheme/ASNAuthorityStatus/Original`
		aj.Asn_indexingStatus.Uri = `http://purl.org/ASN/scheme/ASNIndexingStatus/No`
		aj.Dcterms_rights.Language = "en-au"
		aj.Dcterms_rights.Literal = `©Copyright Australian Curriculum, Assessment and Reporting Authority`
		aj.Dcterms_rightsHolder.Language = "en-au"
		aj.Dcterms_rightsHolder.Literal = `Australian Curriculum, Assessment and Reporting Authority`
		if rstChildren.IsArray() {
			aj.Cls = "folder"
		} else {
			aj.Leaf = "true"
		}
		if gjson.Get(block, "tags").IsObject() {
			aj.Asn_conceptTerm = "SCIENCE_TEACHER_BACKGROUND_INFORMATION"
		}

		// Predicate
		mConnUri := make(map[string]string)
		result := gjson.Get(block, "connections.*")
		if result.IsArray() {
			for _, rUri := range result.Array() {
				uri := rUri.String()
				mConnUri[uri] = mUidTitle[uri]
			}
		}

		nodeType := tool.GetCodeType(code, mCodeChildParent)
		switch nodeType {
		case "GC":
			for uri, title := range mConnUri {
				aj.Asn_skillEmbodied = append(aj.Asn_skillEmbodied, struct {
					Uri       string "json:\"uri\""
					PrefLabel string "json:\"prefLabel\""
				}{
					Uri:       uri,
					PrefLabel: title,
				})
			}

		case "LA":
			for uri, title := range mConnUri {
				aj.Dc_relation = append(aj.Dc_relation, struct {
					Uri       string "json:\"uri\""
					PrefLabel string "json:\"prefLabel\""
				}{
					Uri:       uri,
					PrefLabel: title,
				})
			}

		case "AS":
			for uri, title := range mConnUri {
				aj.Asn_hasLevel = append(aj.Asn_hasLevel, struct {
					Uri       string "json:\"uri\""
					PrefLabel string "json:\"prefLabel\""
				}{
					Uri:       uri,
					PrefLabel: title,
				})
			}

		case "CCP":
			for uri, title := range mConnUri {
				aj.Asn_crossSubjectReference = append(aj.Asn_crossSubjectReference, struct {
					Uri       string "json:\"uri\""
					PrefLabel string "json:\"prefLabel\""
				}{
					Uri:       uri,
					PrefLabel: title,
				})
			}

		default:
			if code != "root" {
				log.Printf("NodeType: '%v' is not one of [GC CCP LA AS], Code is '%v'", nodeType, code)
			}
		}

		////////////////////////////////////////////////////////////////

		if bytes, err := json.Marshal(aj); err == nil {
			parts = append(parts, string(bytes))
		}

		return true
	})

	out = "[" + strings.Join(parts, ",") + "]"       // combine whole
	out = jt.FmtStr(out, "  ")                       // format json
	out = jt.TrimFields(out, true, true, true, true) // remove empty object, string, array

	if !strings.HasSuffix(outName, ".json") {
		outName += ".json"
	}

	outPath := fmt.Sprintf("./%s/%s", outDir, outName)
	fmt.Printf("save at [%s], data length: [%d]\n", outPath, len(out))
	os.WriteFile(outPath, []byte(out), os.ModePerm)
}

/////////////////////////////////////////////////////////////////////////////

// func getIdBlock(js string) (mIdBlock, mIdBlockLeaf map[string]string) {

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	r := strings.NewReader(js)
// 	result, ok := jt.ScanObject(ctx, r, true, true, jt.OUT_FMT)
// 	if !ok {
// 		log.Fatalln("node file is NOT JSON array")
// 	}

// 	mIdBlock = make(map[string]string)
// 	mIdBlockLeaf = make(map[string]string)

// 	for r := range result {
// 		if r.Err != nil {
// 			log.Fatalln(r.Err)
// 		}
// 		id := gjson.Get(r.Obj, "id").String()
// 		mIdBlock[id] = r.Obj

// 		hasChildren := gjson.Get(r.Obj, "children").IsArray()
// 		if !hasChildren {
// 			mIdBlockLeaf[id] = r.Obj
// 		}
// 	}

// 	return
// }

// func childrenId(cBlock string) (cid []string) {
// 	s := strings.Index(cBlock, "[")
// 	e := strings.LastIndex(cBlock, "]")
// 	cBlock = cBlock[s+1 : e]
// 	cBlock = strings.Trim(cBlock, " \n\t")
// 	for _, id := range strings.Split(cBlock, ",") {
// 		cid = append(cid, strings.Trim(id, " \n\t"))
// 	}
// 	return
// }

// func childrenRepl(inpath string, mIdBlock map[string]string) string {

// 	data, err := os.ReadFile(inpath)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	rChildren := regexp.MustCompile(`"children":\s*\[([\n\s]*"http[^"]+",?)+[\n\s]*\]`)
// 	js := string(data)
// 	repl := false

// AGAIN:
// 	repl = false
// 	js = rChildren.ReplaceAllStringFunc(js, func(s string) string {
// 		for _, id := range childrenId(s) {
// 			id = id[1 : len(id)-1]
// 			if block, ok := mIdBlock[id]; ok {
// 				s = strings.ReplaceAll(s, "\""+id+"\"", block)
// 				repl = true
// 			}
// 		}
// 		return s
// 	})

// 	if repl {
// 		goto AGAIN
// 	}

// 	return jt.FmtStr(js, "  ")
// }

// func getRootWholeObject(allNestedSet string) string {

// 	rId := regexp.MustCompile(`"id": "http[^"]+"`)

// 	mIdCnt := make(map[string]int)
// 	rId.ReplaceAllStringFunc(allNestedSet, func(s string) string {
// 		mIdCnt[s]++
// 		return s
// 	})

// 	// fmt.Println(len(mIdCnt))

// 	mIdRootCnt := make(map[string]int)
// 	for idstr, cnt := range mIdCnt {
// 		if cnt == 1 {
// 			mIdRootCnt[idstr] = cnt
// 		}
// 	}

// 	mIdBlock, _ := getIdBlock(allNestedSet)

// 	for idstr := range mIdRootCnt {
// 		s := strings.Index(idstr, "http:")
// 		e := strings.LastIndex(idstr, "\"")
// 		id := idstr[s:e]
// 		return mIdBlock[id]
// 	}

// 	return ""
// }
