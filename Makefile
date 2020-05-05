
.PHONY: test
test:
	@go test -v ./...

.PHONY: dev
dev:
	@go run main.go -c ./test/delve/config/fresher.yaml

.PHONY: release
release:
	@goreleaser release --skip-publish
	@goreleaser

