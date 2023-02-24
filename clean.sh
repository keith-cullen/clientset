#!/bin/bash

kubectl delete item item1 -n itemns
go clean
kubectl delete -f item_crd.yaml
rm -rf ${GOPATH}/github.com/keith-cullen/clientset
rm -rf pkg/client
rm -rf pkg/apis/itemresource/v1/zz_generated.deepcopy.go
