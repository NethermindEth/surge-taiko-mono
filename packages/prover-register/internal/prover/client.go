package prover

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/taikoxyz/taiko-mono/packages/prover-register/internal/logger"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	log        *logger.Logger
}

type GuestData map[string]interface{}

func NewClient(baseURL string, log *logger.Logger) *Client {
	// Ensure baseURL has http:// prefix
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "http://" + baseURL
	}

	return &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		log: log,
	}
}

func (c *Client) GetGuestData(ctx context.Context) (GuestData, error) {
	url := fmt.Sprintf("%s/guest_data", c.baseURL)
	c.log.Debug("fetching guest data", "url", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var data GuestData

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if err := json.Unmarshal(body, &data); err != nil {
		var dataArray []GuestData
		if arrayErr := json.Unmarshal(body, &dataArray); arrayErr != nil {
			return nil, fmt.Errorf("decode response as object or array: object error: %w, array error: %v", err, arrayErr)
		}

		if len(dataArray) == 0 {
			return nil, fmt.Errorf("received empty array")
		}

		data = dataArray[0]
	}

	c.log.Debug("successfully fetched guest data", "keys", getKeys(data))
	return data, nil
}

func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
