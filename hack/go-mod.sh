#!/bin/sh

source ./hack/go-mod-env.sh

if [[ -z ${CI} ]]; then
    go mod tidy
    go mod vendor
else
    go mod vendor -v
fi
