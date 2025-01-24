package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"dockerizer-cli/internal/analyzer"
	"dockerizer-cli/internal/generator"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "dockerizer",
		Usage: "Automatically dockerize your applications",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Value:   ".",
				Usage:   "Path to the project directory",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   ".",
				Usage:   "Output directory for generated files",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "Initialize and analyze the project",
				Action: func(c *cli.Context) error {
					projectPath := c.String("path")
					fmt.Printf("Analyzing project structure in %s...\n", projectPath)

					project, err := analyzer.AnalyzeProject(projectPath)
					if err != nil {
						return fmt.Errorf("failed to analyze project: %w", err)
					}

					fmt.Printf("\nProject Analysis Results:\n")
					fmt.Printf("Language: %s\n", project.Language)
					if project.Framework != "" {
						fmt.Printf("Framework: %s\n", project.Framework)
					}
					if len(project.Dependencies) > 0 {
						fmt.Printf("Dependencies found: %d\n", len(project.Dependencies))
					}

					return nil
				},
			},
			{
				Name:    "generate",
				Aliases: []string{"g"},
				Usage:   "Generate Dockerfile and docker-compose.yml",
				Action: func(c *cli.Context) error {
					projectPath := c.String("path")
					outputPath := c.String("output")

					// Ensure output directory exists
					if err := os.MkdirAll(outputPath, 0755); err != nil {
						return fmt.Errorf("failed to create output directory: %w", err)
					}

					// Analyze project first
					project, err := analyzer.AnalyzeProject(projectPath)
					if err != nil {
						return fmt.Errorf("failed to analyze project: %w", err)
					}

					// Generate Dockerfile
					if err := generator.GenerateDockerfile(project, outputPath); err != nil {
						return fmt.Errorf("failed to generate Dockerfile: %w", err)
					}

					// Generate docker-compose.yml
					if err := generator.GenerateCompose(project, outputPath); err != nil {
						return fmt.Errorf("failed to generate docker-compose.yml: %w", err)
					}

					fmt.Printf("\nSuccessfully generated Docker files in %s\n", outputPath)
					fmt.Println("\nNext steps:")
					fmt.Println("1. Review the generated files")
					fmt.Println("2. Build your container:")
					fmt.Printf("   cd %s && docker-compose build\n", outputPath)
					fmt.Println("3. Run your application:")
					fmt.Println("   docker-compose up")

					return nil
				},
			},
			{
				Name:    "clean",
				Aliases: []string{"c"},
				Usage:   "Remove generated Docker files",
				Action: func(c *cli.Context) error {
					outputPath := c.String("output")
					files := []string{"Dockerfile", "docker-compose.yml"}

					for _, file := range files {
						path := filepath.Join(outputPath, file)
						if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
							return fmt.Errorf("failed to remove %s: %w", file, err)
						}
					}

					fmt.Println("Successfully removed Docker files")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
