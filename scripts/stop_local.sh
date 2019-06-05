#!/bin/bash

if [ -z "$(kind get clusters | grep evntsrc)" ]; then
	echo "!! Nothing to stop"
	exit 0;
fi

echo "Stopping local cluster"
unset KUBECONFIG
kind delete cluster --name=evntsrc