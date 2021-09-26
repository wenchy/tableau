#!/bin/bash

# set -eux
set -e
set -o pipefail

cd "$(git rev-parse --show-toplevel)"

bash ./scripts/gen_pb.sh

test_indir="./cmd/test/protobuf"
test_outdir="./cmd/test/testpb"

# remove *.go
rm -fv $test_outdir/*.go

for item in "$test_indir"/* ; do
    echo "$item"
    if [ -f "$item" ]; then
        protoc \
        --go_out="$test_outdir" \
        --go_opt=paths=source_relative \
        --proto_path="$test_indir" \
        --proto_path="$tableau_indir" \
        "$item"
    fi
done