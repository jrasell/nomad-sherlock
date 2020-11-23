package command

import (
	"fmt"
	"strings"
)

type InspectionListCommand struct {
	Meta
}

func (i *InspectionListCommand) Name() string { return "inspection list" }

func (i *InspectionListCommand) Help() string {
	helpText := `
Usage: nomad-sherlock inspection list <subcommand> [options]
`
	return strings.TrimSpace(helpText)
}

func (i *InspectionListCommand) Synopsis() string {
	return "List all the available Nomad Sherlock inspection results"
}

func (i *InspectionListCommand) Run(args []string) int {

	flags := i.Meta.FlagSet(i.Name(), FlagSetClient)
	flags.Usage = func() { i.Ui.Output(i.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if len(flags.Args()) > 0 {
		i.Ui.Error("Command takes zero arguments")
		return 1
	}

	// Look locally for inspections and list them in the future!
	if i.Meta.address == "" {
		i.Ui.Error("Nomad Sherlock cannot output local inspections currently")
		return 1
	}

	client, err := i.Meta.Client()
	if err != nil {
		i.Ui.Error(fmt.Sprintf("Error initializing client: %v", err))
		return 1
	}

	reports, err := client.Inspection().List(nil)
	if err != nil {
		i.Ui.Error(fmt.Sprintf("Error listing inspections: %v", err))
		return 1
	}

	output := []string{"ID|Region|Start Time|Total Checks|Passed|Failed|Unknown"}

	for _, report := range reports {
		output = append(output, fmt.Sprintf("%s|%s|%s|%v|%v|%v|%v",
			report.ID, report.Region, report.StartTime.String(), report.TotalRun,
			report.TotalPass, report.TotalFail, report.TotalUnknown))
	}
	i.Ui.Output(formatList(output))

	return 0
}
