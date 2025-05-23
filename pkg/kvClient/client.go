package kvClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Amirali-Amirifar/kv/internal/types/api"
	"io"
	"net/http"
)

// Client configuration
type Client struct {
	BaseURL string
	HTTP    *http.Client
}

// NewClient creates a new KV database client
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTP:    &http.Client{},
	}
}

func (c *Client) Connect() (string, error) {
	val, err := c.HTTP.Post(c.BaseURL+"/health", "application/json", nil)
	if err != nil {
		return "", fmt.Errorf("failed to connect to the host %v", err)
	}
	if val.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to connect to the host %v, status code %v", val.StatusCode, val.Status)
	}
	return fmt.Sprintf("Successfully connected to host %s\n", c.BaseURL), nil

}

// Set a new key
func (c *Client) Set(key, value string) error {
	req := api.SetRequest{Key: key, Value: value}
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := c.HTTP.Post(c.BaseURL+"/set", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server error (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// Get the value of a key
func (c *Client) Get(key string) (string, error) {
	req := api.GetRequest{Key: key}
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := c.HTTP.Post(c.BaseURL+"/get", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("server error (%d): %s", resp.StatusCode, string(body))
	}

	var response api.GetResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	return response.Value, nil
}

// Del delete a key
func (c *Client) Del(key string) error {
	req := api.DelRequest{Key: key}
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := c.HTTP.Post(c.BaseURL+"/del", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server error (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}
