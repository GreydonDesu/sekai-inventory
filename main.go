// Package main implements the Sekai Inventory Manager CLI, a command-line tool
// for managing a local inventory of Project SEKAI cards. It exposes commands
// such as init, update, convert, add, remove, search, list, change, and help.
package main

import (
	"fmt"
	"os"
	"sekai-inventory/function"
	"sekai-inventory/tools"
	"strings"
)

// CLI help texts used for consistent messaging across commands.
const (
	filterUsageFormat = "Usage: sekai-inventory %s [--<field> <value>]"
	filterFieldsHelp  = `Valid fields are:
  --character  (Filter by character's given name)
  --rarity     (Filter by card rarity (1, 2, 3, 4, bd))
  --group      (Filter by unit (L/N, MMJ, VBS, WxS, N25, VS))
  --painting   (Filter by painting status (true/false))`

	changeFieldsHelp = `Valid fields are:
  --level       (1-60)
  --skillLevel  (1-5)
  --masterRank  (0-5)
  --sideStory1  (true/false)
  --sideStory2  (true/false)
  --painting    (true/false)`
)

// parseCardIDArgs parses a slice of string arguments into card IDs. It prints
// a usage warning and returns (nil, false) if args is empty or any argument
// fails to parse as a positive integer.
func parseCardIDArgs(cmd string, args []string) ([]int, bool) {
	if len(args) == 0 {
		tools.PrintWarningMessage(fmt.Sprintf("Usage: sekai-inventory %s <cardID> [cardID...]", cmd))
		return nil, false
	}
	ids := make([]int, 0, len(args))
	for _, arg := range args {
		id, err := tools.ParseCardID(arg)
		if err != nil {
			tools.PrintErrorMessage(fmt.Sprintf("Invalid card ID: %s", arg))
			return nil, false
		}
		ids = append(ids, id)
	}
	return ids, true
}

// handleChangeCommand parses and executes the "change" subcommand.
//
// It expects arguments in the form:
//
//	change <cardID> --<field> <value> [--<field> <value> ...]
//
// The function validates the card ID and field/value pairs, delegates the
// update to function.Change, and prints a success message when the card
// has been updated. It returns an error if argument parsing fails or if
// function.Change reports an error.
func handleChangeCommand(args []string) error {
	if len(args) < 3 {
		tools.PrintWarningMessage("Usage: sekai-inventory change <cardID> --<field> <value> [--<field> <value> ...]")
		tools.PrintWarningMessage(changeFieldsHelp)
		return fmt.Errorf("insufficient arguments")
	}

	cardID, err := tools.ParseCardID(args[0])
	if err != nil {
		return err
	}

	updates := make(map[string]string)
	for i := 1; i < len(args)-1; i += 2 {
		field := strings.TrimPrefix(args[i], "--")
		updates[field] = args[i+1]
	}

	if err := function.Change(cardID, updates); err != nil {
		return err
	}

	tools.PrintSuccessMessage(fmt.Sprintf("Card with ID %d updated successfully.", cardID))
	return nil
}

// main is the entry point of the Sekai Inventory Manager CLI. It parses
// command-line arguments and dispatches to the corresponding command
// handlers in the function package.
func main() {
	if len(os.Args) < 2 {
		tools.PrintWarningMessage("Usage: sekai-inventory <command> [arguments]")
		tools.PrintWarningMessage("Run 'sekai-inventory help' for a list of available commands.")
		return
	}

	command := os.Args[1]

	switch command {
	case "init":
		function.Init()

	case "update":
		function.Update()

	case "convert":
		function.Convert()

	case "search", "list":
		filters := tools.ParseFilters(os.Args[2:])
		if filters == nil {
			tools.PrintWarningMessage(fmt.Sprintf(filterUsageFormat, command))
			tools.PrintWarningMessage(filterFieldsHelp)
			return
		}
		if command == "search" {
			function.Search(filters)
		} else {
			function.List(filters)
		}

	case "add":
		if ids, ok := parseCardIDArgs("add", os.Args[2:]); ok {
			function.Add(ids...)
		}

	case "remove":
		if ids, ok := parseCardIDArgs("remove", os.Args[2:]); ok {
			function.Remove(ids...)
		}

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
