
.PHONY: install
install:
	go mod tidy
	go install github.com/golang/mock/mockgen@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/pav5000/smartimports/cmd/smartimports@v0.2.0

.PHONY: generate
generate:
	mkdir -p vendor.protogen
	cp -R ../user-account/api/user_account/ vendor.protogen/
	buf generate
	go run github.com/99designs/gqlgen generate

.PHONY: run
run:
	go run cmd/main.go

.PHONY: format
format:
	smartimports -local "github.com/photo-pixels/user-account/"

.PHONY: lint-full
lint-full:
	goimports -w ./internal/..
	golangci-lint run --config=.golangci.yaml ./...

clear-cache:
	buf mod clear-cache