package server

import (
	"net/http"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/jrasell/nomad-sherlock/internal/agent"
	"github.com/jrasell/nomad-sherlock/sdk"
)

type HTTPServer struct {
	agent *agent.Agent
	log   hclog.Logger

	mux *http.ServeMux
	srv *http.Server

	lock     sync.RWMutex
	regional map[string][]string
	reports  map[string]*sdk.Inspection
}

func NewHTTPServer(log hclog.Logger, agent *agent.Agent) *HTTPServer {

	s := HTTPServer{
		agent:    agent,
		log:      log.Named("http_server"),
		mux:      http.NewServeMux(),
		regional: make(map[string][]string),
		reports:  make(map[string]*sdk.Inspection),
	}

	s.mux.HandleFunc("/v1/inspections", s.wrap(s.inspectionsReq))
	s.mux.HandleFunc("/v1/inspection/", s.wrap(s.inspectionReq))
	s.mux.HandleFunc("/v1/checks", s.wrap(s.checksReq))
	s.mux.HandleFunc("/v1/check/", s.wrap(s.checkReq))

	s.srv = &http.Server{
		Addr:    ":1313",
		Handler: s.mux,
	}
	return &s
}

func (h *HTTPServer) Start() {
	h.log.Info("HTTP server listening", "addr", h.srv.Addr)
	if err := h.srv.ListenAndServe(); err != nil {
		h.log.Info("shutting down HTTP server", "error", err)
	}
}
