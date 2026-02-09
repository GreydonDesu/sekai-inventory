package function

import (
	"fmt"
	"os"
	"sekai-inventory/model"
	"sekai-inventory/tools"
	"time"
)

// Init creates a new empty inventory file if one doesn't already exist.
// This function should be called before using the inventory manager for the first time.
// It performs the following steps:
//  1. Checks if an inventory already exists
//  2. Creates the resources directory if needed
//  3. Initializes an empty inventory with current timestamps
//  4. Saves the inventory to disk
func Init() {
	// Check if the inventory file already exists
	if _, err := os.Stat(tools.InventoryFile); err == nil {
		tools.PrintWarningMessage("Inventory is already initialized.")
		return
	}

	// Ensure the "res" directory exists
	if err := tools.EnsureResDirectory(); err != nil {
		message := fmt.Sprintf("Error: %v\n", err)
		tools.PrintErrorMessage(message)
		return
	}

	// Create an empty inventory
	emptyInventory := &model.Inventory{
		Cards:     []model.CardEntity{}, // Initialize with an empty slice of cards
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save the empty inventory to the inventory.json file
	if err := tools.SaveInventory(emptyInventory); err != nil {
		message := fmt.Sprintf("Error initializing inventory: %v\n", err)
		tools.PrintErrorMessage(message)
		return
	}

	tools.PrintSuccessMessage("Inventory initialized successfully.")
}
