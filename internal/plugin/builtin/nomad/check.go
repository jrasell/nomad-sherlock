package nomad

import (
	"github.com/hashicorp/nomad/api"
	"github.com/jrasell/nomad-sherlock/sdk"
)

var defaultChecks = map[string]*sdk.Check{
	"nomad_1": {
		ID:            "nomad_1",
		Desc:          "access control list should be enabled",
		CheckExecutor: NewACLCheckExecution,
		Detail: `
ACLs provide a fundamental security enhancement for Nomad clusters. ACL
policies allow for fine grained control of access to the Nomad API alongside
short lived tokens to minimize risk in cases of exposure.

Docs: https://www.nomadproject.io/docs/configuration/acl
Learn Guide: https://learn.hashicorp.com/collections/nomad/access-control
`,
	},
	"nomad_2": {
		ID:            "nomad_2",
		Desc:          "Nomad agents should be running the same binary version",
		CheckExecutor: newVersionCheckExecution,
		Detail: `
Ensuring the running version of Nomad is the same across all agents ensures
there are no compatibility issue. Different versions running within the same
cluster can quickly cause issues, fractured configuration and lead to technical
debt as a result of having to rectify the issue.

Docs: https://www.nomadproject.io/downloads
Learn Guide: https://learn.hashicorp.com/tutorials/nomad/get-started-install
`,
	},
	"nomad_3": {
		ID:            "nomad_3",
		Desc:          "Nomad agents should use TLS on the HTTP API",
		CheckExecutor: newTLSHTTPCheckExecution,
		Detail: `
Securing Nomad's cluster communication is not only important for security but
can even ease operations by preventing mistakes and misconfigurations. Nomad
should use mutual TLS (mTLS) for all HTTP communication.

Docs: https://www.nomadproject.io/docs/configuration/tls
Learn Guide: https://learn.hashicorp.com/tutorials/nomad/security-enable-tls
`,
	},
	"nomad_4": {
		ID:            "nomad_4",
		Desc:          "Nomad agents should use TLS for Consul connectivity",
		CheckExecutor: newTLSHTTPCheckExecution,
		Detail: `
Docs: https://www.nomadproject.io/docs/configuration/consul
Learn Guide: https://learn.hashicorp.com/tutorials/consul/tls-encryption-secure
`,
	},
}

type PluginCheckExecutor struct {
	impl []sdk.CheckExecutor
}

func NewAgentInfoExecutor(checks map[string]*sdk.Check) *PluginCheckExecutor {

	executor := PluginCheckExecutor{}

	for _, c := range checks {
		executor.impl = append(executor.impl, c.CheckExecutor(c.ID))
	}
	return &executor
}

func (p *PluginCheckExecutor) Update(info *api.AgentSelf) {
	for _, impl := range p.impl {
		impl.Update(info)
	}
}

func (p *PluginCheckExecutor) Results() []sdk.CheckResult {
	var out []sdk.CheckResult
	for _, impl := range p.impl {
		out = append(out, impl.Results())
	}
	return out
}
