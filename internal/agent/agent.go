package agent

import (
	"sync"

	"github.com/jrasell/nomad-sherlock/internal/plugin/builtin/nomad"
	"github.com/jrasell/nomad-sherlock/sdk"
)

type Agent struct {
	lock         sync.RWMutex
	checkPlugins map[string]sdk.CheckPluginProvider
}

func NewAgent() *Agent {
	return &Agent{
		checkPlugins: make(map[string]sdk.CheckPluginProvider),
	}
}

func (r *Agent) Initialize() error {
	return r.launchPlugins()
}

func (r *Agent) launchPlugins() error {

	r.lock.Lock()
	defer r.lock.Unlock()

	r.checkPlugins["nomad"] = nomad.New()
	return r.checkPlugins["nomad"].SetConfig(nil)
}

func (r *Agent) RunInspection() *sdk.Inspection {

	report := sdk.NewInspection()

	for _, check := range r.checkPlugins {

		res, err := check.ExecuteChecks()
		report.Update(res, err)
	}
	report.Finalize()

	return report
}

func (r *Agent) Plugin() sdk.CheckPluginProvider {
	return r.checkPlugins["nomad"]
}
