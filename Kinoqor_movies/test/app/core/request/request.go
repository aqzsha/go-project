package request

import (
	headercontract "movies/internal/app/core/contracts/microservices/header-contract"
	"movies/test/api"
	"movies/test/app/core/header"
	"net/http"
	"net/url"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

type Request struct {
	T      *testing.T
	Server *api.TestServer
}

func NewRequest(t *testing.T) *Request {
	t.Helper()

	server := api.NewTestServer(t)

	return &Request{
		T:      t,
		Server: server,
	}
}

func (h *Request) Get(path string, status int, body *map[string]any, user headercontract.AuthUser) *httpexpect.Object {
	return h.request(path, http.MethodGet, status, body, user)
}

func (h *Request) Post(path string, status int, body *map[string]any, user headercontract.AuthUser) *httpexpect.Object {
	return h.request(path, http.MethodPost, status, body, user)
}

func (h *Request) Put(path string, status int, body *map[string]any, user headercontract.AuthUser) *httpexpect.Object {
	return h.request(path, http.MethodPut, status, body, user)
}

func (h *Request) Patch(path string, status int, body *map[string]any, user headercontract.AuthUser) *httpexpect.Object {
	return h.request(path, http.MethodPatch, status, body, user)
}

func (h *Request) Delete(path string, status int, body *map[string]any, user headercontract.AuthUser) *httpexpect.Object {
	return h.request(path, http.MethodDelete, status, body, user)
}

func (h *Request) request(path, method string, status int, body *map[string]any, user headercontract.AuthUser) *httpexpect.Object {
	h.T.Helper()
	e := h.Server.Expect()

	parsedPath, err := url.Parse(path)
	if err != nil {
		h.T.Fatalf("failed to parse url %q: %v", path, err)
	}

	var req *httpexpect.Request
	switch method {
	case http.MethodGet:
		if parsedPath.IsAbs() {
			req = e.GET(parsedPath.String())
		} else {
			req = e.GET(parsedPath.Path)
			for k, vals := range parsedPath.Query() {
				for _, v := range vals {
					req = req.WithQuery(k, v)
				}
			}
		}
	case http.MethodPost:
		if parsedPath.IsAbs() {
			req = e.POST(parsedPath.String())
		} else {
			req = e.POST(parsedPath.Path)
			for k, vals := range parsedPath.Query() {
				for _, v := range vals {
					req = req.WithQuery(k, v)
				}
			}
		}
	case http.MethodPut:
		if parsedPath.IsAbs() {
			req = e.PUT(parsedPath.String())
		} else {
			req = e.PUT(parsedPath.Path)
			for k, vals := range parsedPath.Query() {
				for _, v := range vals {
					req = req.WithQuery(k, v)
				}
			}
		}
	case http.MethodPatch:
		if parsedPath.IsAbs() {
			req = e.PATCH(parsedPath.String())
		} else {
			req = e.PATCH(parsedPath.Path)
			for k, vals := range parsedPath.Query() {
				for _, v := range vals {
					req = req.WithQuery(k, v)
				}
			}
		}
	case http.MethodDelete:
		if parsedPath.IsAbs() {
			req = e.DELETE(parsedPath.String())
		} else {
			req = e.DELETE(parsedPath.Path)
			for k, vals := range parsedPath.Query() {
				for _, v := range vals {
					req = req.WithQuery(k, v)
				}
			}
		}
	default:
		h.T.Fatalf("unsupported HTTP method: %s", method)
	}

	req = header.ApplyAuth(
		req,
		user,
	)

	if body != nil {
		req.WithJSON(body)
	}

	return req.Expect().
		Status(status).
		JSON().Object()
}