
.PHONY: test
test:
	@go test -v ./...

.PHONY: delv
delv:
	@go run main.go -c ./test/delve/config/fresher.yaml

.PHONY: dev
dev:
	@go run main.go -c ./test/dev/config/fresher.yaml

.PHONY: release
release:
	@rm -fR dist
	@goreleaser release
