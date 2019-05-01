#!/bin/sh
REGISTRY=quay.io/interconnectedcloud
IMAGE=qdrouterd-operator
TAG=1.0.0-beta2

if [[ -z ${CI} ]]; then
	./hack/go-test.sh
	operator-sdk build ${REGISTRY}/${IMAGE}:${TAG}
else
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -o build/_output/bin/qdrouterd-operator github.com/interconnectedcloud/qdrouterd-operator/cmd/manager
fi
