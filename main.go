package main

import (
	"fmt"
	"os"

	lk "github.com/digisan/logkit"
	"github.com/nsip/mrac-2023/meta"
	"github.com/nsip/mrac-2023/node"
	"github.com/nsip/mrac-2023/tree"
)

const (
	uri4id = "http://vocabulary.curriculum.edu.au/" // "http://rdf.curriculum.edu.au/202110/"
	out    = "data-out"

	metaFile = "./data/Sofia-API-Meta-Data-23022023.json"
	nodeFile = "./data/Sofia-API-Nodes-Data-23022023.json"
	treeFile = "./data/Sofia-API-Tree-Data-23022023.json"
)

func main() {

	os.MkdirAll(fmt.Sprintf("./%s/", out), os.ModePerm)

	bytesMeta, err := os.ReadFile(metaFile)
	lk.FailOnErr("%v", err)
	jsMeta := string(bytesMeta)

	mMetaKeyName, err := meta.Parse(jsMeta, "name")
	lk.FailOnErr("%v", err)
	mMetaKeyPlural, err := meta.Parse(jsMeta, "plural")
	lk.FailOnErr("%v", err)

	for k, v := range mMetaKeyName {
		fmt.Printf("%v: %v\n", k, v)
	}
	fmt.Println("---------------------------------------------")
	for k, v := range mMetaKeyPlural {
		fmt.Printf("%v: %v\n", k, v)
	}

	// fmt.Printf("%d\n", len(mMeta))

	///////////////////////////////////////////////////////////////////

	bytesNode, err := os.ReadFile(nodeFile)
	lk.FailOnErr("%v", err)
	node.Process(bytesNode, uri4id, mMetaKeyName, out)
	// node.GenCodeIdUrlTxt(bytesNode, out) // *** if 'code-url.txt' & 'id-url.txt' exist, DO NOT run this ***

	///////////////////////////////////////////////////////////////////

	bytesTree, err := os.ReadFile(treeFile) // tree.pretty.json
	lk.FailOnErr("%v", err)
	js := string(bytesTree)
	tree.Partition(js, out, mMetaKeyName)

}
