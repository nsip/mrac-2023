package tool

import (
	"fmt"
	"testing"
)

func TestAcScot(t *testing.T) {
	m := getAcScotMap("../../data/SCOT_20231110.txt")
	fmt.Println(len(m))
	fmt.Println(m["AC9M1N06"])
}

//////////////////////////////////////////////////////////////////////////////////////

func TestScanSCOT(t *testing.T) {
	m := scanSCOT("../../data/pp_project_schoolsonlinethesaurus.jsonld")
	fmt.Println(len(m))
	for _, v := range m["http://vocabulary.curriculum.edu.au/scot/2184"] {
		fmt.Println(v)
	}
}

//////////////////////////////////////////////////////////////////////////////////////

// func TestScotJsonLd(t *testing.T) {
// 	// m := scanScotJsonLd("../../data/pp_project_schoolsonlinethesaurus.jsonld")
// 	m := scanScotJsonLd("../../_obsolete/scot.jsonld")
// 	fmt.Println(len(m))
// 	for _, v := range m["http://vocabulary.curriculum.edu.au/scot/2184"] {
// 		fmt.Println(v)
// 	}
// }

func TestGetAsnConceptTerm(t *testing.T) {

	m := GetAsnConceptTerm("../../data/SCOT_20231110.txt", "../../data/pp_project_schoolsonlinethesaurus.jsonld")
	fmt.Println(len(m))
	fmt.Println(m["AC9M1N06"])

	//
	// *** create id-preflabel.txt ***
	//
	// const out = "id-preflabel.txt"
	// const sep = "\t"
	// os.RemoveAll(out)
	// for k, v := range m {
	// 	// fmt.Println(k, v)
	// 	fd.MustAppendFile(out, []byte(strings.Join([]string{k, v}, sep)), true)
	// }
}

func TestLoadIdPrefLbl(t *testing.T) {
	const in = "id-preflabel.txt"
	m := LoadIdPrefLbl(in)
	fmt.Println(len(m))
	code := "AC9M1N06"
	fmt.Printf("%s: %s\n", code, m[code])
}
