package tree

import (
	"fmt"
	"os"
	"testing"

	lk "github.com/digisan/logkit"
	"github.com/nsip/mrac-2023/meta"
)

const (
	out = "../data-out"

	metaFile = "../data/Sofia-API-Meta-Data-13062023.json"
	nodeFile = "../data/Sofia-API-Node-Data-13062023.json"
	treeFile = "../data/Sofia-API-Tree-Data-13062023.json"
)

func TestPartition(t *testing.T) {

	os.MkdirAll(out, os.ModePerm)

	dataMeta, err := os.ReadFile(metaFile)
	lk.FailOnErr("%v", err)
	jsMeta := string(dataMeta)

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

	dataTree, err := os.ReadFile(treeFile) // tree.pretty.json
	lk.FailOnErr("%v", err)
	js := string(dataTree)
	Partition(js, out, mMetaKeyName)
}
