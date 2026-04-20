package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sekai-inventory/model"
	"sekai-inventory/tools"
)

const deprecatedInventoryFile = "res/inventory-DEPRECATED.json"

// Convert migrates an old inventory.json (without the "painting" field) to the
// new schema that includes Painting on each card.
//
// It performs the following steps:
//
//  1. Ensure the "res" directory exists.
//  2. Read res/inventory.json.
//  3. If the "painting" field is already present, do nothing.
//  4. Backup the old file as res/inventory-DEPRECATED.json.
//  5. Unmarshal the inventory, then save it using tools.SaveInventory so that
//     the Painting field is written.
//  6. Update the UpdatedAt timestamp.
//
// If no inventory file is found, Convert prints a warning and returns.
func Convert() {
	// Ensure res directory exists.
	if err := tools.EnsureResDirectory(); err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("Error ensuring res directory: %v", err))
		return
	}

	// Read current inventory.
	data, err := os.ReadFile(tools.InventoryFile)
	if err != nil {
		if os.IsNotExist(err) {
			tools.PrintWarningMessage("No inventory file found to convert. Nothing to do.")
			return
		}
		tools.PrintErrorMessage(fmt.Sprintf("Error reading inventory file: %v", err))
		return
	}

	// Already converted?
	if bytes.Contains(data, []byte(`"painting"`)) {
		tools.PrintWarningMessage("Inventory already uses the latest schema (field 'painting' found). No conversion needed.")
		return
	}

	// Backup old file.
	if err := os.WriteFile(deprecatedInventoryFile, data, 0o600); err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("Error writing backup inventory file: %v", err))
		return
	}

	// Unmarshal into the current Inventory struct.
	var inv model.Inventory
	if err := json.Unmarshal(data, &inv); err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("Error parsing inventory file: %v", err))
		return
	}

	// Save using the new schema (CardEntity now has Painting field).
	if err := tools.SaveInventory(&inv); err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("Error writing converted inventory file: %v", err))
		return
	}

	tools.PrintSuccessMessage("Inventory successfully converted to the new schema (added 'painting' field).")
	tools.PrintWarningMessage("A backup of the previous inventory has been saved as res/inventory-DEPRECATED.json.")

	if err := tools.UpdateTimeSet(); err != nil {
		tools.PrintWarningMessage(fmt.Sprintf("Warning: failed to update inventory timestamp: %v", err))
	}
}
