package node

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	dt "github.com/digisan/gotk/data-type"
	fd "github.com/digisan/gotk/file-dir"
	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
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

func Process(dataNode []byte, uri string, meta map[string]string, outDir string) {

	e := bytes.LastIndexAny(dataNode, "}")
	dataNode = dataNode[:e+1]

	outDir = strings.Trim(outDir, `./\`)
	parts := []string{}
	out := ""

	Scan(dataNode, func(i int, id, block string) bool {

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
	out = jt.FmtStr(out, "  ")

	lk.FailOnErrWhen(!dt.IsJSON([]byte(out)), "%v", errors.New("invalid JSON from node & meta"))

	os.WriteFile(fmt.Sprintf("./%s/node-meta.json", outDir), []byte(out), os.ModePerm)
}

//////////////////////////////////////////////////////////////////////

func MarkUrl(ids, codes []string, mCodeUrl, mIdUrl map[string]string) {

	// for _, code := range codes {
	// 	if url, ok := mCodeUrl[code]; ok {
	// 		for i, code := range codes {
	// 			if NotIn(code, "root", "LA", "AS", "GC", "CCP") {
	// 				mCodeUrl[code] = url
	// 				mIdUrl[ids[i]] = url
	// 			}
	// 		}
	// 		return
	// 	}
	// }

	url := ""
	for i, code := range codes {
		if i == len(codes)-3 {
			switch code {
			case "HAS", "HASS", "ASHAS", "ASHASS":
				url = "http://vocabulary.curriculum.edu.au/MRAC/LA/HASS/"
			case "ENG", "ASENG":
				url = "http://vocabulary.curriculum.edu.au/MRAC/LA/ENG/"
			case "LAN", "ASLAN":
				url = "http://vocabulary.curriculum.edu.au/MRAC/LA/LAN/"
			case "SCI", "ASSCI":
				url = "http://vocabulary.curriculum.edu.au/MRAC/LA/SCI/"
			case "ART", "ASART":
				url = "http://vocabulary.curriculum.edu.au/MRAC/LA/ART/"
			case "HPE", "ASHPE":
				url = "http://vocabulary.curriculum.edu.au/MRAC/LA/HPE/"
			case "MAT", "ASMAT":
				url = "http://vocabulary.curriculum.edu.au/MRAC/LA/MAT/"
			case "TEC", "ASTEC":
				url = "http://vocabulary.curriculum.edu.au/MRAC/LA/TEC/"

			case "CCT":
				url = "http://vocabulary.curriculum.edu.au/MRAC/GC/CCT/"
			case "N":
				url = "http://vocabulary.curriculum.edu.au/MRAC/GC/N/"
			case "DL":
				url = "http://vocabulary.curriculum.edu.au/MRAC/GC/DL/"
			case "L":
				url = "http://vocabulary.curriculum.edu.au/MRAC/GC/L/"
			case "PSC":
				url = "http://vocabulary.curriculum.edu.au/MRAC/GC/PSC/"
			case "IU":
				url = "http://vocabulary.curriculum.edu.au/MRAC/GC/IU/"
			case "EU":
				url = "http://vocabulary.curriculum.edu.au/MRAC/GC/EU/"

			case "AA":
				url = "http://vocabulary.curriculum.edu.au/MRAC/CCP/AA/"
			case "S":
				url = "http://vocabulary.curriculum.edu.au/MRAC/CCP/S/"
			case "A_TSI":
				url = "http://vocabulary.curriculum.edu.au/MRAC/CCP/A_TSI/"
			}
			break
		}
		if i == len(codes)-2 {
			switch code {
			case "AS", "LA":
				url = "http://vocabulary.curriculum.edu.au/MRAC/LA/"
			case "GC":
				url = "http://vocabulary.curriculum.edu.au/MRAC/GC/"
			case "CCP":
				url = "http://vocabulary.curriculum.edu.au/MRAC/CCP/"
			}
			mCodeUrl[code] = url
			mIdUrl[ids[i]] = url
			break
		}
		if i == len(codes)-1 {
			switch code {
			case "root":
				url = "http://vocabulary.curriculum.edu.au/MRAC/"
			}
			mCodeUrl[code] = url
			mIdUrl[ids[i]] = url
			break
		}
	}

	if len(codes) > 2 && url == "" {
		panic("Need Code: " + strings.Join(codes, ","))
	}

	for i, code := range codes {
		if i < len(codes)-2 {
			if _, ok := mCodeUrl[code]; !ok {
				mCodeUrl[code] = url
			}
			if _, ok := mIdUrl[ids[i]]; !ok {
				mIdUrl[ids[i]] = url
			}
		}
	}
}

func TrackCode(ms map[string]string, code string) (codes, ids []string) {
	ID := ""
	for id, valstr := range ms {
		if gjson.Get(valstr, "code").String() == code {
			ID = id
			break
		}
	}
	if ID != "" {
		for _, id := range TrackId(ms, ID) {
			valstr := ms[id]
			codes = append(codes, gjson.Get(valstr, "code").String())
			ids = append(ids, id)
		}
	}
	return // Reverse(codes), Reverse(ids)
}

func TrackId(ms map[string]string, id string) (ids []string) {
	ids = append(ids, id)
	for parent := IsChild(ms, id); len(parent) > 0; parent = IsChild(ms, parent) {
		ids = append(ids, parent)
	}
	return // Reverse(ids)
}

func IsChild(ms map[string]string, childId string) string {
	for id := range ms {
		if HasChild(ms, id, childId) {
			return id
		}
	}
	return ""
}

func HasChild(ms map[string]string, id, childId string) bool {
	valstr := ms[id]
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("err: [%v]\nid: [%v]\nvalstr: [%v]", err, id, valstr)
		}
	}()
	if childrenRst := gjson.Get(valstr, "children"); childrenRst.Type != gjson.Null && childrenRst.IsArray() {
		if children := childrenRst.Array(); len(children) > 0 {
			for _, child := range children {
				// fmt.Println(child)
				if childId == child.String() {
					return true
				}
			}
		}
	}
	return false
}

func Scan2Map(data []byte) map[string]any {
	M := make(map[string]any)
	if err := json.Unmarshal(data, &M); err != nil {
		panic(err)
	}
	return M
}

func Scan2MapStrval(data []byte) map[string]string {
	M := Scan2Map(data)
	ret := make(map[string]string)
	for k, v := range M {
		vdata, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		ret[k] = string(vdata)
	}
	return ret
}

func Scan2Flatmap(data []byte) map[string]map[string]any {
	M := Scan2Map(data)
	ret := make(map[string]map[string]any)
	for k, v := range M {
		vdata, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		mf, err := jt.Flatten(vdata)
		if err != nil {
			panic(err)
		}
		ret[k] = mf
	}
	return ret
}

func GenCodeIdUrlTxt(dataNode []byte, outDir string) {

	mNodeIdBlock := Scan2MapStrval(dataNode)

	// str := ms["92b62493-c251-421f-b774-235cfd597852"]
	// r := gjson.Get(str, "children.#").Int()
	// fmt.Println(r)

	// mf := scan2flatmap(data)
	// m := mf["92b62493-c251-421f-b774-235cfd597852"]
	// v := m["doc.typeName"]
	// fmt.Println(v)
	// v = m["children.10"]
	// fmt.Println(v)

	// trackCode(ms, mf, "92b62493-c251-421f-b774-235cfd597852")

	// fmt.Println(hasChild(ms, "92b62493-c251-421f-b774-235cfd597852", "48e7794a-fd99-4bbd-965e-090c24a4ca00"))

	// fmt.Println(isChild(ms, "48e7794a-fd99-4bbd-965e-090c24a4ca00"))

	////////////////////////////////////////////////

	mCodeUrl := make(map[string]string)
	mIdUrl := make(map[string]string)

	Idx := 0
	for _, block := range mNodeIdBlock {

		code := gjson.Get(block, "code").String()
		codes, ids := TrackCode(mNodeIdBlock, code)
		// for i, code := range codes {
		// 	fmt.Println(code, ids[i])
		// }

		MarkUrl(ids, codes, mCodeUrl, mIdUrl)

		fmt.Printf("processed... %d, %v\n", Idx, codes)
		Idx++

		// if Idx == 1000 {
		// 	break
		// }
	}

	fmt.Println(len(mIdUrl))
	fmt.Println(len(mCodeUrl))

	fPathIdUrl := filepath.Join(outDir, "id-url.txt")
	for id, url := range mIdUrl {
		fd.MustAppendFile(fPathIdUrl, []byte(fmt.Sprintf("%s\t%s", id, url)), true)
	}

	fPathCodeUrl := filepath.Join(outDir, "./code-url.txt")
	for id, url := range mCodeUrl {
		fd.MustAppendFile(fPathCodeUrl, []byte(fmt.Sprintf("%s\t%s", id, url)), true)
	}
}
