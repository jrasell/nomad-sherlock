package main

import (
	"fmt"
	"os"

	"github.com/jrasell/nomad-sherlock/command"
	"github.com/mitchellh/cli"
)

func main() {

	// Create the meta object
	metaPtr := new(command.Meta)

	// The Nomad agent never outputs color
	metaPtr.Ui = &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	versionString := fmt.Sprintf("Nomad Sherlock %s", "hack")
	c := cli.NewCLI("nomad-sherlock", versionString)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"check": func() (cli.Command, error) {
			return &command.CheckCommand{
				Meta: *metaPtr,
			}, nil
		},
		"check info": func() (cli.Command, error) {
			return &command.CheckInfoCommand{
				Meta: *metaPtr,
			}, nil
		},
		"check list": func() (cli.Command, error) {
			return &command.CheckListCommand{
				Meta: *metaPtr,
			}, nil
		},
		"inspection": func() (cli.Command, error) {
			return &command.InspectionCommand{}, nil
		},
		"inspection list": func() (cli.Command, error) {
			return &command.InspectionListCommand{
				Meta: *metaPtr,
			}, nil
		},
		"inspection info": func() (cli.Command, error) {
			return &command.InspectionInfoCommand{
				Meta: *metaPtr,
			}, nil
		},
		"inspection run": func() (cli.Command, error) {
			return &command.InspectionRunCommand{
				Meta: *metaPtr,
			}, nil
		},
		"server": func() (cli.Command, error) {
			return &command.ServerCommand{
				Meta: *metaPtr,
			}, nil
		},
	}

	exitCode, err := c.Run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error executing CLI: %v\n", err)
	}
	os.Exit(exitCode)
}
