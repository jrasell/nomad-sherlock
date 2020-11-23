package nomad

import (
	"github.com/hashicorp/nomad/api"
	"github.com/jrasell/nomad-sherlock/sdk"
)

type consulTLSCheckExecution struct {
	id                         string
	unknown, enabled, disabled int
}

func newConsulTLSCheckExecution(id string) sdk.CheckExecutor {
	return &consulTLSCheckExecution{
		id: id,
	}
}

func (c *consulTLSCheckExecution) Update(i interface{}) {

	agentSelf, ok := i.(*api.AgentSelf)
	if !ok {
		c.unknown++
		return
	}

	if sslEnabled := agentSelf.Config["Consul"].(map[string]interface{})["EnableSSL"].(bool); sslEnabled {
		c.enabled++
	} else {
		c.disabled++
	}
}

func (c *consulTLSCheckExecution) Results() sdk.CheckResult {
	return sdk.GenericEnabledCheckResult(c.unknown, c.enabled, c.disabled, c.id, "Consul TLS enabled")
}
