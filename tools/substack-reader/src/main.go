package main
import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"entext-applications/internal/substack"
)

// getSessionFile returns the path to the session file using XDG conventions
func getSessionFile() (string, error) {
	// Check environment variable first (for testing/custom setups)
	if custom := os.Getenv("SUBSTACK_SESSION_FILE"); custom != "" {
		return custom, nil
	}

	var configDir string
	// XDG Base Directory Specification
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		configDir = filepath.Join(xdgConfig, "substack-reader")
	} else {
		// Default to ~/.config/substack-reader
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot determine home directory: %w", err)
		}
		configDir = filepath.Join(home, ".config", "substack-reader")
	}
	return filepath.Join(configDir, "session.json"), nil
}

// saveSession saves the current session with updated tokens
func saveSession(session *substack.Session, sessionFile string) error {
	// Ensure directory exists
	dir := filepath.Dir(sessionFile)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Save session to file
	data, err := session.Save()
	if err != nil {
		return fmt.Errorf("failed to serialize session: %w", err)
	}

	if err := os.WriteFile(sessionFile, data, 0600); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	// Help should work without session
	if command == "help" || command == "-h" || command == "--help" {
		printUsage()
		os.Exit(0)
	}

	// Check for help flag in command args before loading session
	for _, arg := range os.Args[2:] {
		if arg == "-h" || arg == "--help" {
			// Route to command help without session
			switch command {
			case "inbox":
				runInboxHelp()
			case "article":
				runArticleHelp()
			case "search":
				runSearchHelp()
			case "auth":
				runAuthHelp()
			}
			os.Exit(0)
		}
	}

	// Initialize Substack client
	client := substack.NewClient()

	// Load session from XDG config path
	sessionFile, err := getSessionFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding session file: %v\n", err)
		os.Exit(1)
	}

	session, err := substack.NewSessionFromFile(sessionFile)
	if err != nil {
		// auth command: create or renew a session
		if command == "auth" {
			session = &substack.Session{}
		} else {
			fmt.Fprintf(os.Stderr, "Error loading session: %v\n", err)
			fmt.Fprintf(os.Stderr, "Run setup script or 'substack auth' to authenticate\n")
			os.Exit(1)
		}
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
	case "auth":
		runAuth(client, session, os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}
const toolUsage = `
Substack CLI Tools
Usage: substack <command> [flags]
Commands:
  inbox                   Get chronological inbox posts
  article                 Get article content by post ID
  search                  Search posts with different modes
  auth                    Authenticate with Substack via email link

Examples:
  substack auth
  substack inbox
  substack inbox -after "2024-01-01T00:00:00.000Z"
  substack article -post-id 123456
  substack article -post-id 123456 -base-url "substack.com"
  substack search -query "AI" -mode top
  substack search -query "technology" -mode all -page 1
  substack search -query "newsletter" -mode subscribed -language en

Run 'substack <command> -h' for more information on a command.`

func printUsage() {
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