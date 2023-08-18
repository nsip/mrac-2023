package util

import (
	"fmt"
	"os"
	"os/exec"

	. "github.com/digisan/go-generics/v2"
	fd "github.com/digisan/gotk/file-dir"
	lk "github.com/digisan/logkit"
)

// func FmtJSON(str string) (string, error) {
// 	jsFmt := "fmt.js"
// 	lk.FailOnErrWhen(!fd.FileExists(jsFmt), "%v", fmt.Errorf("%v is not found", jsFmt))

// 	cmd := exec.Command("node", jsFmt, str)
// 	output, err := cmd.Output()
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return "", err
// 	}
// 	return string(output), nil
// }

func FmtJSON(str string) (string, error) {
	fTemp4fmt := "/tmp/json-fmt.json"
	fd.MustWriteFile(fTemp4fmt, StrToConstBytes(str))
	defer os.RemoveAll(fTemp4fmt)
	return FmtJSONFile(fTemp4fmt)
}

func FmtJSONFile(fPath string) (string, error) {
	jsFmtFile := "fmtfile.js"
	lk.FailOnErrWhen(!fd.FileExists(jsFmtFile), "%v", fmt.Errorf("%v is not found", jsFmtFile))

	cmd := exec.Command("node", jsFmtFile, fPath)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return ConstBytesToStr(output), nil
}
