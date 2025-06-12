package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/invopop/jsonschema"

	"github.com/docvibe-ai/api/go/invoice"
)

func main() {
	err := createSchema(
		invoice.Invoice{},
		"github.com/docvibe-ai/api/go/invoice", // importPath
		"invoice.schema.json",                  // schemaFilename
	)
	if err != nil {
		fmt.Println(err)
	}
	err = createSchema(
		invoice.AccountingInvoice{},
		"github.com/docvibe-ai/api/go/invoice", // importPath
		"accounting-invoice.schema.json",       // schemaFilename
	)
	if err != nil {
		fmt.Println(err)
	}
}

func createSchema(val any, importPath, schemaFilename string) error {
	// Get absolute path to repo directory before changing directory
	repoDir, err := filepath.Abs("../..")
	if err != nil {
		return fmt.Errorf("failed to get absolute path of repo directory: %w", err)
	}

	// reflector.AddGoComments needs to be called from the directory of the reflected package
	packageDir := filepath.Join(repoDir, strings.TrimPrefix(importPath, "github.com/docvibe-ai/api/"))
	err = os.Chdir(packageDir)
	if err != nil {
		return fmt.Errorf("failed to change directory to package path %s: %w", packageDir, err)
	}
	reflector := &jsonschema.Reflector{
		Anonymous:      true,
		ExpandedStruct: true,
		DoNotReference: true,
	}
	err = reflector.AddGoComments(importPath, ".")
	if err != nil {
		return fmt.Errorf("failed to parse Go comments: %w", err)
	}

	schema := reflector.Reflect(val)
	schema.ID = jsonschema.ID("https://raw.githubusercontent.com/docvibe-ai/api/refs/heads/master/schema/" + schemaFilename)
	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to generate schema JSON: %w", err)
	}

	outputFile := filepath.Join(repoDir, "schema", schemaFilename)
	err = os.WriteFile(outputFile, schemaJSON, 0644)
	if err != nil {
		return fmt.Errorf("failed to write schema: %w", err)
	}
	fmt.Println("Schema written to", outputFile)
	return nil
}
