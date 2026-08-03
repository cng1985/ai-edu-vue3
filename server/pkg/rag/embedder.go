package rag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Embedder 文本嵌入接口
type Embedder interface {
	Enabled() bool
	Model() string
	Dimensions() int
	Embed(ctx context.Context, texts []string) ([][]float32, error)
	EmbedOne(ctx context.Context, text string) ([]float32, error)
}

// APIEmbedder OpenAI 兼容嵌入 API
type APIEmbedder struct {
	apiKey     string
	baseURL    string
	model      string
	dimensions int
	client     *http.Client
}

func NewAPIEmbedder(apiKey, baseURL, model string, dimensions int) *APIEmbedder {
	if dimensions <= 0 {
		dimensions = 1536
	}
	if model == "" {
		model = "text-embedding-3-small"
	}
	return &APIEmbedder{
		apiKey:     apiKey,
		baseURL:    strings.TrimRight(baseURL, "/"),
		model:      model,
		dimensions: dimensions,
		client:     &http.Client{Timeout: 60 * time.Second},
	}
}

func (e *APIEmbedder) Enabled() bool     { return e.apiKey != "" }
func (e *APIEmbedder) Model() string     { return e.model }
func (e *APIEmbedder) Dimensions() int   { return e.dimensions }

func (e *APIEmbedder) EmbedOne(ctx context.Context, text string) ([]float32, error) {
	vecs, err := e.Embed(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	if len(vecs) == 0 {
		return nil, fmt.Errorf("嵌入结果为空")
	}
	return vecs[0], nil
}

func (e *APIEmbedder) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	if !e.Enabled() {
		return nil, fmt.Errorf("嵌入 API 未配置")
	}
	body, _ := json.Marshal(map[string]interface{}{
		"model": e.model,
		"input": texts,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.baseURL+"/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+e.apiKey)

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("嵌入 API 错误 %d: %s", resp.StatusCode, string(b))
	}
	var result struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	out := make([][]float32, len(result.Data))
	for i, item := range result.Data {
		out[i] = item.Embedding
	}
	return out, nil
}

// LocalEmbedder 本地哈希嵌入（无需外部 API）
type LocalEmbedder struct{}

func NewLocalEmbedder() *LocalEmbedder { return &LocalEmbedder{} }

func (e *LocalEmbedder) Enabled() bool   { return true }
func (e *LocalEmbedder) Model() string   { return "local-hash-v1" }
func (e *LocalEmbedder) Dimensions() int { return LocalEmbeddingDim }

func (e *LocalEmbedder) EmbedOne(ctx context.Context, text string) ([]float32, error) {
	return LocalEmbed(text), nil
}

func (e *LocalEmbedder) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	out := make([][]float32, len(texts))
	for i, t := range texts {
		out[i] = LocalEmbed(t)
	}
	return out, nil
}
