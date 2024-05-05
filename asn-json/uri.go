package main

import (
	"fmt"

	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
)

var (
	// uri4id = "http://vocabulary.curriculum.edu.au/" // "http://rdf.curriculum.edu.au/202110"

	mLaUri = map[string]string{
		"English":                        `http://vocabulary.curriculum.edu.au/framework/E`,
		"The Arts":                       `http://vocabulary.curriculum.edu.au/framework/A`,
		"Health and Physical Education":  `http://vocabulary.curriculum.edu.au/framework/P`,
		"Humanities and Social Sciences": `http://vocabulary.curriculum.edu.au/framework/U`,
		"Languages":                      `http://vocabulary.curriculum.edu.au/framework/L`,
		"Mathematics":                    `http://vocabulary.curriculum.edu.au/framework/M`,
		"Science":                        `http://vocabulary.curriculum.edu.au/framework/S`,
		"Technologies":                   `http://vocabulary.curriculum.edu.au/framework/T`,
		"Work Studies":                   `http://vocabulary.curriculum.edu.au/framework/W`,
	}

	mYrlvlUri = map[string]string{
		"Early years":     `http://vocabulary.curriculum.edu.au/schoolLevel/-`,
		"Foundation Year": `http://vocabulary.curriculum.edu.au/schoolLevel/0`,
		"Year 1":          `http://vocabulary.curriculum.edu.au/schoolLevel/1`,
		"Year 2":          `http://vocabulary.curriculum.edu.au/schoolLevel/2`,
		"Year 3":          `http://vocabulary.curriculum.edu.au/schoolLevel/3`,
		"Year 4":          `http://vocabulary.curriculum.edu.au/schoolLevel/4`,
		"Year 5":          `http://vocabulary.curriculum.edu.au/schoolLevel/5`,
		"Year 6":          `http://vocabulary.curriculum.edu.au/schoolLevel/6`,
		"Year 7":          `http://vocabulary.curriculum.edu.au/schoolLevel/7`,
		"Year 8":          `http://vocabulary.curriculum.edu.au/schoolLevel/8`,
		"Year 9":          `http://vocabulary.curriculum.edu.au/schoolLevel/9`,
		"Year 10":         `http://vocabulary.curriculum.edu.au/schoolLevel/10`,
		"Year 11":         `http://vocabulary.curriculum.edu.au/schoolLevel/11`,
		"Year 12":         `http://vocabulary.curriculum.edu.au/schoolLevel/12`,
	}

	mProglvlUri = map[string]string{
		"1":  `http://vocabulary.curriculum.edu.au/progressionLevel/1`,
		"2":  `http://vocabulary.curriculum.edu.au/progressionLevel/2`,
		"3":  `http://vocabulary.curriculum.edu.au/progressionLevel/3`,
		"4":  `http://vocabulary.curriculum.edu.au/progressionLevel/4`,
		"5":  `http://vocabulary.curriculum.edu.au/progressionLevel/5`,
		"6":  `http://vocabulary.curriculum.edu.au/progressionLevel/6`,
		"7":  `http://vocabulary.curriculum.edu.au/progressionLevel/7`,
		"8":  `http://vocabulary.curriculum.edu.au/progressionLevel/8`,
		"9":  `http://vocabulary.curriculum.edu.au/progressionLevel/9`,
		"10": `http://vocabulary.curriculum.edu.au/progressionLevel/10`,
		"11": `http://vocabulary.curriculum.edu.au/progressionLevel/11`,
		"12": `http://vocabulary.curriculum.edu.au/progressionLevel/12`,
		"13": `http://vocabulary.curriculum.edu.au/progressionLevel/13`,
		"14": `http://vocabulary.curriculum.edu.au/progressionLevel/14`,
	}

	mProglvlABCUri = map[string]string{
		"1a": `http://vocabulary.curriculum.edu.au/progressionLevel/-3`,
		"1b": `http://vocabulary.curriculum.edu.au/progressionLevel/-2`,
		"1c": `http://vocabulary.curriculum.edu.au/progressionLevel/-1`,
	}

	mProglvlABUri = map[string]string{
		"1a": `http://vocabulary.curriculum.edu.au/progressionLevel/-2`,
		"1b": `http://vocabulary.curriculum.edu.au/progressionLevel/-1`,
	}

	mProglvlAUri = map[string]string{
		"1a": `http://vocabulary.curriculum.edu.au/progressionLevel/-1`,
	}
)

func getProLevel(mData map[string]any, path string) string {
	N := 0
AGAIN:
	lk.FailOnErrWhen(N > 500, "%v", fmt.Errorf("getProLevel DeadLoop?"))

	sp := jt.NewSibling(path, "doc.typeName")
	if mData[sp] == "Level" {
		lvlstr := mData[jt.NewSibling(path, "title")].(string)
		tail := ""
		fmt.Sscanf(lvlstr, "Level %s", &tail)
		return tail
	} else {
		path = jt.ParentPath(path)
		if path == "" {
			return ""
		}
		N++
		goto AGAIN
	}
}

func getYears(mData map[string]any, path string) []string {
	N := 0
AGAIN:
	lk.FailOnErrWhen(N > 500, "%v", fmt.Errorf("getYears DeadLoop?"))

	sp := jt.NewSibling(path, "doc.typeName")
	if mData[sp] == "Level" {
		return yearsSplit(mData[jt.NewSibling(path, "title")].(string))
	} else {
		path = jt.ParentPath(path)
		if path == "" {
			return nil
		}
		N++
		goto AGAIN
	}
}
