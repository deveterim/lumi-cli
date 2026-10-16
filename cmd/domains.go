package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/eadafonso/lumi/lumi-cli/pkg/ui"
	"github.com/spf13/cobra"
)

var domainsCmd = &cobra.Command{
	Use:     "domains",
	Aliases: []string{"domain"},
	Short:   "Manage sending domains and DNS verification",
}

var domainsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured sending domains",
	Run: func(cmd *cobra.Command, args []string) {
		c := getClientOrExit()
		ctx := context.Background()

		domains, err := c.ListDomains(ctx)
		if err != nil {
			ui.ExitWithError("Failed to list domains: %v", err)
		}

		if flagJSON {
			_ = ui.PrintJSON(domains)
			return
		}

		if len(domains) == 0 {
			ui.Info("No domains configured. Add one with '%s'", ui.Bold("lumi domains add <domain-name>"))
			return
		}

		fmt.Println(ui.Bold("Sending Domains:"))
		fmt.Printf("%-24s %-26s %-14s %s\n", ui.Gray("ID"), ui.Gray("DOMAIN NAME"), ui.Gray("STATUS"), ui.Gray("REGION"))
		fmt.Println(strings.Repeat("─", 75))

		for _, d := range domains {
			region := d.Region
			if region == "" {
				region = "us-east-1"
			}
			fmt.Printf("%-24s %-26s %-22s %s\n", d.ID, ui.Bold(d.Name), ui.StatusBadge(d.Status), region)
		}
	},
}

var domainsAddCmd = &cobra.Command{
	Use:   "add <domain-name>",
	Short: "Add a new domain for email sending",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClientOrExit()
		ctx := context.Background()
		domainName := args[0]

		domain, err := c.CreateDomain(ctx, domainName)
		if err != nil {
			ui.ExitWithError("Failed to add domain %s: %v", domainName, err)
		}

		if flagJSON {
			_ = ui.PrintJSON(domain)
			return
		}

		ui.Success("Domain %s added successfully!", ui.Bold(domain.Name))
		ui.RenderKeyValue("Domain ID", domain.ID)
		ui.RenderKeyValue("Status", ui.StatusBadge(domain.Status))
		fmt.Println()

		if len(domain.Records) > 0 {
			fmt.Println(ui.Bold("DNS Verification Records:"))
			fmt.Printf("%-10s %-32s %-8s %s\n", ui.Gray("TYPE"), ui.Gray("NAME / HOST"), ui.Gray("STATUS"), ui.Gray("VALUE"))
			fmt.Println(strings.Repeat("─", 90))
			for _, rec := range domain.Records {
				fmt.Printf("%-10s %-32s %-16s %s\n", rec.Type, rec.Name, ui.StatusBadge(rec.Status), rec.Value)
			}
			fmt.Println()
			fmt.Printf("After adding DNS records in your provider, verify with: %s\n", ui.Cyan("lumi domains verify "+domain.ID))
		}
	},
}

var domainsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get DNS records and verification status for a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClientOrExit()
		ctx := context.Background()
		domainID := args[0]

		domain, err := c.GetDomain(ctx, domainID)
		if err != nil {
			ui.ExitWithError("Failed to get domain %s: %v", domainID, err)
		}

		if flagJSON {
			_ = ui.PrintJSON(domain)
			return
		}

		fmt.Println(ui.Bold("Domain Details:"))
		fmt.Println(strings.Repeat("─", 50))
		ui.RenderKeyValue("ID", domain.ID)
		ui.RenderKeyValue("Name", domain.Name)
		ui.RenderKeyValue("Status", ui.StatusBadge(domain.Status))
		ui.RenderKeyValue("Region", domain.Region)
		fmt.Println(strings.Repeat("─", 50))

		if len(domain.Records) > 0 {
			fmt.Println(ui.Bold("\nDNS Records:"))
			fmt.Printf("%-10s %-32s %-16s %s\n", ui.Gray("TYPE"), ui.Gray("NAME / HOST"), ui.Gray("STATUS"), ui.Gray("VALUE"))
			fmt.Println(strings.Repeat("─", 90))
			for _, rec := range domain.Records {
				fmt.Printf("%-10s %-32s %-24s %s\n", rec.Type, rec.Name, ui.StatusBadge(rec.Status), rec.Value)
			}
		}
	},
}

var domainsVerifyCmd = &cobra.Command{
	Use:   "verify <id>",
	Short: "Trigger DNS verification check for a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClientOrExit()
		ctx := context.Background()
		domainID := args[0]

		domain, err := c.VerifyDomain(ctx, domainID)
		if err != nil {
			ui.ExitWithError("Failed to verify domain %s: %v", domainID, err)
		}

		if flagJSON {
			_ = ui.PrintJSON(domain)
			return
		}

		if strings.ToLower(domain.Status) == "verified" {
			ui.Success("Domain %s is VERIFIED and ready to send emails!", ui.Bold(domain.Name))
		} else {
			ui.Warn("Domain %s verification status: %s", domain.Name, ui.StatusBadge(domain.Status))
			fmt.Println("Please verify that all DNS records are correctly propagated.")
		}
	},
}

var domainsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a sending domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClientOrExit()
		ctx := context.Background()
		domainID := args[0]

		if err := c.DeleteDomain(ctx, domainID); err != nil {
			ui.ExitWithError("Failed to delete domain %s: %v", domainID, err)
		}

		ui.Success("Domain %s deleted successfully.", domainID)
	},
}

func init() {
	domainsCmd.AddCommand(domainsListCmd)
	domainsCmd.AddCommand(domainsAddCmd)
	domainsCmd.AddCommand(domainsGetCmd)
	domainsCmd.AddCommand(domainsVerifyCmd)
	domainsCmd.AddCommand(domainsDeleteCmd)
}
