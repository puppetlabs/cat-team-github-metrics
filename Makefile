tag:
	@git tag $(version)
	@git push origin $(version)

lint:
	@golangci-lint run

build:
	@WORKINGDIR=$(pwd) goreleaser build --snapshot --rm-dist --single-target

build-all: lint
	@WORKINGDIR=$(pwd) goreleaser build --snapshot --rm-dist
	@docker build -t chelnak/puppet-github-metrics .
