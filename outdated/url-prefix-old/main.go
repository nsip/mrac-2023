package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	fd "github.com/digisan/gotk/file-dir"
	"github.com/tidwall/gjson"
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

// Tree Path
func FetchTime(fPath string) (yyyy, mm string) {
	data, err := os.ReadFile(fPath)
	if err != nil {
		log.Fatal(err)
	}
	layout := "2006-01-02T15:04:05.000Z"
	ts := gjson.Get(string(data), "created_at").String()
	t, err := time.Parse(layout, ts)
	if err != nil {
		log.Fatal(err)
	}
	ts = t.Format("2006-01-02")
	ss := strings.Split(ts, "-")
	return ss[0], ss[1]
}

func main() {

	// rUUID := regexp.MustCompile(`[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`)

	mCodeUrl := LoadUrl("../data/code-url.txt")
	mIdUrl := LoadUrl("../data/id-url.txt")

	fmt.Println(len(mCodeUrl))
	fmt.Println(len(mIdUrl))

	///////////////////////////////////////

	r := regexp.MustCompile(`http://vocabulary.curriculum.edu.au/+[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`)

	base := "http://vocabulary.curriculum.edu.au/MRAC"
	yyyy, mm := FetchTime("../data/Sofia-API-Tree-Data-18072023.json")
	prefix := fmt.Sprintf("%s/%s/%s", base, yyyy, mm)
	urlModify := func(url string) string {
		return strings.Replace(url, base, prefix, 1)
	}

	///////////////////////////////////////

	dirIn := "../data-out/asn-json"
	dirOut := "../data-out/asn-json/url"
	fd.MustCreateDir(dirOut)

	fs, err := os.ReadDir(dirIn)
	if err != nil {
		panic(err)
	}
	for _, f := range fs {

		if f.IsDir() || filepath.Ext(f.Name()) != ".json" {
			continue
		}
		// if f.Name() != "la-Languages.json" {
		// 	continue
		// }

		fPath := filepath.Join(dirIn, f.Name())
		// fmt.Println(fPath)

		data, err := os.ReadFile(fPath)
		if err != nil {
			panic(err)
		}
		str := string(data)

		fd := r.FindAllString(str, -1)
		fmt.Println(fPath, len(fd))

		str = r.ReplaceAllStringFunc(str, func(s string) string {
			id := s[len(s)-36:]
			if url, ok := mIdUrl[id]; ok {
				url = urlModify(url) + id
				return url
			}
			return s
		})

		fPath = filepath.Join(dirOut, f.Name())
		if err := os.WriteFile(fPath, []byte(str), os.ModePerm); err != nil {
			panic(err)
		}
	}

	///////////////////////////////////////

	fmt.Println("-------------------------------------------------------")

	dirIn = "../data-out/asn-json-ld"
	dirOut = "../data-out/asn-json-ld/url"
	fd.MustCreateDir(dirOut)

	fs, err = os.ReadDir(dirIn)
	if err != nil {
		panic(err)
	}
	for _, f := range fs {

		if f.IsDir() || filepath.Ext(f.Name()) != ".json" {
			continue
		}
		// if f.Name() != "la-Languages.json" {
		// 	continue
		// }

		fPath := filepath.Join(dirIn, f.Name())
		// fmt.Println(fPath)

		data, err := os.ReadFile(fPath)
		if err != nil {
			panic(err)
		}
		str := string(data)

		fd := r.FindAllString(str, -1)
		fmt.Println(fPath, len(fd))

		str = r.ReplaceAllStringFunc(str, func(s string) string {
			id := s[len(s)-36:]
			if url, ok := mIdUrl[id]; ok {
				url = urlModify(url) + id
				return url
			}
			return s
		})

		fPath = filepath.Join(dirOut, f.Name())
		if err := os.WriteFile(fPath, []byte(str), os.ModePerm); err != nil {
			panic(err)
		}
	}
}
