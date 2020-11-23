package server

import (
	"net/http"
	"strings"
)

func (h *HTTPServer) checksReq(_ http.ResponseWriter, req *http.Request) (interface{}, error) {
	if req.Method != http.MethodGet {
		return nil, newCodedError(http.StatusMethodNotAllowed, errInvalidMethod)
	}
	out := h.agent.Plugin().GetChecks()
	return out, nil
}

func (h *HTTPServer) checkReq(_ http.ResponseWriter, req *http.Request) (interface{}, error) {
	if req.Method != http.MethodGet {
		return nil, newCodedError(http.StatusMethodNotAllowed, errInvalidMethod)
	}

	check := strings.TrimPrefix(req.URL.Path, "/v1/check/")
	split := strings.Split(check, "_")

	if len(split) != 2 {
		return nil, newCodedError(http.StatusBadRequest, "Invalid check format")
	}

	if found := h.agent.Plugin().GetCheck(check); found != nil {
		return found, nil
	}
	return nil, newCodedError(http.StatusNotFound, "Check not found")
}
