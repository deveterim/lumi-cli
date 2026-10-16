package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/eadafonso/lumi/lumi-cli/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	apiKeyName       string
	apiKeyPermission string
)

var apiKeysCmd = &cobra.Command{
	Use:     "api-keys",
	Aliases: []string{"keys", "key"},
	Short:   "Manage workspace API Keys",
}

var apiKeysListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all API keys in the workspace",
	Run: func(cmd *cobra.Command, args []string) {
		c := getClientOrExit()
		ctx := context.Background()

		keys, err := c.ListApiKeys(ctx)
		if err != nil {
			ui.ExitWithError("Failed to list API keys: %v", err)
		}

		if flagJSON {
			_ = ui.PrintJSON(keys)
			return
		}

		if len(keys) == 0 {
			ui.Info("No API keys found. Create one with '%s'", ui.Bold("lumi api-keys create --name <name>"))
			return
		}

		fmt.Println(ui.Bold("API Keys:"))
		fmt.Printf("%-24s %-24s %-16s %s\n", ui.Gray("ID"), ui.Gray("NAME"), ui.Gray("PERMISSION"), ui.Gray("CREATED AT"))
		fmt.Println(strings.Repeat("─", 75))

		for _, k := range keys {
			fmt.Printf("%-24s %-24s %-16s %s\n", k.ID, ui.Bold(k.Name), k.Permission, k.CreatedAt)
		}
	},
}

var apiKeysCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new API key",
	Run: func(cmd *cobra.Command, args []string) {
		c := getClientOrExit()
		ctx := context.Background()

		if apiKeyName == "" {
			ui.ExitWithError("Key name (--name) is required.")
		}

		key, err := c.CreateApiKey(ctx, apiKeyName, apiKeyPermission)
		if err != nil {
			ui.ExitWithError("Failed to create API key: %v", err)
		}

		if flagJSON {
			_ = ui.PrintJSON(key)
			return
		}

		ui.Success("API Key created successfully!")
		ui.RenderKeyValue("ID", key.ID)
		ui.RenderKeyValue("Name", key.Name)
		if key.Token != "" {
			ui.RenderKeyValue("Secret Token", ui.Green(key.Token))
			ui.Warn("Make sure to copy your secret token now. You won't be able to see it again!")
		}
	},
}

var apiKeysRevokeCmd = &cobra.Command{
	Use:     "revoke <id>",
	Aliases: []string{"delete", "remove"},
	Short:   "Revoke and delete an API key",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClientOrExit()
		ctx := context.Background()
		keyID := args[0]

		if err := c.DeleteApiKey(ctx, keyID); err != nil {
			ui.ExitWithError("Failed to revoke API key %s: %v", keyID, err)
		}

		ui.Success("API key %s revoked successfully.", keyID)
	},
}

func init() {
	apiKeysCreateCmd.Flags().StringVar(&apiKeyName, "name", "", "Name or label for the API Key")
	apiKeysCreateCmd.Flags().StringVar(&apiKeyPermission, "permission", "full_access", "Permission scope: full_access or sending_only")

	apiKeysCmd.AddCommand(apiKeysListCmd)
	apiKeysCmd.AddCommand(apiKeysCreateCmd)
	apiKeysCmd.AddCommand(apiKeysRevokeCmd)
}
