package node

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	dt "github.com/digisan/gotk/data-type"
	lk "github.com/digisan/logkit"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	jt "github.com/digisan/json-tool"
)

func Scan(data []byte, f func(i int, id, block string) bool) {

	js := string(data)
	r := regexp.MustCompile(`"[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}":`)
	pGrp := r.FindAllStringIndex(js, -1)
	// fmt.Println(len(pGrp), js[pGrp[0][0]:pGrp[0][1]])

	for i := 0; i < len(pGrp); i++ {

		p, pn := pGrp[i], []int{}
		if i < len(pGrp)-1 {
			p, pn = pGrp[i], pGrp[i+1]
		}

		ids, ide := p[0]+1, p[1]-2
		id := js[ids:ide]
		// fmt.Println(id)

		blkStart, blkEnd := 0, 0
		if i < len(pGrp)-1 {
			blkStart, blkEnd = p[1], pn[0]-1
		} else {
			blkStart, blkEnd = p[1], len(js)-1
		}

		block := js[blkStart:blkEnd]
		block = strings.TrimSuffix(strings.TrimSpace(block), ",")

		if !f(i, id, block) {
			break
		}
	}
}

func Process(data []byte, uri string, meta map[string]string, outDir string) {

	e := bytes.LastIndexAny(data, "}")
	data = data[:e+1]

	outDir = strings.Trim(outDir, `./\`)
	parts := []string{}
	out := ""

	Scan(data, func(i int, id, block string) bool {

		// "uuid": {id} => "id": "http://abc/def/{id}"
		newIdVal := fmt.Sprintf("%s%s", uri, gjson.Get(block, "uuid").String())
		block, _ = sjson.Set(block, "uuid", newIdVal)
		block = strings.Replace(block, `"uuid"`, `"id"`, 1)

		m := make(map[string]interface{})
		json.Unmarshal([]byte(gjson.Get(block, "connections").String()), &m)

		for k, v := range m {
			// "abcdeft" => "Levels" etc.
			block = strings.Replace(block, k, meta[k], 1)
			// "abc-def" => "http://abc/def/{id}"
			for _, a := range v.([]interface{}) {
				block = strings.Replace(block, a.(string), fmt.Sprintf("%s%s", uri, a), 1)
			}
		}

		part := fmt.Sprintf(`"%s": %s`, id, block)
		parts = append(parts, part)
		return true
	})

	out = "{" + strings.Join(parts, ",") + "}"
	out = jt.FmtStr(out, "    ")

	lk.FailOnErrWhen(!dt.IsJSON([]byte(out)), "%v", errors.New("invalid JSON from node & meta"))

	os.WriteFile(fmt.Sprintf("./%s/node-meta.json", outDir), []byte(out), os.ModePerm)
}
