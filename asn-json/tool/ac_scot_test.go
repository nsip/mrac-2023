package tool

import (
	"fmt"
	"testing"
)

func TestAcScot(t *testing.T) {
	m := getAcScotMap("../data/SCOT_20231110.txt")
	fmt.Println(len(m))
	fmt.Println(m["AC9ADA10C01"])
}

func TestScotJsonLd(t *testing.T) {
	m := scanScotJsonLd("../data/pp_project_schoolsonlinethesaurus.jsonld")
	fmt.Println(len(m))
	for _, v := range m["http://vocabulary.curriculum.edu.au/scot/2184"] {
		fmt.Println(v)
	}
}

func TestGetAsnConceptTerm(t *testing.T) {
	m := GetAsnConceptTerm("../data/SCOT_20231110.txt", "../data/pp_project_schoolsonlinethesaurus.jsonld")
	fmt.Println(m["AC9HC10K05"])
}
