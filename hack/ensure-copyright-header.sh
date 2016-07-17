#!/bin/bash

# Copyright (C) 2016 wikiwi.io
#
# This software may be modified and distributed under the terms
# of the MIT license. See the LICENSE file for details.

shopt -s globstar

read -d '' ASTERISK_HEADER <<EOL
/*
 * Copyright (C) $(date +%Y) wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */
EOL

read -d '' SHARP_HEADER <<EOL
# Copyright (C) $(date +%Y) wikiwi.io
#
# This software may be modified and distributed under the terms
# of the MIT license. See the LICENSE file for details.
EOL

for f in pkg/**/*.go cmd/**/*.go test/**/*.go; do
  if ! grep -q "Copyright (C)" "$f"
  then
    echo "Adding license notice to $f"
    echo -e "${ASTERISK_HEADER}\n" > $f.new && cat "$f" >> "$f.new" && mv "$f.new" "$f"
  fi
done

for f in Makefile*; do
  if ! grep -q "Copyright (C)" "$f"
  then
    echo "Adding license notice to $f"
    echo -e "${SHARP_HEADER}\n" > $f.new && cat "$f" >> "$f.new" && mv "$f.new" "$f"
  fi
done

for f in hack/**/*.sh;  do
  if ! grep -q "Copyright (C)" "$f"
  then
    echo "Adding license notice to $f"
    head -n1 "$f" > $f.new && echo -e "\n${SHARP_HEADER}\n" >> $f.new && tail -n +2 "$f" >> "$f.new" && mv "$f.new" "$f"
  fi
done

