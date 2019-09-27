#!/bin/bash

selectors="storer websocks stsmetrics streams passport users apigw interconnect streamauth ingress billing wui emails";

for selector in $selectors 
do
    kubectl get pods | grep $selector | awk '{print $1}' | xargs kubectl delete pod
done