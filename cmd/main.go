package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"dockerizer-cli/internal/analyzer"
	"dockerizer-cli/internal/generator"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func main() {
	app := &cli.App{
		Name:  "dockerizer",
		Usage: "Automatically dockerize your applications",
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "Initialize and analyze the project",
				Action: func(c *cli.Context) error {
					fmt.Println("🔍 Analyzing project structure...")

					// Analyze project
					project, err := analyzer.AnalyzeProject(".")
					if err != nil {
						return fmt.Errorf("failed to analyze project: %w", err)
					}

					if project.Language == "" {
						fmt.Println("❌ Could not automatically detect the project language.")
						return selectLanguageManually(project)
					}

					// Confirm language detection
					fmt.Printf("✨ Detected %s project\n", project.Language)
					fmt.Print("Is this correct? [Y/n]: ")
					var response string
					fmt.Scanln(&response)
					if response != "" && strings.ToLower(response) != "y" {
						return selectLanguageManually(project)
					}

					// Framework detection
					if project.Framework == "" {
						fmt.Println("❌ Could not automatically detect the framework.")
						return selectFrameworkManually(project)
					}

					fmt.Printf("✨ Detected %s framework\n", project.Framework)
					fmt.Print("Is this correct? [Y/n]: ")
					fmt.Scanln(&response)
					if response != "" && strings.ToLower(response) != "y" {
						return selectFrameworkManually(project)
					}

					// Port configuration
					if len(project.Ports) > 0 {
						defaultPort := project.Ports[0]
						fmt.Printf("✨ Default port for %s is %s\n", project.Framework, defaultPort)
						fmt.Print("Would you like to use a different port? [y/N]: ")
						fmt.Scanln(&response)
						if strings.ToLower(response) == "y" {
							var portStr string
							for {
								fmt.Print("Enter port number: ")
								fmt.Scanln(&portStr)
								port, err := strconv.Atoi(portStr)
								if err != nil || port < 1 || port > 65535 {
									fmt.Println("Please enter a valid port number between 1 and 65535")
									continue
								}
								project.Ports[0] = portStr
								break
							}
						}
					}

					// Database selection
					fmt.Print("Does your project need a database? [Y/n]: ")
					fmt.Scanln(&response)
					if response == "" || strings.ToLower(response) == "y" {
						// Load database options from config
						dbConfig, err := loadDatabaseConfig()
						if err != nil {
							return err
						}

						var dbOptions []string
						for dbName := range dbConfig.Databases {
							dbOptions = append(dbOptions, dbName)
						}

						fmt.Println("\nAvailable databases:")
						for i, db := range dbOptions {
							fmt.Printf("%d) %s\n", i+1, db)
						}

						var dbIndex int
						for {
							fmt.Print("Select database (enter number): ")
							fmt.Scanln(&response)
							index, err := strconv.Atoi(response)
							if err != nil || index < 1 || index > len(dbOptions) {
								fmt.Println("Please enter a valid number")
								continue
							}
							dbIndex = index - 1
							break
						}

						dbType := dbOptions[dbIndex]
						project.Database = dbType

						// Ask for database port
						dbInfo := dbConfig.Databases[dbType]
						defaultDBPort := fmt.Sprintf("%d", dbInfo.Port)
						fmt.Printf("✨ Default port for %s is %s\n", dbType, defaultDBPort)
						fmt.Print("Would you like to use a different port? [y/N]: ")
						fmt.Scanln(&response)
						if strings.ToLower(response) == "y" {
							var portStr string
							for {
								fmt.Print("Enter database port number: ")
								fmt.Scanln(&portStr)
								port, err := strconv.Atoi(portStr)
								if err != nil || port < 1 || port > 65535 {
									fmt.Println("Please enter a valid port number between 1 and 65535")
									continue
								}
								dbInfo.Port = port
								dbConfig.Databases[dbType] = dbInfo
								break
							}
						}
					}

					fmt.Println("\n📦 Generating Docker files...")

					// Generate Dockerfile
					if err := generator.GenerateDockerfile(project, "."); err != nil {
						fmt.Printf("⚠️  Warning: %v\n", err)
						fmt.Println("Continuing with docker-compose.yml generation...")
					} else {
						fmt.Println("✅ Successfully generated Dockerfile")
					}

					// Generate docker-compose.yml
					if err := generator.GenerateCompose(project, "."); err != nil {
						return fmt.Errorf("failed to generate docker-compose.yml: %w", err)
					}
					fmt.Println("✅ Successfully generated docker-compose.yml")

					fmt.Println("\nNext steps:")
					fmt.Println("1. Review the generated files")
					fmt.Println("2. Build and run your container:")
					fmt.Println("   docker-compose up --build")

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func selectLanguageManually(project *analyzer.ProjectType) error {
	// Get available languages from config files
	files, err := filepath.Glob("supported/*.yaml")
	if err != nil {
		return fmt.Errorf("failed to read supported languages: %w", err)
	}

	var languages []string
	for _, file := range files {
		if file == "supported/databases.yaml" {
			continue
		}

		config, err := loadLanguageConfig(file)
		if err != nil {
			continue
		}
		languages = append(languages, config.Name)
	}

	selectPrompt := promptui.Select{
		Label: "Select your project's language",
		Items: languages,
	}

	_, language, err := selectPrompt.Run()
	if err != nil {
		return fmt.Errorf("language selection failed: %w", err)
	}

	project.Language = language
	return selectFrameworkManually(project)
}

func selectFrameworkManually(project *analyzer.ProjectType) error {
	// Find the language config file
	files, err := filepath.Glob("supported/*.yaml")
	if err != nil {
		return fmt.Errorf("failed to read supported languages: %w", err)
	}

	var frameworks []string
	for _, file := range files {
		if file == "supported/databases.yaml" {
			continue
		}

		config, err := loadLanguageConfig(file)
		if err != nil {
			continue
		}

		if config.Name == project.Language {
			for name := range config.Frameworks {
				frameworks = append(frameworks, name)
			}
			break
		}
	}

	selectPrompt := promptui.Select{
		Label: "Select your project's framework",
		Items: frameworks,
	}

	_, framework, err := selectPrompt.Run()
	if err != nil {
		return fmt.Errorf("framework selection failed: %w", err)
	}

	project.Framework = framework
	return nil
}

type DatabaseConfig struct {
	Databases map[string]struct {
		Name        string   `yaml:"name"`
		Image       string   `yaml:"image"`
		Port        int      `yaml:"port"`
		Environment []string `yaml:"environment"`
	} `yaml:"databases"`
}

func loadDatabaseConfig() (*DatabaseConfig, error) {
	data, err := ioutil.ReadFile("supported/databases.yaml")
	if err != nil {
		return nil, err
	}

	var config DatabaseConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func loadLanguageConfig(file string) (*analyzer.LanguageConfig, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var config analyzer.LanguageConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
