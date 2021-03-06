package main

import (
	"github.com/hashicorp/terraform/plugin"

	"github.com/rickardl/terraform-artifactory-provider/pkg/artifactory"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: artifactory.Provider,
	})
}
