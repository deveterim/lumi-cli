package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/eadafonso/lumi/lumi-cli/pkg/config"
	"github.com/eadafonso/lumi/lumi-cli/pkg/ui"
	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:     "whoami",
	Aliases: []string{"status"},
	Short:   "Show current authentication and workspace status",
	Run: func(cmd *cobra.Command, args []string) {
		if cfg.APIKey == "" {
			ui.Warn("No API key configured. Run 'lumi login' to authenticate.")
			return
		}

		c := getClientOrExit()
		ctx := context.Background()

		configPath, _ := config.GetConfigPath()

		maskedKey := cfg.APIKey
		if len(maskedKey) > 10 {
			maskedKey = maskedKey[:7] + "..." + maskedKey[len(maskedKey)-4:]
		}

		domains, err := c.ListDomains(ctx)
		domainsCount := 0
		if err == nil {
			domainsCount = len(domains)
		}

		if flagJSON {
			_ = ui.PrintJSON(map[string]any{
				"api_key":       maskedKey,
				"base_url":      cfg.BaseURL,
				"config_path":   configPath,
				"domains_count": domainsCount,
				"authenticated": err == nil,
			})
			return
		}

		ui.Banner()
		fmt.Println(ui.Bold("Authentication & Environment Status:"))
		fmt.Println(strings.Repeat("─", 45))
		ui.RenderKeyValue("API Key", maskedKey)
		ui.RenderKeyValue("Base URL", cfg.BaseURL)
		ui.RenderKeyValue("Config Path", configPath)

		if err != nil {
			ui.RenderKeyValue("Connection", ui.Red(fmt.Sprintf("Failed (%v)", err)))
		} else {
			ui.RenderKeyValue("Connection", ui.Green("Connected & Active"))
			ui.RenderKeyValue("Total Domains", domainsCount)
		}
		fmt.Println(strings.Repeat("─", 45))
	},
}
