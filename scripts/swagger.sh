#!/bin/bash

cd ./pkg;

files=`find . -type f -name "*.proto" -not -path "./event/*" -not -path "./utils/*"`

for file in $files;
do
	protoc \
		-I. \
		-I$GOPATH/src \
		-I../vendor \
		-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--swagger_out=logtostderr=true:../api \
		$file;
done