package utils

import (
	"fmt"
	"time"

	"github.com/chud-lori/go-boilerplate/config"
)

func Banner(cfg *config.AppConfig) {
	banner := `
		************************************************
		* *
		* Go Net/HTTP Boilerplate              *
		* *
		* Server is now running...             *
		* *
		************************************************
		`
	// Print the banner
	fmt.Println(banner)
	fmt.Printf("           Version: %s\n", cfg.Version)                                   // Example: Add a version
	fmt.Printf("           Current Time: %s\n", time.Now().Format("2006-01-02 15:04:05")) // Example: Add current time
	fmt.Printf("           Listening on port: %s\n", cfg.ServerPort)
	fmt.Println() // Add an empty line for spacing
}
