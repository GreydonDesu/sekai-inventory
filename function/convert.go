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

// Convert migrates an old inventory.json (without the "painting" field)
// to the new schema.
//
// Steps:
//  1. Read res/inventory.json
//  2. If "painting" is already present, do nothing
//  3. Backup the old file as res/inventory-DEPRECATED.json
//  4. Unmarshal into model.Inventory (Painting defaults to false)
//  5. Save the inventory again via tools.SaveInventory (now with "painting")
//  6. Update the updated_at timestamp
func Convert() {
	// Ensure res directory exists
	if err := tools.EnsureResDirectory(); err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("Error ensuring res directory: %v", err))
		return
	}

	// Read current inventory
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

	// Backup old file
	if err := os.WriteFile(deprecatedInventoryFile, data, 0644); err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("Error writing backup inventory file: %v", err))
		return
	}

	// Unmarshal into the current Inventory struct
	var inv model.Inventory
	if err := json.Unmarshal(data, &inv); err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("Error parsing inventory file: %v", err))
		return
	}

	// Save using the new schema (CardEntity now has 'Painting' field)
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
