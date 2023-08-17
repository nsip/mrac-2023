package main

import (
	"fmt"
	"testing"
)

func TestYearsSplit(t *testing.T) {
	fmt.Println(yearsSplit("Foundation Year"), len(yearsSplit("Foundation Year")))
	fmt.Println(yearsSplit("Level 2 (Years 1- 2)"), len(yearsSplit("Level 2 (Years 1- 2)")))
	fmt.Println(yearsSplit("Level 2 (Years 1 - 2)"), len(yearsSplit("Level 2 (Years 1 - 2)")))
	fmt.Println(yearsSplit("Level 2 (Years 1-2)"), len(yearsSplit("Level 2 (Years 1-2)")))
	fmt.Println(yearsSplit("Year 9 and 10"), len(yearsSplit("Year 9 and 10")))
}
