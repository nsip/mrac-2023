package main

var (
	context = `"@context": {
		  "asn": "http://purl.org/ASN/schema/core/",
		  "dc": "http://purl.org/dc/terms/",
		  "gem": "http://purl.org/gem/qualifiers/",
		  "esa": "https://www.esa.edu.au/",
		  "skos": "http://www.w3.org/2004/02/skos/core#",
		  "xsd": "http://www.w3.org/2001/XMLSchema#",
		  "@language": "en-au"
		}`

	mPrefNamespace = map[string]string{
		"asn": "http://purl.org/ASN/schema/core/",
		"deo": "http://purl.org/spar/deo",
		"esa": "https://www.esa.edu.au/",
		"dc":  "http://purl.org/dc/terms/",
		"gem": "http://purl.org/gem/qualifiers/",
	}

	mPrefRepl = map[string]string{
		"dc_":      "dc:",
		"dcterms_": "dc:",
		"asn_":     "asn:",
	}

	mFieldRepl = map[string]string{
		"text":        "dc:text",
		"description": "dc:description",
		"children":    "gem:hasChild",
		"id":          "@id",
	}

	mFieldRemove = map[string]struct{}{
		"cls":  {},
		"leaf": {},
	}
)
