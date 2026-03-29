package tools

import (
	"encoding/json"
	"os"
	"time"
)

// readMetadata reads tools.MetadataFile and decodes it into a Metadata value.
//
// It returns an error if the file cannot be opened or decoded.
func ReadMetadata() (*Metadata, error) {
	file, err := os.Open(MetadataFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var metadata Metadata
	if err := json.NewDecoder(file).Decode(&metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// formatTime converts a timestamp string to a more readable format.
//
// It attempts several common timestamp layouts (RFC1123, RFC3339, and related
// variants). If the timestamp is empty, it returns "Not available". If no
// known format matches, it returns "Unknown format (<original>)".
func FormatTime(timestamp string) string {
	if timestamp == "" {
		return "Not available"
	}

	// Try different time formats.
	formats := []string{
		time.RFC1123,                    // Standard HTTP header format.
		"Mon, 02 Jan 2006 15:04:05 MST", // RFC1123 without GMT.
		time.RFC1123Z,                   // RFC1123 with numeric zone.
		time.RFC3339,                    // ISO8601/RFC3339.
		"2006-01-02T15:04:05Z",          // Basic ISO8601.
		time.RFC822,                     // Another common format.
		time.RFC822Z,                    // RFC822 with numeric zone.
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timestamp); err == nil {
			return t.Format("02 Jan 2006 15:04:05")
		}
	}

	// If we cannot parse the timestamp, return a user-friendly message.
	return "Unknown format (" + timestamp + ")"
}
