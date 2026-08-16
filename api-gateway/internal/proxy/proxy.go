package proxy

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

type Proxy struct {
	client *http.Client
}

func New() *Proxy {
	return &Proxy{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}
func (p *Proxy) Forward(baseURL string) fiber.Handler {
	baseURL = strings.TrimRight(baseURL, "/")
	return func(c fiber.Ctx) error {
		var body io.Reader
		target := baseURL + c.OriginalURL()
		if len(c.Body()) > 0 {
			body = bytes.NewReader(c.Body())
		}
		req, err := http.NewRequestWithContext(c.Context(), c.Method(), target, body)
		if err != nil {
			return fiber.NewError(fiber.StatusBadGateway, "failed to create request")
		}
		c.Request().Header.VisitAll(func(key, value []byte) {
			req.Header.Set(string(key), string(value))
		})
		requestID := requestid.FromContext(c)
		if requestID != "" {
			req.Header.Set("X-Request-ID", requestID)
		}
		resp, err := p.client.Do(req)
		if err != nil {
			return fiber.NewError(fiber.StatusBadGateway, "failed to send request")
		}
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fiber.NewError(fiber.StatusBadGateway, "failed to read response body")
		}
		for k, vv := range resp.Header {
			for _, v := range vv {
				c.Set(k, v)
			}
		}
		return c.Status(resp.StatusCode).SendString(string(respBody))
	}
}
