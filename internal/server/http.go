package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/hashicorp/go-msgpack/codec"
)

func decodeBody(req *http.Request, out interface{}) error {

	if req.Body == http.NoBody {
		return errors.New("request body is empty")
	}

	dec := json.NewDecoder(req.Body)
	return dec.Decode(&out)
}

func (h *HTTPServer) wrap(handler func(w http.ResponseWriter, r *http.Request) (interface{}, error)) func(w http.ResponseWriter, r *http.Request) {
	f := func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		// Defer a function which allows us to log the time taken to fulfill
		// the HTTP request.
		defer func() {
			h.log.Info("request complete", "method", r.Method,
				"path", r.URL, "duration", time.Since(start))
		}()

		// Handle the request, allowing us to the get response object and any
		// error from the endpoint.
		obj, err := handler(w, r)
		if err != nil {
			h.handleHTTPError(w, r, err)
			return
		}

		// If we have a response object, encode it.
		if obj != nil {
			var buf bytes.Buffer

			enc := codec.NewEncoder(&buf, &codec.JsonHandle{HTMLCharsAsIs: true})

			// Encode the object. If we fail to do this, handle the error so
			// that this can be passed to the operator.
			err := enc.Encode(obj)
			if err != nil {
				h.handleHTTPError(w, r, err)
				return
			}

			//  Set the content type header and write the data to the HTTP
			//  reply.
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(buf.Bytes())
		}
	}

	return f
}

// handleHTTPError is used to handle HTTP handler errors within the wrap func.
// It sets response headers where required and ensure appropriate errors are
// logged.
func (h *HTTPServer) handleHTTPError(w http.ResponseWriter, r *http.Request, err error) {

	// Start with a default internal server error and the error message
	// that was returned.
	code := http.StatusInternalServerError
	errMsg := err.Error()

	// If the error was a custom codedError update the response code to
	// that of the wrapped error.
	if codedErr, ok := err.(codedError); ok {
		code = codedErr.Code()
	}

	// Write the status code header.
	w.WriteHeader(code)

	// Write the response body. If we get an error, log this as it will
	// provide some operator insight if this happens regularly.
	if _, wErr := w.Write([]byte(errMsg)); wErr != nil {
		h.log.Error("failed to write response error", "error", wErr)
	}
	h.log.Error("request failed", "method", r.Method, "path", r.URL, "error", errMsg, "code", code)
}
