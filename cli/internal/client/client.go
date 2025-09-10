package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

// Client API 客户端
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	userAgent  string
}

// NewClient 创建新的 API 客户端
func NewClient(baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	u, _ := url.Parse(baseURL)

	return &Client{
		baseURL:    u,
		httpClient: httpClient,
		userAgent:  "Prism-CLI/1.0",
	}
}

// Get 发送 GET 请求
func (c *Client) Get(endpoint string) (*http.Response, error) {
	return c.do("GET", endpoint, nil)
}

// Post 发送 POST 请求
func (c *Client) Post(endpoint string, body interface{}) (*http.Response, error) {
	return c.do("POST", endpoint, body)
}

// Put 发送 PUT 请求
func (c *Client) Put(endpoint string, body interface{}) (*http.Response, error) {
	return c.do("PUT", endpoint, body)
}

// Delete 发送 DELETE 请求
func (c *Client) Delete(endpoint string) (*http.Response, error) {
	return c.do("DELETE", endpoint, nil)
}

// do 执行 HTTP 请求
func (c *Client) do(method, endpoint string, body interface{}) (*http.Response, error) {
	// 构建 URL
	u := *c.baseURL
	u.Path = path.Join(u.Path, endpoint)

	// 准备请求体
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	// 创建请求
	req, err := http.NewRequest(method, u.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", c.userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	return resp, nil
}

// GetJSON 发送 GET 请求并解析 JSON 响应
func (c *Client) GetJSON(endpoint string, result interface{}) error {
	resp, err := c.Get(endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

// PostJSON 发送 POST 请求并解析 JSON 响应
func (c *Client) PostJSON(endpoint string, body interface{}, result interface{}) error {
	resp, err := c.Post(endpoint, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %s", resp.Status)
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}

	return nil
}