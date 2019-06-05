#!/bin/bash

if ! [ -x "$(command -v kind)" ]; then
	echo "=== Kind not found; Installing ==="
	GO111MODULE="on" go get sigs.k8s.io/kind@v0.3.0
fi

if [ -z "$(kind get clusters | grep evntsrc)" ]; then
	echo "=== Starting local cluster ==="
	kind create cluster --name=evntsrc
else 
	echo "=== Using existing local cluster ==="
fi

export KUBECONFIG="$(kind get kubeconfig-path --name="evntsrc")"