package main

import (
	"fmt"
	"os"

	"entext-applications/internal/koyfin"
)

func runTranscript(client *koyfin.Client, session *koyfin.Session, args []string) {
	fs := newFlagSet("transcript")
	action := fs.String("action", "list", "Action: list, get, summary")
	kid := fs.String("kid", "", "Koyfin identifier for the stock (required for list)")
	transcriptID := fs.Int("transcript-id", 0, "Key development identifier (required for get/summary)")
	limit := fs.Int("limit", 10, "Maximum results for list action (1-64)")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	switch *action {
	case "list":
		if *kid == "" {
			fmt.Fprintf(os.Stderr, "Error: -kid flag is required for list action\n")
			fs.Usage()
			os.Exit(1)
		}
		if *limit <= 0 || *limit > 64 {
			fmt.Fprintf(os.Stderr, "Error: -limit must be between 1 and 64\n")
			os.Exit(1)
		}

		req := koyfin.TranscriptListRequest{
			KID:   *kid,
			Limit: *limit,
		}

		transcripts, err := client.ListTranscripts(session, &req)
		if err != nil {
			exitWithError(err)
		}

		printJSON(map[string]any{"data": transcripts})

	case "get":
		if *transcriptID == 0 {
			fmt.Fprintf(os.Stderr, "Error: -transcript-id flag is required for get action\n")
			fs.Usage()
			os.Exit(1)
		}

		transcript, err := client.GetTranscript(session, *transcriptID)
		if err != nil {
			exitWithError(err)
		}

		printJSON(map[string]any{"data": transcript})

	case "summary":
		if *transcriptID == 0 {
			fmt.Fprintf(os.Stderr, "Error: -transcript-id flag is required for summary action\n")
			fs.Usage()
			os.Exit(1)
		}

		summary, err := client.GetTranscriptSummary(session, *transcriptID)
		if err != nil {
			exitWithError(err)
		}

		printJSON(map[string]any{"data": summary})

	default:
		fmt.Fprintf(os.Stderr, "Error: unknown action '%s'. Use: list, get, summary\n", *action)
		fs.Usage()
		os.Exit(1)
	}
}
