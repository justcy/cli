package services

import (
	"fmt"
	"sort"

	mcli "github.com/justcy/cli/cmd"
	"github.com/urfave/cli/v2"
)

func init() {
	mcli.Register(&cli.Command{
		Name:   "services",
		Usage:  "List services in the registry",
		Action: List,
	})
}

// List fetches running services from the registry and lists them. Exits on
// error.
func List(ctx *cli.Context) error {
	r := *mcli.DefaultOptions().Registry
	srvs, err := r.ListServices()
	if err != nil {
		return err
	}

	var services []string
	for _, srv := range srvs {
		services = append(services, srv.Name)
	}

	sort.Strings(services)
	for _, srv := range services {
		fmt.Println(srv)
	}

	return nil
}
