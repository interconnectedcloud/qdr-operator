#!/bin/sh
REGISTRY=quay.io/interconnectedcloud
IMAGE=qdr-operator
TAG=1.0.0-beta6

if [[ -z ${CI} ]]; then
	./hack/go-test.sh
	operator-sdk build ${REGISTRY}/${IMAGE}:${TAG}
else
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -o build/_output/bin/qdr-operator github.com/interconnectedcloud/qdr-operator/cmd/manager
fi
