package function

import (
	"fmt"
	"os"
	"sekai-inventory/model"
	"sekai-inventory/tools"
	"time"
)

// Init creates a new empty inventory file if one does not already exist.
//
// It is intended to be called before using the inventory manager for the
// first time. Init performs the following steps:
//
//  1. Check if an inventory already exists; if so, print a warning and return.
//  2. Ensure the "res" directory exists.
//  3. Initialize an empty inventory with current timestamps.
//  4. Save the inventory to disk.
//
// Any file or directory errors are reported as error messages.
func Init() {
	// Check if the inventory file already exists.
	if _, err := os.Stat(tools.InventoryFile); err == nil {
		tools.PrintWarningMessage("Inventory is already initialized.")
		return
	}

	// Ensure the "res" directory exists.
	if err := tools.EnsureResDirectory(); err != nil {
		message := fmt.Sprintf("Error: %v\n", err)
		tools.PrintErrorMessage(message)
		return
	}

	// Create an empty inventory.
	emptyInventory := &model.Inventory{
		Cards:     []model.CardEntity{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save the empty inventory to the inventory.json file.
	if err := tools.SaveInventory(emptyInventory); err != nil {
		message := fmt.Sprintf("Error initializing inventory: %v\n", err)
		tools.PrintErrorMessage(message)
		return
	}

	tools.PrintSuccessMessage("Inventory initialized successfully.")
}
