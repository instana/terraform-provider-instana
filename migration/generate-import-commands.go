//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TerraformState represents the structure of a terraform.tfstate file
type TerraformState struct {
	Version          int        `json:"version"`
	TerraformVersion string     `json:"terraform_version"`
	Resources        []Resource `json:"resources"`
}

// Resource represents a resource in the state file
type Resource struct {
	Module    string     `json:"module,omitempty"`
	Mode      string     `json:"mode"`
	Type      string     `json:"type"`
	Name      string     `json:"name"`
	Provider  string     `json:"provider"`
	Instances []Instance `json:"instances"`
}

// Instance represents an instance of a resource
type Instance struct {
	IndexKey      interface{}            `json:"index_key,omitempty"`
	SchemaVersion int                    `json:"schema_version"`
	Attributes    map[string]interface{} `json:"attributes"`
}

// ImportCommand represents a Terraform import CLI command
type ImportCommand struct {
	Module       string
	ResourceType string
	ResourceName string
	IndexKey     string
	ID           string
}

func main() {
	// Default values
	stateFilePath := "terraform.tfstate"
	outputFilePath := "import-commands.sh"

	// Parse command line arguments
	args := os.Args[1:]
	if len(args) > 0 {
		stateFilePath = args[0]
	}
	if len(args) > 1 {
		outputFilePath = args[1]
	}

	// Check if state file exists
	if _, err := os.Stat(stateFilePath); os.IsNotExist(err) {
		fmt.Printf("Error: State file '%s' not found.\n", stateFilePath)
		fmt.Printf("\nUsage: go run generate-import-commands.go [state_file_path] [output_file_path]\n")
		fmt.Printf("Example: go run generate-import-commands.go terraform.tfstate import-commands.sh\n")
		os.Exit(1)
	}

	fmt.Printf("Reading state file: %s\n", stateFilePath)
	fmt.Printf("Output file: %s\n\n", outputFilePath)

	// Generate import commands
	if err := generateImportCommands(stateFilePath, outputFilePath); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func generateImportCommands(stateFilePath, outputFilePath string) error {
	// Read the state file
	data, err := os.ReadFile(stateFilePath)
	if err != nil {
		return fmt.Errorf("failed to read state file: %w", err)
	}

	// Parse JSON
	var state TerraformState
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("failed to parse state file JSON: %w", err)
	}

	// Check if there are resources
	if len(state.Resources) == 0 {
		fmt.Println("No resources found in the state file.")
		return nil
	}

	// Collect import commands
	var importCommands []ImportCommand

	for _, resource := range state.Resources {
		// Only process managed resources
		if resource.Mode != "managed" {
			continue
		}

		// Only process Instana resources (resources with type starting with "instana_")
		if len(resource.Type) < 8 || resource.Type[:8] != "instana_" {
			fmt.Printf("Skipping non-Instana resource: %s.%s\n", resource.Type, resource.Name)
			continue
		}

		// Process each instance
		for _, instance := range resource.Instances {
			// Get the resource ID
			id, ok := instance.Attributes["id"]
			if !ok {
				fmt.Printf("Warning: No ID found for resource %s.%s\n", resource.Type, resource.Name)
				continue
			}

			idStr, ok := id.(string)
			if !ok {
				fmt.Printf("Warning: ID is not a string for resource %s.%s\n", resource.Type, resource.Name)
				continue
			}

			// Handle index_key (for resources with for_each or count)
			indexKey := ""
			if instance.IndexKey != nil {
				switch v := instance.IndexKey.(type) {
				case string:
					indexKey = v
				case float64:
					indexKey = fmt.Sprintf("%d", int(v))
				case int:
					indexKey = fmt.Sprintf("%d", v)
				}
			}

			importCommand := ImportCommand{
				Module:       resource.Module,
				ResourceType: resource.Type,
				ResourceName: resource.Name,
				IndexKey:     indexKey,
				ID:           idStr,
			}
			importCommands = append(importCommands, importCommand)

			// Build resource address for logging
			resourceAddr := buildResourceAddress(resource.Module, resource.Type, resource.Name, indexKey)
			fmt.Printf("Generated import command for: %s (ID: %s)\n", resourceAddr, idStr)
		}
	}

	// Write import commands to file
	if len(importCommands) == 0 {
		fmt.Println("No import commands generated.")
		return nil
	}

	// Create output file
	file, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Write shell script header
	header := `#!/bin/bash
# Terraform Import Commands
# Generated from terraform.tfstate
# 
# Usage:
#   chmod +x import-commands.sh
#   ./import-commands.sh
#
# Or run commands individually as needed
#
# Note: For module resources, you may need to be in the correct directory
# or use the -chdir flag if your modules are in subdirectories

set -e  # Exit on error

echo "Starting Terraform import process..."
echo ""

`
	if _, err := file.WriteString(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Group commands by module for better organization
	moduleGroups := make(map[string][]ImportCommand)
	for _, cmd := range importCommands {
		moduleKey := cmd.Module
		if moduleKey == "" {
			moduleKey = "root"
		}
		moduleGroups[moduleKey] = append(moduleGroups[moduleKey], cmd)
	}

	// Write import commands grouped by module
	for modulePath, commands := range moduleGroups {
		// Write module header
		if modulePath == "root" {
			if _, err := file.WriteString("# Root module resources\n"); err != nil {
				return fmt.Errorf("failed to write module header: %w", err)
			}
		} else {
			moduleHeader := fmt.Sprintf("# Module: %s\n", modulePath)
			if _, err := file.WriteString(moduleHeader); err != nil {
				return fmt.Errorf("failed to write module header: %w", err)
			}
		}

		// Write commands for this module
		for _, cmd := range commands {
			// Build the resource address
			resourceAddr := buildResourceAddress(cmd.Module, cmd.ResourceType, cmd.ResourceName, cmd.IndexKey)

			// Escape any special characters in the ID if needed
			escapedID := escapeShellString(cmd.ID)

			// Generate the terraform import command
			importCmd := fmt.Sprintf("terraform import %s %s\n", resourceAddr, escapedID)

			if _, err := file.WriteString(importCmd); err != nil {
				return fmt.Errorf("failed to write import command: %w", err)
			}
		}

		// Add blank line between module groups
		if _, err := file.WriteString("\n"); err != nil {
			return fmt.Errorf("failed to write newline: %w", err)
		}
	}

	// Write footer
	footer := `
echo ""
echo "✓ Import process completed successfully!"
echo "✓ Total resources imported: ` + fmt.Sprintf("%d", len(importCommands)) + `"
`
	if _, err := file.WriteString(footer); err != nil {
		return fmt.Errorf("failed to write footer: %w", err)
	}

	fmt.Printf("\n✓ Successfully generated %d import command(s)\n", len(importCommands))
	fmt.Printf("✓ Import commands written to: %s\n", outputFilePath)

	// Get absolute path for better user experience
	absPath, err := filepath.Abs(outputFilePath)
	if err == nil {
		fmt.Printf("✓ Full path: %s\n", absPath)
	}

	fmt.Printf("\nTo execute the import commands:\n")
	fmt.Printf("  chmod +x %s\n", outputFilePath)
	fmt.Printf("  ./%s\n", outputFilePath)
	fmt.Printf("\nOr run individual commands from the file as needed.\n")

	return nil
}

// buildResourceAddress constructs the full resource address including module path and index key
func buildResourceAddress(module, resourceType, resourceName, indexKey string) string {
	// Start with the resource type and name
	addr := fmt.Sprintf("%s.%s", resourceType, resourceName)

	// Add index key if present (for resources with for_each or count)
	if indexKey != "" {
		// Check if it's a string key (for_each) or numeric (count)
		// String keys need quotes, numeric keys don't
		if _, err := fmt.Sscanf(indexKey, "%d", new(int)); err == nil {
			// Numeric index (count)
			addr = fmt.Sprintf("%s[%s]", addr, indexKey)
		} else {
			// String index (for_each) - escape quotes in the key
			escapedKey := strings.ReplaceAll(indexKey, `"`, `\"`)
			addr = fmt.Sprintf(`%s["%s"]`, addr, escapedKey)
		}
	}

	// Prepend module path if present
	if module != "" {
		addr = fmt.Sprintf("%s.%s", module, addr)
	}

	return addr
}

// escapeShellString escapes special characters in strings for shell scripts
func escapeShellString(s string) string {
	// If the string contains spaces or special characters, quote it
	if strings.ContainsAny(s, " \t\n\"'$`\\!*?[](){};<>|&") {
		// Escape backslashes and double quotes
		s = strings.ReplaceAll(s, `\`, `\\`)
		s = strings.ReplaceAll(s, `"`, `\"`)
		return fmt.Sprintf(`"%s"`, s)
	}
	return s
}

// Made with Bob
