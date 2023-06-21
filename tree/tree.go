package tree

import (
	"fmt"
	"os"
	"strings"

	dt "github.com/digisan/gotk/data-type"
	fd "github.com/digisan/gotk/file-dir"
	"github.com/digisan/gotk/strs"
	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
	. "github.com/nsip/mrac-2023/tree/sub"
)

func LoadUrl(fPath string) map[string]string {
	m := make(map[string]string)
	fd.FileLineScan(fPath, func(line string) (bool, string) {
		ss := strings.Split(line, "\t")
		m[ss[0]] = ss[1]
		return true, ""
	}, "")
	return m
}

func Partition(js, outDir string, mMeta map[string]string) {

	fileContent := CCP(js, outDir)
	err := os.WriteFile(fmt.Sprintf("./%s/ccp-%s.json", outDir, "Cross-curriculum Priorities"), []byte(fileContent), os.ModePerm)
	lk.FailOnErr("%v", err)

	for gc, fileContent := range GC(js) {
		err = os.WriteFile(fmt.Sprintf("./%s/gc-%s.json", outDir, gc), []byte(fileContent), os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}

	for la, fileContent := range LA(js) {
		err := os.WriteFile(fmt.Sprintf("./%s/la-%s.json", outDir, la), []byte(fileContent), os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}

	//////////////////////////////////////////////////////////////

	var (
		mProcFlag = map[string]bool{
			"la-English":                        true,
			"la-Humanities and Social Sciences": true,
			"la-Health and Physical Education":  true,
			"la-Languages":                      true,
			"la-Mathematics":                    true,
			"la-Science":                        true,
			"la-The Arts":                       true,
			"la-Technologies":                   true,
		}

		// mCodeUrl = LoadUrl("../data/code-url.txt")
		mIdUrl = LoadUrl("../data/id-url.txt")

		// mUrlID = map[string]string{
		// 	"la-English":                        "http://vocabulary.curriculum.edu.au/MRAC/LA/ENG/",
		// 	"la-Humanities and Social Sciences": "http://vocabulary.curriculum.edu.au/MRAC/LA/HASS/",
		// 	"la-Health and Physical Education":  "http://vocabulary.curriculum.edu.au/MRAC/LA/HPE/",
		// 	"la-Languages":                      "http://vocabulary.curriculum.edu.au/MRAC/LA/LAN/",
		// 	"la-Mathematics":                    "http://vocabulary.curriculum.edu.au/MRAC/LA/MAT/",
		// 	"la-Science":                        "http://vocabulary.curriculum.edu.au/MRAC/LA/SCI/",
		// 	"la-The Arts":                       "http://vocabulary.curriculum.edu.au/MRAC/LA/ART/",
		// 	"la-Technologies":                   "http://vocabulary.curriculum.edu.au/MRAC/LA/TEC/",
		// }
	)

	for fName, proc := range mProcFlag {
		if !proc {
			continue
		}

		in := fmt.Sprintf("./%s/%s.json", outDir, fName)
		lk.Log("Processing... %s", in)

		data, err := os.ReadFile(in)
		lk.WarnOnErr("%v", err)
		if err != nil {
			return
		}

		js := ReStruct(string(data))

		js = ConnectionsFieldMapping(js, mIdUrl, mMeta)

		if len(js) > 0 {
			lk.FailOnErrWhen(!dt.IsJSON([]byte(js)), "%v", fmt.Errorf("invalid JSON from [ReStruct %s]", fName))
			js = jt.FmtStr(js, "  ")

			// remove unwanted line
			js, err = strs.StrLineScan(js, func(line string) (bool, string) {
				trimmed := strings.TrimSpace(line)
				if strs.HasAnySuffix(trimmed, `: [],`, `: []`, `: "",`, `: ""`) {
					return false, ""
				}
				return true, line
			})
			if err != nil {
				return
			}

			out := fmt.Sprintf("./%s/restructure/%s.json", outDir, fName)
			fd.MustWriteFile(out, []byte(js))
		}
	}
}
