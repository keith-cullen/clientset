#!/bin/bash

touch boilerplate.txt
source "${GOPATH}/src/k8s.io/code-generator/kube_codegen.sh"
kube::codegen::gen_helpers \
    --boilerplate boilerplate.txt \
    ./pkg
kube::codegen::gen_client \
    --boilerplate boilerplate.txt \
    --output-pkg "github.com/keith-cullen/clientset/pkg/client" \
    --output-dir "./pkg/client" \
    --with-watch \
    ./pkg/apis
