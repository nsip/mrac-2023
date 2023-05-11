package sub

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	dt "github.com/digisan/gotk/data-type"
	"github.com/digisan/gotk/strs"
	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func LA(js string) (mOut map[string]string) {

	mOut = make(map[string]string)

	mNameLaAs := map[string]string{
		"English":      "English",
		"HASS":         "Humanities and Social Sciences",
		"HPE":          "Health and Physical Education",
		"Languages":    "Languages",
		"Mathematics":  "Mathematics",
		"Science":      "Science",
		"Technologies": "Technologies",
		"The Arts":     "The Arts",
	}

	mRoot := map[string]interface{}{
		"code":       gjson.Get(js, "code").String(),
		"uuid":       gjson.Get(js, "uuid").String(),
		"type":       gjson.Get(js, "type").String(),
		"created_at": gjson.Get(js, "created_at").String(),
		"title":      gjson.Get(js, "title").String(),
		"children":   nil,
	}

	mASfield := map[string]interface{}{
		"code":       "",
		"uuid":       "",
		"type":       "",
		"created_at": "",
		"title":      "",
		"children":   nil,
	}

	mLAfield := map[string]interface{}{
		"code":       "",
		"uuid":       "",
		"type":       "",
		"created_at": "",
		"title":      "",
		"children":   nil,
	}

	mAS := make(map[string][]string)
	mLA := make(map[string][]string)

	valueC1 := gjson.Get(js, "children")
	if valueC1.IsArray() {
		for _, r1 := range valueC1.Array() {
			if r1.IsObject() { // "Achievement Standards", "Cross-curriculum Priorities", "General Capabilities", "Learning Areas"
				block1 := r1.String()
				valueTitle1 := gjson.Get(block1, "title")
				title1str := valueTitle1.String()
				fmt.Println(title1str, ":")

				valueC2 := gjson.Get(block1, "children")
				if valueC2.IsArray() {
					for _, r2 := range valueC2.Array() {
						if r2.IsObject() { // "English", "Mathematics", etc.
							block2 := r2.String()
							valueTitle2 := gjson.Get(block2, "title")
							title2str := valueTitle2.String()
							fmt.Println("	", title2str)

							switch title1str {
							case "Achievement Standards":
								mAS[title2str] = append(mAS[title2str], block2)
							case "Learning Areas":
								mLA[title2str] = append(mLA[title2str], block2)
							}
						}
					}
				}

				switch title1str {
				case "Achievement Standards":
					mASfield["code"] = gjson.Get(block1, "code").String()
					mASfield["uuid"] = gjson.Get(block1, "uuid").String()
					mASfield["type"] = gjson.Get(block1, "type").String()
					mASfield["created_at"] = gjson.Get(block1, "created_at").String()
					mASfield["title"] = gjson.Get(block1, "title").String()
					mASfield["children"] = mAS
				case "Learning Areas":
					mLAfield["code"] = gjson.Get(block1, "code").String()
					mLAfield["uuid"] = gjson.Get(block1, "uuid").String()
					mLAfield["type"] = gjson.Get(block1, "type").String()
					mLAfield["created_at"] = gjson.Get(block1, "created_at").String()
					mLAfield["title"] = gjson.Get(block1, "title").String()
					mLAfield["children"] = mLA
				}
			}
		}
	}

	if len(mAS) != len(mLA) {
		log.Println("[Achievement Standards] children count is NOT same as [Learning Areas] children count")
		if len(mLA) < len(mAS) {
			log.Fatalln("[Learning Areas] children count less than [Achievement Standards] children count")
		}
	}

NEXT_LA:
	for la, blockLA := range mLA {
		out := ""
		for as, blockAS := range mAS {
			if la == as || mNameLaAs[la] == as {

				out, _ = sjson.Set(out, "code", mRoot["code"])
				out, _ = sjson.Set(out, "uuid", mRoot["uuid"])
				out, _ = sjson.Set(out, "type", mRoot["type"])
				out, _ = sjson.Set(out, "created_at", mRoot["created_at"])
				out, _ = sjson.Set(out, "title", mRoot["title"])

				out, _ = sjson.Set(out, "children.0.code", mASfield["code"])
				out, _ = sjson.Set(out, "children.0.uuid", mASfield["uuid"])
				out, _ = sjson.Set(out, "children.0.type", mASfield["type"])
				out, _ = sjson.Set(out, "children.0.created_at", mASfield["created_at"])
				out, _ = sjson.Set(out, "children.0.title", mASfield["title"])
				for i, bAS := range blockAS {
					path := fmt.Sprintf("children.0.children.%d", i)
					out, _ = sjson.SetRaw(out, path, bAS)
				}

				out, _ = sjson.Set(out, "children.1.code", mLAfield["code"])
				out, _ = sjson.Set(out, "children.1.uuid", mLAfield["uuid"])
				out, _ = sjson.Set(out, "children.1.type", mLAfield["type"])
				out, _ = sjson.Set(out, "children.1.created_at", mLAfield["created_at"])
				out, _ = sjson.Set(out, "children.1.title", mLAfield["title"])
				for i, bLA := range blockLA {
					path := fmt.Sprintf("children.1.children.%d", i)
					out, _ = sjson.SetRaw(out, path, bLA)
				}

				out = jt.FmtStr(out, "  ")
				lk.FailOnErrWhen(!dt.IsJSON([]byte(out)), "%v", errors.New("invalid JSON from [la]"))
				mOut[la] = out
				continue NEXT_LA
			}
		}
	}

	return
}

