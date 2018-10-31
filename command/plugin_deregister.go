package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*PluginDeregisterCommand)(nil)
var _ cli.CommandAutocomplete = (*PluginDeregisterCommand)(nil)

type PluginDeregisterCommand struct {
	*BaseCommand
}

func (c *PluginDeregisterCommand) Synopsis() string {
	return "Deregister an existing plugin in the catalog"
}

func (c *PluginDeregisterCommand) Help() string {
	helpText := `
Usage: vault plugin deregister [options] NAME TYPE

  Deregister an existing plugin in the catalog. If the plugin does not exist,
  no action is taken (the command is idempotent). The last argument of type
  takes "auth", "database", or "secret".

  Deregister the plugin named my-custom-plugin:

      $ vault plugin deregister my-custom-plugin auth

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PluginDeregisterCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *PluginDeregisterCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultPlugins(consts.PluginTypeUnknown)
}

func (c *PluginDeregisterCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PluginDeregisterCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 2:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1, got %d)", len(args)))
		return 1
	case len(args) > 2:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	pluginName := strings.TrimSpace(args[0])
	pluginType, err := consts.ParsePluginType(strings.TrimSpace(args[1]))
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	if err := client.Sys().DeregisterPlugin(&api.DeregisterPluginInput{
		Name: pluginName,
		Type: pluginType,
	}); err != nil {
		c.UI.Error(fmt.Sprintf("Error deregistering plugin named %s: %s", pluginName, err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Deregistered plugin (if it was registered): %s", pluginName))
	return 0
}
