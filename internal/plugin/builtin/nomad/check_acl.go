package nomad

import (
	"strconv"

	"github.com/hashicorp/nomad/api"
	"github.com/jrasell/nomad-sherlock/sdk"
)

type aclCheckExecution struct {
	id                         string
	unknown, enabled, disabled int
}

func NewACLCheckExecution(id string) sdk.CheckExecutor {
	return &aclCheckExecution{
		id: id,
	}
}

func (a *aclCheckExecution) Update(i interface{}) {

	agentSelf, ok := i.(*api.AgentSelf)
	if !ok {
		a.unknown++
		return
	}

	if enabled := agentSelf.Config["ACL"].(map[string]interface{})["Enabled"].(bool); enabled {
		a.enabled++
	} else {
		a.disabled++
	}
}

func (a *aclCheckExecution) Results() sdk.CheckResult {

	res := sdk.CheckResult{
		ID: a.id,
		Meta: map[string]string{
			"unknown":  strconv.Itoa(a.unknown),
			"enabled":  strconv.Itoa(a.enabled),
			"disabled": strconv.Itoa(a.disabled),
		},
	}

	if a.unknown > 0 {
		res.Msg = "failed to interrogate all Nomad agents"
		res.State = sdk.CheckResultStateUnknown
	} else if a.disabled > 0 {
		res.Msg = "ACLs are disabled"
		res.State = sdk.CheckResultStateFail
	} else {
		res.Msg = "ACLs are enabled"
		res.State = sdk.CheckResultStatePass
	}

	return res
}
