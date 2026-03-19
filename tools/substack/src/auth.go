package main

import (
	"fmt"
	"os"
	"strings"

	"entext-applications/internal/substack"
	"entext-applications/internal/validator"
)

const profileCmdHelp = `
Set Substack email address for authentication.

Usage: substack profile [flags]

Flags:
  -email string
		Substack email address

Examples:
  # Set email address
  substack profile -email "user@example.com"

  # Interactive mode (prompts for email)
  substack profile`

const authCmdHelp = `
Complete Substack authentication with email link.

Usage: substack auth [flags]

Flags:
  -auth-link string
		Authentication link from Substack email

Examples:
  # Complete auth with link (for automation)
  substack auth -auth-link "https://substack.com/auth?token=..."

  # Interactive mode (prompts for link)
  substack auth

Flow:
1. First run: substack profile -email <email>
2. Check your email for the authentication link
3. Run: substack auth -auth-link <link>`

// runProfile sets the email address on the session
func runProfile(client *substack.Client, session *substack.Session, args []string) {
	var emailAddr string
	fs := newFlagSet("email_login")
	fs.StringVar(&emailAddr, "email", "", "Substack email address")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	emailAddr = strings.TrimSpace(emailAddr)
	v := validator.New()
	v.Check(emailAddr != "", "email", "email must not be empty")
	v.Check(validator.Matches(emailAddr, validator.EmailRX), "email", "must be a valid email address")

	// Validate email format
	if !v.Valid() {
		for k, v := range v.Errors {
			fmt.Fprintf(os.Stderr, "invalid input: %s - %s\n", k, v)
		}
		os.Exit(1)
	}

	// Set email on session
	session.Email = emailAddr
	client.StartEmailLinkLogin(emailAddr)

	fmt.Printf("✓ Email saved!\nEmail set to: %s\n", emailAddr)
	fmt.Println("Next step: Request authentication link\n\tsubstack auth")
}

// runAuth completes authentication with the link from Substack
func runAuth(client *substack.Client, session *substack.Session, args []string) {
	// check valid session before authenticate
	checkValidSession(session)

	fs := newFlagSet("auth")
	var authString string
	fs.StringVar(&authString, "auth_string", "", "Authentication link or code from Substack email")

	err := fs.Parse(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	// Request authentication link if not provided
	authString = strings.TrimSpace(authString)
	v := validator.New()
	v.Check(authString != "", "auth_string", "auth_string must not be empty")

	// validate input
	if !v.Valid() {
		for k, v := range v.Errors {
			fmt.Fprintf(os.Stderr, "validation error: %s - %s\n", k, v)
		}
		os.Exit(1)
	}

	// Complete authentication flow
	fmt.Println("\nCompleting authentication...")
	if strings.HasPrefix(authString, "https://") {
		err = client.AuthenticateFromResponse(session, authString)
	} else {
		err = client.AuthorizationCodeComplete(session, authString)
	}

	if err != nil {
		exitWithError(err)
	}

	fmt.Println("✓ Authentication successful!\n")
	fmt.Printf("Session saved for: %s\n", session.Email)
}
