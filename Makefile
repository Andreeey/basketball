.PHONY: vendor
vendor:
	go mod vendor

.PHONY: lint
lint:
	golangci-lint run

.PHONY: run
run:
	go run main.go
