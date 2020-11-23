package command

import (
	"fmt"
	"strings"

	"github.com/jrasell/nomad-sherlock/internal/agent"
	"github.com/jrasell/nomad-sherlock/sdk"
)

type CheckInfoCommand struct {
	Meta
}

func (c *CheckInfoCommand) Name() string { return "check info" }

func (c *CheckInfoCommand) Help() string {
	helpText := `
Usage: nomad-sherlock check info [options] <id>
`
	return strings.TrimSpace(helpText)
}

func (c *CheckInfoCommand) Synopsis() string {
	return "Detail an individual Nomad Sherlock check"
}

func (c *CheckInfoCommand) Run(args []string) int {

	flags := c.Meta.FlagSet(c.Name(), FlagSetClient)
	flags.Usage = func() { c.Ui.Output(c.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if len(flags.Args()) < 1 || len(flags.Args()) > 1 {
		c.Ui.Error("Command takes one argument: <id>")
		return 1
	}

	if c.Meta.address == "" {

		a := agent.NewAgent()

		if err := a.Initialize(); err != nil {
			c.Ui.Error(fmt.Sprintf("Error initializing agent: %v", err))
			return 1
		}
		check := a.Plugin().GetCheck(args[0])
		c.outputCheck(check)
		return 0
	}

	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %v", err))
		return 1
	}

	check, err := client.Check().Info(flags.Args()[0], nil)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error detailing check: %v", err))
		return 1
	}

	c.outputCheck(check)
	return 0
}

func (c *CheckInfoCommand) outputCheck(check *sdk.CheckInfo) {
	c.Ui.Output(fmt.Sprintf("%s: %s\n%s", check.ID, check.Desc, check.Detail))
}