func ConnFieldMapping(js, uri string, meta map[string]string) string {
	r1 := regexp.MustCompile(`"connections":\s*\{[^{}]*\},?`)
	r2 := regexp.MustCompile(`"[\d\w]{40}":\s*\[([\n\s]*"[\d\w-]+",?[\n\s]*)+\],?`)
	return r1.ReplaceAllStringFunc(js, func(s string) string {
		return r2.ReplaceAllStringFunc(s, func(ss string) string {
			starts, _ := strs.IndexAll(ss, "\"")
			code := ss[starts[0]+1 : starts[1]]
			ssMeta := strings.ReplaceAll(ss, code, meta[code])
			m := make(map[string]string)
			for i, pos := range starts {
				if i > 1 && i%2 == 1 {
					id := ss[starts[i-1]+1 : pos]
					m[id] = uri + id
				}
			}
			for id, uri := range m {
				ssMeta = strings.ReplaceAll(ssMeta, id, uri)
			}
			return ssMeta
		})
	})
}

// only test for restruct_***
// func laRestructure(js string, I int) string {

// 	for i := 0; i < 100; i++ {

// 		path := fmt.Sprintf("children.%d.children.%d.code", I, i)
// 		testcode := gjson.Get(js, path).String()
// 		if testcode == "" {
// 			break
// 		}
// 		fmt.Println(testcode)

// 		for j := 0; j < 100; j++ {
// 			path := fmt.Sprintf("children.%d.children.%d.children.%d.code", I, i, j)
// 			testcode := gjson.Get(js, path).String()
// 			if testcode == "" {
// 				break
// 			}
// 			fmt.Printf("\t%s\n", testcode)

// 			for k := 0; k < 100; k++ {
// 				path := fmt.Sprintf("children.%d.children.%d.children.%d.children.%d.code", I, i, j, k)
// 				testcode := gjson.Get(js, path).String()
// 				if testcode == "" {
// 					break
// 				}
// 				fmt.Printf("\t\t%s\n", testcode)

// 				for l := 0; l < 100; l++ {
// 					path := fmt.Sprintf("children.%d.children.%d.children.%d.children.%d.children.%d.code", I, i, j, k, l)
// 					testcode := gjson.Get(js, path).String()
// 					if testcode == "" {
// 						break
// 					}
// 					fmt.Printf("\t\t\t%s\n", testcode)
// 				}
// 			}
// 		}
// 	}

// 	return ""
// }
