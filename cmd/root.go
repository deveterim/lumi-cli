package cmd

import (
	"fmt"
	"os"

	"github.com/eadafonso/lumi/lumi-cli/pkg/client"
	"github.com/eadafonso/lumi/lumi-cli/pkg/config"
	"github.com/eadafonso/lumi/lumi-cli/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	flagAPIKey  string
	flagBaseURL string
	flagJSON    bool

	cfg        *config.Config
	mailClient *client.Client
)

var rootCmd = &cobra.Command{
	Use:   "lumi",
	Short: "Official CLI for Lumi Email Platform (Resend-like)",
	Long: ui.Cyan(`
  _      _    _  __  __  _____ 
 | |    | |  | ||  \/  ||_   _|
 | |    | |  | || \  / |  | |  
 | |    | |  | || |\/| |  | |  
 | |____| |__| || |  | | _| |_ 
 |______|\____/ |_|  |_||_____|
`) + "\nOfficial Command Line Interface for Lumi Email Platform:\nDisparar e-mails transacionais, gerir domínios, chaves de API e testar webhooks locais.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		cfg, err = config.Load()
		if err != nil {
			cfg = &config.Config{BaseURL: config.DefaultBaseURL}
		}

		if flagAPIKey != "" {
			cfg.APIKey = flagAPIKey
		}
		if flagBaseURL != "" {
			cfg.BaseURL = flagBaseURL
		}

		mailClient = client.New(cfg.APIKey, cfg.BaseURL)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagAPIKey, "api-key", "", "Lumi API Key (lu_live_... overrides config & LUMI_API_KEY)")
	rootCmd.PersistentFlags().StringVar(&flagBaseURL, "base-url", "", "Lumi Base URL (overrides default https://api.lumi.email/v1)")
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "Output results in formatted JSON")

	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(whoamiCmd)
	rootCmd.AddCommand(emailsCmd)
	rootCmd.AddCommand(domainsCmd)
	rootCmd.AddCommand(apiKeysCmd)
	rootCmd.AddCommand(templatesCmd)
	rootCmd.AddCommand(webhooksCmd)
	rootCmd.AddCommand(versionCmd)
}

func getClientOrExit() *client.Client {
	if cfg.APIKey == "" {
		ui.ExitWithError("Nenhuma API Key configurada. Execute '%s' ou passe '--api-key lu_live_...'.", ui.Bold("lumi login"))
	}
	return mailClient
}
