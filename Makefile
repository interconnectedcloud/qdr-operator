export REGISTRY=quay.io/interconnectedcloud
export IMAGE=qdr-operator
export TAG=1.4.0-beta1
export DOCKER=docker


.PHONY: all
all: build

.PHONY: mod
mod:
	./scripts/go-mod.sh

.PHONY: format
format:
	./hack/go-fmt.sh

.PHONY: go-generate
go-generate: mod
	./hack/go-gen.sh

.PHONY: sdk-generate
sdk-generate: mod
	./hack/go-gen.sh

.PHONY: vet
vet:
	./hack/go-vet.sh

.PHONY: test
test:
	./hack/go-test.sh

.PHONY: cluster-test
cluster-test:
	go test --timeout=30m --count=1 -v "./test/e2e"

.PHONY: build
build:
	./hack/go-build.sh

.PHONY: docker-push
docker-push:
	${DOCKER} push ${REGISTRY}/${IMAGE}:${TAG}

.PHONY: clean
clean:
	rm -rf build/_output
	rm -rf vendor
