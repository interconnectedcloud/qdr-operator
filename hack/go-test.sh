#!/bin/sh

if [[ -z ${CI} ]]; then
    ./hack/go-vet.sh
    ./hack/go-fmt.sh
    ./hack/catalog-source.sh
fi

#local test
GOCACHE=off go test `go list ./test/... | grep -v e2e`
