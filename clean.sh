#!/bin/bash

set -e

find . -type f -executable -exec sh -c "file -i '{}' | grep -q 'x-executable; charset=binary'" \; -print | xargs rm -f

# delete others
# find ./ -type f \( -name "*.log" -o -name "*.temp" \) -exec rm {} \;
find ./ \( -type f \( -name "*.log" -o -name "*.temp" \) -o -type d -name "fatal" \) -exec rm -rf {} +

rm -f ./util/*.json

if [ -n "$1" ] && [ $1 == "all" ]; then
    rm -rf ./data-out ./package
    echo "executables & /data-out & /package are deleted"
else
    echo "only deleted executables"
fi