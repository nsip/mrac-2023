package main

import (
	"fmt"
	"os"

	"github.com/digisan/gotk/strs"
	jt "github.com/digisan/json-tool"
)

func main() {
	de, err := os.ReadDir("./")
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	for _, f := range de {
		fName := f.Name()
		if strs.HasAnySuffix(fName, ".json", ".jsonld") {
			bytes, err := os.ReadFile(fName)
			if err != nil {
				fmt.Printf("%v", err)
				return
			}
			bytes = jt.Fmt(bytes, "  ")
			err = os.WriteFile(fName, bytes, os.ModePerm)
			if err != nil {
				fmt.Printf("%v", err)
				return
			}
		}
	}
}
