package generate

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	mcli "github.com/justcy/cli/cmd"
	"github.com/justcy/cli/generator"
	tmpl "github.com/justcy/cli/generator/template"
	"github.com/urfave/cli/v2"
)

func init() {
	mcli.Register(&cli.Command{
		Name:  "generate",
		Usage: "Generate project template files after the fact",
		Subcommands: []*cli.Command{
			{
				Name:   "kubernetes",
				Usage:  "Generate Kubernetes resource template files",
				Action: Kubernetes,
			},
			{
				Name:   "skaffold",
				Usage:  "Generate Skaffold template files",
				Action: Skaffold,
			},
			{
				Name:   "sqlc",
				Usage:  "Generate sqlc resources",
				Action: Sqlc,
			},
		},
	})
}

// Kubernetes generates Kubernetes resource template files in the current
// working directory. Exits on error.
func Kubernetes(ctx *cli.Context) error {
	service, err := getService()
	if err != nil {
		return err
	}

	vendor, err := getServiceVendor(service)
	if err != nil {
		return err
	}

	g := generator.New(
		generator.Service(service),
		generator.Vendor(vendor),
		generator.Directory("."),
		generator.Client(strings.HasSuffix(service, "-client")),
	)

	files := []generator.File{
		{Path: "plugins.go", Template: tmpl.Plugins},
		{Path: "resources/clusterrole.yaml", Template: tmpl.KubernetesClusterRole},
		{Path: "resources/configmap.yaml", Template: tmpl.KubernetesEnv},
		{Path: "resources/deployment.yaml", Template: tmpl.KubernetesDeployment},
		{Path: "resources/rolebinding.yaml", Template: tmpl.KubernetesRoleBinding},
	}

	g.Generate(files)

	return nil
}

// Skaffold generates Skaffold template files in the current working directory.
// Exits on error.
func Skaffold(ctx *cli.Context) error {
	service, err := getService()
	if err != nil {
		return err
	}

	vendor, err := getServiceVendor(service)
	if err != nil {
		return err
	}

	g := generator.New(
		generator.Service(service),
		generator.Vendor(vendor),
		generator.Directory("."),
		generator.Client(strings.HasSuffix(service, "-client")),
		generator.Skaffold(true),
	)

	files := []generator.File{
		{Path: ".dockerignore", Template: tmpl.DockerIgnore},
		{Path: "go.mod", Template: tmpl.Module},
		{Path: "plugins.go", Template: tmpl.Plugins},
		{Path: "resources/clusterrole.yaml", Template: tmpl.KubernetesClusterRole},
		{Path: "resources/configmap.yaml", Template: tmpl.KubernetesEnv},
		{Path: "resources/deployment.yaml", Template: tmpl.KubernetesDeployment},
		{Path: "resources/rolebinding.yaml", Template: tmpl.KubernetesRoleBinding},
		{Path: "skaffold.yaml", Template: tmpl.SkaffoldCFG},
	}

	if err := g.Generate(files); err != nil {
		return err
	}

	fmt.Println("skaffold project template files generated")

	return nil
}

// Sqlc generates sqlc files in the current working directory.
// Exits on error.
func Sqlc(ctx *cli.Context) error {
	service, err := getService()
	if err != nil {
		return err
	}

	vendor, err := getServiceVendor(service)
	if err != nil {
		return err
	}

	g := generator.New(
		generator.Service(service),
		generator.Vendor(vendor),
		generator.Directory("."),
		generator.Client(strings.HasSuffix(service, "-client")),
		generator.Sqlc(true),
	)

	files := []generator.File{
		{Path: "postgres/queries/example.sql", Template: tmpl.QueryExample},
		{Path: "postgres/migrations/", Template: ""},
	}

	path := "postgres/postgres.go"
	if _, err := os.Stat("./" + path); err != nil {
		files = append(files, generator.File{Path: path, Template: tmpl.Postgres})
	}

	path = "postgres/sqlc.yaml"
	if _, err := os.Stat("./" + path); err != nil {
		files = append(files, generator.File{Path: path, Template: tmpl.Sqlc})
	}

	if err := g.Generate(files); err != nil {
		return err
	}

	fmt.Println("Sqlc project template files generated")

	return nil
}

func getService() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir[strings.LastIndex(dir, "/")+1:], nil
}

func getServiceVendor(s string) (string, error) {
	f, err := os.Open("go.mod")
	if err != nil {
		return "", err
	}
	defer f.Close()

	line := ""
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "module ") {
			line = scanner.Text()
			break

		}
	}
	if line == "" {
		return "", nil
	}

	module := line[strings.LastIndex(line, " ")+1:]
	if module == s {
		return "", nil
	}

	return module[:strings.LastIndex(module, "/")] + "/", nil
}
