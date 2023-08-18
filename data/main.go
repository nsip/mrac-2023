package main

import (
	"fmt"
	"os"

	"github.com/digisan/gotk/strs"
	u "github.com/nsip/mrac-2023/util"
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
			if formatted, err := u.FmtJSONFile(fName); err == nil {
				if err := os.WriteFile(fName, []byte(formatted), os.ModePerm); err != nil {
					fmt.Printf("%v", err)
					return
				}
			}
			// bytes, err := os.ReadFile(fName)
			// if err != nil {
			// 	fmt.Printf("%v", err)
			// 	return
			// }
			// if formatted, err := u.FmtJSON(string(bytes)); err == nil {
			// 	if err := os.WriteFile(fName, []byte(formatted), os.ModePerm); err != nil {
			// 		fmt.Printf("%v", err)
			// 		return
			// 	}
			// }
		}
	}
}
