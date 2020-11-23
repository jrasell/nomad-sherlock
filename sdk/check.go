package sdk

import "strconv"

type CheckInfo struct {
	ID, Desc, Detail string
}

type Check struct {
	ID, Desc, Detail string
	Results          CheckResult
	CheckExecutor    CheckExecutorSetupFunc
}

type CheckResultState string

const (
	CheckResultStatePass    CheckResultState = "Pass"
	CheckResultStateFail    CheckResultState = "Fail"
	CheckResultStateUnknown CheckResultState = "Unknown"
)

type CheckResult struct {
	ID    string            `json:"id"`
	State CheckResultState  `json:"state"`
	Msg   string            `json:"msg"`
	Meta  map[string]string `json:"meta"`
}

type CheckExecutor interface {
	Update(i interface{})
	Results() CheckResult
}

type CheckExecutorSetupFunc func(id string) CheckExecutor

func FilterChecks(disabled []string, checks map[string]*Check) map[string]*Check {
	out := make(map[string]*Check)
	for _, id := range disabled {
		if c, ok := checks[id]; ok {
			continue
		} else {
			out[id] = c
		}
	}
	return out
}

func GenericEnabledCheckResult(unknown, enabled, disabled int, id, msg string) CheckResult {
	res := CheckResult{
		ID: id,
		Meta: map[string]string{
			"unknown":  strconv.Itoa(unknown),
			"enabled":  strconv.Itoa(enabled),
			"disabled": strconv.Itoa(disabled),
		},
	}

	var msgPrefix string

	if unknown > 0 {
		msgPrefix = "failed to check all Nomad agents for "
		res.State = CheckResultStateUnknown
	} else if disabled > 0 {
		msgPrefix = "not all Nomad agents using "
		res.State = CheckResultStateFail
	} else {
		res.Msg = "all Nomad agents using "
		res.State = CheckResultStatePass
	}
	res.Msg = msgPrefix + msg

	return res
}
