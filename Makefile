TARGET_LINUX = GOARCH=amd64 GOOS=linux
CONTAINER_REGISTRY = eu.gcr.io/txsvc-dev

.PHONY: all
all: test

.PHONY: test
test:
	cd pkg/api && go test
	
.PHONY: test_coverage
test_coverage:
	go test `go list ./...` -coverprofile=coverage.txt -covermode=atomic

.PHONY: test_build
test_build:
	cd cmd/cdn && go build main.go && rm main

.PHONY: build_cdn
build_cdn:
	cd cmd/cdn && ${TARGET_LINUX} go build -o svc main.go && docker build -t ${CONTAINER_REGISTRY}/cdn . && docker push ${CONTAINER_REGISTRY}/cdn

.PHONY: example_api_gcp
example_api:
	cd examples/api && gcloud app deploy . --quiet