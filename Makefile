tag:
	@git tag $(version)
	@git push origin $(version)

lint:
	@golangci-lint run

build: lint
	@WORKINGDIR=$(pwd) goreleaser build --snapshot --rm-dist

release: lint
	@WORKINGDIR=$(pwd) goreleaser release --snapshot --rm-dist
