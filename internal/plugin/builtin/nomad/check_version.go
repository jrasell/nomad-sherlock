package nomad

import (
	"fmt"
	"github.com/hashicorp/nomad/api"
	"github.com/jrasell/nomad-sherlock/sdk"
	"strconv"
)

type versionCheckExecution struct {
	id       string
	unknown  int
	versions map[string]int
}

func newVersionCheckExecution(id string) sdk.CheckExecutor {
	return &versionCheckExecution{
		id:       id,
		versions: make(map[string]int),
	}
}

func (v *versionCheckExecution) Update(i interface{}) {

	agentSelf, ok := i.(*api.AgentSelf)
	if !ok {
		v.unknown++
		return
	}

	versionConfig, ok := agentSelf.Config["Version"]
	if !ok {
		v.unknown++
		return
	}

	versionConfigImpl, ok := versionConfig.(map[string]interface{})
	if !ok {
		v.unknown++
		return
	}

	version, ok := versionConfigImpl["Version"].(string)
	if !ok {
		v.unknown++
		return
	}

	if preRelease := versionConfigImpl["VersionPrerelease"]; preRelease != "" {
		version += fmt.Sprintf("-%s", preRelease.(string))
	}
	v.versions[version] ++
}

func (v *versionCheckExecution) Results() sdk.CheckResult {

	res := sdk.CheckResult{
		ID: v.id,
		Meta: map[string]string{
			"unknown": strconv.Itoa(v.unknown),
		},
	}

	for v, c := range v.versions {
		res.Meta[v] = strconv.Itoa(c)
	}

	if v.unknown > 0 {
		res.Msg = "failed to interrogate all Nomad agents"
		res.State = sdk.CheckResultStateUnknown
	} else if len(v.versions) > 1 {
		res.Msg = "inconsistent Nomad agent versions found"
		res.State = sdk.CheckResultStateFail
	} else {
		res.Msg = "all Nomad agent versions match"
		res.State = sdk.CheckResultStatePass
	}

	return res
}
