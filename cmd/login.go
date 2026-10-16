package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/eadafonso/lumi/lumi-cli/pkg/config"
	"github.com/eadafonso/lumi/lumi-cli/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	loginKey     string
	loginBaseURL string
)

var loginCmd = &cobra.Command{
	Use:     "login",
	Aliases: []string{"configure", "auth"},
	Short:   "Authenticate and configure the Lumi CLI",
	Long:    "Set up your Lumi API Key and default endpoint settings stored securely in ~/.lumi/config.json",
	Run: func(cmd *cobra.Command, args []string) {
		ui.Banner()

		reader := bufio.NewReader(os.Stdin)

		apiKey := loginKey
		if apiKey == "" {
			fmt.Printf("%s: ", ui.Bold("Enter your Lumi API Key (e.g. ma_live_...)"))
			input, _ := reader.ReadString('\n')
			apiKey = strings.TrimSpace(input)
		}

		if apiKey == "" {
			ui.ExitWithError("API Key cannot be empty.")
		}

		baseURL := loginBaseURL
		if baseURL == "" {
			defaultURL := config.DefaultBaseURL
			if cfg.BaseURL != "" {
				defaultURL = cfg.BaseURL
			}
			fmt.Printf("%s [%s]: ", ui.Bold("Lumi API URL"), ui.Gray(defaultURL))
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if input != "" {
				baseURL = input
			} else {
				baseURL = defaultURL
			}
		}

		cfg.APIKey = apiKey
		cfg.BaseURL = baseURL

		if err := cfg.Save(); err != nil {
			ui.ExitWithError("Failed to save configuration: %v", err)
		}

		configPath, _ := config.GetConfigPath()
		ui.Success("Successfully authenticated and configured Lumi CLI!")
		fmt.Printf("  %s %s\n", ui.Gray("Config saved to:"), ui.Bold(configPath))
		fmt.Printf("  %s %s\n", ui.Gray("API Endpoint:"), ui.Cyan(cfg.BaseURL))
		fmt.Println()
		fmt.Printf("Try sending a test email with: %s\n", ui.Cyan("lumi emails send --to user@example.com --from hello@yourdomain.com --subject 'Test' --text 'Hello!'"))
	},
}

func init() {
	loginCmd.Flags().StringVar(&loginKey, "key", "", "Lumi API Key")
	loginCmd.Flags().StringVar(&loginBaseURL, "url", "", "Lumi API URL")
}
