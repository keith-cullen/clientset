#!/bin/bash

go clean
kubectl delete item item1 -n itemns
kubectl delete -f item_crd.yaml
rm -rf pkg/client
rm -rf pkg/apis/itemresource/v1/zz_generated.defaults.go
rm -rf pkg/apis/itemresource/v1/zz_generated.deepcopy.go
