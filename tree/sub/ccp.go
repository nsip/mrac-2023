package sub

import (
	"errors"
	"fmt"

	dt "github.com/digisan/gotk/data-type"
	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func CCP(js, outDir string) string {

	mRoot := map[string]interface{}{
		"code":       gjson.Get(js, "code").String(),
		"uuid":       gjson.Get(js, "uuid").String(),
		"type":       gjson.Get(js, "type").String(),
		"created_at": gjson.Get(js, "created_at").String(),
		"title":      gjson.Get(js, "title").String(),
		"children":   nil,
	}

	valueC1 := gjson.Get(js, "children")
	if valueC1.IsArray() {
		for _, r1 := range valueC1.Array() {
			if r1.IsObject() { // "Achievement Standards", "Cross-curriculum Priorities", "General Capabilities", "Learning Areas"
				block1 := r1.String()
				valueTitle1 := gjson.Get(block1, "title")
				title1str := valueTitle1.String()
				fmt.Println(title1str, ":")
				if title1str == "Cross-curriculum Priorities" {
					mRoot["children"] = block1
				}
			}
		}
	}

	out := ""
	out, _ = sjson.Set(out, "code", mRoot["code"])
	out, _ = sjson.Set(out, "uuid", mRoot["uuid"])
	out, _ = sjson.Set(out, "type", mRoot["type"])
	out, _ = sjson.Set(out, "created_at", mRoot["created_at"])
	out, _ = sjson.Set(out, "title", mRoot["title"])
	out, _ = sjson.SetRaw(out, fmt.Sprintf("children.%d", 0), mRoot["children"].(string))

	out = jt.FmtStr(out, "    ")
	lk.FailOnErrWhen(!dt.IsJSON([]byte(out)), "%v", errors.New("invalid JSON from [ccp]"))
	return out
}
