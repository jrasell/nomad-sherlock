package nomad

import (
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/nomad/api"
	"github.com/jrasell/nomad-sherlock/sdk"
)

type Plugin struct {
	cfg           map[string]string
	enabledChecks map[string]*sdk.Check
	client        *api.Client
}

func New() sdk.CheckPluginProvider {
	return &Plugin{}
}

func (p *Plugin) SetConfig(cfg map[string]string) error {
	p.cfg = cfg
	p.client, _ = api.NewClient(api.DefaultConfig())

	p.enabledChecks = defaultChecks

	if ids := cfg[sdk.PluginConfigKeyDisabledChecks]; ids != "" {
		p.enabledChecks = sdk.FilterChecks(strings.Split(ids, ","), defaultChecks)
	}
	return nil
}

func (p *Plugin) ExecuteChecks() ([]sdk.CheckResult, error) {

	nodes, _, err := p.client.Nodes().List(nil)
	if err != nil {
		return nil, err
	}

	var mErr *multierror.Error

	executors := NewAgentInfoExecutor(p.enabledChecks)

	for _, node := range nodes {
		nc, err := p.client.GetNodeClient(node.ID, nil)
		if err != nil {
			mErr = multierror.Append(mErr, err)
			continue
		}

		agentSelf, err := nc.Agent().Self()
		if err != nil {
			mErr = multierror.Append(mErr, err)
			continue
		}

		executors.Update(agentSelf)
	}

	return executors.Results(), mErr.ErrorOrNil()
}

func (p *Plugin) GetChecks() []*sdk.CheckInfo {
	out := []*sdk.CheckInfo{}

	for _, check := range p.enabledChecks {
		out = append(out, &sdk.CheckInfo{
			ID:     check.ID,
			Desc:   check.Desc,
			Detail: check.Detail,
		})
	}

	return out
}

func (p *Plugin) GetCheck(id string) *sdk.CheckInfo {
	if check := defaultChecks[id]; check != nil {
		return &sdk.CheckInfo{
			ID:     check.ID,
			Desc:   check.Desc,
			Detail: check.Detail,
		}
	}
	return nil
}
