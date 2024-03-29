tag:
	@git tag -a $(version) -m "Release $(version)"
	@git push --follow-tags

lint:
	@golangci-lint run ./...

build: lint
	@WORKINGDIR=$(pwd) goreleaser build --snapshot --clean

snapshot:
	@WORKINGDIR=$(pwd) goreleaser release --snapshot --clean

release: lint snapshot
	@docker push ghcr.io/puppetlabs/cat-team-github-metrics:dev
