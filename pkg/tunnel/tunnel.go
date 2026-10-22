package tunnel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/deveterim/lumi-cli/pkg/ui"
)

type ListenerOptions struct {
	Port      int
	ForwardTo string
	Secret    string
}

type WebhookEvent struct {
	ID        string         `json:"id"`
	Type      string         `json:"type"`
	CreatedAt string         `json:"created_at"`
	Data      map[string]any `json:"data"`
}

func StartListener(ctx context.Context, opts ListenerOptions) error {
	forwardURL := opts.ForwardTo

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}

		var event WebhookEvent
		_ = json.Unmarshal(body, &event)

		eventType := event.Type
		if eventType == "" {
			eventType = "custom.event"
		}

		fmt.Printf("\n%s %s %s\n",
			ui.Gray(time.Now().Format("15:04:05")),
			ui.Cyan("➔ WEBHOOK RECEIVED"),
			ui.Bold(eventType),
		)

		if event.ID != "" {
			fmt.Printf("  %s %s\n", ui.Gray("Event ID:"), ui.Bold(event.ID))
		}

		// Forwarding
		if forwardURL != "" {
			req, err := http.NewRequestWithContext(r.Context(), r.Method, forwardURL, bytes.NewReader(body))
			if err != nil {
				ui.Error("Failed to create forward request: %v", err)
				http.Error(w, "Internal forward error", http.StatusInternalServerError)
				return
			}

			// Copy original headers
			for k, v := range r.Header {
				req.Header[k] = v
			}

			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Do(req)
			latency := time.Since(start)

			if err != nil {
				ui.Error("Failed to forward to %s: %v (%s)", forwardURL, err, latency)
				w.WriteHeader(http.StatusBadGateway)
				_, _ = w.Write([]byte(fmt.Sprintf(`{"status":"error","forward_error":"%s"}`, err.Error())))
				return
			}
			defer resp.Body.Close()

			forwardBody, _ := io.ReadAll(resp.Body)

			statusColor := ui.Green
			if resp.StatusCode >= 400 {
				statusColor = ui.Red
			}

			fmt.Printf("  %s %s ➔ %s (%s)\n",
				ui.Gray("Forwarded:"),
				statusColor(fmt.Sprintf("[%d %s]", resp.StatusCode, http.StatusText(resp.StatusCode))),
				forwardURL,
				latency,
			)

			w.WriteHeader(resp.StatusCode)
			_, _ = w.Write(forwardBody)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"received"}`))
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", opts.Port),
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	fmt.Println(ui.Green("● Webhook listener ready!"))
	fmt.Printf("  %s http://localhost:%d\n", ui.Gray("Listening on:"), opts.Port)
	if forwardURL != "" {
		fmt.Printf("  %s %s\n", ui.Gray("Forwarding to:"), ui.Cyan(forwardURL))
	}
	fmt.Printf("  %s Press %s to exit\n\n", ui.Gray("Quit:"), ui.Bold("Ctrl+C"))

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
