package main

import (
	"fmt"
	"os"

	fd "github.com/digisan/gotk/file-dir"
	"github.com/digisan/gotk/strs"
	lk "github.com/digisan/logkit"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func main() {

	// asn json
	{
		outDir := "./asn-json-ccp/"
		fd.MustCreateDir(outDir)

		dataCcp, err := os.ReadFile("../data-out/asn-json/ccp-Cross-curriculum Priorities.json")
		lk.FailOnErr("%v", err)

		jsCcp := string(dataCcp)
		parts := []string{
			gjson.Get(jsCcp, "children.0.children.0").Raw,
			gjson.Get(jsCcp, "children.0.children.1").Raw,
			gjson.Get(jsCcp, "children.0.children.2").Raw,
		}

		mPathLiteral := make(map[string]string)
		mLiteralPart := make(map[string]string)
		for i, part := range parts {

			literal := gjson.Get(part, "asn_statementNotation.literal").String()
			fmt.Println(literal)

			mLiteralPart[literal] = part

			switch i {
			case 0:
				mPathLiteral["children.0.children.0"] = literal
			case 1:
				mPathLiteral["children.0.children.1"] = literal
			case 2:
				mPathLiteral["children.0.children.2"] = literal
			default:
				lk.FailOnErr("%v", fmt.Errorf("only 3 children under CCP. If not, change here"))
			}
		}

		for path, literal := range mPathLiteral {

			part := mLiteralPart[literal]
			outPart, err := sjson.Delete(jsCcp, strs.TrimTailFromLast(path, "."))
			lk.FailOnErr("%v", err)

			pathOne := "children.0.children.0"
			out, err := sjson.SetRaw(outPart, pathOne, part)
			lk.FailOnErr("%v", err)

			os.WriteFile(outDir+literal+".json", []byte(out), os.ModePerm)
		}
	}

	// asn-json-ld
	{
		outDir := "./asn-json-ld-ccp/"
		fd.MustCreateDir(outDir)

		dataCcp, err := os.ReadFile("../data-out/asn-json-ld/ccp-Cross-curriculum Priorities.json")
		lk.FailOnErr("%v", err)

		jsCcp := string(dataCcp)
		parts := []string{
			gjson.Get(jsCcp, "gem:hasChild.0.gem:hasChild.0").Raw,
			gjson.Get(jsCcp, "gem:hasChild.0.gem:hasChild.1").Raw,
			gjson.Get(jsCcp, "gem:hasChild.0.gem:hasChild.2").Raw,
		}

		mPathLiteral := make(map[string]string)
		mLiteralPart := make(map[string]string)
		for i, part := range parts {

			literal := gjson.Get(part, "asn:statementNotation").String()
			fmt.Println(literal)

			mLiteralPart[literal] = part

			switch i {
			case 0:
				mPathLiteral["gem:hasChild.0.gem:hasChild.0"] = literal
			case 1:
				mPathLiteral["gem:hasChild.0.gem:hasChild.1"] = literal
			case 2:
				mPathLiteral["gem:hasChild.0.gem:hasChild.2"] = literal
			default:
				lk.FailOnErr("%v", fmt.Errorf("only 3 children under CCP. If not, change here"))
			}
		}

		for path, literal := range mPathLiteral {

			part := mLiteralPart[literal]
			outPart, err := sjson.Delete(jsCcp, strs.TrimTailFromLast(path, "."))
			lk.FailOnErr("%v", err)

			pathOne := "gem:hasChild.0.gem:hasChild.0"
			out, err := sjson.SetRaw(outPart, pathOne, part)
			lk.FailOnErr("%v", err)

			os.WriteFile(outDir+literal+".json", []byte(out), os.ModePerm)
		}
	}
}
