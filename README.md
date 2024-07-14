# mrac-2023

'scot.jsonld' from <http://vocabulary.curriculum.edu.au/scot/export/scot.jsonld>

0. 2023-Nov

# convert Retool export mapping-*.csv, which is the database of mappings of ScOT to Curriculum Content Descriptions
# into SCOT.txt, mapping of content descriptions to SCoT URL

# Input 2023-Nov/mapping-*.csv
# Output data/mapping-*.csv

1. /node2/node2_test.go (if there are existing node-meta.json & code-url.txt &
   id-url.txt, can ignore this step)

# input files specified in node2_test.go:
# ../data/Sofia-API-Meta-Data.json
# ../data/Sofia-API-Node-Data.json
# ../data/Sofia-API-Tree-Data.json

# code-url.txt is commented out
# cd node2; go test; files output to data/

# data/code-url.txt : Statement Notation => URL prefix, e.g.
# ASLANKORF10_5605	http://vocabulary.curriculum.edu.au/MRAC/2024/04/LA/LAN/

# data/id-url.txt : GUID => URL prefix, e.g.
# 6ecea5f2-6b36-41a5-ae53-274ae0de6268	http://vocabulary.curriculum.edu.au/MRAC/2024/04/LA/LAN/

# data/node-meta.json : map GUID to node struct, based on Sofia-API-Node-Data-*.json,
# interleave into it data from Sofia-API-Node-Meta-*.json
# seems to be just http://vocabulary.curriculum.edu.au prefix on GUIDs,
# not even with the child subdirectories



2. /asn-json/tool/ac_scot_test.go [in TestGetAsnConceptTerm] (if there is
   existing id-preflabel.txt, can ignore this step)

# input files:
# data/SCOT.txt
# data/pp_project_schoolsonlinethesaurus.jsonld

# data/id-preflabel.txt: Content Description Code => ScOT URIs and labels
# code to generate had been disabled (!)
 
# cd asn-json/tool; go test; files output to data/

3. /tree/tree_test.go (then copy gc-_.json & ccp-_.json, from /data-out/ into
   /data-out/restructure/ before next step)

# cd tree; go test
# outputs data-out/* and data-out/restructure/*
# data-out/ld* appears to be updated to data-out/restructure/ld*
# cp data-out/ccp* data-out/restructure
# cp data-out/gc* data-out/restructure

4. /asn-json (rewriting... /asn-json-v2)

# cd asn-json; go test; generates data-out/asn-json
# these are the files with the keys we expect. But ignores Languages. Seems to be obsolete code.

# cd asn-json-v2; go test; generates ./asn-json-v2
# these are the files with the keys we expect. 

# was dumping files locally, moving to ../data-out/asn-json

5. /task-force-array-2023-06-21

# forces asn_hasLevel to be array. Should be redundant, already done in asn-json-v2
# ... but isn't: modifies data-out/asn-json/* in situ.

<!-- 5. /task-remove-fields-2023-06-22 -->

# remove "asn_conceptTerm" with value ["[]", "SCIENCE_TEACHER_BACKGROUND_INFORMATION"]

7. /asn-json-ld

# reads data-out/asn-json, generates data-out/asn-json-ld
# Both directories still contain ScOT

8. /task-remove-fields-2023-06-23

# processes JSON-LD: strips dc:text, adds @language @value to dc:description
# updates asn-json-ld in situ

9. /task-split-ccp (outdir for ccp here)

# breaks up data-out/asn-json-ld/ccp-Cross-curriculum Priorities.json
# into task-split-ccp-2022-11-12/asn-json-ccp/*.json, task-split-ccp-2022-11-12/asn-json-ld-ccp/*.json

10. /task-structure-improvement (indir only for ccp here)

# tidy of ccp only (?): date-time, skos:prefLabel  on curriculum statemnents, 
# language literal position to @language @value @asn:listID
# make hasLevel be array
# input task-split-ccp-2022-11-12/asn-json-ld-ccp/*
# output task-structure-improvement-2023-03-23/asn-json-ld-ccp/*

11. `mkdir ./release && cp -rf ./data-out/asn-json ./release`
    `cp -rf ./data-out/asn-json-ld ./release`
    `cp -rf ./split-ccp/asn-json-ccp ./release` *
    `cp -rf ./structure-improvement/asn-json-ld-ccp ./release` *

12. /task-add-dctitle-2023-10-10 (only for json-ld) (/release must exist)

# processes release/asn-json-ld/* in situ
# copies  dc:title to dc:description (not for asn-json-ld-ccp !!)
