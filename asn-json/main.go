package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	fd "github.com/digisan/gotk/file-dir"
	"github.com/digisan/gotk/track"
	jt "github.com/digisan/json-tool"
	"github.com/nsip/mrac-2023/asn-json/tool"
)

const (
	metaFile = "../data/Sofia-API-Meta-Data-06062023.json"
	nodeFile = "../data/Sofia-API-Node-Data-06062023.json"
	treeFile = "../data/Sofia-API-Tree-Data-06062023.json"
	nmFile   = "../data-out/node-meta.json"
)

var mEscStr = map[string]string{
	`\n`: "*LF*",
	`\"`: "*DQ*",
}

func removeEsc(js string) string {
	for esc, str := range mEscStr {
		js = strings.ReplaceAll(js, esc, str)
	}
	return js
}

func restoreEsc(js string) string {
	for esc, str := range mEscStr {
		js = strings.ReplaceAll(js, str, esc)
	}
	return js
}

func main() {
	defer track.TrackTime(time.Now())

	{
		outDir := "../data-out/asn-json"
		outFile := "asn-node.json"
		os.MkdirAll(outDir, os.ModePerm)
		outPath := filepath.Join(outDir, outFile)

		if !fd.FileExists(outPath) {
			nmData, err := os.ReadFile(nmFile)
			if err != nil {
				panic(err)
			}
			nodeProc(nmData, outDir, outFile, treeFile, uri4id)
		}

		// 	// 	/////

		// 	data, err := os.ReadFile(outpath)
		// 	if err != nil {
		// 		log.Fatalln(err)
		// 	}

		// 	mIdBlock, _ := getIdBlock(string(data))

		// 	inpath4exp := outpath
		// 	outexp := childrenRepl(inpath4exp, mIdBlock)
		// 	// os.WriteFile("./out/asnexp.json", []byte(outexp), os.ModePerm)

		// 	rootWholeBlock := getRootWholeObject(outexp)
		// 	os.WriteFile("./out/asn-node-one.json", []byte(rootWholeBlock), os.ModePerm)

	}

	// return

	//////////////////////////////////////////////////////////////////////

	{
		mInputLa := map[string]string{
			"la-Languages.json":                      "Languages",
			"la-English.json":                        "English",
			"la-Humanities and Social Sciences.json": "Humanities and Social Sciences",
			"la-Health and Physical Education.json":  "Health and Physical Education",
			"la-Mathematics.json":                    "Mathematics",
			"la-Science.json":                        "Science",
			"la-Technologies.json":                   "Technologies",
			"la-The Arts.json":                       "The Arts",
			"ccp-Cross-curriculum Priorities.json":   "CCP",
			"gc-Critical and Creative Thinking.json": "GC-CCT",
			"gc-Digital Literacy.json":               "GC-DL",
			"gc-Ethical Understanding.json":          "GC-EU",
			"gc-Intercultural Understanding.json":    "GC-IU",
			"gc-Literacy.json":                       "GC-L",
			"gc-Numeracy.json":                       "GC-N",
			"gc-Personal and Social capability.json": "GC-PSC",
		}

		dataTree, err := os.ReadFile(treeFile)
		if err != nil {
			log.Fatalln(err)
		}
		mCodeParent := tool.GetCodeParentMap(dataTree)

		dataNode, err := os.ReadFile(nmFile)
		if err != nil {
			log.Fatalln(err)
		}
		// mUidTitle := scanNodeIdTitle(dataNode) // title should be node title

		mNodeData, err := jt.Flatten(dataNode)
		if err != nil {
			log.Fatalln(err)
		}

		wg := sync.WaitGroup{}
		wg.Add(len(mInputLa))

		for file, la := range mInputLa {

			go func(file, la string) {

				if file != "la-Languages.json" {
					wg.Done()
					return
				}

				fmt.Printf("----- %s ----- %s\n", file, la)

				var (
					prevDocTypePath = ""
					retEL           = `` // used by 'Level' & its descendants
					retPL           = `` // used by 'Level' & its descendants
					progLvlABC      = "" // indicate Level 1a, 1b or 1c
				)

				data, err := os.ReadFile(filepath.Join(`../data-out/restructure`, file))
				if err != nil {
					log.Fatalln(err)
				}
				js := removeEsc(string(data))

				///
				switch {
				case strings.Contains(js, `"Level 1c"`):
					progLvlABC = "1c"
				case strings.Contains(js, `"Level 1b"`):
					progLvlABC = "1b"
				case strings.Contains(js, `"Level 1a"`):
					progLvlABC = "1a"
				}
				///

				paths, _ := jt.GetLeavesPathOrderly(js)

				js = treeProc3(
					[]byte(js),
					la,
					mCodeParent,
					mNodeData,
					paths,
					&prevDocTypePath,
					&retEL,
					&retPL,
					progLvlABC,
				)

				js = restoreEsc(js)

				os.WriteFile(filepath.Join(`../data-out/asn-json`, file), []byte(js), os.ModePerm)

				wg.Done()

			}(file, la)
		}

		wg.Wait()
	}
}
