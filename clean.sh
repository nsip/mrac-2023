#!/bin/bash

set -e

find . -type f -executable -exec sh -c "file -i '{}' | grep -q 'x-executable; charset=binary'" \; -print | xargs rm -f

# delete others
# find ./ -type f \( -name "*.log" -o -name "*.temp" \) -exec rm {} \;
find ./ \( -type f \( -name "*.log" -o -name "*.temp" \) -o -type d -name "fatal" \) -exec rm -rf {} +

rm -rf ./data-out