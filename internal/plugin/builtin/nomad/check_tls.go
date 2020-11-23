package nomad

import (
	"github.com/hashicorp/nomad/api"
	"github.com/jrasell/nomad-sherlock/sdk"
)

type tlsHTTPCheckExecution struct {
	id                         string
	unknown, enabled, disabled int
}

func newTLSHTTPCheckExecution(id string) sdk.CheckExecutor {
	return &tlsHTTPCheckExecution{
		id: id,
	}
}

func (t *tlsHTTPCheckExecution) Update(i interface{}) {

	agentSelf, ok := i.(*api.AgentSelf)
	if !ok {
		t.unknown++
		return
	}

	if httpEnabled := agentSelf.Config["TLSConfig"].(map[string]interface{})["EnableHTTP"].(bool); httpEnabled {
		t.enabled++
	} else {
		t.disabled++
	}
}

func (t *tlsHTTPCheckExecution) Results() sdk.CheckResult {
	return sdk.GenericEnabledCheckResult(t.unknown, t.enabled, t.disabled, t.id, "HTTP API TLS enabled")
}
