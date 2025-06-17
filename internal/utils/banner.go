package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/chud-lori/go-boilerplate/config"
)

func Banner(cfg *config.AppConfig) {
	const (
		Cyan    = "\x1b[36m"
		Green   = "\x1b[32m"
		Magenta = "\x1b[35m"
		Yellow  = "\x1b[33m"
		Bold    = "\x1b[1m"
		Reset   = "\x1b[0m"
	)

	appEnv := strings.ToUpper(cfg.AppEnv)
	logLevel := strings.ToUpper(cfg.LogLevel)

	banner := fmt.Sprintf(`
%s%s╔════════════════════════════════════════════════════════════════════╗%s
%s%s║%s                      Go Net/HTTP Boilerplate                       %s║%s
%s%s╟────────────────────────────────────────────────────────────────────╢%s
%s%s║%s [>] Status       : %s%-48s%s║%s
%s%s║%s [#] Environment  : %s%-48s%s║%s
%s%s║%s [*] Log Level    : %s%-48s%s║%s
%s%s║%s [@] Author       : %s%-48s%s║%s
%s%s╚════════════════════════════════════════════════════════════════════╝%s
`,
		Green, Bold, Reset, // Top border
		Green, Bold, Reset, Green, Reset, // Title
		Green, Bold, Reset, // Divider
		Green, Bold, Reset, Green, "Running", Green, Reset, // Status
		Green, Bold, Reset, Yellow, appEnv, Green, Reset, // Environment
		Green, Bold, Reset, Yellow, logLevel, Green, Reset, // Log Level
		Green, Bold, Reset, Magenta, "@chud_lori", Green, Reset, // Author
		Green, Bold, Reset) // Bottom border

	fmt.Println(banner)

	fmt.Printf("      %s%sVersion%s      : %s%s%s\n", Bold, Cyan, Reset, Green, cfg.Version, Reset)
	fmt.Printf("      %s%sCurrent Time%s : %s%s%s\n", Bold, Cyan, Reset, Green, time.Now().Format("2006-01-02 15:04:05"), Reset)
	fmt.Printf("      %s%sListening on%s : %s%d%s\n", Bold, Cyan, Reset, Green, cfg.ServerPort, Reset)
	fmt.Println()
}
