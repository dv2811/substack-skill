// Koyfin CLI tools - replicates MCP tools as command-line utilities
// Usage: koyfin <command> [flags]
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"entext-applications/internal/koyfin"
)

// getSessionFile returns the path to the session file using XDG conventions
func getSessionFile() (string, error) {
	// Check environment variable first (for testing/custom setups)
	if custom := os.Getenv("KOYFIN_SESSION_FILE"); custom != "" {
		return custom, nil
	}

	var configDir string

	// XDG Base Directory Specification
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		configDir = filepath.Join(xdgConfig, "koyfin")
	} else {
		// Default to ~/.config/koyfin
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot determine home directory: %w", err)
		}
		configDir = filepath.Join(home, ".config", "koyfin")
	}

	return filepath.Join(configDir, "session.json"), nil
}

// saveSession saves the current session with updated tokens
func saveSession(session *koyfin.Session, sessionFile string) error {
	// Ensure directory exists
	dir := filepath.Dir(sessionFile)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Save session to file (includes updated tokens, timestamps, cookies)
	if err := session.SaveToFile(sessionFile); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
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

	// Initialize Koyfin client
	client := koyfin.NewClient()

	// Load session from XDG config path
	sessionFile, err := getSessionFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding session file: %v\n", err)
		os.Exit(1)
	}

	session, err := koyfin.NewSessionFromFile(sessionFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading session: %v\n", err)
		fmt.Fprintf(os.Stderr, "Run setup script or create session file at: %s\n", sessionFile)
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
	case "search":
		runSearch(client, session, os.Args[2:])
	case "snapshot":
		runSnapshot(client, session, os.Args[2:])
	case "ticker-data":
		runTickerData(client, session, os.Args[2:])
	case "transcript":
		runTranscript(client, session, os.Args[2:])
	case "schema":
		runSchema(client, session, os.Args[2:])
	case "etf-holdings":
		runETFHoldings(client, session, os.Args[2:])
	case "screener-schema":
		runScreenerSchema(client, session, os.Args[2:])
	case "screener":
		runScreener(client, session, os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Koyfin CLI Tools

Usage: koyfin <command> [flags]

Commands:
  search                  Search for stocks/tickers by name
  snapshot                Get snapshot data for tickers
  ticker-data             Get time series data for a ticker
  transcript              Earnings call transcripts (list, get, summary)
  schema                  Get indicator schema
  etf-holdings            Get ETF holdings
  screener-schema         Get screener filter schema
  screener                Run stock screener

Examples:
  koyfin search -q "Apple"
  koyfin snapshot -kids <list_of_koyfin_ids> -category Equity
  koyfin ticker-data -kid <koyfin_id> -key "p_candle_range" -date-from "2024-01-01"
  koyfin transcript -action list -kid <koyfin_id> -limit 5
  koyfin transcript -action get -transcript-id 12345
  koyfin transcript -action summary -transcript-id 12345
  koyfin schema -asset-type Equity -indicator-type financials
  koyfin etf-holdings -kids <list_of_koyfin_ids> -category ETF
  koyfin screener-schema -asset-type Equity
  koyfin screener -filters '[{"key":"mkt","min":1000,"max":10000}]' -page-size 50

Run 'koyfin <command> -h' for more information on a command.`)
}

// findProjectRoot finds the root of the project by looking for go.mod
func findProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "."
}

// printJSON prints data as formatted JSON
func printJSON(data any) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
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
