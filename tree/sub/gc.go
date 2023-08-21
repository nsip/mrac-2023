package sub

import (
	"errors"
	"fmt"

	. "github.com/digisan/go-generics/v2"
	dt "github.com/digisan/gotk/data-type"

	// jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
	u "github.com/nsip/mrac-2023/util"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func GC(js string) map[string]string {

	mOut := make(map[string]string)

	var (
		// gcTitles = []string{
		// 	"Literacy",
		// 	"Numeracy",
		// 	"Critical and Creative Thinking",
		// 	"Personal and Social capability",
		// 	"Digital Literacy",
		// 	"Intercultural Understanding",
		// 	"Ethical Understanding",
		// }
		gcCodes = []string{
			"CCT",
			"DL",
			"EU",
			"IU",
			"L",
			"N",
			"PSC",
		}
	)

	// L0
	mRoot := map[string]any{
		"code":       gjson.Get(js, "code").String(),
		"uuid":       gjson.Get(js, "uuid").String(),
		"type":       gjson.Get(js, "type").String(),
		"created_at": gjson.Get(js, "created_at").String(),
		"title":      gjson.Get(js, "title").String(),
		"children":   nil,
	}

	// L1
	mGC := map[string]any{
		"code":       "",
		"uuid":       "",
		"type":       "",
		"created_at": "",
		"title":      "",
		"children":   nil,
	}

	// L2
	mL := map[string]any{
		"code":         "",
		"uuid":         "",
		"type":         "",
		"created_at":   "",
		"title":        "",
		"doc.typeName": "",
		"children":     nil,
	}

	// L2
	mN := map[string]any{
		"code":         "",
		"uuid":         "",
		"type":         "",
		"created_at":   "",
		"title":        "",
		"doc.typeName": "",
		"children":     nil,
	}

	// L2
	mCCT := map[string]any{
		"code":         "",
		"uuid":         "",
		"type":         "",
		"created_at":   "",
		"title":        "",
		"doc.typeName": "",
		"children":     nil,
	}

	// L2
	mPSC := map[string]any{
		"code":         "",
		"uuid":         "",
		"type":         "",
		"created_at":   "",
		"title":        "",
		"doc.typeName": "",
		"children":     nil,
	}

	// L2
	mDL := map[string]any{
		"code":         "",
		"uuid":         "",
		"type":         "",
		"created_at":   "",
		"title":        "",
		"doc.typeName": "",
		"children":     nil,
	}

	// L2
	mIU := map[string]any{
		"code":         "",
		"uuid":         "",
		"type":         "",
		"created_at":   "",
		"title":        "",
		"doc.typeName": "",
		"children":     nil,
	}

	// L2
	mEU := map[string]any{
		"code":         "",
		"uuid":         "",
		"type":         "",
		"created_at":   "",
		"title":        "",
		"doc.typeName": "",
		"children":     nil,
	}

	var (
		mL2s = []map[string]any{mL, mN, mCCT, mPSC, mDL, mIU, mEU}
	)

	valueC1 := gjson.Get(js, "children")
	if valueC1.IsArray() {
		for _, r1 := range valueC1.Array() {
			if r1.IsObject() { // "Achievement Standards", "Cross-curriculum Priorities", "General Capabilities", "Learning Areas"
				block1 := r1.String()
				valueTitle1 := gjson.Get(block1, "title")
				title1str := valueTitle1.String()
				fmt.Println(title1str, ":")

				if title1str == "General Capabilities" {
					// mRoot["children"] = block1

					mGC["code"] = gjson.Get(block1, "code").String()
					mGC["uuid"] = gjson.Get(block1, "uuid").String()
					mGC["type"] = gjson.Get(block1, "type").String()
					mGC["created_at"] = gjson.Get(block1, "created_at").String()
					mGC["title"] = gjson.Get(block1, "title").String()

					valueC2 := gjson.Get(block1, "children")
					if valueC2.IsArray() {
						for _, r2 := range valueC2.Array() {
							if r2.IsObject() {
								block2 := r2.String()
								title2str := gjson.Get(block2, "title").String()
								code2str := gjson.Get(block2, "code").String()

								// if In(title2str, gcTitles...) {
								// 	fmt.Println("  ", title2str)

								if In(code2str, gcCodes...) {
									fmt.Println("  ", title2str, "  ", code2str)

									var m map[string]any

									// switch {
									// case strings.EqualFold("Literacy", title2str):
									// 	m = mNLLP
									// case strings.EqualFold("Numeracy", title2str):
									// 	m = mNNLP
									// case strings.EqualFold("Critical and Creative Thinking", title2str):
									// 	m = mCCT
									// case strings.EqualFold("Personal and Social capability", title2str):
									// 	m = mPSC
									// case strings.EqualFold("Digital Literacy", title2str):
									// 	m = mDL
									// case strings.EqualFold("Intercultural Understanding", title2str):
									// 	m = mIU
									// case strings.EqualFold("Ethical Understanding", title2str):
									// 	m = mEU
									// }

									switch code2str {
									case "CCT":
										m = mCCT
									case "DL":
										m = mDL
									case "EU":
										m = mEU
									case "IU":
										m = mIU
									case "L":
										m = mL
									case "N":
										m = mN
									case "PSC":
										m = mPSC
									}

									m["code"] = gjson.Get(block2, "code").String()
									m["uuid"] = gjson.Get(block2, "uuid").String()
									m["type"] = gjson.Get(block2, "type").String()
									m["created_at"] = gjson.Get(block2, "created_at").String()
									m["title"] = gjson.Get(block2, "title").String()
									m["doc.typeName"] = gjson.Get(block2, "doc.typeName").String()
									m["children"] = gjson.Get(block2, "children").String()
								}
							}
						}
					}
				}
			}
		}
	}

	// fmt.Println(mRoot["title"])
	// fmt.Println(mGC["title"])
	// fmt.Println(mIU["title"])
	// fmt.Println(mEU["title"])

	for _, L2 := range mL2s {

		if MapAllValAreEmpty(L2) {
			continue
		}

		out := ""
		out, _ = sjson.Set(out, "code", mRoot["code"])
		out, _ = sjson.Set(out, "uuid", mRoot["uuid"])
		out, _ = sjson.Set(out, "type", mRoot["type"])
		out, _ = sjson.Set(out, "created_at", mRoot["created_at"])
		out, _ = sjson.Set(out, "title", mRoot["title"])
		out, _ = sjson.Set(out, "children.0.code", mGC["code"])
		out, _ = sjson.Set(out, "children.0.uuid", mGC["uuid"])
		out, _ = sjson.Set(out, "children.0.type", mGC["type"])
		out, _ = sjson.Set(out, "children.0.created_at", mGC["created_at"])
		out, _ = sjson.Set(out, "children.0.title", mGC["title"])
		out, _ = sjson.Set(out, "children.0.children.0.code", L2["code"])
		out, _ = sjson.Set(out, "children.0.children.0.uuid", L2["uuid"])
		out, _ = sjson.Set(out, "children.0.children.0.type", L2["type"])
		out, _ = sjson.Set(out, "children.0.children.0.created_at", L2["created_at"])
		out, _ = sjson.Set(out, "children.0.children.0.title", L2["title"])
		out, _ = sjson.Set(out, "children.0.children.0.doc.typeName", L2["doc.typeName"])
		out, _ = sjson.SetRaw(out, "children.0.children.0.children", L2["children"].(string))

		out, err := u.FmtJSON(out)
		lk.FailOnErr("%v", err)

		lk.FailOnErrWhen(!dt.IsJSON([]byte(out)), "%v", errors.New("invalid JSON from [gc]"))
		mOut[fmt.Sprint(L2["title"])] = out
	}

	return mOut
}
