package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const Version = "1.0.0"

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Attachment struct {
	Content     string `json:"content"` // base64
	Filename    string `json:"filename"`
	ContentType string `json:"content_type,omitempty"`
}

type SendEmailParams struct {
	From        string            `json:"from"`
	To          any               `json:"to"` // string or []string
	Subject     string            `json:"subject"`
	Html        string            `json:"html,omitempty"`
	Text        string            `json:"text,omitempty"`
	Cc          []string          `json:"cc,omitempty"`
	Bcc         []string          `json:"bcc,omitempty"`
	ReplyTo     string            `json:"reply_to,omitempty"`
	TemplateId  string            `json:"template_id,omitempty"`
	Variables   map[string]any    `json:"variables,omitempty"`
	ScheduledAt string            `json:"scheduled_at,omitempty"`
	Attachments []Attachment      `json:"attachments,omitempty"`
	Tags        []Tag             `json:"tags,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
}

type SendEmailResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type Email struct {
	ID        string   `json:"id"`
	From      string   `json:"from"`
	To        []string `json:"to"`
	Subject   string   `json:"subject"`
	Html      string   `json:"html,omitempty"`
	Text      string   `json:"text,omitempty"`
	Status    string   `json:"status"`
	CreatedAt string   `json:"created_at"`
	LastEvent string   `json:"last_event,omitempty"`
}

type DomainRecord struct {
	Record   string `json:"record"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	Status   string `json:"status"`
	Priority int    `json:"priority,omitempty"`
}

type Domain struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Status    string         `json:"status"`
	Region    string         `json:"region,omitempty"`
	Records   []DomainRecord `json:"records"`
	CreatedAt string         `json:"created_at"`
}

type Template struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Subject   string `json:"subject"`
	Html      string `json:"html,omitempty"`
	CreatedAt string `json:"created_at"`
}

type ApiKey struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Token      string `json:"token,omitempty"`
	MaskedKey  string `json:"masked_key,omitempty"`
	Permission string `json:"permission"`
	CreatedAt  string `json:"created_at"`
}

func New(apiKey, baseURL string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) SetTimeout(t time.Duration) {
	c.httpClient.Timeout = t
}

func (c *Client) request(ctx context.Context, method, path string, body any, out any) error {
	if c.apiKey == "" {
		return errors.New("missing API Key. Run 'mailifica login' or provide --api-key / MAILIFICA_API_KEY")
	}

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(data)
	}

	url := c.baseURL + "/" + strings.TrimLeft(path, "/")
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "mailifica-cli/"+Version)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		var errData struct {
			Message string `json:"message"`
			Error   string `json:"error"`
		}
		_ = json.Unmarshal(respBytes, &errData)
		msg := errData.Message
		if msg == "" {
			msg = errData.Error
		}
		if msg == "" {
			msg = fmt.Sprintf("HTTP %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		}
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, msg)
	}

	if out != nil && len(respBytes) > 0 {
		return json.Unmarshal(respBytes, out)
	}

	return nil
}

// Emails
func (c *Client) SendEmail(ctx context.Context, params *SendEmailParams) (*SendEmailResponse, error) {
	var res struct {
		Data   SendEmailResponse `json:"data"`
		ID     string            `json:"id"`
		Status string            `json:"status"`
	}
	if err := c.request(ctx, "POST", "/emails", params, &res); err != nil {
		return nil, err
	}
	if res.Data.ID != "" {
		return &res.Data, nil
	}
	return &SendEmailResponse{ID: res.ID, Status: res.Status}, nil
}

func (c *Client) GetEmail(ctx context.Context, id string) (*Email, error) {
	var res struct {
		Data Email `json:"data"`
	}
	if err := c.request(ctx, "GET", "/emails/"+id, nil, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (c *Client) ListEmails(ctx context.Context) ([]Email, error) {
	var res struct {
		Data []Email `json:"data"`
	}
	if err := c.request(ctx, "GET", "/emails", nil, &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

// Domains
func (c *Client) ListDomains(ctx context.Context) ([]Domain, error) {
	var res struct {
		Data []Domain `json:"data"`
	}
	if err := c.request(ctx, "GET", "/domains", nil, &res); err != nil {
		// Fallback if returned directly as array
		var direct []Domain
		if err2 := c.request(ctx, "GET", "/domains", nil, &direct); err2 == nil {
			return direct, nil
		}
		return nil, err
	}
	return res.Data, nil
}

func (c *Client) CreateDomain(ctx context.Context, name string) (*Domain, error) {
	var res Domain
	if err := c.request(ctx, "POST", "/domains", map[string]string{"name": name}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetDomain(ctx context.Context, id string) (*Domain, error) {
	var res Domain
	if err := c.request(ctx, "GET", "/domains/"+id, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) VerifyDomain(ctx context.Context, id string) (*Domain, error) {
	var res Domain
	if err := c.request(ctx, "POST", "/domains/"+id+"/verify", nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) DeleteDomain(ctx context.Context, id string) error {
	return c.request(ctx, "DELETE", "/domains/"+id, nil, nil)
}

// Templates
func (c *Client) ListTemplates(ctx context.Context) ([]Template, error) {
	var res struct {
		Data []Template `json:"data"`
	}
	if err := c.request(ctx, "GET", "/templates", nil, &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (c *Client) GetTemplate(ctx context.Context, id string) (*Template, error) {
	var res struct {
		Data Template `json:"data"`
	}
	if err := c.request(ctx, "GET", "/templates/"+id, nil, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

// API Keys
func (c *Client) ListApiKeys(ctx context.Context) ([]ApiKey, error) {
	var res struct {
		Data []ApiKey `json:"data"`
	}
	if err := c.request(ctx, "GET", "/api-keys", nil, &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (c *Client) CreateApiKey(ctx context.Context, name string, permission string) (*ApiKey, error) {
	var res ApiKey
	payload := map[string]string{
		"name":       name,
		"permission": permission,
	}
	if err := c.request(ctx, "POST", "/api-keys", payload, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) DeleteApiKey(ctx context.Context, id string) error {
	return c.request(ctx, "DELETE", "/api-keys/"+id, nil, nil)
}
