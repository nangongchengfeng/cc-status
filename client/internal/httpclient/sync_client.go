package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Report struct {
	RequestID           string `json:"request_id"`
	AppType             string `json:"app_type"`
	Model               string `json:"model"`
	InputTokens         uint32 `json:"input_tokens"`
	OutputTokens        uint32 `json:"output_tokens"`
	CacheReadTokens     uint32 `json:"cache_read_tokens"`
	CacheCreationTokens uint32 `json:"cache_creation_tokens"`
	CreatedAt           int64  `json:"created_at"`
	SessionID           string `json:"session_id"`
	DataSource          string `json:"data_source"`
}

type SyncResponse struct {
	Code           int    `json:"code"`
	Message        string `json:"message"`
	AcceptedCount  int    `json:"accepted_count"`
	DuplicateCount int    `json:"duplicate_count"`
}

type SyncClient struct {
	baseURL    string
	authToken  string
	httpClient *http.Client
}

func NewSyncClient(baseURL string, authToken string, timeoutSeconds int) *SyncClient {
	timeout := time.Duration(timeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	return &SyncClient{
		baseURL:   strings.TrimRight(baseURL, "/"),
		authToken: authToken,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (client *SyncClient) Sync(clientID string, reports []Report) (SyncResponse, error) {
	payload := struct {
		ClientID string   `json:"client_id"`
		Reports  []Report `json:"reports"`
	}{
		ClientID: clientID,
		Reports:  reports,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return SyncResponse{}, fmt.Errorf("marshal sync payload: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, client.baseURL+"/api/v1/sync", bytes.NewReader(body))
	if err != nil {
		return SyncResponse{}, fmt.Errorf("build sync request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	if client.authToken != "" {
		request.Header.Set("Authorization", "Bearer "+client.authToken)
	}

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		response, err := client.httpClient.Do(request.Clone(request.Context()))
		if err != nil {
			lastErr = fmt.Errorf("send sync request: %w", err)
			// 只对网络层瞬时错误做指数退避，业务失败交由上层按整批失败处理。
			if attempt < 2 {
				time.Sleep(time.Duration(1<<attempt) * 10 * time.Millisecond)
				continue
			}
			return SyncResponse{}, lastErr
		}

		body, readErr := io.ReadAll(response.Body)
		_ = response.Body.Close()
		if readErr != nil {
			return SyncResponse{}, fmt.Errorf("read sync response: %w", readErr)
		}

		if response.StatusCode >= 500 {
			lastErr = fmt.Errorf("unexpected status code: %d", response.StatusCode)
			if attempt < 2 {
				time.Sleep(time.Duration(1<<attempt) * 10 * time.Millisecond)
				continue
			}
			return SyncResponse{}, lastErr
		}

		if response.StatusCode < 200 || response.StatusCode >= 300 {
			return SyncResponse{}, fmt.Errorf("unexpected status code: %d", response.StatusCode)
		}

		var decoded SyncResponse
		if err := json.Unmarshal(body, &decoded); err != nil {
			return SyncResponse{}, fmt.Errorf("decode sync response: %w", err)
		}

		return decoded, nil
	}

	return SyncResponse{}, lastErr
}
