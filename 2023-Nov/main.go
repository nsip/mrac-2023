package main

import (
	"fmt"

	ct "github.com/digisan/csv-tool"
	fd "github.com/digisan/gotk/file-dir"
	lk "github.com/digisan/logkit"
)

// make scot txt, "mrac\turl"
func main() {
	const url = "http://vocabulary.curriculum.edu.au/scot/"
	const mapFile = "mapping-20231110.csv"
	const outFile = "./SCOT_20231110.txt"
	ct.ScanFile(mapFile, func(i, n int, headers, items []string) (ok bool, hdr string, row string) {
		scot, mrac := "", ""
		for i, hdr := range headers {
			if hdr == "scot" {
				// fmt.Println(i, items[i])
				scot = items[i]
			}
			if hdr == "mrac" {
				// fmt.Println(i, items[i])
				mrac = items[i]
			}
		}
		lk.FailOnErrWhen(len(scot) == 0 || len(mrac) == 0, "%v", fmt.Errorf("scot or mrac missing"))
		fd.MustAppendFile(outFile, []byte(fmt.Sprintf("%s\t%s", mrac, url+scot)), true)
		return true, "", ""
	}, false, "")
}
