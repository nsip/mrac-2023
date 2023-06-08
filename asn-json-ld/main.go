package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/digisan/gotk/strs"
	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
	"github.com/nsip/mrac-2023/asn-json/tool"
	"github.com/tidwall/gjson"
)

func init() {
	lk.WarnDetail(false)
}

func addContext(js, ctx string) string {
	return strings.Replace(js, "{", "{"+context+",", 1)
}

func replace(js string) string {
	r := regexp.MustCompile(`"(dc_|dcterms_|asn_)[^"]+"`)
	js = r.ReplaceAllStringFunc(js, func(s string) string {
		s = strings.Trim(s, "\"")
		ss := strings.Split(s, "_")
		p1 := mPrefRepl[ss[0]+"_"]
		return "\"" + p1 + ss[1] + "\""
	})
	return js
}

func findIdLinkage(js string, mFamilyTree map[string][]string) (mIdLink2P, mIdLink2C map[string][]string) {
	mIdLink2C = make(map[string][]string)
	mIdLink2P = make(map[string][]string)
	for _, children := range mFamilyTree {
		for _, child := range children {
			if strings.HasSuffix(child, ".id") {
				id := gjson.Get(js, child).String()
				pid := gjson.Get(js, jt.NewUncle(child, "id")).String()
				mIdLink2P[id] = append(mIdLink2P[id], pid)
				mIdLink2C[pid] = append(mIdLink2C[pid], id)
			}
		}
	}
	return
}

func cvt2jsonld(asnpath string) {

	data, err := os.ReadFile(asnpath)
	if err != nil {
		panic(err)
	}
	js := string(data)
	// fmt.Println(len(js))

	_, mFamilyTree := jt.FamilyTree(js)
	// fmt.Println(len(mLvlSiblings))
	mIdLink2P, _ := findIdLinkage(js, mFamilyTree)
	// fmt.Println(len(mIdLink2P), len(mIdLink2C))

	// fmt.Println(mIdLink2P[uri4id + "649c9d14-75b7-41e3-ac5f-c4c86fd8f57c"])
	// fmt.Println(mIdLink2C[uri4id + "649c9d14-75b7-41e3-ac5f-c4c86fd8f57c"])
	// fmt.Println(mIdLink2C[uri4id + "652a716a-67c2-4174-9dbd-79977ba3f049"])

	for oldPref, newPref := range mPrefRepl {
		js = strings.ReplaceAll(js, "\""+oldPref, "\""+newPref)
	}

	for oldField, newField := range mFieldRepl {
		js = strings.ReplaceAll(js, "\""+oldField+"\"", "\""+newField+"\"")
	}

	rRm := regexp.MustCompile(`("cls":\s*"\w+",?)|("leaf":\s*"\w+",?)`)
	js = rRm.ReplaceAllStringFunc(js, func(s string) string {
		return ""
	})

	rId := regexp.MustCompile(`"@id":\s*"http:[^"]+",?`)
	js = rId.ReplaceAllStringFunc(js, func(s string) (ret string) {

		id := tool.FetchValue(s, "|")

		pids := []string{}
		for _, pid := range mIdLink2P[id] {
			pids = append(pids, fmt.Sprintf(`{ "@id": "%s" }`, pid))
		}
		pidstr := ""
		if len(pids) > 0 {
			pidstr = fmt.Sprintf(`"gem:isChildOf": [%s]`, strings.Join(pids, ","))
		}
		if pidstr != "" {
			ret = pidstr + "," + s
		} else {
			ret = s
		}

		typestr := `"@type": [ "asn:Statement", "skos:Concept" ],`
		ret = typestr + ret

		return
	})

	rModified := regexp.MustCompile(`"dc:modified":\s*\{[^{}]+\},?`)
	js = rModified.ReplaceAllStringFunc(js, func(s string) string {
		suffix := ""
		if strings.HasSuffix(s, ",") {
			suffix = ","
		}
		str := tool.FetchValue(s, "|")
		// return fmt.Sprintf(`"dc:modified": { "@value": "%s", "@type": [ "xsd:dateTime", "skos:Concept" ] }%s`, str, suffix)
		return fmt.Sprintf(`"dc:modified": { "@value": "%s", "@type": "xsd:dateTime" }%s`, str, suffix)
	})

	rLangLit := regexp.MustCompile(`\{[\s\n]*"language":\s*"[^"]+",?[\s\n]*"literal":\s*"[^"]+"[\s\n]*\},?`)
	js = rLangLit.ReplaceAllStringFunc(js, func(s string) string {
		// fmt.Println(s)

		suffix := ""
		if strings.HasSuffix(s, ",") {
			suffix = ","
		}

		starts, _ := strs.IndexAll(s, "\"")
		lang := s[starts[2]+1 : starts[3]]
		lit := s[starts[6]+1 : starts[7]]
		// fmt.Println(lang, lit)

		if lang == "en-au" {
			return fmt.Sprintf(`"%s"%s`, lit, suffix)
		} else {
			return fmt.Sprintf(`{ "@language": "%s", "@value": "%s" }%s`, lang, lit, suffix)
		}
	})

	rUri := regexp.MustCompile(`"uri":`)
	js = rUri.ReplaceAllStringFunc(js, func(s string) string {
		return `"@id":`
	})

	rPrefLabel := regexp.MustCompile(`"prefLabel":`)
	js = rPrefLabel.ReplaceAllStringFunc(js, func(s string) string {
		return `"skos:prefLabel":`
	})

	rYrLvl := regexp.MustCompile(`"dc:educationLevel":\s*\[[^\[\]]+\],?`)
	js = rYrLvl.ReplaceAllStringFunc(js, func(s string) string {

		start, end := strings.Index(s, "["), strings.LastIndex(s, "]")
		block := s[start+1 : end]
		starts, _ := strs.IndexAll(block, "\"")

		year := ""
		years := []string{}
		for i := 0; i < strings.Count(block, "\"@id\""); i++ {
			idx := i * 8
			// atid := block[starts[idx+2]+1 : starts[idx+3]]
			prefLabel := block[starts[idx+6]+1 : starts[idx+7]]
			years = append(years, prefLabel)
		}
		if len(years) > 0 {
			sort.Strings(years)
			year = years[0]
		}

		ret := fmt.Sprintf(`"esa:nominalYearLevel": "%s",`, year)
		ret += fmt.Sprintf(`"dc:isPartOf": { "@id": "%s" },`, "TBD")
		return ret + s
	})

	js = addContext(js, context)

	name := filepath.Base(asnpath)
	jsonldpath := filepath.Join("../data-out/asn-json-ld", name)
	os.WriteFile(jsonldpath, []byte(js), os.ModePerm)
}

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

	wg := sync.WaitGroup{}
	wg.Add(len(mInputLa))

	for file := range mInputLa {
		go func(file string) {
			if file != "la-Languages.json" {
				wg.Done()
				return
			}
			cvt2jsonld(filepath.Join("../data-out/asn-json", file))
			wg.Done()

		}(file)
	}
	wg.Wait()

	lk.Warn("FORMAT each out file, then process extra duplicated line in 'main_test.go'")
}