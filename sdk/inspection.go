package sdk

import (
	"github.com/hashicorp/go-multierror"
	"time"
)

type Inspection struct {
	ID     string `json:"id"`
	Region string `json:"region"`

	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Duration  string    `json:"duration"`

	Warnings []string `json:"warnings"`

	TotalRun     int `json:"totalRun"`
	TotalFail    int `json:"totalFail"`
	TotalUnknown int `json:"totalUnknown"`
	TotalPass    int `json:"totalPass"`

	Results []RunCheckResult `json:"results"`
}

type InspectionSubmission struct {
	ID string `json:"id"`
}

type RunCheckResult struct {
	ID     string      `json:"id"`
	Result CheckResult `json:"result"`
}

func NewInspection() *Inspection {
	return &Inspection{
		Region:    "global",
		StartTime: time.Now(),
	}
}

func (r *Inspection) AddWarning(warn string) {
	r.Warnings = append(r.Warnings, warn)
}

func (r *Inspection) Update(res []CheckResult, err error) {

	if err != nil {
		if mErr, ok := err.(*multierror.Error); ok {
			for _, unwrappedErr := range mErr.Errors {
				r.Warnings = append(r.Warnings, unwrappedErr.Error())
			}
		} else {
			r.AddWarning(err.Error())
		}
	}

	for _, result := range res {
		r.TotalRun++

		switch result.State {
		case CheckResultStatePass:
			r.TotalPass++
		case CheckResultStateFail:
			r.TotalFail++
		case CheckResultStateUnknown:
			r.TotalUnknown++
		}

		rr := RunCheckResult{
			ID:     result.ID,
			Result: result,
		}
		r.Results = append(r.Results, rr)
	}
}

func (r *Inspection) Finalize() {
	r.EndTime = time.Now()
	r.Duration = r.EndTime.Sub(r.StartTime).String()
}
