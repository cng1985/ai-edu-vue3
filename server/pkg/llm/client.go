package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	apiKey   string
	baseURL  string
	model    string
	authType string
	enabled  bool
	client   *http.Client
}

func NewClient(apiKey, baseURL, model string) *Client {
	return NewClientWithAuth(apiKey, baseURL, model, "Bearer")
}

// NewClientWithAuth 创建 OpenAI 兼容客户端，并支持常见厂商认证头。
func NewClientWithAuth(apiKey, baseURL, model, authType string) *Client {
	return &Client{
		apiKey:   apiKey,
		baseURL:  strings.TrimRight(baseURL, "/"),
		model:    model,
		authType: strings.TrimSpace(authType),
		enabled:  apiKey != "",
		client:   &http.Client{Timeout: 120 * time.Second},
	}
}

func (c *Client) Enabled() bool { return c.enabled }

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Message 对话消息（对外暴露）
type Message = chatMessage

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type streamDelta struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

// StreamMessages 多轮对话流式输出
func (c *Client) StreamMessages(ctx context.Context, messages []chatMessage, onToken func(string)) (string, error) {
	if !c.enabled {
		return "", fmt.Errorf("LLM 未配置")
	}
	body, _ := json.Marshal(chatRequest{
		Model:    c.model,
		Messages: messages,
		Stream:   true,
	})
	return c.doStream(ctx, body, onToken)
}

// Complete 非流式补全，用于结构化 JSON 输出
func (c *Client) Complete(ctx context.Context, messages []chatMessage) (string, error) {
	if !c.enabled {
		return "", fmt.Errorf("LLM 未配置")
	}
	body, _ := json.Marshal(chatRequest{
		Model:    c.model,
		Messages: messages,
		Stream:   false,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	c.applyAuth(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("LLM API 错误 %d: %s", resp.StatusCode, string(b))
	}
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("LLM 返回为空")
	}
	return result.Choices[0].Message.Content, nil
}

// StreamChat 调用 OpenAI 兼容接口，通过 onToken 回调流式输出
func (c *Client) StreamChat(ctx context.Context, systemPrompt, userPrompt string, onToken func(string)) (string, error) {
	return c.StreamMessages(ctx, []chatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}, onToken)
}

func (c *Client) doStream(ctx context.Context, body []byte, onToken func(string)) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	c.applyAuth(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("LLM API 错误 %d: %s", resp.StatusCode, string(b))
	}

	var full strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}
		var chunk streamDelta
		if json.Unmarshal([]byte(data), &chunk) != nil || len(chunk.Choices) == 0 {
			continue
		}
		token := chunk.Choices[0].Delta.Content
		if token == "" {
			continue
		}
		full.WriteString(token)
		if onToken != nil {
			onToken(token)
		}
	}
	return full.String(), scanner.Err()
}

func (c *Client) applyAuth(req *http.Request) {
	switch strings.ToLower(c.authType) {
	case "api-key", "apikey":
		req.Header.Set("api-key", c.apiKey)
	case "x-api-key":
		req.Header.Set("x-api-key", c.apiKey)
	default:
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
}
