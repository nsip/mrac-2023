package meta

import (
	"errors"
	"fmt"

	. "github.com/digisan/go-generics"
	lk "github.com/digisan/logkit"
	"github.com/tidwall/gjson"
)

func Parse(js, field string) (map[string]string, error) {
	if NotIn(field, "name", "plural") {
		return nil, fmt.Errorf("[%s] is NOT a field in Sofia-API-Meta-Data.json or NOT supported", field)
	}
	mMeta := make(map[string]string)
	if r := gjson.Get(js, "fields"); r.Type != gjson.Null && r.IsArray() {
		for _, ra := range r.Array() {
			if ra.IsObject() {
				if rKey := ra.Get("key"); rKey.Type != gjson.Null {
					if rField := ra.Get(field); rField.Type != gjson.Null {
						mMeta[rKey.Str] = rField.Str
					}
				} else {
					lk.FailOnErr("field 'key' is missing %v", errors.New(""))
				}
			}
		}
	}
	return mMeta, nil
}
