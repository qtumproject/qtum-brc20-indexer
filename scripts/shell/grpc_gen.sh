#!/usr/bin/env bash

# 获取当前工作目录的父级目录名和当前目录名
APP_RELATIVE_PATH="$(basename "$(dirname "$(pwd)")")/$(basename "$(pwd)")"
# 使用双引号括起来以避免空格问题
INTERNAL_PROTO_FILES=$(find internal -name *.proto)
API_PROTO_FILES=$(cd "../../../api/${APP_RELATIVE_PATH}" && find . -name *.proto)
if [ -z "$API_PROTO_FILES" ]; then
    echo "API_PROTO_FILES is empty, skip gen proto files"
else
    echo "gen proto files: ${API_PROTO_FILES}"
    cd "../../../api/${APP_RELATIVE_PATH}" && protoc --proto_path=. \
               --proto_path=../../../third_party \
               --go_out=paths=source_relative:. \
               --go-grpc_out=paths=source_relative:. \
               --go-gin_out . --go-gin_opt=paths=source_relative \
               "${API_PROTO_FILES}" && protoc-go-inject-tag -input="./*/*.pb.go"
fi
