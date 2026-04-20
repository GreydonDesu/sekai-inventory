package tools

import "time"

// FormatTime converts a timestamp string to a human-readable format.
//
// It attempts several common timestamp layouts (RFC1123, RFC3339, and related
// variants). If the timestamp is empty, it returns "Not available". If no
// known format matches, it returns "Unknown format (<original>)".
func FormatTime(timestamp string) string {
	if timestamp == "" {
		return "Not available"
	}

	formats := []string{
		time.RFC1123,
		"Mon, 02 Jan 2006 15:04:05 MST",
		time.RFC1123Z,
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		time.RFC822,
		time.RFC822Z,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timestamp); err == nil {
			return t.Format("02 Jan 2006 15:04:05")
		}
	}

	return "Unknown format (" + timestamp + ")"
}
