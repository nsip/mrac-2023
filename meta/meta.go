package meta

import (
	"fmt"

	. "github.com/digisan/go-generics/v2"
	"github.com/tidwall/gjson"
)

func Parse(js, field string) (map[string]string, error) {
	if NotIn(field, "name", "plural") {
		return nil, fmt.Errorf("[%s] is NOT a field in Sofia-API-Meta-Data.json or still NOT supported", field)
	}
	mMeta := make(map[string]string)
	r := gjson.Get(js, "fields")
	if r.IsArray() {
		for _, ra := range r.Array() {
			if ra.IsObject() {
				mMeta[ra.Get("key").String()] = ra.Get(field).String()
			}
		}
	}
	return mMeta, nil
}
