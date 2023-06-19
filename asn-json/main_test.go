package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/digisan/gotk/strs"
)

func TestMain(t *testing.T) {
	main()
}

func TestFixEmptyArray(t *testing.T) {
	fPath := `../data-out/asn-json/la-Languages.json`
	data, err := os.ReadFile(fPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(data))

	out, err := strs.StrLineScan(string(data), func(line string) (bool, string) {
		return !strs.HasAnySuffix(line, "[]", "[],"), line
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := os.WriteFile(`../data-out/asn-json/la-Languages-fix.json`, []byte(out), os.ModePerm); err != nil {
		fmt.Println(err)
		return
	}
}

func TestFixDup_asn_hasLevel(t *testing.T) {

	const (
		// aim = `"asn_hasLevel": [`
		aim = `"asn_crossSubjectReference": [`
	)

	func() {

		fPath := `../data-out/asn-json/la-Languages-fix.json`
		data, err := os.ReadFile(fPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(len(data))

		out, err := strs.StrLineScanEx(string(data), 20, 20, "*****", func(line string, cache []string) (bool, string) {
			if ln := strings.TrimSpace(line); ln == "]," {
				if c11 := strings.TrimSpace(cache[21]); c11 == aim {
					if above := strings.Join(cache[:20], "\n"); strings.Count(above, aim) == 1 {
						// fmt.Println(c11)
						return false, ""
					}
				}
			}
			return true, line
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := os.WriteFile(`../data-out/asn-json/la-Languages-fix-temp.json`, []byte(out), os.ModePerm); err != nil {
			fmt.Println(err)
			return
		}
	}()

	/////////////////////////////////////////////////////////////////////////////////////////////

	func() {

		fPath := `../data-out/asn-json/la-Languages-fix-temp.json`
		data, err := os.ReadFile(fPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer os.RemoveAll(fPath)

		fmt.Println(len(data))

		out, err := strs.StrLineScanEx(string(data), 20, 20, "*****", func(line string, cache []string) (bool, string) {
			if ln := strings.TrimSpace(line); ln == aim {
				if above := strings.Join(cache[:20], "\n"); strings.Count(above, aim) == 1 {
					return true, ","
				}
			}
			return true, line
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		// out = jt.FmtStr(out, "  ") // Fmt overwrites duplicated fields
		if err = os.WriteFile(`../data-out/asn-json/la-Languages-fix.json`, []byte(out), os.ModePerm); err != nil {
			fmt.Println(err)
			return
		}
	}()
}
