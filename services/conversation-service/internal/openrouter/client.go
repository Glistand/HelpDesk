package openrouter

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

type Client struct {
	APIKey  string
	Model   string
	BaseURL string
	HTTP    *http.Client
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

const systemPrompt = `Ты — виртуальный помощник support desk.
Отвечай кратко и по делу на языке пользователя.
Если не уверен в ответе или пользователь просит человека — скажи, что можно нажать «Нужен человек», чтобы связаться с агентом.
Не выдумывай факты о компании.`

func (c *Client) Reply(ctx context.Context, history []ChatMessage) (string, error) {
	return c.ReplyWithInstructions(ctx, history, "")
}

func (c *Client) ReplyWithInstructions(ctx context.Context, history []ChatMessage, instructions string) (string, error) {
	if c.APIKey == "" {
		return "", fmt.Errorf("OPENROUTER_API_KEY is not set")
	}
	base := strings.TrimRight(c.BaseURL, "/")
	model := c.Model
	if model == "" {
		model = "inclusionai/ling-3.0-flash-vl:free"
	}
	prompt := systemPrompt
	if strings.TrimSpace(instructions) != "" {
		prompt += "\n\nКонтекст проекта:\n" + strings.TrimSpace(instructions)
	}
	msgs := make([]ChatMessage, 0, len(history)+1)
	msgs = append(msgs, ChatMessage{Role: "system", Content: prompt})
	msgs = append(msgs, history...)

	body, _ := json.Marshal(chatRequest{Model: model, Messages: msgs})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://github.com/Glistand/HelpDesk")
	req.Header.Set("X-Title", "HelpDesk Support Widget")

	httpClient := c.HTTP
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 45 * time.Second}
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode == http.StatusTooManyRequests || res.StatusCode >= 500 {
		// one fallback free model
		return c.replyWithModel(ctx, history, "openrouter/free", prompt)
	}
	if res.StatusCode >= 300 {
		return "", fmt.Errorf("openrouter %d: %s", res.StatusCode, string(raw))
	}
	var parsed chatResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", err
	}
	if parsed.Error != nil && parsed.Error.Message != "" {
		return "", fmt.Errorf("openrouter: %s", parsed.Error.Message)
	}
	if len(parsed.Choices) == 0 || strings.TrimSpace(parsed.Choices[0].Message.Content) == "" {
		return "", fmt.Errorf("empty openrouter response")
	}
	return strings.TrimSpace(parsed.Choices[0].Message.Content), nil
}

func (c *Client) replyWithModel(ctx context.Context, history []ChatMessage, model, prompt string) (string, error) {
	orig := c.Model
	c.Model = model
	defer func() { c.Model = orig }()
	base := strings.TrimRight(c.BaseURL, "/")
	msgs := append([]ChatMessage{{Role: "system", Content: prompt}}, history...)
	body, _ := json.Marshal(chatRequest{Model: model, Messages: msgs})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	httpClient := c.HTTP
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 45 * time.Second}
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 300 {
		return "", fmt.Errorf("openrouter fallback %d: %s", res.StatusCode, string(raw))
	}
	var parsed chatResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", err
	}
	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("empty fallback response")
	}
	return strings.TrimSpace(parsed.Choices[0].Message.Content), nil
}

func HistoryFromRepo(role, body string) ChatMessage {
	switch role {
	case "visitor":
		return ChatMessage{Role: "user", Content: body}
	case "bot", "agent":
		return ChatMessage{Role: "assistant", Content: body}
	default:
		return ChatMessage{Role: "user", Content: body}
	}
}
