package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type HTTPUserClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewHTTPClient(baseURL string, apiKey string, timeout time.Duration) *HTTPUserClient {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &HTTPUserClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}
func (c *HTTPUserClient) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("create user req", err)
	}
	url := fmt.Sprintf("%s/users/", c.baseURL)
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create user HTTP req", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpRequest.Header.Set("X-API-KEY", c.apiKey)
	}
	log.Printf("USER CLIENT CREATE URL: %s", url)
	resp, err := c.httpClient.Do(httpRequest)
	if err != nil {
		log.Printf("USER CLIENT CREATE URL %s failed: %v", url, err)
		return nil, ErrUserService
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body", err)
	}
	switch resp.StatusCode {
	case http.StatusCreated:
		var created User
		if err := json.Unmarshal(respBody, &created); err != nil {
			return nil, fmt.Errorf("unmarshal response body: %v", ErrInvalidResponse, err)
		}
		if created.ID == uuid.Nil {
			return nil, ErrInvalidResponse
		}
		return &created, nil
	case http.StatusConflict:
		return nil, ErrUserAlreadyExists
	default:
		return nil, fmt.Errorf("status=%d body=%s", ErrUserService, resp.StatusCode, string(respBody))
	}

}
func (c *HTTPUserClient) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if userID == uuid.Nil {
		return ErrUserNotFound
	}
	url := fmt.Sprintf("%s/api/users/%s", c.baseURL, userID.String())
	deleteRequest, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("create delete user HTTP request", err)
	}
	if c.apiKey != "" {
		deleteRequest.Header.Set("X-API-KEY", c.apiKey)
	}
	resp, err := c.httpClient.Do(deleteRequest)
	if err != nil {
		return ErrUserService
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return ErrUserNotFound
	default:
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete user HTTP request: status=%d body=%s", resp.StatusCode, string(body))
	}
}
func (c *HTTPUserClient) GetuserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	if id == uuid.Nil {
		return nil, ErrUserNotFound
	}
	url := fmt.Sprintf("%s/api/users/%s", c.baseURL, id.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create get user HTTP request: %w", err)
	}
	if c.apiKey != "" {
		req.Header.Set("X-API-KEY", c.apiKey)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, ErrUserService
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body", err)
	}
	log.Printf(
		"USER CLIENT: url=%s status=%d body=%s",
		url,
		resp.StatusCode,
		string(body),
	)
	switch resp.StatusCode {
	case http.StatusOK:
		var user User
		if err := json.Unmarshal(body, &user); err != nil {
			return nil, fmt.Errorf("unmarshal response body: %v", ErrInvalidResponse, err)
		}
		if user.ID == uuid.Nil {
			return nil, ErrInvalidResponse
		}
		return &user, nil
	case http.StatusNotFound:
		return nil, ErrUserNotFound
	default:
		return nil, fmt.Errorf("status=%d body=%s", ErrUserService, resp.StatusCode, string(body))
	}

}
