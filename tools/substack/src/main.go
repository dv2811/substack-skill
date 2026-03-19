package main

import (
	"encoding/json"
	"entext-applications/internal/substack"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

// getSessionFile returns the path to the session file in the binary's directory
func getSessionFile() (string, error) {
	// Check environment variable first (for testing/custom setups)
	if custom := os.Getenv("SUBSTACK_SESSION_FILE"); custom != "" {
		return custom, nil
	}

	// Get the directory where the binary is located
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("cannot determine executable path: %w", err)
	}
	binDir := filepath.Dir(execPath)

	return filepath.Join(binDir, "session.json"), nil
}

// saveSession saves the current session with updated tokens
func saveSession(session *substack.Session, sessionFile string) error {
	// Save session to file
	data, err := session.Save()
	if err != nil {
		return fmt.Errorf("failed to serialize session: %w", err)
	}

	if err := os.WriteFile(sessionFile, data, 0644); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		printUsage(substackToolUsage)
		os.Exit(1)
	}

	command := os.Args[1]
	// Help should work without session
	if command == "help" || command == "-h" || command == "--help" {
		printUsage(substackToolUsage)
		os.Exit(0)
	}

	// Check for help flag in command args before loading session
	for _, arg := range os.Args[2:] {
		if arg == "-h" || arg == "--help" {
			// Route to command help without session
			switch command {
			case "inbox":
				printUsage(InboxCmdHelp)
			case "article":
				printUsage(articleHelp)
			case "search":
				printUsage(SearchCmdHelp)
			case "auth":
				printUsage(authCmdHelp)
			case "profile":
				printUsage(profileCmdHelp)
			}
			os.Exit(0)
		}
	}

	// Initialize Substack client
	client := substack.NewClient()

	// Load session from binary directory
	sessionFile, err := getSessionFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding session file: %v\n", err)
		os.Exit(1)
	}

	session := &substack.Session{}
	// Load session if not in initiation flow
	if command != "profile" {
		file, err := os.OpenFile(sessionFile, os.O_RDONLY, 0644)
		if err == nil {
			err = session.LoadFromFile(file)
			// no defer just close file explicitly
			file.Close()
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading session: %v\n", err)
		fmt.Fprintf(os.Stderr, "Run 'substack profile -email <email>' then 'substack auth' to authenticate\n")
		os.Exit(1)
	}

	// Set up signal handling to save session on interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		// Save session before exit on signal
		if saveErr := saveSession(session, sessionFile); saveErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to save session: %v\n", saveErr)
		}
		os.Exit(130)
	}()

	// Defer session save to persist any token updates after command completes
	defer func() {
		if err := saveSession(session, sessionFile); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to save session: %v\n", err)
		}
	}()

	// Route to appropriate command
	switch command {
	case "inbox":
		runInbox(client, session, os.Args[2:])
	case "article":
		runArticle(client, session, os.Args[2:])
	case "search":
		runSearch(client, session, os.Args[2:])
	case "profile":
		runProfile(client, session, os.Args[2:])
	case "auth":
		runAuth(client, session, os.Args[2:])
	case "help", "-h", "--help":
		printUsage(substackToolUsage)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage(substackToolUsage)
		os.Exit(1)
	}
}

const substackToolUsage = `
Substack CLI Tools
Usage: substack <command> [flags]
Commands:
  inbox    	Get chronological inbox posts
  article   Get article content by post ID
  search    Search posts with different modes
  profile   Set Substack email address
  auth      Authenticate with Substack via email link

Examples:
  substack profile -email "user@example.com"
  substack auth -auth_string "https://substack.com/auth?token=..."
  substack inbox
  substack inbox -after "2024-01-01T00:00:00.000Z"
  substack article -post-id 123456
  substack article -post-id 123456 -base-url "substack.com"
  substack search -query "AI" -mode top
  substack search -query "technology" -mode all -page 1
  substack search -query "newsletter" -mode subscribed -language en

Run 'substack <command> -h' for more information on a command.`

func printUsage(toolUsage string) {
	fmt.Println(toolUsage)
}

// printJSON prints data as formatted JSON
func printJSON(data any) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", " ")
	if err := encoder.Encode(data); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}

// exitWithError prints error and exits
func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

// Common flag sets for different commands
func newFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	return fs
}

// check if the session is valid
func checkValidSession(session *substack.Session) {
	// If email not set, prompt for it first
	if session == nil || session.Email == "" {
		fmt.Fprintf(os.Stderr, "existing session with valid email must be provided\n")
		os.Exit(1)
	}
}
