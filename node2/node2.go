package node2

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	. "github.com/digisan/go-generics/v2"
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

func GenNodeIdBlockMap(data []byte) map[string]string {
	m := make(map[string]string)
	Scan(data, func(i int, id, block string) bool {
		m[id] = block
		return true
	})
	return m
}

func GenNodeCodeBlockMap(data []byte) map[string]string {
	m := make(map[string]string)
	Scan(data, func(i int, id, block string) bool {
		if r := gjson.Get(block, "code"); r.Type != gjson.Null {
			m[r.Str] = block
		} else {
			lk.FailOnErr("%v has no [code]\n%v", id, errors.New(""))
		}
		return true
	})
	return m
}

func GetCodeById(id string, mIdBlock map[string]string) string {
	if block, ok := mIdBlock[id]; ok {
		if rCode := gjson.Get(block, "code"); rCode.Type != gjson.Null {
			return rCode.Str
		}
		lk.FailOnErr("%v", fmt.Errorf("ID: '%s' has no 'code' field", id))
	}
	return ""
}

func GetIdByCode(code string, mCodeBlock map[string]string) string {
	if block, ok := mCodeBlock[code]; ok {
		if rID := gjson.Get(block, "uuid"); rID.Type != gjson.Null {
			return rID.Str
		}
		if rID := gjson.Get(block, "id"); rID.Type != gjson.Null {
			return rID.Str
		}
		lk.FailOnErr("%v", fmt.Errorf("code: '%s' has no 'uuid' or 'id' field", code))
	}
	return ""
}

// not including map[root]***
func GenChildParentMap(data []byte, mIdBlock map[string]string) (mIDChildParent map[string]string, mCodeChildParent map[string]string) {
	mIDChildParent = make(map[string]string)
	mCodeChildParent = make(map[string]string)
	Scan(data, func(i int, id, block string) bool {
		if r := gjson.Get(block, "children"); r.Type != gjson.Null && r.IsArray() {
			for _, rChild := range r.Array() {

				// id: parent ID; idChild: child ID
				idChild := rChild.Str
				mIDChildParent[idChild] = id

				// code in parent
				if rCodeP := gjson.Get(block, "code"); r.Type != gjson.Null {
					if cBlock, ok := mIdBlock[idChild]; ok {
						// code in child
						if rCodeC := gjson.Get(cBlock, "code"); r.Type != gjson.Null {
							mCodeChildParent[rCodeC.Str] = rCodeP.Str
						} else {
							lk.FailOnErr("%v has no [code]\n%v", idChild, errors.New(""))
						}
					} else {
						lk.FailOnErr("%v has no content\n%v", idChild, errors.New(""))
					}
				} else {
					lk.FailOnErr("%v has no [code]\n%v", id, errors.New(""))
				}
			}
		} else {
			// lk.Log("%v has no [children]", id)
		}
		return true
	})
	return
}

func RetrieveAncestryID(id string, mIDChildParent map[string]string) []string {
	Ancestry := []string{id}
AGAIN:
	if pID, ok := mIDChildParent[id]; ok {
		Ancestry = append(Ancestry, pID)
		id = pID
		goto AGAIN
	}
	return Reverse(Ancestry)
}

func RetrieveAncestryCode(code string, mCodeChildParent map[string]string) []string {
	Ancestry := []string{code}
AGAIN:
	if pCode, ok := mCodeChildParent[code]; ok {
		Ancestry = append(Ancestry, pCode)
		code = pCode
		goto AGAIN
	}
	return Reverse(Ancestry)
}

func IsAncestorID(id, ancestor string, mIDChildParent map[string]string) bool {
	ancestry := RetrieveAncestryID(id, mIDChildParent)
	return In(ancestor, ancestry...) && IdxOf(id, ancestry...) > IdxOf(ancestor, ancestry...)
}

func IsAncestorCode(code, ancestor string, mCodeChildParent map[string]string) bool {
	ancestry := RetrieveAncestryCode(code, mCodeChildParent)
	return In(ancestor, ancestry...) && IdxOf(code, ancestry...) > IdxOf(ancestor, ancestry...)
}

//////////////////////////////////////////////////////////////////////////

func UpdateNodeWithMeta(dataNode []byte, URI string, meta map[string]string, outPath string) {
	e := bytes.LastIndexAny(dataNode, "}")
	dataNode = dataNode[:e+1]

	outPath = strings.TrimSuffix(outPath, ".json") + ".json"
	parts := []string{}
	out := ""

	Scan(dataNode, func(i int, id, block string) bool {

		// "uuid": {id} => "id": "http://abc/def/{id}"
		newIdVal := fmt.Sprintf("%s%s", URI, gjson.Get(block, "uuid").String())
		block, _ = sjson.Set(block, "uuid", newIdVal)
		block = strings.Replace(block, `"uuid"`, `"id"`, 1)

		m := make(map[string]interface{})
		json.Unmarshal([]byte(gjson.Get(block, "connections").String()), &m)

		for k, v := range m {
			// "abcdeft" => "Levels" etc.
			block = strings.Replace(block, k, meta[k], 1)
			// "abc-def" => "http://abc/def/{id}"
			for _, a := range v.([]interface{}) {
				block = strings.Replace(block, a.(string), fmt.Sprintf("%s%s", URI, a), 1)
			}
		}

		part := fmt.Sprintf(`"%s": %s`, id, block)
		parts = append(parts, part)
		return true
	})

	out = "{" + strings.Join(parts, ",") + "}"
	out = jt.FmtStr(out, "  ")

	lk.FailOnErrWhen(!dt.IsJSON([]byte(out)), "%v", errors.New("invalid JSON from node & meta"))
	os.WriteFile(outPath, []byte(out), os.ModePerm)
}

//////////////////////////////////////////////////////////////////////////

func MakeIdUrlText(mIdBlock, mCodeBlock, mIDChildParent, mCodeChildParent map[string]string, outPath4IdUrl, outPath4CodeUrl string) {

	var (
		mIdUrl   = make(map[string]string)
		mCodeUrl = make(map[string]string)
	)

	for code := range mCodeBlock {

		url := ""
		ancestors := RetrieveAncestryCode(code, mCodeChildParent)

		switch {
		case len(ancestors) >= 3:
			code := ancestors[2]
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

			default:
				lk.Warn("code '%v' is missing its url (2)")
			}

		case len(ancestors) == 2:
			code := ancestors[1]
			switch code {
			case "AS", "LA":
				url = "http://vocabulary.curriculum.edu.au/MRAC/LA/"
			case "GC":
				url = "http://vocabulary.curriculum.edu.au/MRAC/GC/"
			case "CCP":
				url = "http://vocabulary.curriculum.edu.au/MRAC/CCP/"
			default:
				lk.Warn("code '%v' is missing its url (1)")
			}

		case len(ancestors) == 1:
			code := ancestors[0]
			switch code {
			case "root":
				url = "http://vocabulary.curriculum.edu.au/MRAC/"
			default:
				lk.Warn("code '%v' is missing its url (0)")
			}

		default:
		}

		mCodeUrl[code] = url
		mIdUrl[GetIdByCode(code, mCodeBlock)] = url
	}

	for id, url := range mIdUrl {
		fd.MustAppendFile(outPath4IdUrl, []byte(fmt.Sprintf("%s\t%s", id, url)), true)
	}
	for code, url := range mCodeUrl {
		fd.MustAppendFile(outPath4CodeUrl, []byte(fmt.Sprintf("%s\t%s", code, url)), true)
	}
}
