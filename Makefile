.PHONY: all
all: test

.PHONY: test
test:
	cd pkg/api && go test
	
.PHONY: test_coverage
test_coverage:
	go test `go list ./...` -coverprofile=coverage.txt -covermode=atomic
