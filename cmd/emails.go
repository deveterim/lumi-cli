package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/eadafonso/lumi/lumi-cli/pkg/client"
	"github.com/eadafonso/lumi/lumi-cli/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	emailTo          []string
	emailFrom        string
	emailSubject     string
	emailHtml        string
	emailText        string
	emailReplyTo     string
	emailCc          []string
	emailBcc         []string
	emailAttachments []string
	emailTemplateID  string
	emailVariables   string
	emailScheduledAt string
	emailTags        []string
)

var emailsCmd = &cobra.Command{
	Use:     "emails",
	Aliases: []string{"email", "mail"},
	Short:   "Send and manage transactional emails",
}

var emailsSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a transactional email",
	Example: `  lumi emails send --to user@example.com --from hello@yourdomain.com --subject "Hello" --text "Hello world"
  lumi emails send --to user@example.com --from hello@yourdomain.com --subject "Invoice" --html "<h1>Your Invoice</h1>" --attach ./invoice.pdf
  lumi emails send --to user@example.com --from hello@yourdomain.com --subject "Welcome" --template tmpl_welcome --variables '{"name":"John"}'`,
	Run: func(cmd *cobra.Command, args []string) {
		c := getClientOrExit()
		ctx := context.Background()

		if len(emailTo) == 0 {
			ui.ExitWithError("Recipient (--to) is required.")
		}
		if emailFrom == "" {
			if cfg.DefaultFrom != "" {
				emailFrom = cfg.DefaultFrom
			} else {
				ui.ExitWithError("Sender (--from) is required.")
			}
		}
		if emailSubject == "" && emailTemplateID == "" {
			ui.ExitWithError("Subject (--subject) or Template (--template) is required.")
		}

		htmlContent := emailHtml
		if strings.HasPrefix(htmlContent, "@") {
			filePath := strings.TrimPrefix(htmlContent, "@")
			data, err := os.ReadFile(filePath)
			if err != nil {
				ui.ExitWithError("Failed to read HTML file %s: %v", filePath, err)
			}
			htmlContent = string(data)
		}

		textContent := emailText
		if strings.HasPrefix(textContent, "@") {
			filePath := strings.TrimPrefix(textContent, "@")
			data, err := os.ReadFile(filePath)
			if err != nil {
				ui.ExitWithError("Failed to read text file %s: %v", filePath, err)
			}
			textContent = string(data)
		}

		var variables map[string]any
		if emailVariables != "" {
			if err := json.Unmarshal([]byte(emailVariables), &variables); err != nil {
				ui.ExitWithError("Invalid JSON in --variables: %v", err)
			}
		}

		var attachments []client.Attachment
		for _, attachPath := range emailAttachments {
			data, err := os.ReadFile(attachPath)
			if err != nil {
				ui.ExitWithError("Failed to read attachment %s: %v", attachPath, err)
			}
			attachments = append(attachments, client.Attachment{
				Filename: filepath.Base(attachPath),
				Content:  base64.StdEncoding.EncodeToString(data),
			})
		}

		var tags []client.Tag
		for _, tagStr := range emailTags {
			if strings.Contains(tagStr, "=") {
				parts := strings.SplitN(tagStr, "=", 2)
				tags = append(tags, client.Tag{Name: parts[0], Value: parts[1]})
			} else {
				tags = append(tags, client.Tag{Name: tagStr, Value: "true"})
			}
		}

		params := &client.SendEmailParams{
			From:        emailFrom,
			To:          emailTo,
			Subject:     emailSubject,
			Html:        htmlContent,
			Text:        textContent,
			ReplyTo:     emailReplyTo,
			Cc:          emailCc,
			Bcc:         emailBcc,
			TemplateId:  emailTemplateID,
			Variables:   variables,
			ScheduledAt: emailScheduledAt,
			Attachments: attachments,
			Tags:        tags,
		}

		resp, err := c.SendEmail(ctx, params)
		if err != nil {
			ui.ExitWithError("Failed to send email: %v", err)
		}

		if flagJSON {
			_ = ui.PrintJSON(resp)
			return
		}

		ui.Success("Email queued successfully!")
		ui.RenderKeyValue("Email ID", resp.ID)
		ui.RenderKeyValue("Status", ui.StatusBadge(resp.Status))
		ui.RenderKeyValue("To", strings.Join(emailTo, ", "))
		ui.RenderKeyValue("From", emailFrom)
		if emailSubject != "" {
			ui.RenderKeyValue("Subject", emailSubject)
		}
	},
}

var emailsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get email status and delivery details by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClientOrExit()
		ctx := context.Background()
		emailID := args[0]

		email, err := c.GetEmail(ctx, emailID)
		if err != nil {
			ui.ExitWithError("Failed to get email %s: %v", emailID, err)
		}

		if flagJSON {
			_ = ui.PrintJSON(email)
			return
		}

		fmt.Println(ui.Bold("Email Delivery Information:"))
		fmt.Println(strings.Repeat("─", 45))
		ui.RenderKeyValue("ID", email.ID)
		ui.RenderKeyValue("Status", ui.StatusBadge(email.Status))
		ui.RenderKeyValue("Subject", email.Subject)
		ui.RenderKeyValue("From", email.From)
		ui.RenderKeyValue("To", strings.Join(email.To, ", "))
		ui.RenderKeyValue("Created At", email.CreatedAt)
		if email.LastEvent != "" {
			ui.RenderKeyValue("Last Event", email.LastEvent)
		}
		fmt.Println(strings.Repeat("─", 45))
	},
}

var emailsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List recently sent emails",
	Run: func(cmd *cobra.Command, args []string) {
		c := getClientOrExit()
		ctx := context.Background()

		emails, err := c.ListEmails(ctx)
		if err != nil {
			ui.ExitWithError("Failed to list emails: %v", err)
		}

		if flagJSON {
			_ = ui.PrintJSON(emails)
			return
		}

		if len(emails) == 0 {
			ui.Info("No emails found in this workspace.")
			return
		}

		fmt.Println(ui.Bold("Recent Emails:"))
		fmt.Printf("%-24s %-14s %-28s %s\n", ui.Gray("ID"), ui.Gray("STATUS"), ui.Gray("RECIPIENT"), ui.Gray("SUBJECT"))
		fmt.Println(strings.Repeat("─", 80))

		for _, em := range emails {
			recip := strings.Join(em.To, ", ")
			if len(recip) > 26 {
				recip = recip[:23] + "..."
			}
			subj := em.Subject
			if len(subj) > 30 {
				subj = subj[:27] + "..."
			}
			fmt.Printf("%-24s %-22s %-28s %s\n", em.ID, ui.StatusBadge(em.Status), recip, subj)
		}
	},
}

func init() {
	emailsSendCmd.Flags().StringSliceVar(&emailTo, "to", nil, "Recipient email address(es) (comma-separated or multiple --to)")
	emailsSendCmd.Flags().StringVar(&emailFrom, "from", "", "Sender email address (e.g. hello@yourdomain.com)")
	emailsSendCmd.Flags().StringVar(&emailSubject, "subject", "", "Email subject line")
	emailsSendCmd.Flags().StringVar(&emailHtml, "html", "", "Email HTML body (or @path/to/file.html)")
	emailsSendCmd.Flags().StringVar(&emailText, "text", "", "Email Plain Text body (or @path/to/file.txt)")
	emailsSendCmd.Flags().StringVar(&emailReplyTo, "reply-to", "", "Reply-to address")
	emailsSendCmd.Flags().StringSliceVar(&emailCc, "cc", nil, "CC email address(es)")
	emailsSendCmd.Flags().StringSliceVar(&emailBcc, "bcc", nil, "BCC email address(es)")
	emailsSendCmd.Flags().StringSliceVar(&emailAttachments, "attach", nil, "Path to file(s) to attach")
	emailsSendCmd.Flags().StringVar(&emailTemplateID, "template", "", "Template ID")
	emailsSendCmd.Flags().StringVar(&emailVariables, "variables", "", "Template variables in JSON format")
	emailsSendCmd.Flags().StringVar(&emailScheduledAt, "scheduled-at", "", "ISO-8601 timestamp for scheduled delivery")
	emailsSendCmd.Flags().StringSliceVar(&emailTags, "tag", nil, "Tags in key=value format (e.g. --tag env=prod)")

	emailsCmd.AddCommand(emailsSendCmd)
	emailsCmd.AddCommand(emailsGetCmd)
	emailsCmd.AddCommand(emailsListCmd)
}
