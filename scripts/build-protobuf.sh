#!/usr/bin/env bash

API_VERSION=0.1.3

rm -rf protobuf veidemann_api
mkdir protobuf veidemann_api
cd protobuf
wget -q https://github.com/google/protobuf/releases/download/v3.5.1/protoc-3.5.1-linux-x86_64.zip
unzip protoc-3.5.1-linux-x86_64.zip
rm protoc-3.5.1-linux-x86_64.zip
wget -O - -q https://github.com/nlnwa/veidemann-api/archive/${API_VERSION}.tar.gz | tar --strip-components=2 -zx
go get github.com/golang/protobuf/proto
go get github.com/golang/protobuf/protoc-gen-go

bin/protoc -I. --go_out=plugins=grpc:../veidemann_api *.proto
