package function

import (
	"errors"
	"fmt"
	"sekai-inventory/tools"
	"strings"
	"time"
)

// Update fetches the latest card, character, and skill data from the remote database.
//
// It compares the current Git commit ID of the Sekai-World master data
// repository with the commit ID stored in the local metadata. If the commit
// has not changed, Update skips downloading new data and reports that the
// local data is already up to date. Otherwise it downloads the latest
// cards.json, gameCharacters.json, skills.json, updates metadata, and prints a summary.
//
// Progress is reported as a simple text progress bar on stdout.
func Update() {
	// Check if metadata file exists and read last update time.
	var lastUpdate time.Time
	if metadata, err := tools.LoadMetadata(); err == nil {
		if t, err := time.Parse(time.RFC3339, metadata.Timestamp); err == nil {
			lastUpdate = t
		}
	}

	// Show update start message.
	tools.PrintSuccessMessage("Starting database update...")
	if !lastUpdate.IsZero() {
		fmt.Printf("Last update was on: %s\n", lastUpdate.Format("02 Jan 2006 15:04:05"))
	}
	fmt.Println()

	// Variable to store last stage for progress display.
	var lastStage string

	// Fetch and save data with progress reporting.
	err := tools.FetchAndSaveData(func(stage string, progress float64) {
		// Only clear line and print stage if it changed.
		if stage != lastStage {
			fmt.Print("\r" + strings.Repeat(" ", 80))
			lastStage = stage
		}

		// Progress callback.
		fmt.Printf("\r%s... [", stage)
		bars := int(progress * 20)
		bar := strings.Repeat("=", bars)
		if bars < 20 {
			bar += ">"
			bar += strings.Repeat(" ", 19-bars)
		}
		fmt.Printf("%s] %.0f%%", bar, progress*100)
	})
	if err != nil {
		// Special case: data already up to date.
		if errors.Is(err, tools.ErrNoUpdateNeeded) {
			fmt.Println() // Finish the progress line if any.
			tools.PrintSuccessMessage("Your local data is already up to date. No update needed.")
			return
		}

		tools.PrintErrorMessage(fmt.Sprintf("\nUpdate failed: %v", err))
		return
	}

	fmt.Println() // Add newline after progress bar.

	// Read and display update summary.
	if metadata, err := tools.LoadMetadata(); err == nil {
		fmt.Println()
		tools.PrintSuccessMessage("Update completed successfully!")
		fmt.Printf("\nUpdate Summary:\n")
		fmt.Printf("  Cards database updated:      %s\n", tools.FormatTime(metadata.CardsLastUpdate))
		fmt.Printf("  Characters database updated: %s\n", tools.FormatTime(metadata.CharsLastUpdate))
		fmt.Printf("  Skills database updated:     %s\n", tools.FormatTime(metadata.SkillsLastUpdate))
		fmt.Printf("  Data version:                %s\n", metadata.GitCommitID[:7])
	} else {
		tools.PrintSuccessMessage("\nUpdate completed successfully!")
	}
}
