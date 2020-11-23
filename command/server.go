package command

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/jrasell/nomad-sherlock/internal/agent"
	"github.com/jrasell/nomad-sherlock/internal/server"
)

type ServerCommand struct {
	Meta
}

func (s *ServerCommand) Help() string {
	return ""
}

func (s *ServerCommand) Synopsis() string {
	return "Run a Nomad Sherlock server"
}

func (s *ServerCommand) Name() string { return "server" }

func (s *ServerCommand) Run(_ []string) int {

	a := agent.NewAgent()

	if err := a.Initialize(); err != nil {
		s.Ui.Error(fmt.Sprintf("Error initializing agent: %v", err))
		return 1
	}

	srv := server.NewHTTPServer(hclog.New(hclog.DefaultOptions), a)
	srv.Start()
	return 0
}
