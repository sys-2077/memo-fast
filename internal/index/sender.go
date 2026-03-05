package index

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// IndexRequest is the payload sent to the memo-fast memory_index MCP tool.
type IndexRequest struct {
	Files   []FilePayload   `json:"files"`
	Commits []CommitPayload `json:"commits"`
	Config  ConfigPayload   `json:"config"`
}

// FilePayload represents a single file in the index request.
type FilePayload struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// CommitPayload represents a single commit in the index request.
type CommitPayload struct {
	SHA     string   `json:"sha"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	Date    string   `json:"date"`
	Files   []string `json:"files"`
}

// ConfigPayload carries the collection name for indexing.
type ConfigPayload struct {
	Collection string `json:"collection"`
}

// IndexResponse is the response from the memory_index tool.
type IndexResponse struct {
	IndexedFiles   int      `json:"indexed_files"`
	IndexedCommits int      `json:"indexed_commits"`
	Entities       int      `json:"entities"`
	Errors         []string `json:"errors"`
}

type rpcToolCallRequest struct {
	JSONRPC string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  rpcToolCallParams `json:"params"`
	ID      int               `json:"id"`
}

type rpcToolCallParams struct {
	Name      string              `json:"name"`
	Arguments rpcToolCallArgsBody `json:"arguments"`
}

type rpcToolCallArgsBody struct {
	Files      []FilePayload   `json:"files"`
	Collection string          `json:"collection"`
	Commits    []CommitPayload `json:"commits"`
}

type rpcResponse struct {
	Error  *rpcError  `json:"error"`
	Result *rpcResult `json:"result"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type rpcResult struct {
	Content []rpcContent `json:"content"`
	// FastMCP also includes structuredContent.result for convenience.
	StructuredContent struct {
		Result string `json:"result"`
	} `json:"structuredContent"`
}

type rpcContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Send posts the index request via MCP JSON-RPC tools/call(memory_index).
func Send(apiURL, apiKey string, req IndexRequest) (*IndexResponse, error) {
	rpcReq := rpcToolCallRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params: rpcToolCallParams{
			Name: "memory_index",
			Arguments: rpcToolCallArgsBody{
				Files:      req.Files,
				Collection: req.Config.Collection,
				Commits:    req.Commits,
			},
		},
		ID: 1,
	}

	body, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	url := mcpEndpoint(apiURL)
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("sending request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var rpcResp rpcResponse
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("RPC error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}
	if rpcResp.Result == nil {
		return nil, fmt.Errorf("RPC response missing result")
	}

	payloadText := extractPayloadText(*rpcResp.Result)
	if strings.TrimSpace(payloadText) == "" {
		return nil, fmt.Errorf("RPC response missing tool payload")
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(payloadText), &payload); err != nil {
		return nil, fmt.Errorf("parsing tool payload: %w", err)
	}
	if toolErr, ok := payload["error"]; ok && strings.TrimSpace(fmt.Sprint(toolErr)) != "" {
		return nil, fmt.Errorf("memory_index failed: %v", toolErr)
	}

	var result IndexResponse
	if err := json.Unmarshal([]byte(payloadText), &result); err != nil {
		return nil, fmt.Errorf("parsing index response payload: %w", err)
	}

	return &result, nil
}

func extractPayloadText(result rpcResult) string {
	for _, item := range result.Content {
		if item.Type == "text" && strings.TrimSpace(item.Text) != "" {
			return item.Text
		}
	}
	return result.StructuredContent.Result
}

func mcpEndpoint(apiURL string) string {
	url := strings.TrimSpace(apiURL)
	url = strings.TrimRight(url, "/")
	if strings.HasSuffix(url, "/api/index") {
		return strings.TrimSuffix(url, "/api/index") + "/mcp"
	}
	if strings.HasSuffix(url, "/mcp") {
		return url
	}
	return url + "/mcp"
}
