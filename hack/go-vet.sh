#!/usr/bin/env bash

if [[ -z ${CI} ]]; then
    ./hack/go-dep.sh
    operator-sdk generate k8s
    ./hack/update-codegen.sh
fi
go vet ./...
