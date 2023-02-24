package main

import (
	"github.com/justcy/cli/cmd"

	// register commands
	_ "github.com/justcy/cli/cmd/call"
	_ "github.com/justcy/cli/cmd/completion"
	_ "github.com/justcy/cli/cmd/describe"
	_ "github.com/justcy/cli/cmd/generate"
	_ "github.com/justcy/cli/cmd/new"
	_ "github.com/justcy/cli/cmd/run"
	_ "github.com/justcy/cli/cmd/services"
	_ "github.com/justcy/cli/cmd/stream"

	// plugins
	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
)

func main() {
	cmd.Run()
}
