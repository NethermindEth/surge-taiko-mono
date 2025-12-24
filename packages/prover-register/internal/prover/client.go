package prover

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	log        *zap.SugaredLogger
}

type GuestData map[string]interface{}

func NewClient(baseURL string, log *zap.SugaredLogger) *Client {
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

// GetGuestData fetches guest data from the original REST endpoint.
func (c *Client) GetGuestData(ctx context.Context) (GuestData, error) {
	url := fmt.Sprintf("%s/guest_data", c.baseURL)
	c.log.Debugw("fetching guest data", "url", url)

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

	c.log.Debugw("successfully fetched guest data", "keys", getKeys(data))
	return data, nil
}

// GetGuestDataFromNethermind fetches guest data from Nethermind's JSON-RPC endpoint.
func (c *Client) GetGuestDataFromNethermind(ctx context.Context) (GuestData, error) {
	c.log.Debugw("fetching guest data from Nethermind via JSON-RPC", "url", c.baseURL)

	rpcRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "taiko_getTdxGuestInfo",
		"params":  []interface{}{},
		"id":      1,
	}

	requestBody, err := json.Marshal(rpcRequest)
	if err != nil {
		return nil, fmt.Errorf("marshal JSON-RPC request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var rpcResponse struct {
		JSONRPC string    `json:"jsonrpc"`
		ID      int       `json:"id"`
		Result  GuestData `json:"result"`
		Error   *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&rpcResponse); err != nil {
		return nil, fmt.Errorf("decode JSON-RPC response: %w", err)
	}

	if rpcResponse.Error != nil {
		return nil, fmt.Errorf("JSON-RPC error %d: %s", rpcResponse.Error.Code, rpcResponse.Error.Message)
	}

	if rpcResponse.Result == nil {
		return nil, fmt.Errorf("received null result from JSON-RPC")
	}

	c.log.Debugw("successfully fetched guest data from Nethermind", "keys", getKeys(rpcResponse.Result))
	return rpcResponse.Result, nil
}

func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
