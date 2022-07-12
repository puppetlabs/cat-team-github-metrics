package main

import (
	"os"

	"github.com/puppetlabs/cat-team-github-metrics/cmd"
)

func main() {
	os.Exit(cmd.Execute())
}
