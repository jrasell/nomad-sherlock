package command

import (
	"github.com/mitchellh/cli"
	"strings"
)

type InspectionCommand struct{}

func (i *InspectionCommand) Help() string {
	helpText := `
Usage: nomad-sherlock inspection <subcommand> [options]

  This command groups subcommands for performing inspection runs as well as
  interacting with them.

  Execute an inspection:

      $ nomad-sherlock inspection run

  Please see the individual subcommand help for detailed usage information.
`
	return strings.TrimSpace(helpText)
}

func (i *InspectionCommand) Synopsis() string {
	return "Interact and run Nomad Sherlock inspections"
}

func (i *InspectionCommand) Run(_ []string) int {
	return cli.RunResultHelp
}
