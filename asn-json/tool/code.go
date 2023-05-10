package tool

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

func findCodeParent(code string, start int, mParent map[string]int) string {
	pCode, pDis := "", math.MaxInt32
	for code, pStart := range mParent {
		if dis := start - pStart; dis > 0 && dis < pDis {
			pCode = code
			pDis = dis
		}
	}
	return pCode
}

// input json MUST be well formatted, indent is 2 spaces.
func GetCodeParentMap(data []byte) map[string]string {

	js := string(data)

	// check
	if !strings.HasPrefix(js, "{\n  \"code\": \"root\",") &&
		!strings.HasPrefix(js, "{\r\n  \"code\": \"root\",") {
		panic("input json MUST be well formatted, indent is 2 spaces.")
	}

	reCodes := []*regexp.Regexp{}
	for i := 6; i <= 54; i += 4 {
		reCodes = append(reCodes, regexp.MustCompile(fmt.Sprintf(`\n[ ]{%d}"code":(\s)*"[^"]+"`, i)))
	}

	mCodeChildren := make(map[string][]string)
	mCodeParent := make(map[string]string)
	mLvlStartGrp := []map[string]int{}

	for _, r := range reCodes {

		locGrp := r.FindAllStringIndex(js, -1)

		starts := []int{}
		codes := []string{}

		for _, loc := range locGrp {
			s, e := loc[0], loc[1]
			codeln := strings.Trim(js[s:e], "\n ")
			codeln = codeln[:len(codeln)-1]
			n := strings.LastIndex(codeln, "\"")
			code := codeln[n+1:]
			codes = append(codes, code)
			starts = append(starts, s)
		}

		starts = append(starts, len(js))

		mLvlStart := make(map[string]int)
		for i, code := range codes {
			mLvlStart[code] = starts[i]
		}

		mLvlStartGrp = append(mLvlStartGrp, mLvlStart)
	}

	for i := len(mLvlStartGrp) - 2; i >= 0; i-- {
		mThis, mParent := mLvlStartGrp[i+1], mLvlStartGrp[i]
		if len(mThis) == 0 {
			continue
		}
		for code, v := range mThis {
			pCode := findCodeParent(code, v, mParent)
			mCodeChildren[pCode] = append(mCodeChildren[pCode], code)
			mCodeParent[code] = pCode
		}

		// dump top levels 'code' to find related 'title'
		// if i == 0 || i == 1 {
		// 	keys, _ := tsi.Map2KVs(mParent, func(i, j string) bool { return i < j }, nil)
		// 	fmt.Println(keys)
		// 	fmt.Println("------------------------------------------------------------")
		// }
	}

	return mCodeParent
}

func GetCodeAncestors(mCodeParent map[string]string, code string) (ancestors []string) {
	for pCode, ok := mCodeParent[code]; ok; pCode, ok = mCodeParent[pCode] {
		ancestors = append(ancestors, pCode)
	}
	return
}

func GetCodeAncestor(mCodeParent map[string]string, code string, level int) string {
	ancestors := GetCodeAncestors(mCodeParent, code)
	ancestors = append([]string{code}, ancestors...)
	index := len(ancestors) - level - 1
	return ancestors[index]
}

func GetAncestorTitle(mCodeParent map[string]string, code, group string) string {

	m := mCodeTitle1
	switch group {
	case "LA":
		m = mLATitle
	case "AS":
		m = mASTitle
	case "CCP":
		m = mCCPTitle
	case "GC":
		m = mGCTitle
	default:
		// log.Fatalln("[group] can only be 'LA', 'AS', 'CCP', 'GC'")
	}

	ancestors := GetCodeAncestors(mCodeParent, code)
	for _, a := range append(ancestors, code) {
		if title, ok := m[a]; ok {
			return title
		}
	}
	return ""
}

var (
	mCodeTitle0 = map[string]string{
		"AS":  "Achievement Standards",
		"CCP": "Cross-curriculum Priorities",
		"GC":  "General Capabilities",
		"LA":  "Learning Areas",
	}

	mCodeTitle1 = map[string]string{
		"AA":    "Asia and Australia's Engagement with Asia",
		"ART":   "The Arts",
		"ASART": "The Arts",
		"ASENG": "English",
		"ASHAS": "Humanities and Social Sciences",
		"ASHPE": "Health and Physical Education",
		"ASLAN": "Languages",
		"ASMAT": "Mathematics",
		"ASSCI": "Science",
		"ASTEC": "Technologies",
		"A_TSI": "Aboriginal and Torres Strait Islander Histories and Cultures",
		"CCT":   "Critical and Creative Thinking",
		"DL":    "Digital Literacy",
		"ENG":   "English",
		"EU":    "Ethical Understanding",
		"HAS":   "Humanities and Social Sciences",
		"HPE":   "Health and Physical Education",
		"IU":    "Intercultural Understanding",
		"L":     "National Literacy Learning Progression",
		"LAN":   "Languages",
		"MAT":   "Mathematics",
		"N":     "National Numeracy Learning Progression",
		"PSC":   "Personal and Social Capability",
		"S":     "Sustainability",
		"SCI":   "Science",
		"TEC":   "Technologies",
	}

	mLATitle = map[string]string{
		"HAS": mCodeTitle1["HAS"],
		"HPE": mCodeTitle1["HPE"],
		"ART": mCodeTitle1["ART"],
		"LAN": mCodeTitle1["LAN"],
		"ENG": mCodeTitle1["ENG"],
		"SCI": mCodeTitle1["SCI"],
		"TEC": mCodeTitle1["TEC"],
		"MAT": mCodeTitle1["MAT"],
	}

	mASTitle = map[string]string{
		"ASLAN": mCodeTitle1["ASLAN"],
		"ASHPE": mCodeTitle1["ASHPE"],
		"ASART": mCodeTitle1["ASART"],
		"ASENG": mCodeTitle1["ASENG"],
		"ASHAS": mCodeTitle1["ASHAS"],
		"ASMAT": mCodeTitle1["ASMAT"],
		"ASSCI": mCodeTitle1["ASSCI"],
		"ASTEC": mCodeTitle1["ASTEC"],
	}

	mCCPTitle = map[string]string{
		"A_TSI": mCodeTitle1["A_TSI"],
		"AA":    mCodeTitle1["AA"],
		"S":     mCodeTitle1["S"],
	}

	mGCTitle = map[string]string{
		"DL":  mCodeTitle1["DL"],
		"IU":  mCodeTitle1["IU"],
		"CCT": mCodeTitle1["CCT"],
		"L":   mCodeTitle1["L"],
		"PSC": mCodeTitle1["PSC"],
		"EU":  mCodeTitle1["EU"],
		"N":   mCodeTitle1["N"],
	}
)
