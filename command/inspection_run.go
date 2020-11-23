package command

import (
	"encoding/json"
	"fmt"
	"github.com/jrasell/nomad-sherlock/internal/agent"
)

type InspectionRunCommand struct {
	Meta
}

func (i *InspectionRunCommand) Help() string {
	return ""
}

func (i *InspectionRunCommand) Synopsis() string {
	return "Runs a Nomad Sherlock inspection"
}

func (i *InspectionRunCommand) Name() string { return "inspection run" }

func (i *InspectionRunCommand) Run(args []string) int {

	flags := i.Meta.FlagSet(i.Name(), FlagSetClient)
	flags.Usage = func() { i.Ui.Output(i.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	a := agent.NewAgent()

	if err := a.Initialize(); err != nil {
		i.Ui.Error(fmt.Sprintf("Error initializing agent: %v", err))
		return 1
	}

	result := a.RunInspection()

	if i.Meta.address == "" {
		out, err := json.Marshal(result)
		if err != nil {
			i.Ui.Error(fmt.Sprintf("Error marshalling report: %v", err))
			return 1
		}
		i.Ui.Output(string(out))
		return 0
	}

	client, err := i.Meta.Client()
	if err != nil {
		i.Ui.Error(fmt.Sprintf("Error initializing client: %v", err))
		return 1
	}

	id, err := client.Inspection().Submit(result)
	if err != nil {
		i.Ui.Error(fmt.Sprintf("Error submittied report: %v", err))
		return 1
	}
	i.Ui.Output(fmt.Sprintf("Inspection ID: %s", id.ID))
	return 0
}
