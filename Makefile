tag:
	@git tag -a $(version) -m "Release $(version)"
	@git push --follow-tags

lint:
	@golangci-lint run ./...

build: lint
	@WORKINGDIR=$(pwd) goreleaser build --snapshot --rm-dist

release: lint
	@WORKINGDIR=$(pwd) goreleaser release --snapshot --rm-dist
	@docker push ghcr.io/puppetlabs/cat-team-github-metrics:dev

.PHONY: workflow
workflow:
	@cd workflow
	@go run . > workflow.yml
	@relay workflow save cat-github-metrics -f workflow.yml
	@cd -
