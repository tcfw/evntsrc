#!/bin/bash

files=`find ./internal -type f -name "*.proto"`

dir=""

for file in $files;
do
	dir=`dirname $file`
	mkdir -p $dir/../protos

	protoc \
		--gofast_out=Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api,plugins=grpc:$dir/ \
		--grpc-gateway_out=Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api,logtostderr=true:$dir/ \
		--js_out=import_style=commonjs,binary:web/src/protos \
		--js_out=library=evntsrc,binary:web/src/protos \
		$file \
		-I $dir \
		-I . \
		-I $GOPATH/src \
		-I ./vendor \
		-I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis
done

mkdir -p web/src/protos/google/api/;
mkdir -p web/src/protos/github.com/gogo/protobuf/gogoproto/;
# GOGO proto
protoc \
	--js_out=library=gogo_pb.js,binary:web/src/protos/github.com/gogo/protobuf/gogoproto/ \
	-I $dir \
	-I . \
	-I $GOPATH/src \
	-I ./vendor \
	./vendor/github.com/gogo/protobuf/gogoproto/gogo.proto

# Google annotations
protoc \
	--js_out=library=annotations_pb.js,binary:web/src/protos/google/api/ \
	-I $dir \
	-I . \
	-I $GOPATH/src \
	-I ./vendor \
	-I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api/annotations.proto