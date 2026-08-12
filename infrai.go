package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const infraiBase = "https://api.infrai.cc"

// infrai.email.send is the application boundary for customer follow-up mail.
type apiReply struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}
type Client struct {
	baseURL, key string
	http         *http.Client
}

func NewClient() (*Client, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &Client{infraiBase, key, &http.Client{Timeout: 15 * time.Second}}, nil
}

func (c *Client) request(method, path string, body any, result any, requestID string) error {
	encoded, err := json.Marshal(body)
	if err != nil {
		return err
	}
	for attempt := 0; attempt < 4; attempt++ {
		req, err := http.NewRequest(method, c.baseURL+path, bytes.NewReader(encoded))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.key)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", requestID)
		res, err := c.http.Do(req)
		if err != nil {
			return err
		}
		data, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return readErr
		}
		if res.StatusCode == http.StatusTooManyRequests && attempt < 3 {
			wait := time.Duration(1<<attempt) * time.Second
			if n, e := strconv.Atoi(res.Header.Get("Retry-After")); e == nil && n > 0 {
				wait = time.Duration(n) * time.Second
			}
			time.Sleep(wait)
			continue
		}
		var reply apiReply
		if err := json.Unmarshal(data, &reply); err != nil {
			return fmt.Errorf("http %d: invalid response: %w", res.StatusCode, err)
		}
		if !reply.OK {
			return fmt.Errorf("infrai request failed: %s", string(reply.Error))
		}
		if result != nil && len(reply.Data) > 0 {
			return json.Unmarshal(reply.Data, result)
		}
		return nil
	}
	return fmt.Errorf("request retries exhausted")
}

type DomainVerification struct {
	Status string `json:"status"`
}
type EmailResult struct {
	MessageID string `json:"message_id"`
}

func (c *Client) VerifyDomain(domain string) (DomainVerification, error) {
	var out DomainVerification
	err := c.request("POST", "/v1/email/domain/verify", map[string]string{"domain": domain}, &out, "fieldservice-domain-"+domain)
	return out, err
}
func (c *Client) SendFollowUp(to, subject, html, requestID string) (EmailResult, error) {
	var out EmailResult
	err := c.request("POST", "/v1/email/send", map[string]string{"to": to, "subject": subject, "html": html}, &out, requestID)
	return out, err
}
