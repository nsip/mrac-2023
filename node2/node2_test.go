package node2

import (
	"fmt"
	"os"
	"testing"

	lk "github.com/digisan/logkit"
)

func TestNode2(t *testing.T) {

	nodeData, err := os.ReadFile("../data/Sofia-API-Node-Data-06062023.json")
	lk.FailOnErr("%v", err)

	mIdBlock := GenNodeIdBlockMap(nodeData)
	mCodeBlock := GenNodeCodeBlockMap(nodeData)
	mIDChildParent, mCodeChildParent := GenChildParentMap(nodeData, mIdBlock)

	fmt.Printf("Total: %d - %d - %d - %d\n", len(mIdBlock), len(mCodeBlock), len(mIDChildParent), len(mCodeChildParent))

	ancestryID := RetrieveAncestryID("ffdaf9d5-514b-4f0d-873c-130ffbde52f4", mIDChildParent)
	fmt.Println(ancestryID)

	ancestryCode := RetrieveAncestryID("AC9LIN10C03_E3", mCodeChildParent)
	fmt.Println(ancestryCode)

	fmt.Println(IsAncestorCode("AC9LIN10C03_E3", "Indicator", mCodeChildParent))

	fmt.Println(GetIdByCode("AC9LIN10C03_E3", mCodeBlock))
	fmt.Println(GetCodeById("ffdaf9d5-514b-4f0d-873c-130ffbde52f4", mIdBlock))
}

func TestMakeUrlText(t *testing.T) {

	nodeData, err := os.ReadFile("../data/Sofia-API-Node-Data-06062023.json")
	lk.FailOnErr("%v", err)

	mIdBlock := GenNodeIdBlockMap(nodeData)
	mCodeBlock := GenNodeCodeBlockMap(nodeData)
	mIDChildParent, mCodeChildParent := GenChildParentMap(nodeData, mIdBlock)

	MakeIdUrlText(mIdBlock, mCodeBlock, mIDChildParent, mCodeChildParent, "../data/id-url.txt", "../data/code-url.txt")
}
