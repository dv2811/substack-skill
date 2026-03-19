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

// getSessionFile returns the path to the session file in the binary's directory
func getSessionFile() (string, error) {
	// Check environment variable first (for testing/custom setups)
	if custom := os.Getenv("KOYFIN_SESSION_FILE"); custom != "" {
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
func saveSession(session *koyfin.Session, sessionFile string) error {
	// Save session to file (includes updated tokens, timestamps, cookies)
	if err := session.SaveToFile(sessionFile); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		printUsage(koyfinCLIUsage)
		os.Exit(1)
	}

	command := os.Args[1]

	// Help should work without session
	if command == "help" || command == "-h" || command == "--help" {
		printUsage(koyfinCLIUsage)
		os.Exit(0)
	}

	// Check for help flag in command args before loading session
	for _, arg := range os.Args[2:] {
		if arg == "-h" || arg == "--help" {
			// Route to command help without session
			switch command {
			case "search":
				// search flags will handle its own help
			case "snapshot":
				// snapshot flags will handle its own help
			case "ticker-data":
				// ticker-data flags will handle its own help
			case "transcript":
				// transcript flags will handle its own help
			case "schema":
				// schema flags will handle its own help
			case "etf-holdings":
				// etf-holdings flags will handle its own help
			case "screener-schema":
				// screener-schema flags will handle its own help
			case "screener":
				// screener flags will handle its own help
			case "auth":
				printUsage(koyfinCLIUsage)
				fmt.Println("\nAuth Command:")
				fmt.Println(authCmdHelp)
				os.Exit(0)
			}
			break
		}
	}

	// Initialize Koyfin client
	client := koyfin.NewClient()

	// Load session from binary directory
	sessionFile, err := getSessionFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding session file: %v\n", err)
		os.Exit(1)
	}

	session, err := koyfin.NewSessionFromFile(sessionFile)
	if err != nil {
		// auth command: create or renew a session
		if command == "auth" {
			session = &koyfin.Session{}
		} else {
			fmt.Fprintf(os.Stderr, "Error loading session: %v\n", err)
			fmt.Fprintf(os.Stderr, "Run 'koyfin auth' to authenticate\n")
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
	case "auth":
		runAuth(client, session, os.Args[2:])
	case "help", "-h", "--help":
		printUsage(koyfinCLIUsage)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage(koyfinCLIUsage)
		os.Exit(1)
	}
}

const koyfinCLIUsage = `
Koyfin CLI Tools

Usage: koyfin <command> [flags]

Commands:
  search          Search for stocks/tickers by name
  snapshot        Get snapshot data for tickers
  ticker-data     Get time series data for a ticker
  transcript      Earnings call transcripts (list, get, summary)
  schema          Get indicator schema
  etf-holdings    Get ETF holdings
  screener-schema Get screener filter schema
  screener        Run stock screener
  auth            Authenticate with Koyfin (email/password)

Examples:
  koyfin auth -email <user email> -password <password for Koyfin account>
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

Run 'koyfin <command> -h' for more information on a command.`

// common utility for showing tool usage
func printUsage(usage string) {
	fmt.Println(usage)
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
