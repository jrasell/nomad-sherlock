package command

import (
	"fmt"
	"strings"

	"github.com/jrasell/nomad-sherlock/internal/agent"
	"github.com/jrasell/nomad-sherlock/sdk"
)

type CheckListCommand struct {
	Meta
}

func (c *CheckListCommand) Name() string { return "check list" }

func (c *CheckListCommand) Help() string {
	helpText := `
Usage: nomad-sherlock check list <subcommand> [options]
`
	return strings.TrimSpace(helpText)
}

func (c *CheckListCommand) Synopsis() string {
	return "List all the available Nomad Sherlock checks"
}

func (c *CheckListCommand) Run(args []string) int {

	flags := c.Meta.FlagSet(c.Name(), FlagSetClient)
	flags.Usage = func() { c.Ui.Output(c.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if len(flags.Args()) > 0 {
		c.Ui.Error("Command takes zero arguments")
		return 1
	}

	if c.Meta.address == "" {

		a := agent.NewAgent()

		if err := a.Initialize(); err != nil {
			c.Ui.Error(fmt.Sprintf("Error initializing agent: %v", err))
			return 1
		}
		check := a.Plugin().GetChecks()
		c.outputChecks(check)
		return 0
	}

	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %v", err))
		return 1
	}

	check, err := client.Check().List(nil)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing checks: %v", err))
		return 1
	}

	c.outputChecks(check)
	return 0
}

func (c *CheckListCommand) outputChecks(checks []*sdk.CheckInfo) {
	output := []string{"ID|Description"}
	for _, check := range checks {
		output = append(output, fmt.Sprintf("%s|%s", check.ID, check.Desc))
	}

	c.Ui.Output(formatList(output))
}
