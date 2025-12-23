package acceptance

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)
type apiDriver interface {
	call(method string, path string, payload any, headers map[string]string) (apiResponse, error)
}
type apiDriverImpl struct {
	client  *http.Client
	baseURL string
}

type apiResponse struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

func (d apiDriverImpl) call(method string, path string, payload any, headers map[string]string) (apiResponse, error) {
	base, err := requireBaseURL(d.baseURL)
	if err != nil {
		return apiResponse{}, err
	}

	json, err := parseToJSON(payload)
	if err != nil {
		return apiResponse{}, err
	}

	fullURL := buildFullURL(base, path)

	client := d.client
	if client == nil {
		client = http.DefaultClient
	}

	req, err := http.NewRequest(method, fullURL, json)
	if err != nil {
		return apiResponse{}, err
	}

	applyPayloadHeaders(req, payload)
	applyHeaders(req, headers)

	resp, err := client.Do(req)
	if err != nil {
		return apiResponse{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return apiResponse{}, err
	}

	return apiResponse{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Header:     resp.Header,
	}, nil
}

func requireBaseURL(baseURL string) (string, error) {
	if baseURL == "" {
		return "", errors.New("api driver base URL is empty")
	}
	return strings.TrimRight(baseURL, "/"), nil
}

func parseToJSON(payload any)  (io.Reader, error) {
	if payload == nil {
		return nil, nil
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}

func buildFullURL(base, route string) string {
	if route == "" {
		route = "/"
	}
	if !strings.HasPrefix(route, "/") {
		route = "/" + route
	}
	return base + route
}

func applyPayloadHeaders(req *http.Request, payload any) {
	if payload == nil {
		return
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
}

func applyHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}
