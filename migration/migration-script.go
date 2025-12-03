//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// TerraformState represents the structure of a terraform.tfstate file
type TerraformState struct {
	Version          int        `json:"version"`
	TerraformVersion string     `json:"terraform_version"`
	Resources        []Resource `json:"resources"`
}

// Resource represents a resource in the state file
type Resource struct {
	Mode      string     `json:"mode"`
	Type      string     `json:"type"`
	Name      string     `json:"name"`
	Provider  string     `json:"provider"`
	Instances []Instance `json:"instances"`
}

// Instance represents an instance of a resource
type Instance struct {
	SchemaVersion int                    `json:"schema_version"`
	Attributes    map[string]interface{} `json:"attributes"`
}

// ImportBlock represents a Terraform import block
type ImportBlock struct {
	ResourceType string
	ResourceName string
	ID           string
}

func main() {
	// Default values
	stateFilePath := "terraform.tfstate"
	outputFilePath := "import.tf"

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
		fmt.Printf("\nUsage: go run generate_import_blocks.go [state_file_path] [output_file_path]\n")
		fmt.Printf("Example: go run generate_import_blocks.go terraform.tfstate import.tf\n")
		os.Exit(1)
	}

	fmt.Printf("Reading state file: %s\n", stateFilePath)
	fmt.Printf("Output file: %s\n\n", outputFilePath)

	// Generate import blocks
	if err := generateImportBlocks(stateFilePath, outputFilePath); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func generateImportBlocks(stateFilePath, outputFilePath string) error {
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

	// Collect import blocks
	var importBlocks []ImportBlock

	for _, resource := range state.Resources {
		// Only process managed resources
		if resource.Mode != "managed" {
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

			importBlock := ImportBlock{
				ResourceType: resource.Type,
				ResourceName: resource.Name,
				ID:           idStr,
			}
			importBlocks = append(importBlocks, importBlock)

			fmt.Printf("Generated import block for: %s.%s (ID: %s)\n",
				resource.Type, resource.Name, idStr)
		}
	}

	// Write import blocks to file
	if len(importBlocks) == 0 {
		fmt.Println("No import blocks generated.")
		return nil
	}

	// Create output file
	file, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Write import blocks
	for i, block := range importBlocks {
		importBlockStr := fmt.Sprintf(`import {
  to = %s.%s
  id = "%s"
}
`, block.ResourceType, block.ResourceName, block.ID)

		if _, err := file.WriteString(importBlockStr); err != nil {
			return fmt.Errorf("failed to write import block: %w", err)
		}

		// Add newline between blocks (except for the last one)
		if i < len(importBlocks)-1 {
			if _, err := file.WriteString("\n"); err != nil {
				return fmt.Errorf("failed to write newline: %w", err)
			}
		}
	}

	fmt.Printf("\n✓ Successfully generated %d import block(s)\n", len(importBlocks))
	fmt.Printf("✓ Import blocks written to: %s\n", outputFilePath)

	// Get absolute path for better user experience
	absPath, err := filepath.Abs(outputFilePath)
	if err == nil {
		fmt.Printf("✓ Full path: %s\n", absPath)
	}

	return nil
}
