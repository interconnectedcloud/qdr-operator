#!/usr/bin/env bash

if [[ -z ${CI} ]]; then
    ./hack/go-vet.sh
    ./hack/go-fmt.sh
    ./hack/catalog-source.sh
fi

#local test
go test -count=1 `go list ./test/... | grep -v e2e`
