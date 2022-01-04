.PHONY: build

build:
	goreleaser build --snapshot --rm-dist
