package microservices

import (
	"context"
	"encoding/json"
	"fmt"
	"gateway/internal/app/core/contracts/headercontract"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Response struct {
	Data       interface{}
	StatusCode int
	Error      error
}

type RequestOption struct {
	Body   io.Reader
	Query  url.Values
	Header map[string]string
}

type RequestHandler struct {
	client Client
}

func NewRequestHandler(client Client) *RequestHandler {
	return &RequestHandler{client: client}
}

func (h *RequestHandler) Get(ctx context.Context, url string, option RequestOption) Response {
	return h.request(ctx, http.MethodGet, url, option)
}

func (h *RequestHandler) Post(ctx context.Context, url string, option RequestOption) Response {
	return h.request(ctx, http.MethodPost, url, option)
}

func (h *RequestHandler) Put(ctx context.Context, url string, option RequestOption) Response {
	return h.request(ctx, http.MethodPut, url, option)
}

func (h *RequestHandler) Patch(ctx context.Context, url string, option RequestOption) Response {
	return h.request(ctx, http.MethodPatch, url, option)
}

func (h *RequestHandler) Delete(ctx context.Context, url string, option RequestOption) Response {
	return h.request(ctx, http.MethodDelete, url, option)
}

func (h *RequestHandler) request(ctx context.Context, method, url string, option RequestOption) Response {
	req, err := http.NewRequestWithContext(ctx, method, h.client.FullPath(url), option.Body)
	if err != nil {
		return Response{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Errorf("failed to create request: %w", err),
		}
	}

	if option.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	authUser, ok := ctx.Value(headercontract.AuthUserKey{}).(headercontract.AuthUser)
	if ok {
		req.Header.Set(headercontract.UserIdKey, strconv.FormatInt(authUser.ID, 10))
	}

	for k, v := range option.Header {
		req.Header.Set(k, v)
	}

	req.URL.RawQuery = option.Query.Encode()

	resp, err := h.client.Do(ctx, req)
	if err != nil {
		return Response{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Errorf("failed to send request: %w", err),
		}
	}
	defer resp.Body.Close()

	var response interface{}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return Response{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Errorf("failed to decode response: %w, status: %s", err, resp.Status),
		}
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		return Response{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Errorf("failed to decode response, status: %s", resp.Status),
		}
	}

	return Response{
		Data:       response,
		StatusCode: resp.StatusCode,
		Error:      nil,
	}
}
