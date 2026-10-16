package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	Cyan    = color.New(color.FgCyan, color.Bold).SprintFunc()
	Green   = color.New(color.FgGreen, color.Bold).SprintFunc()
	Red     = color.New(color.FgRed, color.Bold).SprintFunc()
	Yellow  = color.New(color.FgYellow, color.Bold).SprintFunc()
	Magenta = color.New(color.FgMagenta, color.Bold).SprintFunc()
	Blue    = color.New(color.FgBlue, color.Bold).SprintFunc()
	Gray    = color.New(color.FgHiBlack).SprintFunc()
	Bold    = color.New(color.Bold).SprintFunc()
)

func Banner() {
	banner := `
  __  __       _ _ _  __ _           
 |  \/  |     (_) (_)/ _(_)          
 | \  / | __ _ _| |_| |_ _  ___ __ _ 
 | |\/| |/ _` + "`" + ` | | |  _| |/ __/ _` + "`" + ` |
 | |  | | (_| | | | | | | | (_| (_| |
 |_|  |_|\__,_|_|_|_|_| |_|\___\__,_|
`
	fmt.Println(Cyan(banner))
	fmt.Println(Gray(" Official CLI for Lumi Email Infrastructure Platform"))
	fmt.Println()
}

func Success(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s %s\n", Green("✔"), msg)
}

func Info(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s %s\n", Blue("ℹ"), msg)
}

func Warn(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s %s\n", Yellow("⚠"), msg)
}

func Error(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s %s\n", Red("✖"), msg)
}

func PrintJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func RenderKeyValue(key string, val any) {
	fmt.Printf("  %s %v\n", Gray(fmt.Sprintf("%-16s:", key)), Bold(fmt.Sprintf("%v", val)))
}

func StatusBadge(status string) string {
	lower := strings.ToLower(status)
	switch lower {
	case "verified", "delivered", "active", "success":
		return Green(fmt.Sprintf("[%s]", strings.ToUpper(status)))
	case "pending", "queued", "sending":
		return Yellow(fmt.Sprintf("[%s]", strings.ToUpper(status)))
	case "failed", "bounced", "revoked", "error", "complained":
		return Red(fmt.Sprintf("[%s]", strings.ToUpper(status)))
	default:
		return Gray(fmt.Sprintf("[%s]", strings.ToUpper(status)))
	}
}

func ExitWithError(format string, a ...any) {
	Error(format, a...)
	os.Exit(1)
}
