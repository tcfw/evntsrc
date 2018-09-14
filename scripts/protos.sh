#!/bin/bash

files=`find ./pkg -type f -name "*.proto"`

dir=""

for file in $files;
do
	dir=`dirname $file`
	mkdir -p $dir/../protos

	protoc \
		--gofast_out=Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api,plugins=grpc:$dir/ \
		--grpc-gateway_out=Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api,logtostderr=true:$dir/ \
		$file \
		-I $dir \
		-I . \
		-I $GOPATH/src \
		-I ./vendor \
		-I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis
done