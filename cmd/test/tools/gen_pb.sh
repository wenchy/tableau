#!/bin/bash

output_dir="../"

tableau_dir="../../../pkg/protobuf/"
for item in "$tableau_dir"/* ; do
    echo "$item"
    if [ -f "$item" ]; then
        protoc -I"$tableau_dir" --go_out="$output_dir" "$item"
    fi
done

test_dir="../protobuf"
for item in "$test_dir"/* ; do
    echo "$item"
    if [ -f "$item" ]; then
        protoc -I"$test_dir" -I"$tableau_dir" --go_out="$output_dir" "$item"
    fi
done
