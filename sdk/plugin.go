package sdk

type CheckPluginProvider interface {
	SetConfig(cfg map[string]string) error
	ExecuteChecks() ([]CheckResult, error)
	GetChecks() []*CheckInfo
	GetCheck(id string) *CheckInfo
}

const (
	PluginConfigKeyDisabledChecks = "disabled_checks"
)
