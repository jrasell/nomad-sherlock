package api

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

type Client struct {
	addr       string
	httpClient *http.Client
}

type QueryOptions struct {
	Params map[string]string
}

func NewClient(addr string) (*Client, error) {

	if _, err := url.Parse(addr); err != nil {
		return nil, fmt.Errorf("invalid address '%s': %v", addr, err)
	}

	client := &Client{
		addr:       addr,
		httpClient: defaultHttpClient(),
	}
	return client, nil
}

func defaultHttpClient() *http.Client {
	httpClient := cleanhttp.DefaultClient()
	transport := httpClient.Transport.(*http.Transport)
	transport.TLSHandshakeTimeout = 10 * time.Second
	transport.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	return httpClient
}

func (c *Client) query(endpoint string, out interface{}, q *QueryOptions) error {
	r, err := c.newRequest(http.MethodGet, endpoint)
	if err != nil {
		return err
	}
	r.setQueryOptions(q)
	_, resp, err := requireOK(c.doRequest(r))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := decodeBody(resp, out); err != nil {
		return err
	}
	return nil
}

func (c *Client) write(endpoint string, in, out interface{}) error {
	r, err := c.newRequest("PUT", endpoint)
	if err != nil {
		return err
	}
	r.obj = in
	_, resp, err := requireOK(c.doRequest(r))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if out != nil {
		if err := decodeBody(resp, &out); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) doRequest(r *request) (time.Duration, *http.Response, error) {
	req, err := r.toHTTP()
	if err != nil {
		return 0, nil, err
	}
	start := time.Now()
	resp, err := c.httpClient.Do(req)
	diff := time.Now().Sub(start)

	// If the response is compressed, we swap the body's reader.
	if resp != nil && resp.Header != nil {
		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			greader, err := gzip.NewReader(resp.Body)
			if err != nil {
				return 0, nil, err
			}

			// The gzip reader doesn't close the wrapped reader so we use
			// multiCloser.
			reader = &multiCloser{
				reader:       greader,
				inorderClose: []io.Closer{greader, resp.Body},
			}
		default:
			reader = resp.Body
		}
		resp.Body = reader
	}

	return diff, resp, err
}

func (c *Client) newRequest(method, path string) (*request, error) {
	base, _ := url.Parse(c.addr)
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	r := &request{
		addr:   c.addr,
		method: method,
		url: &url.URL{
			Scheme:  base.Scheme,
			User:    base.User,
			Host:    base.Host,
			Path:    u.Path,
			RawPath: u.RawPath,
		},
		params: make(map[string][]string),
	}

	// Add in the query parameters, if any
	for key, values := range u.Query() {
		for _, value := range values {
			r.params.Add(key, value)
		}
	}
	return r, nil
}

type multiCloser struct {
	reader       io.Reader
	inorderClose []io.Closer
}

func (m *multiCloser) Close() error {
	for _, c := range m.inorderClose {
		if err := c.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (m *multiCloser) Read(p []byte) (int, error) {
	return m.reader.Read(p)
}

type request struct {
	addr   string
	method string
	url    *url.URL
	params url.Values
	body   io.Reader
	obj    interface{}
}

func (r *request) setQueryOptions(q *QueryOptions) {
	if q == nil {
		return
	}
	for k, v := range q.Params {
		r.params.Set(k, v)
	}
}

func (r *request) toHTTP() (*http.Request, error) {
	// Encode the query parameters
	r.url.RawQuery = r.params.Encode()

	// Check if we should encode the body
	if r.body == nil && r.obj != nil {
		if b, err := encodeBody(r.obj); err != nil {
			return nil, err
		} else {
			r.body = b
		}
	}

	// Create the HTTP request
	req, err := http.NewRequest(r.method, r.url.RequestURI(), r.body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept-Encoding", "gzip")

	req.URL.Host = r.url.Host
	req.URL.Scheme = r.url.Scheme
	req.Host = r.url.Host
	return req, nil
}

func requireOK(d time.Duration, resp *http.Response, e error) (time.Duration, *http.Response, error) {
	if e != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return d, nil, e
	}
	if resp.StatusCode != 200 {
		var buf bytes.Buffer
		io.Copy(&buf, resp.Body)
		resp.Body.Close()
		return d, nil, fmt.Errorf("unexpected response code: %d (%s)", resp.StatusCode, buf.Bytes())
	}
	return d, resp, nil
}

func encodeBody(obj interface{}) (io.Reader, error) {
	if reader, ok := obj.(io.Reader); ok {
		return reader, nil
	}

	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(obj); err != nil {
		return nil, err
	}
	return buf, nil
}

func decodeBody(resp *http.Response, out interface{}) error {
	switch resp.ContentLength {
	case 0:
		if out == nil {
			return nil
		}
		return errors.New("got 0 byte response with non-nil decode object")
	default:
		dec := json.NewDecoder(resp.Body)
		return dec.Decode(out)
	}
}
