package api

import (
	"errors"

	"github.com/jrasell/nomad-sherlock/sdk"
)

type Inspection struct {
	client *Client
}

func (c *Client) Inspection() *Inspection {
	return &Inspection{client: c}
}

func (i *Inspection) List(q *QueryOptions) ([]*sdk.Inspection, error) {
	var resp []*sdk.Inspection
	err := i.client.query("/v1/inspections", &resp, q)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *Inspection) Info(id string, q *QueryOptions) (*sdk.Inspection, error) {
	var resp *sdk.Inspection
	err := i.client.query("/v1/inspection/"+id, &resp, q)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *Inspection) Submit(report *sdk.Inspection) (*sdk.InspectionSubmission, error) {
	var resp *sdk.InspectionSubmission
	if report == nil || report.Region == "" {
		return nil, errors.New("missing inspection region identifier")
	}
	if err :=  i.client.write("/v1/inspections", report, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
