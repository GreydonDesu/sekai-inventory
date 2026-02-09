package function

import (
	"encoding/json"
	"os"
	"sekai-inventory/tools"
	"time"
)

// readMetadata reads the metadata file and returns its contents
func readMetadata() (*tools.Metadata, error) {
	file, err := os.Open(tools.MetadataFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var metadata tools.Metadata
	if err := json.NewDecoder(file).Decode(&metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// formatTime converts a timestamp to a more readable format.
// Handles various timestamp formats including RFC1123, RFC3339, and ISO8601.
// Returns "Not available" if the timestamp is empty or cannot be parsed.
func formatTime(timestamp string) string {
	if timestamp == "" {
		return "Not available"
	}

	// Try different time formats
	formats := []string{
		time.RFC1123,      // Standard HTTP header format
		"Mon, 02 Jan 2006 15:04:05 MST",  // RFC1123 without GMT
		time.RFC1123Z,     // RFC1123 with numeric zone
		time.RFC3339,      // ISO8601/RFC3339
		"2006-01-02T15:04:05Z",           // Basic ISO8601
		time.RFC822,       // Another common format
		time.RFC822Z,      // RFC822 with numeric zone
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timestamp); err == nil {
			return t.Format("02 Jan 2006 15:04:05")
		}
	}

	// If we can't parse the timestamp, return a user-friendly message
	return "Unknown format (" + timestamp + ")"
}
