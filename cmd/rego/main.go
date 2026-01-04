package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rebolo",
	Short: "ReboloLang - A modern Go web framework inspired by Rebolo, Barranquilla",
	Long:  `ReboloLang is a batteries-included web framework for Go with Bun.js asset pipeline, hot reload, and Rails-like conventions.`,
}

var newCmd = &cobra.Command{
	Use:   "new [app-name]",
	Short: "Create a new ReboloLang application",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		fmt.Printf("Creating new ReboloLang app: %s\n", appName)
		generateApp(appName)
	},
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start development server with hot reload",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting ReboloLang development server...")
		startDevServer()
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate resources, models, controllers",
	Aliases: []string{"g"},
}

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database operations",
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running database migrations...")
		runMigrations()
	},
}

var resourceCmd = &cobra.Command{
	Use:   "resource [name] [fields...]",
	Short: "Generate a complete resource (model, controller, views, routes)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resourceName := args[0]
		fields := args[1:]
		fmt.Printf("Generating resource: %s with fields: %v\n", resourceName, fields)
		generateResource(resourceName, fields)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(devCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(dbCmd)
	
	generateCmd.AddCommand(resourceCmd)
	dbCmd.AddCommand(migrateCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Functions implemented in separate files
