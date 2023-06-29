#!/bin/bash

set -e

find . -type f -executable -exec sh -c "file -i '{}' | grep -q 'x-executable; charset=binary'" \; -print | xargs rm -f

# delete others
# find ./ -type f \( -name "*.log" -o -name "*.temp" \) -exec rm {} \;
find ./ \( -type f \( -name "*.log" -o -name "*.temp" \) -o -type d -name "fatal" \) -exec rm -rf {} +

if [ -n "$1" ] && [ $1 == "all" ]; then
    rm -rf ./data-out
    echo "executables & /data-out are deleted"
else
    echo "only deleted executables"
fi