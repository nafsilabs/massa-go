#!/bin/bash

set -e

# Create output directory
mkdir -p client/proto

# Clean up any previous downloads
rm -rf proto-gen

mkdir -p proto-gen

# Get protos from massa-proto repo
echo "Downloading massa-proto..."
curl -sSL https://github.com/massalabs/massa-proto/archive/refs/heads/main.zip -o proto-gen/massa-proto.zip
cd proto-gen

unzip -q massa-proto.zip

cd massa-proto-main

# Update go_package in proto files
echo "Updating go_package paths..."
arr=("./proto/abis/massa/abi/v1" "./proto/apis/massa/api/v1" "./proto/commons/massa/model/v1")

for dir in "${arr[@]}"
do
    for filename in "$dir"/*.proto; do
        if [ -f "$filename" ]; then
            sed -i '' 's|option go_package = "github.com/massalabs|option go_package = "github.com/jwmdev/massa-go/client/proto|' "$filename"
        fi
    done
done

# Generate Go code using buf workspace
echo "Generating Go code..."
buf generate --template ../../buf.gen.yaml

# Copy generated files to the target directory
echo "Copying generated files..."
cp -r client/proto/* ../../client/proto/

cd ../..

echo "Cleaning up..."
rm -rf proto-gen/

echo "Done! Generated files are in client/proto/"