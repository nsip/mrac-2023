package main

import (
	"fmt"
	"strings"

	fd "github.com/digisan/gotk/file-dir"
	lk "github.com/digisan/logkit"
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
