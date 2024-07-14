package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	main()
}

/*
func TestCompareWithDataOut(t *testing.T) {

	dir1 := "../data-out/asn-json"
	dir2 := "./"

	des, err := os.ReadDir(dir1)
	if err != nil {
		log.Fatalln(err)
	}

	for _, de := range des {

		if de.IsDir() {
			continue
		}

		name := de.Name()
		if name == "asn-node.json" {
			continue
		}

		file1 := filepath.Join(dir1, name)
		file2 := filepath.Join(dir2, name)

		data2, err := os.ReadFile(file2)
		if err != nil {
			log.Fatalln(err)
		}

		data1, err := os.ReadFile(file1)
		if err != nil {
			log.Fatalln(err)
		}

		if string(data1) != string(data2) {
			log.Fatalf("Not Equal: %v - %v", file1, file2)
		}
	}
}
*/
