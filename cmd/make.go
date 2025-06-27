package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

//go:embed migration_template.go.tmpl
var migrationTemplate string

//go:embed add_to_migration_template.go.tmpl
var addToMigrationTemplate string

var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "Generate application components",
	Long:  `Generate various application components like migrations, models, controllers, etc.`,
}

var makeMigrationCmd = &cobra.Command{
	Use:   "migration [name]",
	Short: "Create a new migration file",
	Long:  `Create a new migration file in the database/migrations directory`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		migrationName := args[0]
		createMigration(migrationName)
	},
}

var makeMigrationAliasCmd = &cobra.Command{
	Use:   "make:migration [name]",
	Short: "Create a new migration file",
	Long:  `Create a new migration file in the database/migrations directory`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		migrationName := args[0]
		createMigration(migrationName)
	},
}

func init() {
	rootCmd.AddCommand(makeCmd)
	rootCmd.AddCommand(makeMigrationAliasCmd)
	makeCmd.AddCommand(makeMigrationCmd)
}

func createMigration(name string) {
	// Convert name to struct name (PascalCase)
	structName := toPascalCase(name)

	// Extract table name from migration name
	tableName := extractTableName(name)

	// Create timestamp prefix
	timestamp := time.Now().Format("2006_01_02_150405")

	// Create file name
	fileName := fmt.Sprintf("%s_%s.go", timestamp, name)

	// Ensure migrations directory exists
	migrationsDir := "database/migrations"
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		fmt.Printf("Error creating migrations directory: %v\n", err)
		return
	}

	// Create the migration file
	filePath := filepath.Join(migrationsDir, fileName)

	// Select template based on migration name
	selectedTemplate := migrationTemplate
	if strings.HasPrefix(name, "add_") && strings.Contains(name, "_to_") {
		selectedTemplate = addToMigrationTemplate
	}

	// Parse and execute template
	tmpl, err := template.New("migration").Parse(selectedTemplate)
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating migration file: %v\n", err)
		return
	}
	defer file.Close()

	data := struct {
		StructName string
		TableName  string
	}{
		StructName: structName,
		TableName:  tableName,
	}

	if err := tmpl.Execute(file, data); err != nil {
		fmt.Printf("Error writing migration: %v\n", err)
		return
	}

	fmt.Printf("Migration created successfully: %s\n", filePath)
}

func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}
	return strings.Join(parts, "")
}

func extractTableName(migrationName string) string {
	// Handle common patterns
	if strings.HasPrefix(migrationName, "create_") && strings.HasSuffix(migrationName, "_table") {
		// create_users_table -> users
		return strings.TrimSuffix(strings.TrimPrefix(migrationName, "create_"), "_table")
	} else if strings.HasPrefix(migrationName, "add_") && strings.Contains(migrationName, "_to_") {
		// add_column_to_users -> users
		parts := strings.Split(migrationName, "_to_")
		if len(parts) == 2 {
			if strings.HasSuffix(parts[1], "_table") {
				parts[1] = strings.TrimSuffix(parts[1], "_table")
			}
			return parts[1]
		}
	}

	// Default: use the migration name as table name
	return migrationName
}
