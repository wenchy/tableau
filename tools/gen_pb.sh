#!/bin/bash
input_dir="../proto"
output_dir="../"
for item in "$input_dir"/* ; do
    echo "$item"
    if [ -f "$item" ]; then
        protoc -I"$input_dir" --go_out="$output_dir" "$item"
    fi
done
