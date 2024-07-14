# inputs: 
# data/mapping.csv from Retool
# data/pp_project_schoolsonlinethesaurus.jsonld from PoolParty
# data/Sofia-API-Meta-Data-*.json
# data/Sofia-API-Node-Data-*.json
# data/Sofia-API-Tree-Data-*.json

cd data
cp mapping-*.csv mapping.csv
cp Sofia-API-Meta-Data-*.json Sofia-API-Meta-Data.json
cp Sofia-API-Node-Data-*.json Sofia-API-Node-Data.json
cp Sofia-API-Tree-Data-*.json Sofia-API-Tree-Data.json
cd ..
echo "SCOT mapping preprocessing..."
cd 2023-Nov
go run main.go
cd ..
echo "Node preprocessing..."
cd node2
go test
echo "ScOT preprocessing..."
cd ../asn-json/tool
go test
echo "Node restructuring..."
cd ../../tree
go test
cd ..
cp data-out/ccp* data-out/restructure
cp data-out/gc* data-out/restructure
echo "Translation to ASN JSON..."
cd asn-json-v2
go test
echo "Arrays on asn_hasLevel..."
cd ../task-force-array-2023-06-21
go run main.go
echo "Translation to ASN JSON LD..."
cd ../asn-json-ld
go test
echo "remove dc:text, add dc:description@language, @value..."
cd ../task-remove-fields-2023-06-23
go run main.go
echo "Break up CCP..."
cd ../task-split-ccp-2022-11-12
go run main.go
echo "Restructure CCP..."
cd ../task-structure-improvement-2023-03-23
go run main.go
cd ..
rm -rf ./release
mkdir ./release
cp -rf ./data-out/asn-json ./release
cp -rf ./data-out/asn-json-ld ./release
cp -rf ./task-split-ccp-2022-11-12/asn-json-ccp ./release
cp -rf ./task-structure-improvement-2023-03-23/asn-json-ld-ccp ./release
echo "Copy dc:title to dc:description..."
cd task-add-dctitle-2023-10-10
go run main.go
cd ..
echo "Validate..."
cp validate/validate.rb release
cp data/pp_project_schoolsonlinethesaurus.jsonld release/asn-json-ld
cd release
ruby validate.rb
