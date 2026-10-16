package cmd

import (
	"fmt"
	"runtime"

	"github.com/eadafonso/lumi/lumi-cli/pkg/client"
	"github.com/eadafonso/lumi/lumi-cli/pkg/ui"
	"github.com/spf13/cobra"
)

var Version = client.Version

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of lumi CLI",
	Run: func(cmd *cobra.Command, args []string) {
		if flagJSON {
			_ = ui.PrintJSON(map[string]string{
				"version":    Version,
				"go_version": runtime.Version(),
				"os":         runtime.GOOS,
				"arch":       runtime.GOARCH,
			})
			return
		}

		fmt.Printf("lumi version %s (%s/%s, built with %s)\n",
			ui.Bold(Version),
			runtime.GOOS,
			runtime.GOARCH,
			runtime.Version(),
		)
	},
}
