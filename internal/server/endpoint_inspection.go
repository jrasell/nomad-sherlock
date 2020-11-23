package server

import (
	"net/http"
	"strings"

	"github.com/hashicorp/nomad/helper/uuid"
	"github.com/jrasell/nomad-sherlock/sdk"
)

func (h *HTTPServer) inspectionReq(_ http.ResponseWriter, req *http.Request) (interface{}, error) {
	if req.Method != http.MethodGet {
		return nil, newCodedError(http.StatusMethodNotAllowed, errInvalidMethod)
	}

	id := strings.TrimPrefix(req.URL.Path, "/v1/inspection/")

	h.lock.Lock()
	defer h.lock.Unlock()

	inspection, ok := h.reports[id]
	if !ok {
		return nil, newCodedError(http.StatusNotFound, "Inspection not found")
	}
	return inspection, nil
}

func (h *HTTPServer) inspectionsReq(res http.ResponseWriter, req *http.Request) (interface{}, error) {
	switch req.Method {
	case http.MethodGet:
		return h.inspectionsReqGet(res, req)
	case http.MethodPost, http.MethodPut:
		return h.inspectionsReqPut(res, req)
	default:
		return nil, newCodedError(http.StatusMethodNotAllowed, errInvalidMethod)
	}
}

func (h *HTTPServer) inspectionsReqGet(_ http.ResponseWriter, req *http.Request) (interface{}, error) {
	region := req.URL.Query().Get("region")

	h.lock.Lock()
	defer h.lock.Unlock()

	out := []*sdk.Inspection{}

	if region != "" {
		stored := h.regional[region]

		for _, id := range stored {
			if r := h.reports[id]; r != nil {
				out = append(out, r)
			}
		}
	} else {
		for _, stored := range h.reports {
			out = append(out, stored)
		}
	}
	return out, nil
}

func (h *HTTPServer) inspectionsReqPut(_ http.ResponseWriter, req *http.Request) (interface{}, error) {

	var report sdk.Inspection
	if err := decodeBody(req, &report); err != nil {
		return nil, newCodedError(400, err.Error())
	}

	if report.Region == "" {
		return nil, newCodedError(400, "region identifier must be provided")
	}

	h.lock.Lock()
	defer h.lock.Unlock()

	id := uuid.Generate()
	report.ID = id

	h.regional[report.Region] = append(h.regional[report.Region], id)
	h.reports[id] = &report
	return &sdk.InspectionSubmission{ID: id}, nil
}
