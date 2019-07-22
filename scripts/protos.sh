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
		--js_out=import_style=commonjs_strict,binary:web/src/protos \
		--js_out=library=evntsrc,binary:web/src/protos \
		$file \
		-I $dir \
		-I /usr/local/include \
		-I . \
		-I $GOPATH/src \
		-I vendor/github.com/gogo/googleapis \
		-I vendor
done

mkdir -p web/src/protos/google/api/;
mkdir -p web/src/protos/github.com/gogo/protobuf/gogoproto/;
# GOGO proto
protoc \
	--js_out=library=gogo_pb,binary:web/src/protos/github.com/gogo/protobuf/gogoproto/ \
	-I $dir \
	-I . \
	-I $GOPATH/src \
	-I ./vendor \
	./vendor/github.com/gogo/protobuf/gogoproto/gogo.proto

# Google annotations
protoc \
	--js_out=library=annotations_pb,binary:web/src/protos/google/api/ \
	-I $dir \
	-I . \
	-I $GOPATH/src \
	-I vendor \
	-I vendor/github.com/gogo/googleapis \
	./vendor/github.com/gogo/googleapis/google/api/annotations.proto