package util

import (
	"fmt"
	"os"
	"testing"

	. "github.com/digisan/go-generics"
)

func TestFmtJSON(t *testing.T) {
	js := `{"abc":true, "def": 1000.2}`
	jsFmt, err := FmtJSON(js)
	if err == nil {
		fmt.Println(jsFmt)
	}
}

func TestFmtJSONFile(t *testing.T) {
	fName := "Sofia-API-Tree-Data-18072023"
	fIn := fmt.Sprintf("../data/original/%s.json", fName)
	fOut := fmt.Sprintf("./%s-fmt.json", fName)

	jsFmt, err := FmtJSONFile(fIn)
	if err == nil {
		os.WriteFile(fOut, StrToConstBytes(jsFmt), os.ModePerm)
		fmt.Println(len(jsFmt))
	}
}
