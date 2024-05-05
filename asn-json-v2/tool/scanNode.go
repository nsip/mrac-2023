package tool

import (
	"regexp"
	"strings"
)

func ScanNode(data []byte, f func(i int, id, block string) bool) {

	js := string(data)
	r := regexp.MustCompile(`"[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}":`)
	pGrp := r.FindAllStringIndex(js, -1)
	// fmt.Println(len(pGrp), js[pGrp[0][0]:pGrp[0][1]])

	for i := 0; i < len(pGrp); i++ {

		p, pn := pGrp[i], []int{}
		if i < len(pGrp)-1 {
			p, pn = pGrp[i], pGrp[i+1]
		}

		ids, ide := p[0]+1, p[1]-2
		id := js[ids:ide]
		// fmt.Println(id)

		blks, blke := 0, 0
		if i < len(pGrp)-1 {
			blks, blke = p[1], pn[0]-1
		} else {
			blks, blke = p[1], len(js)-1
		}

		block := js[blks:blke]
		block = strings.TrimSuffix(block, " ")
		block = strings.TrimSuffix(block, "\n")
		block = strings.TrimSuffix(block, ",")

		if !f(i, id, block) {
			break
		}
	}
}
