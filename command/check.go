package command

import (
	"github.com/mitchellh/cli"
	"strings"
)

type CheckCommand struct {
	Meta
}

func (c *CheckCommand) Name() string { return "check" }

func (c *CheckCommand) Help() string {
	helpText := `
Usage: nomad-sherlock check <subcommand> [options] [args]

  This command groups subcommands for detailing the checks that are currently
  available to Nomad Sherlock.

  View all currently enabled checks:

      $ nomad-sherlock check list

  View the detail of an individual check:

      $ nomad-sherlock check info <id>

  Please see the individual subcommand help for detailed usage information.
`
	return strings.TrimSpace(helpText)
}

func (c *CheckCommand) Synopsis() string {
	return "Detail the available Nomad Sherlock checks"
}

func (c *CheckCommand) Run(_ []string) int {
	return cli.RunResultHelp
}
