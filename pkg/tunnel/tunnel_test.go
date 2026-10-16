package tunnel

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListenerForwarding(t *testing.T) {
	forwardReceived := false
	forwardServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		forwardReceived = true
		body, _ := io.ReadAll(r.Body)
		if !bytes.Contains(body, []byte("email.delivered")) {
			t.Errorf("expected body to contain email.delivered")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"received":true}`))
	}))
	defer forwardServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	port := 14892
	go func() {
		_ = StartListener(ctx, ListenerOptions{
			Port:      port,
			ForwardTo: forwardServer.URL,
		})
	}()

	// Wait for listener to boot
	time.Sleep(150 * time.Millisecond)

	payload := `{"id":"evt_1","type":"email.delivered","data":{"recipient":"user@example.com"}}`
	resp, err := http.Post(
		"http://localhost:14892/",
		"application/json",
		bytes.NewReader([]byte(payload)),
	)

	if err != nil {
		t.Fatalf("failed to send webhook: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if !forwardReceived {
		t.Errorf("expected destination server to receive forwarded webhook")
	}
}
