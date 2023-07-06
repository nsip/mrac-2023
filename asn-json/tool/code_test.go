package tool

import (
	"fmt"
	"os"
	"testing"

	lk "github.com/digisan/logkit"
	"github.com/nsip/mrac-2023/node2"
)

func TestGetAllCode(t *testing.T) {

	nodeData, err := os.ReadFile("../../data/Sofia-API-Node-Data-04072023.json")
	lk.FailOnErr("%v", err)

	mIdBlock := node2.GenNodeIdBlockMap(nodeData)
	_, mCodeChildParent := node2.GenChildParentMap(nodeData, mIdBlock)

	fmt.Println(GetAncestorTitle("LSLiS5.6", "GC", mCodeChildParent))
}
