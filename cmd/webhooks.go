package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/deveterim/lumi-cli/pkg/tunnel"
	"github.com/deveterim/lumi-cli/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	webhookForwardTo string
	webhookPort      int
	webhookSecret    string
)

var webhooksCmd = &cobra.Command{
	Use:     "webhooks",
	Aliases: []string{"webhook"},
	Short:   "Manage and test webhooks locally",
}

var webhooksListenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Start a local webhook tunnel listener and forward events to your local app",
	Example: `  # Listen on port 4000 and forward to local Laravel backend
  lumi webhooks listen --forward-to http://localhost:8000/lumi/webhook

  # Listen on custom port 8080 and forward to Next.js API
  lumi webhooks listen --port 8080 --forward-to http://localhost:3000/api/webhook`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.Banner()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-sigChan
			ui.Info("\nShutting down webhook listener...")
			cancel()
		}()

		opts := tunnel.ListenerOptions{
			Port:      webhookPort,
			ForwardTo: webhookForwardTo,
			Secret:    webhookSecret,
		}

		if err := tunnel.StartListener(ctx, opts); err != nil {
			ui.ExitWithError("Listener error: %v", err)
		}
	},
}

func init() {
	webhooksListenCmd.Flags().StringVarP(&webhookForwardTo, "forward-to", "f", "", "Local destination URL to forward webhook events (e.g. http://localhost:8000/lumi/webhook)")
	webhooksListenCmd.Flags().IntVarP(&webhookPort, "port", "p", 4000, "Port for the local listener")
	webhooksListenCmd.Flags().StringVar(&webhookSecret, "secret", "", "Webhook signing secret for verification")

	webhooksCmd.AddCommand(webhooksListenCmd)
}
