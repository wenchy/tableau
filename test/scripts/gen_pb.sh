#!/bin/bash

# set -eux
set -e
set -o pipefail

cd "$(git rev-parse --show-toplevel)"

bash ./scripts/gen_pb.sh

tableau_proto="./proto"
test_indir="./test/protoconf"
test_outdir="./test/testpb"

# remove *.go
rm -fv $test_outdir/*.go

gen_pb()
{
    for item in "$1"/* ; do
        echo "$item"
        if [ -f "$item" ]; then
            protoc \
            --go_out="$test_outdir" \
            --go_opt=paths=source_relative \
            --proto_path="$1" \
            --proto_path="$tableau_proto" \
            "$item"
        fi
    done
}
gen_pb $test_indir/xlsx
gen_pb $test_indir/xml