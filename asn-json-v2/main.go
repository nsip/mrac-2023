package main

import (
	"fmt"
	"strings"

	. "github.com/digisan/go-generics"
	fd "github.com/digisan/gotk/file-dir"
	jt "github.com/digisan/json-tool/scan"
)

func loadUrl(fPath string) map[string]string {
	m := make(map[string]string)
	fd.FileLineScan(fPath, func(line string) (bool, string) {
		ss := strings.Split(line, "\t")
		m[ss[0]] = ss[1]
		return true, ""
	}, "")
	return m
}

var (
	mIdUrl = loadUrl("../data/id-url.txt")
)

func main() {

	fPath := "../data-out/restructure/la-Languages.json"
	fOut := "./la-Languages.json"

	opt := jt.OptLineProc{
		Fn_KV:          nil,
		Fn_KV_Str:      proc_kv_str,
		Fn_KV_Obj_Open: nil,
		Fn_KV_Arr_Open: nil,
		Fn_Obj:         nil,
		Fn_Arr:         nil,
		Fn_Elem:        nil,
		Fn_Elem_Str:    nil,
	}

	fmt.Println(jt.ScanJsonLine(fPath, fOut, opt))
}

///////////////////////////////////////////////////////////////

func proc_kv(I int, path, k string, v any) (bool, string) {
	return true, fmt.Sprintf(`"%v": %v`, k, v)
}

func proc_kv_str(I int, path, k, v string) (bool, string) {

	if k == "code" {
		dcterms := ""
		if In(v, "root", "LA") {
			dcterms = `
			  "dcterms_rights": { "language": "en-au", "literal": "©Copyright Australian Curriculum, Assessment and Reporting Authority" },
			  "dcterms_rightsHolder": { "language": "en-au", "literal": "Australian Curriculum, Assessment and Reporting Authority" },`
		}
		return true, dcterms + fmt.Sprintf(`"asn_statementNotation": { "language": "en-au", "literal": "%v" },
		"asn_authorityStatus": { "uri": "http://purl.org/ASN/scheme/ASNAuthorityStatus/Original" },
		"asn_indexingStatus": { "uri": "http://purl.org/ASN/scheme/ASNIndexingStatus/No" }`, v)
	}

	if k == "uuid" {
		return true, fmt.Sprintf(`"id": "%s%s"`, mIdUrl[v], v) // mIdUrl[value] already append with '/'
	}

	if k == "type" {
		return false, ""
	}

	if k == "created_at" {
		return true, fmt.Sprintf(`"dcterms_modified": { "literal": "%v" }`, v)
	}

	if k == "title" {
		// with 'text' sibling

		// without 'text' sibling
		return true, fmt.Sprintf(`"dcterms_title": { "language": "en-au", "literal": "%s" }, "text": "%s"`, v, v)
	}

	if k == "text" {
		return true, fmt.Sprintf(`"dcterms_description": { "language": "en-au", "literal": "%s" }, "text": "%s"`, v, v)
	}

	if k == "position" {
		return true, fmt.Sprintf(`"asn_listID": "%v"`, v)
	}

	return true, fmt.Sprintf(`"%v": "%v"`, k, v)
}

func proc_kv_oo(k string, v any) string {
	return fmt.Sprintf(`"%v": %v`, k, v)
}

func proc_kv_ao(k string, v any) string {
	return fmt.Sprintf(`"%v": %v`, k, v)
}

func proc_obj(v any) string {
	return v.(string)
}

func proc_arr(v any) string {
	return v.(string)
}

func proc_elem(v any) string {
	return v.(string)
}

func proc_elem_s(v any) string {
	return fmt.Sprintf(`"%v"`, v)
}
