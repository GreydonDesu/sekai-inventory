package main

import (
	"fmt"
	"os"
	"sekai-inventory/function"
	"sekai-inventory/tools"
	"strings"
)

// Common help texts for consistent messaging across commands
const (
	filterUsageFormat = "Usage: sekai-inventory %s [--<field> <value>]"
	filterFieldsHelp  = `Valid fields are:
  --character  (Filter by character's given name)
  --rarity     (Filter by card rarity (1, 2, 3, 4, bd))
  --group      (Filter by unit (L/N, MMJ, VBS, WxS, N25, VS))`
	
	changeFieldsHelp = `Valid fields are:
  --level       (1-60)
  --skillLevel  (1-5)
  --masterRank  (0-5)
  --sideStory1  (true/false)
  --sideStory2  (true/false)`
)

// handleChangeCommand processes the change command and its arguments
func handleChangeCommand(args []string) error {
	if len(args) < 3 {
		tools.PrintWarningMessage("Usage: sekai-inventory change <cardID> --<field> <value> [--<field> <value> ...]")
		tools.PrintWarningMessage(changeFieldsHelp)
		return fmt.Errorf("insufficient arguments")
	}

	// Parse the card ID
	cardID, err := tools.ParseCardID(args[0])
	if err != nil {
		return err
	}

	// Parse the fields and values into a map
	updates := make(map[string]string)
	for i := 1; i < len(args)-1; i += 2 {
		field := strings.TrimPrefix(args[i], "--")
		if i+1 >= len(args) {
			return fmt.Errorf("missing value for field '%s'", field)
		}
		updates[field] = args[i+1]
	}

	if err := function.Change(cardID, updates); err != nil {
		return err
	}

	tools.PrintSuccessMessage(fmt.Sprintf("Card with ID %d updated successfully.", cardID))
	return nil
}

func main() {
	if len(os.Args) < 2 {
		tools.PrintWarningMessage("Usage: sekai-inventory <command> [arguments]")
		tools.PrintWarningMessage("Run 'sekai-inventory help' for a list of available commands.")
		return
	}

	// Parse the command
	command := os.Args[1]

	// Handle commands
	switch command {

	case "init":
		function.Init()

	case "update":
		function.Update()

	case "search", "list":
		// Parse filters for search/list commands
		filters := tools.ParseFilters(os.Args[2:])
		if filters == nil {
			// Invalid usage
			tools.PrintWarningMessage(fmt.Sprintf(filterUsageFormat, command))
			tools.PrintWarningMessage(filterFieldsHelp)
			return
		}

		// Call the appropriate function
		if command == "search" {
			function.Search(filters)
		} else {
			function.List(filters)
		}

	case "add":
		if len(os.Args) < 3 {
			tools.PrintWarningMessage("Usage: sekai-inventory add <cardID> [cardID...]")
			return
		}

		// Parse all cardIDs
		var cardIDs []int
		for _, arg := range os.Args[2:] {
			cardID, err := tools.ParseCardID(arg)
			if err != nil {
				tools.PrintErrorMessage(fmt.Sprintf("Invalid card ID: %s", arg))
				return
			}
			cardIDs = append(cardIDs, cardID)
		}

		// Call the Add function with all parsed cardIDs
		function.Add(cardIDs...)

	case "remove":
		if len(os.Args) < 3 {
			tools.PrintWarningMessage("Usage: sekai-inventory remove <cardID> [cardID...]")
			return
		}

		// Parse all cardIDs
		var cardIDs []int
		for _, arg := range os.Args[2:] {
			cardID, err := tools.ParseCardID(arg)
			if err != nil {
				tools.PrintErrorMessage(fmt.Sprintf("Invalid card ID: %s", arg))
				return
			}
			cardIDs = append(cardIDs, cardID)
		}

		// Call the Remove function with all parsed cardIDs
		function.Remove(cardIDs...)

	case "change":
		if err := handleChangeCommand(os.Args[2:]); err != nil {
			tools.PrintErrorMessage(err.Error())
		}

	case "help":
		function.Help()

	default:
		tools.PrintErrorMessage("Unknown command: " + command)
		tools.PrintWarningMessage("Run 'sekai-inventory help' for a list of available commands.")
	}
}
