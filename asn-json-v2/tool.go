package main

import (
	"fmt"
	"strings"

	. "github.com/digisan/go-generics"
	fd "github.com/digisan/gotk/file-dir"
	lk "github.com/digisan/logkit"
	"github.com/nsip/mrac-2023/node2"
)

func LoadIdPrefLbl(fPath string) map[string]string {

	lk.FailOnErrWhen(!fd.FileExists(fPath), "%v", fmt.Errorf("id-preflabel.txt (%s) [made from scot.txt & .jsonld] doesn't exist", fPath))

	rt := make(map[string]string)
	fd.FileLineScan(fPath, func(line string) (bool, string) {
		kv := strings.SplitN(line, "\t", 2)
		rt[kv[0]] = kv[1]
		return true, ""
	}, "")
	return rt
}

func kvStrJoin(kvStrGrp ...string) string {
	nonEmptyStrGrp := []string{}
	for _, kvStr := range kvStrGrp {
		if strings.Trim(kvStr, " \t\n") != "" {
			nonEmptyStrGrp = append(nonEmptyStrGrp, kvStr)
		}
	}
	return strings.Join(nonEmptyStrGrp, ",")
}

func loadUrl(fPath string) map[string]string {
	m := make(map[string]string)
	fd.FileLineScan(fPath, func(line string) (bool, string) {
		ss := strings.Split(line, "\t")
		m[ss[0]] = ss[1]
		return true, ""
	}, "")
	return m
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

func GetAncestorTitle(code, group string, mCodeChildParent map[string]string) string {

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

	ancestors := node2.RetrieveAncestry(code, mCodeChildParent)
	for _, a := range append(ancestors, code) {
		if title, ok := m[a]; ok {
			return title
		}
	}
	return ""
}

func GetCodeType(code string, mCodeChildParent map[string]string) string {
	ancestors := node2.RetrieveAncestry(code, mCodeChildParent)
	for _, ancestor := range ancestors {
		if In(ancestor, "LA", "AS", "GC", "CCP") {
			return ancestor
		}
	}
	return ""
}

func GetIdType(id string, mIdBlock, mCodeChildParent map[string]string) string {
	code := node2.GetCodeById(id, mIdBlock)
	return GetCodeType(code, mCodeChildParent)
}
