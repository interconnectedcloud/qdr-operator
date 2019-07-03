
.PHONY: all
all: build

.PHONY: dep
dep:
	./hack/go-dep.sh

.PHONY: format
format:
	./hack/go-fmt.sh

.PHONY: sdk-generate
sdk-generate: dep
	operator-sdk generate k8s

.PHONY: vet
vet:
	./hack/go-vet.sh

.PHONY: test
test:
	./hack/go-test.sh

.PHONY: cluster-test
cluster-test:
	go test --count=1 -v "./test/e2e"

.PHONY: build
build:
	./hack/go-build.sh
