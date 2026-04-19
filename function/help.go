package function

import (
	"fmt"
	"sekai-inventory/tools"

	"github.com/fatih/color"
)

// commandHelp holds the syntax, description, and optional field list for a
// single CLI command, used to render the help output.
type commandHelp struct {
	syntax      string
	description string
	fields      []fieldHelp
}

// fieldHelp describes a command field's name and its human-readable description.
type fieldHelp struct {
	name        string
	description string
}

// Help prints a colorized overview of available commands, their syntax, and
// their purpose. It is the implementation behind the "help" CLI command.
func Help() {
	// Define command documentation.
	commands := map[string]commandHelp{
		"System Commands": {
			"init", "Initialize the application by creating an empty inventory file", nil,
		},
		"Data Commands": {
			"update", "Fetch the latest data from the server and update local files", nil,
		},
		"Query Commands": {
			"search [--<field> <value>]", "List cards that are available but not in the inventory",
			[]fieldHelp{
				{"character", "Search by character's given name"},
				{"rarity", "Search by card rarity (1, 2, 3, 4, bd)"},
				{"group", "Search by unit (L/N, MMJ, VBS, WxS, N25, VS)"},
			},
		},
		"Inventory Commands": {
			"list [--<field> <value>]", "Display the inventory as a list",
			[]fieldHelp{
				{"character", "Filter by character's given name"},
				{"rarity", "Filter by card rarity (1, 2, 3, 4, bd)"},
				{"group", "Filter by unit (L/N, MMJ, VBS, WxS, N25, VS)"},
				{"painting", "Filter by painting status (true/false)"},
			},
		},
		"Card Management": {
			"add <cardID> [...]", "Add one or more cards to the inventory with default values", nil,
		},
		"": {
			"remove <cardID> [...]", "Remove one or more cards from the inventory", nil,
		},
		" ": {
			"change <cardID> --<field> <value>", "Modify fields of a card in the inventory",
			[]fieldHelp{
				{"level", "Card level (1-60)"},
				{"skillLevel", "Skill level (1-4)"},
				{"masteryRank", "Mastery rank (0-5)"},
				{"sideStory1", "Side story 1 unlock status (true/false)"},
				{"sideStory2", "Side story 2 unlock status (true/false)"},
				{"painting", "Painting status (true/false)"},
			},
		},
		"Migration": {
			"convert", "Convert existing inventory.json to the latest schema (adds 'painting' and creates a backup)", nil,
		},
		"Help": {
			"help", "Display this help message", nil,
		},
	}

	tools.PrintSuccessMessage("Project Sekai Inventory Manager - Help")
	fmt.Println()

	// Print each command section.
	lastSection := ""
	for section, cmd := range commands {
		if section != "" && section != " " {
			if lastSection != "" {
				fmt.Println()
			}
			color.New(color.FgHiCyan, color.Bold).Printf("== %s ==\n", section)
			lastSection = section
		}

		// Command syntax and description.
		bold := color.New(color.Bold)
		fmt.Printf("  ")
		bold.Printf("%-25s", cmd.syntax)
		fmt.Printf(" %s\n", cmd.description)

		// Print fields if any.
		if cmd.fields != nil {
			fmt.Println("    Available fields:")
			for _, field := range cmd.fields {
				fmt.Printf("      ")
				color.HiYellow("--%-12s", field.name)
				fmt.Printf(" %s\n", field.description)
			}
		}
	}

	// Print footer.
	fmt.Println()
	tools.PrintWarningMessage("Note: Use square brackets [] for optional parameters and <> for required parameters.")
}
