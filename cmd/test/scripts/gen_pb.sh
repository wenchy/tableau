#!/bin/bash

# set -eux
set -e
set -o pipefail

cd "$(git rev-parse --show-toplevel)"

output_dir="./cmd/test/"

tableau_dir="./pkg/protobuf/"
for item in "$tableau_dir"/* ; do
    echo "$item"
    if [ -f "$item" ]; then
        protoc -I"$tableau_dir" --go_out="$output_dir" "$item"
    fi
done

test_dir="./cmd/test/protobuf"
for item in "$test_dir"/* ; do
    echo "$item"
    if [ -f "$item" ]; then
        protoc -I"$test_dir" -I"$tableau_dir" --go_out="$output_dir" "$item"
    fi
done

tableaupb_dir="./cmd/test/github.com/Wenchy/tableau/pkg/tableaupb"
testpb_dir="./cmd/test/github.com/Wenchy/tableau/cmd/test/testpb"

# update tableaupb
rsync -avz "$tableaupb_dir" "./pkg/"
# update testpb
rsync -avz "$testpb_dir" "./cmd/test/"

# remove
rm -rf "./cmd/test/github.com"