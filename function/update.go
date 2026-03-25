package function

import (
	"fmt"
	"sekai-inventory/tools"
	"strings"
	"time"
)

// Update fetches the latest card and character data from the game's database.
// This function should be called periodically to ensure the local data is up to date,
// especially when new cards are released in the game.
func Update() {
	// Check if metadata file exists and read last update time
	var lastUpdate time.Time
	if metadata, err := readMetadata(); err == nil {
		if t, err := time.Parse(time.RFC3339, metadata.Timestamp); err == nil {
			lastUpdate = t
		}
	}

	// Show update start message
	tools.PrintSuccessMessage("Starting database update...")
	if !lastUpdate.IsZero() {
		fmt.Printf("Last update was on: %s\n", lastUpdate.Format("02 Jan 2006 15:04:05"))
	}
	fmt.Println()

	// Variable to store last stage for progress display
	var lastStage string

	// Fetch and save data with progress reporting
	if err := tools.FetchAndSaveData(func(stage string, progress float64) {
		// Only clear line and print stage if it changed
		if stage != lastStage {
			fmt.Print("\r" + strings.Repeat(" ", 80))
			lastStage = stage
		}

		// Progress callback
		fmt.Printf("\r%s... [", stage)
		bars := int(progress * 20)
		bar := strings.Repeat("=", bars)
		if bars < 20 {
			bar += ">"
			bar += strings.Repeat(" ", 19-bars)
		}
		fmt.Printf("%s] %.0f%%", bar, progress*100)
	}); err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("\nUpdate failed: %v", err))
		return
	}
	fmt.Println() // Add newline after progress bar

	// Read and display update summary
	if metadata, err := readMetadata(); err == nil {
		fmt.Println()
		tools.PrintSuccessMessage("Update completed successfully!")
		fmt.Printf("\nUpdate Summary:\n")
		fmt.Printf("  Cards database updated:      %s\n", formatTime(metadata.CardsLastUpdate))
		fmt.Printf("  Characters database updated: %s\n", formatTime(metadata.CharsLastUpdate))
		fmt.Printf("  Data version:                %s\n", metadata.GitCommitID[:7])
	} else {
		tools.PrintSuccessMessage("\nUpdate completed successfully!")
	}
}
