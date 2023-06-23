#!/bin/bash

set -e

find . -type f -executable -exec sh -c "file -i '{}' | grep -q 'x-executable; charset=binary'" \; -print | xargs rm -f

# delete others
for f in "`find ./ -name '*.log' -or -name '*.temp'`"; do rm "$f"; done

rm -rf ./data-out