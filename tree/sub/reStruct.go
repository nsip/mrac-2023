package sub

import (
	"errors"
	"fmt"
	"strings"

	. "github.com/digisan/go-generics/v2"
	lk "github.com/digisan/logkit"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func init() {
	lk.WarnDetail(false)
	lk.Log2F(true, "reStruct.log")
}

func ReStruct(js string) string {

	fSf := fmt.Sprintf

	var (
		mLA2PathLvl = map[string]string{}
		mLA2PathLA  = map[string]string{}
		mAS         = make(map[string]string)
		asCodes     = []string{}
		laCodes     = []string{}
		laCodesRem  = []string{}
	)

	// 0:
	for I := 0; I < 2; I++ {
		const P0 = "children.%d."
		path := fSf(P0+"code", I)
		code0 := gjson.Get(js, path).String()
		if code0 == "" {
			break
		}
		fmt.Println(code0)

		// 1:
		for i := 0; i < 100; i++ {
			var P1 = strings.Repeat(P0, 2)
			path := fSf(P1+"code", I, i)
			code1 := gjson.Get(js, path).String()
			if code1 == "" {
				break
			}
			// fmt.Printf("\t%s\n", code1)

			path = fSf(P1+"doc.typeName", I, i)
			if typeName := gjson.Get(js, path).String(); typeName == "Level" {
				fmt.Printf("\t%s - ok\n", code1)

				if code0 == "AS" {
					asCodes = append(asCodes, code1)
					mAS[code1] = gjson.Get(js, fSf(P1+"children.0", I, i)).String()
				}
				if code0 == "LA" {
					laCodes = append(laCodes, code1)
					mLA2PathLvl[code1] = fSf(P1+"asn_hasLevel", I, i)
				}

			} else if typeName == "Learning Area" {

				if code0 == "LA" {
					laCodesRem = append(laCodesRem, code1)
					mLA2PathLA[code1] = fSf(P1+"asn_hasLevel", I, i)
				}
			}

			// 2:
			for j := 0; j < 100; j++ {
				var P2 = strings.Repeat(P0, 3)
				path := fSf(P2+"code", I, i, j)
				code2 := gjson.Get(js, path).String()
				if code2 == "" {
					break
				}
				// fmt.Printf("\t\t%s\n", code2)

				path = fSf(P2+"doc.typeName", I, i, j)
				if typeName := gjson.Get(js, path).String(); typeName == "Level" {
					fmt.Printf("\t\t%s - ok\n", code2)

					if code0 == "AS" {
						asCodes = append(asCodes, code2)
						mAS[code2] = gjson.Get(js, fSf(P2+"children.0", I, i, j)).String()
					}
					if code0 == "LA" {
						laCodes = append(laCodes, code2)
						mLA2PathLvl[code2] = fSf(P2+"asn_hasLevel", I, i, j)
					}

				} else if typeName == "Learning Area" {

					if code0 == "LA" {
						laCodesRem = append(laCodesRem, code2)
						mLA2PathLA[code2] = fSf(P2+"asn_hasLevel", I, i, j)
					}
				}

				// 3:
				for k := 0; k < 100; k++ {
					var P3 = strings.Repeat(P0, 4)
					path := fSf(P3+"code", I, i, j, k)
					code3 := gjson.Get(js, path).String()
					if code3 == "" {
						break
					}
					// fmt.Printf("\t\t\t%s", code3)

					path = fSf(P3+"doc.typeName", I, i, j, k)
					if typeName := gjson.Get(js, path).String(); typeName == "Level" {
						fmt.Printf("\t\t\t%s - ok\n", code3)

						if code0 == "AS" {
							asCodes = append(asCodes, code3)
							mAS[code3] = gjson.Get(js, fSf(P3+"children.0", I, i, j, k)).String()
						}
						if code0 == "LA" {
							laCodes = append(laCodes, code3)
							mLA2PathLvl[code3] = fSf(P3+"asn_hasLevel", I, i, j, k)
						}

					} else if typeName == "Learning Area" {

						if code0 == "LA" {
							laCodesRem = append(laCodesRem, code3)
							mLA2PathLA[code3] = fSf(P3+"asn_hasLevel", I, i, j, k)
						}
					}

					// 4:
					for l := 0; l < 100; l++ {
						var P4 = strings.Repeat(P0, 5)
						path := fSf(P4+"code", I, i, j, k, l)
						code4 := gjson.Get(js, path).String()
						if code4 == "" {
							break
						}
						// fmt.Printf("\t\t\t\t%s", code4)

						path = fSf(P4+"doc.typeName", I, i, j, k, l)
						if typeName := gjson.Get(js, path).String(); typeName == "Level" {
							fmt.Printf("\t\t\t\t%s - ok\n", code4)

							if code0 == "AS" {
								asCodes = append(asCodes, code4)
								mAS[code4] = gjson.Get(js, fSf(P4+"children.0", I, i, j, k, l)).String()
							}
							if code0 == "LA" {
								laCodes = append(laCodes, code4)
								mLA2PathLvl[code4] = fSf(P4+"asn_hasLevel", I, i, j, k, l)
							}

						} else if typeName == "Learning Area" {

							if code0 == "LA" {
								laCodesRem = append(laCodesRem, code4)
								mLA2PathLA[code4] = fSf(P4+"asn_hasLevel", I, i, j, k, l)
							}
						}

						// 5:
						for m := 0; m < 100; m++ {
							var P5 = strings.Repeat(P0, 6)
							path := fSf(P5+"code", I, i, j, k, l, m)
							code5 := gjson.Get(js, path).String()
							if code5 == "" {
								break
							}
							// fmt.Printf("\t\t\t\t\t%s", code5)

							path = fSf(P5+"doc.typeName", I, i, j, k, l, m)
							if typeName := gjson.Get(js, path).String(); typeName == "Level" {
								fmt.Printf("\t\t\t\t\t%s - ok\n", code5)

								if code0 == "AS" {
									asCodes = append(asCodes, code5)
									mAS[code5] = gjson.Get(js, fSf(P5+"children.0", I, i, j, k, l, m)).String()
								}
								if code0 == "LA" {
									laCodes = append(laCodes, code5)
									mLA2PathLvl[code5] = fSf(P5+"asn_hasLevel", I, i, j, k, l, m)
								}

							} else if typeName == "Learning Area" {

								if code0 == "LA" {
									laCodesRem = append(laCodesRem, code5)
									mLA2PathLA[code5] = fSf(P5+"asn_hasLevel", I, i, j, k, l, m)
								}
							}
						}
					}
				}
			}
		}
	}

	// checking unused asCode
	asRemCont := []string{}
	for _, as := range asCodes {
		la := strings.TrimPrefix(as, "AS")
		la = strings.TrimSuffix(la, "L")
		if NotIn(la, laCodes...) {
			lk.Warn("AS has [%s], BUT LA has no [%s]", as, la)
			asRemCont = append(asRemCont, mAS[as])
		}
	}

	// checking which laCode has no asCode
	for _, la := range laCodes {
		as := "AS" + la + "L"
		if NotIn(as, asCodes...) {
			lk.Log("LA has [%s], BUT AS has no [%s], ignore [%s]", la, as, as)
		}
	}

	for laCode, path := range mLA2PathLvl {
		// path += fmt.Sprintf(".%d", len(gjson.Get(js, path).Array())) // modify path, append to the last child
		if content, ok := mAS["AS"+laCode+"L"]; ok {
			if content != "" {
				fmt.Println("PATH: --->", path)
				js, _ = sjson.SetRaw(js, path, content)
			}
		} else {
			// keey `"asn_hasLevel": ""` ?
			// js, _ = sjson.Set(js, path, "")
		}
	}

	lk.FailOnErrWhen(len(mLA2PathLA) > 1, "%v", errors.New("Learning Area's count MUST be One"))
	for _, path := range mLA2PathLA {
		js, _ = sjson.SetRaw(js, path, "["+strings.Join(asRemCont, ",")+"]")
	}

	// remove AS part
	js, _ = sjson.Delete(js, "children.0")
	return js
}
