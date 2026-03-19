package main

import (
	"fmt"
	"os"
	"time"

	"entext-applications/internal/koyfin"
	"entext-applications/internal/validator"
)

const authCmdHelp = `
Authenticate with Koyfin using email and password.

Usage: koyfin auth [flags]

Flags:
  -email string
        Koyfin email address
  -password string
        Koyfin password

Examples:
  # Interactive mode (prompts for credentials)
  koyfin auth

  # Non-interactive mode (for automation)
  koyfin auth -email "user@example.com" -password "secret"
`

// runAuth authenticates with Koyfin using email/password and obtains access tokens
func runAuth(client *koyfin.Client, session *koyfin.Session, args []string) {
	var email, pwd string
	fs := newFlagSet("auth")
	fs.StringVar(&email, "email", "", "Koyfin email address")
	fs.StringVar(&pwd, "password", "", "Koyfin password")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	v := validator.New()
	v.Check(email != "", "enail", "user email for Koyfin account must not be empty")
	v.Check(pwd != "", "password", "password for Koyfin account must not be empty")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")

	// Validate email format
	if !v.Valid() {
		for k, v := range v.Errors {
			fmt.Fprintf(os.Stderr, "invalid input: %s - %s\n", k, v)
		}
		os.Exit(1)
	}

	// Set credentials on session
	session.UserName = email
	session.Password = pwd

	// Perform login to get access tokens
	fmt.Println("\nAuthenticating with Koyfin...")
	currentTimeStamp := time.Now().Unix()
	err := client.Login(session, currentTimeStamp)
	if err != nil {
		exitWithError(err)
	}

	fmt.Printf("✓ Authentication successful!\nSession saved for: %s\n", email)
}
