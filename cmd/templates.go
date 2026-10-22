package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/deveterim/lumi-cli/pkg/ui"
	"github.com/spf13/cobra"
)

var templatesCmd = &cobra.Command{
	Use:     "templates",
	Aliases: []string{"template"},
	Short:   "View and manage email templates",
}

var templatesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all email templates in the workspace",
	Run: func(cmd *cobra.Command, args []string) {
		c := getClientOrExit()
		ctx := context.Background()

		templates, err := c.ListTemplates(ctx)
		if err != nil {
			ui.ExitWithError("Failed to list templates: %v", err)
		}

		if flagJSON {
			_ = ui.PrintJSON(templates)
			return
		}

		if len(templates) == 0 {
			ui.Info("No templates found in this workspace.")
			return
		}

		fmt.Println(ui.Bold("Email Templates:"))
		fmt.Printf("%-24s %-28s %s\n", ui.Gray("ID"), ui.Gray("NAME"), ui.Gray("SUBJECT"))
		fmt.Println(strings.Repeat("─", 75))

		for _, t := range templates {
			subj := t.Subject
			if len(subj) > 35 {
				subj = subj[:32] + "..."
			}
			fmt.Printf("%-24s %-28s %s\n", t.ID, ui.Bold(t.Name), subj)
		}
	},
}

var templatesGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get details and HTML content of a template",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClientOrExit()
		ctx := context.Background()
		templateID := args[0]

		template, err := c.GetTemplate(ctx, templateID)
		if err != nil {
			ui.ExitWithError("Failed to get template %s: %v", templateID, err)
		}

		if flagJSON {
			_ = ui.PrintJSON(template)
			return
		}

		fmt.Println(ui.Bold("Template Details:"))
		fmt.Println(strings.Repeat("─", 50))
		ui.RenderKeyValue("ID", template.ID)
		ui.RenderKeyValue("Name", template.Name)
		ui.RenderKeyValue("Subject", template.Subject)
		ui.RenderKeyValue("Created At", template.CreatedAt)
		fmt.Println(strings.Repeat("─", 50))

		if template.Html != "" {
			fmt.Println(ui.Bold("\nHTML Content:"))
			fmt.Println(template.Html)
		}
	},
}

func init() {
	templatesCmd.AddCommand(templatesListCmd)
	templatesCmd.AddCommand(templatesGetCmd)
}
