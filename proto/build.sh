#!/bin/bash
set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
PROTO_FILE="$SCRIPT_DIR/api.proto"
OUT_DIR="$SCRIPT_DIR/pb"

mkdir -p "$OUT_DIR"

protoc \
  --proto_path="$SCRIPT_DIR" \
  --go_out="$OUT_DIR" --go_opt=paths=source_relative \
  --go-grpc_out="$OUT_DIR" --go-grpc_opt=paths=source_relative \
  "$(basename "$PROTO_FILE")"

echo "gRPC code generated in $OUT_DIR"