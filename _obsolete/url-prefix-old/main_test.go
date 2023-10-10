package main

import (
	"fmt"
	"testing"
)

func TestFetchTime(t *testing.T) {
	fmt.Println(FetchTime("../data/Sofia-API-Tree-Data-18072023.json"))
}

func TestMain(t *testing.T) {
	main()
}
