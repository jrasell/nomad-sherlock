package command

import (
	"fmt"
	"strings"
)

type InspectionInfoCommand struct {
	Meta
}

func (i *InspectionInfoCommand) Name() string { return "inspection info" }

func (i *InspectionInfoCommand) Help() string {
	helpText := `
Usage: nomad-sherlock check info [options] <id>
`
	return strings.TrimSpace(helpText)
}

func (i *InspectionInfoCommand) Synopsis() string {
	return "Detail an individual inspection result"
}

func (i *InspectionInfoCommand) Run(args []string) int {

	flags := i.Meta.FlagSet(i.Name(), FlagSetClient)
	flags.Usage = func() { i.Ui.Output(i.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if len(flags.Args()) < 1 || len(flags.Args()) > 1 {
		i.Ui.Error("Command takes one argument: <id>")
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

	inspection, err := client.Inspection().Info(flags.Args()[0], nil)
	if err != nil {
		i.Ui.Error(fmt.Sprintf("Error detailing inspection: %v", err))
		return 1
	}

	info := []string{
		fmt.Sprintf("ID|%s", inspection.ID),
		fmt.Sprintf("Region|%s", inspection.Region),
		fmt.Sprintf("Start Time|%s", inspection.StartTime.String()),
		fmt.Sprintf("End Time|%s", inspection.EndTime.String()),
		fmt.Sprintf("Duration|%s", inspection.Duration),
		fmt.Sprintf("Total Run|%v", inspection.TotalRun),
		fmt.Sprintf("Total Pass|%v", inspection.TotalPass),
		fmt.Sprintf("Total Fail|%v", inspection.TotalFail),
		fmt.Sprintf("Total Unknown|%v", inspection.TotalUnknown),
	}

	results := []string{"Check ID|State|Message|Meta"}
	for _, res := range inspection.Results {
		results = append(results, fmt.Sprintf(
			"%s|%s|%s|%s", res.ID, res.Result.State, res.Result.Msg, res.Result.Meta,
		))
	}

	i.Ui.Output(formatKV(info))
	i.Ui.Output("\n")
	i.Ui.Output(formatList(results))
	return 0
}
