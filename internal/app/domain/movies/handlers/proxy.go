package handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProxyHandler struct {
	client  *http.Client
	baseURL string
}

func NewProxyHandler(client *http.Client, baseURL string) *ProxyHandler {
	return &ProxyHandler{
		client:  client,
		baseURL: baseURL,
	}
}

func (h *ProxyHandler) Handle(c *gin.Context) {

	targetURL := h.baseURL + c.Request.RequestURI

	req, err := http.NewRequest(
		c.Request.Method,
		targetURL,
		c.Request.Body,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Header = c.Request.Header.Clone()

	resp, err := h.client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	c.Status(resp.StatusCode)

	_, _ = io.Copy(c.Writer, resp.Body)
}
