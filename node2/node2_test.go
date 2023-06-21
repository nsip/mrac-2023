package node2

import (
	"fmt"
	"os"
	"testing"

	lk "github.com/digisan/logkit"
	"github.com/nsip/mrac-2023/meta"
)

func TestNode2(t *testing.T) {

	dataNode, err := os.ReadFile("../data/Sofia-API-Node-Data-13062023.json")
	lk.FailOnErr("%v", err)

	mIdBlock := GenNodeIdBlockMap(dataNode)
	mCodeBlock := GenNodeCodeBlockMap(dataNode)
	mIdChildParent, mCodeChildParent := GenChildParentMap(dataNode, mIdBlock)

	fmt.Printf("Total: %d - %d - %d - %d\n", len(mIdBlock), len(mCodeBlock), len(mIdChildParent), len(mCodeChildParent))

	ancestryID := RetrieveAncestry("3ae876a8-10b9-44c7-9b2e-13b2ba08e217", mIdChildParent)
	fmt.Println(ancestryID)

	ancestryCode := RetrieveAncestryAsCodeById("3ae876a8-10b9-44c7-9b2e-13b2ba08e217", mIdChildParent, mIdBlock)
	fmt.Println(ancestryCode)

	ancestryCode = RetrieveAncestry("AC9TDI4P07_E5", mCodeChildParent)
	fmt.Println(ancestryCode)

	// fmt.Println(IsAncestorCode("AC9LIN10C03_E3", "Indicator", mCodeChildParent))

	// fmt.Println(GetIdByCode("AC9LIN10C03_E3", mCodeBlock))
	// fmt.Println(GetCodeById("ffdaf9d5-514b-4f0d-873c-130ffbde52f4", mIdBlock))
}

const (
	uri4id   = "http://vocabulary.curriculum.edu.au/" // "http://rdf.curriculum.edu.au/202110/"
	metaFile = "../data/Sofia-API-Meta-Data-13062023.json"
	nodeFile = "../data/Sofia-API-Node-Data-13062023.json"
	treeFile = "../data/Sofia-API-Tree-Data-13062023.json"
)

// *** //
func TestUpdateNodeWithMeta(t *testing.T) {

	dataNode, err := os.ReadFile(nodeFile)
	lk.FailOnErr("%v", err)

	dataMeta, err := os.ReadFile(metaFile)
	lk.FailOnErr("%v", err)
	jsMeta := string(dataMeta)
	mMetaKeyName, err := meta.Parse(jsMeta, "name")
	lk.FailOnErr("%v", err)
	// mMetaKeyPlural, err := meta.Parse(jsMeta, "plural")
	// lk.FailOnErr("%v", err)

	UpdateNodeWithMeta(dataNode, uri4id, mMetaKeyName, "../data/node-meta.json")
}

// *** //
func TestMakeUrlText(t *testing.T) {

	dataNode, err := os.ReadFile(nodeFile)
	lk.FailOnErr("%v", err)

	mIdBlock := GenNodeIdBlockMap(dataNode)
	mCodeBlock := GenNodeCodeBlockMap(dataNode)
	mIDChildParent, mCodeChildParent := GenChildParentMap(dataNode, mIdBlock)

	MakeIdUrlText(mIdBlock, mCodeBlock, mIDChildParent, mCodeChildParent, "../data/id-url.txt", "../data/code-url.txt")
}
