package node2

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

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
func GenChildParentMap(dataNode []byte, mIdBlock map[string]string) (mIDChildParent map[string]string, mCodeChildParent map[string]string) {
	mIDChildParent = make(map[string]string)
	mCodeChildParent = make(map[string]string)
	Scan(dataNode, func(i int, id, block string) bool {
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

// func RetrieveAncestryID(id string, mIDChildParent map[string]string) []string {
// 	Ancestry := []string{id}
// AGAIN:
// 	if pID, ok := mIDChildParent[id]; ok {
// 		Ancestry = append(Ancestry, pID)
// 		id = pID
// 		goto AGAIN
// 	}
// 	return Reverse(Ancestry)
// }

// func RetrieveAncestryCode(code string, mCodeChildParent map[string]string) []string {
// 	Ancestry := []string{code}
// AGAIN:
// 	if pCode, ok := mCodeChildParent[code]; ok {
// 		Ancestry = append(Ancestry, pCode)
// 		code = pCode
// 		goto AGAIN
// 	}
// 	return Reverse(Ancestry)
// }

func RetrieveAncestry(IdOrCode string, mIdOrCodeChildParent map[string]string) []string {
	Ancestry := []string{IdOrCode}
AGAIN:
	if pIdOrCode, ok := mIdOrCodeChildParent[IdOrCode]; ok {
		Ancestry = append(Ancestry, pIdOrCode)
		IdOrCode = pIdOrCode
		goto AGAIN
	}
	return Reverse(Ancestry)
}

func RetrieveAncestryAsCodeById(Id string, mIdChildParent, mIdBlock map[string]string) (codes []string) {
	ancestors := RetrieveAncestry(Id, mIdChildParent)
	for _, ancestor := range ancestors {
		codes = append(codes, GetCodeById(ancestor, mIdBlock))
	}
	return
}

// func IsAncestorID(id, ancestor string, mIDChildParent map[string]string) bool {
// 	ancestry := RetrieveAncestryID(id, mIDChildParent)
// 	return In(ancestor, ancestry...) && IdxOf(id, ancestry...) > IdxOf(ancestor, ancestry...)
// }

// func IsAncestorCode(code, ancestor string, mCodeChildParent map[string]string) bool {
// 	ancestry := RetrieveAncestryCode(code, mCodeChildParent)
// 	return In(ancestor, ancestry...) && IdxOf(code, ancestry...) > IdxOf(ancestor, ancestry...)
// }

func IsAncestorCode(IdOrCode, ancestor string, mIdOrCodeChildParent map[string]string) bool {
	ancestry := RetrieveAncestry(IdOrCode, mIdOrCodeChildParent)
	return In(ancestor, ancestry...) && IdxOf(IdOrCode, ancestry...) > IdxOf(ancestor, ancestry...)
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

var (
	yyyy string
	mm   string
)

func fetchTS(fPathOfTree string) (yyyy, mm string) {
	data, err := os.ReadFile(fPathOfTree)
	if err != nil {
		log.Fatal(err)
	}
	layout := "2006-01-02T15:04:05.000Z"
	ts := gjson.Get(string(data), "created_at").String()
	t, err := time.Parse(layout, ts)
	if err != nil {
		log.Fatal(err)
	}
	ts = t.Format("2006-01-02")
	ss := strings.Split(ts, "-")
	return ss[0], ss[1]
}

func SetTS4Url(fPathOfTree string) {
	yyyy, mm = fetchTS(fPathOfTree)
}

func MakeIdUrlText(mIdBlock, mCodeBlock, mIDChildParent, mCodeChildParent map[string]string, outPath4IdUrl, outPath4CodeUrl string) {

	lk.FailOnErrWhen(len(yyyy) == 0 || len(mm) == 0, "%v", errors.New("TimeStamp is empty for URL, 'SetTS4Url' before 'MakeIdUrlText'"))

	var (
		mIdUrl   = make(map[string]string)
		mCodeUrl = make(map[string]string)
	)

	for code := range mCodeBlock {

		url := ""
		ancestors := RetrieveAncestry(code, mCodeChildParent)

		switch {
		case len(ancestors) >= 3:
			code := ancestors[2]
			switch code {
			case "HAS", "HASS", "ASHAS", "ASHASS":
				url = fmt.Sprintf("http://vocabulary.curriculum.edu.au/MRAC/%s/%s/LA/HASS/", yyyy, mm)
			case "ENG", "ASENG":
				url = fmt.Sprintf("http://vocabulary.curriculum.edu.au/MRAC/%s/%s/LA/ENG/", yyyy, mm)
			case "LAN", "ASLAN":
				url = fmt.Sprintf("http://vocabulary.curriculum.edu.au/MRAC/%s/%s/LA/LAN/", yyyy, mm)
			case "SCI", "ASSCI":
				url = fmt.Sprintf("http://vocabulary.curriculum.edu.au/MRAC/%s/%s/LA/SCI/", yyyy, mm)
			case "ART", "ASART":
				url = fmt.Sprintf("http://vocabulary.curriculum.edu.au/MRAC/%s/%s/LA/ART/", yyyy, mm)
			case "HPE", "ASHPE":
				url = fmt.Sprintf("http://vocabulary.curriculum.edu.au/MRAC/%s/%s/LA/HPE/", yyyy, mm)
			case "MAT", "ASMAT":
				url = fmt.Sprintf("http://vocabulary.curriculum.edu.au/MRAC/%s/%s/LA/MAT/", yyyy, mm)
			case "TEC", "ASTEC":
				url = fmt.Sprintf("http://vocabulary.curriculum.edu.au/MRAC/%s/%s/LA/TEC/", yyyy, mm)

			case "CCT":
				url = fmt.Sprintf("http://vocabulary.curriculum.edu.au/MRAC/%s/%s/GC/CCT/", yyyy, mm)
			case "N":
				url = fmt.Sprintf("http://vocabulary.curriculum.edu.au/MRAC/%s/%s/GC/N/", yyyy, mm)
			case "DL":
				url = fmt.Sprintf("http://vocabulary.curriculum.edu.au/MRAC/%s/%s/GC/DL/", yyyy, mm)
			case "L":
				url = fmt.Sprintf("http://vocabulary.curriculum.edu.au/MRAC/%s/%s/GC/L/", yyyy, mm)
			case "PSC":
				url = fmt.Sprintf("http://vocabulary.curriculum.edu.au/MRAC/%s/%s/GC/PSC/", yyyy, mm)
			case "IU":
				url = fmt.Sprintf("http://vocabulary.curriculum.edu.au/MRAC/%s/%s/GC/IU/", yyyy, mm)
			case "EU":
				url = fmt.Sprintf("http://vocabulary.curriculum.edu.au/MRAC/%s/%s/GC/EU/", yyyy, mm)

			case "AA":
				url = fmt.Sprintf("http://vocabulary.curriculum.edu.au/MRAC/%s/%s/CCP/AA/", yyyy, mm)
			case "S":
				url = fmt.Sprintf("http://vocabulary.curriculum.edu.au/MRAC/%s/%s/CCP/S/", yyyy, mm)
			case "A_TSI":
				url = fmt.Sprintf("http://vocabulary.curriculum.edu.au/MRAC/%s/%s/CCP/A_TSI/", yyyy, mm)

			default:
				lk.Warn("code '%v' is missing its url (2)")
			}

		case len(ancestors) == 2:
			code := ancestors[1]
			switch code {
			case "AS", "LA":
				url = fmt.Sprintf("http://vocabulary.curriculum.edu.au/MRAC/%s/%s/LA/", yyyy, mm)
			case "GC":
				url = fmt.Sprintf("http://vocabulary.curriculum.edu.au/MRAC/%s/%s/GC/", yyyy, mm)
			case "CCP":
				url = fmt.Sprintf("http://vocabulary.curriculum.edu.au/MRAC/%s/%s/CCP/", yyyy, mm)
			default:
				lk.Warn("code '%v' is missing its url (1)")
			}

		case len(ancestors) == 1:
			code := ancestors[0]
			switch code {
			case "root":
				url = fmt.Sprintf("http://vocabulary.curriculum.edu.au/MRAC/%s/%s/", yyyy, mm)
			default:
				lk.Warn("code '%v' is missing its url (0)")
			}

		default:
		}

		mCodeUrl[code] = url
		mIdUrl[GetIdByCode(code, mCodeBlock)] = url
	}

	os.RemoveAll(outPath4IdUrl)
	for id, url := range mIdUrl {
		fd.MustAppendFile(outPath4IdUrl, []byte(fmt.Sprintf("%s\t%s", id, url)), true)
	}

	os.RemoveAll(outPath4CodeUrl)
	for code, url := range mCodeUrl {
		fd.MustAppendFile(outPath4CodeUrl, []byte(fmt.Sprintf("%s\t%s", code, url)), true)
	}
}
