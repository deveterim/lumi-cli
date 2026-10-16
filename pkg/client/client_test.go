package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientSendEmail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test_key_123" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if r.URL.Path != "/emails" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		var payload SendEmailParams
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":     "email_mock_789",
			"status": "queued",
		})
	}))
	defer server.Close()

	c := New("test_key_123", server.URL)
	resp, err := c.SendEmail(context.Background(), &SendEmailParams{
		From:    "sender@example.com",
		To:      []string{"receiver@example.com"},
		Subject: "Test Subject",
		Html:    "<p>Hello</p>",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != "email_mock_789" {
		t.Errorf("expected id email_mock_789, got %s", resp.ID)
	}
	if resp.Status != "queued" {
		t.Errorf("expected status queued, got %s", resp.Status)
	}
}

func TestClientDomains(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" && r.URL.Path == "/domains" {
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []Domain{
					{ID: "dom_1", Name: "mail.example.com", Status: "verified"},
				},
			})
			return
		}
		if r.Method == "POST" && r.URL.Path == "/domains" {
			_ = json.NewEncoder(w).Encode(Domain{
				ID:     "dom_2",
				Name:   "news.example.com",
				Status: "pending",
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	c := New("test_key_123", server.URL)
	domains, err := c.ListDomains(context.Background())
	if err != nil {
		t.Fatalf("failed to list domains: %v", err)
	}
	if len(domains) != 1 || domains[0].Name != "mail.example.com" {
		t.Errorf("unexpected domains: %+v", domains)
	}

	created, err := c.CreateDomain(context.Background(), "news.example.com")
	if err != nil {
		t.Fatalf("failed to create domain: %v", err)
	}
	if created.ID != "dom_2" || created.Name != "news.example.com" {
		t.Errorf("unexpected created domain: %+v", created)
	}
}
