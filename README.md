# mrac-2023

'scot.jsonld' from <http://vocabulary.curriculum.edu.au/scot/export/scot.jsonld>

1. /node2/node2_test.go (if there are existing node-meta.json & code-url.txt &
   id-url.txt, can ignore this step)

2. /asn-json/tool/ac_scot_test.go [in TestGetAsnConceptTerm] (if there is
   existing id-preflabel.txt, can ignore this step)

3. /tree/tree_test.go (then copy gc-_.json & ccp-_.json, from data-out/ into
   data-out/restructure/ before next step)

4. /asn-json

5. /task-force-array-2023-06-21

<!-- 5. /task-remove-fields-2023-06-22 -->

7. /asn-json-ld

8. /task-remove-fields-2023-06-23

9. /task-split-ccp (outdir for ccp here)

10. /task-structure-improvement

11. `mkdir release && cp -rf /data-out/asn-json /release`
    `cp -rf /data-out/asn-json-ld /release`
    `cp -rf /split-ccp/asn-json-ccp /release`
    `cp -rf /structure-improvement/asn-json-ld-ccp /release`

12. /task-add-dctitle-2023-10-10 (only to deal with json-ld) (/release must
    exist)
