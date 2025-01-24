package analyzer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ProjectType represents the type of project detected
type ProjectType struct {
	Language     string
	Framework    string
	BaseImage    string
	Dependencies []string
	Ports        []string
	Database     string
	Environment  []string
}

// LanguageConfig represents a language configuration from YAML
type LanguageConfig struct {
	Name           string                     `yaml:"name"`
	FileIndicators []string                   `yaml:"file_indicators"`
	BaseImage      string                     `yaml:"base_image"`
	BuildFlags     []string                   `yaml:"build_flags,omitempty"`
	Frameworks     map[string]FrameworkConfig `yaml:"frameworks"`
}

// FrameworkConfig represents a framework configuration from YAML
type FrameworkConfig struct {
	Name            string   `yaml:"name"`
	Dependencies    []string `yaml:"dependencies"`
	Port            int      `yaml:"port"`
	BuildCommand    string   `yaml:"build_command,omitempty"`
	StartCommand    string   `yaml:"start_command"`
	DevCommand      string   `yaml:"dev_command,omitempty"`
	DatabaseOptions []string `yaml:"database_options,omitempty"`
	Environment     []string `yaml:"environment,omitempty"`
	FilePermissions []string `yaml:"file_permissions,omitempty"`
}

// loadLanguageConfig loads language configuration from YAML file
func loadLanguageConfig(langFile string) (*LanguageConfig, error) {
	data, err := ioutil.ReadFile(langFile)
	if err != nil {
		return nil, err
	}

	var config LanguageConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// AnalyzeProject analyzes the given directory and returns project information
func AnalyzeProject(path string) (*ProjectType, error) {
	project := &ProjectType{}

	// Get supported languages configurations
	supportedLangs, err := filepath.Glob("supported/*.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read supported languages: %w", err)
	}

	// Check project files
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// Create a map of found files for quick lookup
	foundFiles := make(map[string]bool)
	for _, file := range files {
		foundFiles[file.Name()] = true
	}

	// Check each language configuration
	for _, langFile := range supportedLangs {
		if langFile == "supported/databases.yaml" {
			continue
		}

		config, err := loadLanguageConfig(langFile)
		if err != nil {
			continue
		}

		// Check if any of the file indicators exist
		for _, indicator := range config.FileIndicators {
			if foundFiles[indicator] {
				project.Language = config.Name
				project.BaseImage = config.BaseImage

				// Update base image based on detected version
				if err := UpdateBaseImage(project); err != nil {
					// Log error but continue with default version
					fmt.Printf("Warning: Could not detect version, using default: %v\n", err)
				}

				// Detect framework
				if err := detectFramework(path, project, config); err != nil {
					return nil, err
				}

				return project, nil
			}
		}
	}

	return project, nil
}

func detectFramework(path string, project *ProjectType, config *LanguageConfig) error {
	switch project.Language {
	case "Node.js":
		return detectNodeFramework(path, project, config)
	case "Python":
		return detectPythonFramework(path, project, config)
	case "Go":
		return detectGoFramework(path, project, config)
	case "PHP":
		return detectPHPFramework(path, project, config)
	}
	return nil
}

func detectNodeFramework(path string, project *ProjectType, config *LanguageConfig) error {
	packageJSONPath := filepath.Join(path, "package.json")
	data, err := ioutil.ReadFile(packageJSONPath)
	if err != nil {
		return err
	}

	var packageJSON struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}

	if err := yaml.Unmarshal(data, &packageJSON); err != nil {
		return err
	}

	// Combine dependencies
	deps := make(map[string]string)
	for k, v := range packageJSON.Dependencies {
		deps[k] = v
	}
	for k, v := range packageJSON.DevDependencies {
		deps[k] = v
	}

	// Check each framework
	for name, framework := range config.Frameworks {
		for _, dep := range framework.Dependencies {
			if _, ok := deps[dep]; ok {
				project.Framework = name
				if framework.Port != 0 {
					project.Ports = []string{fmt.Sprintf("%d", framework.Port)}
				}
				return nil
			}
		}
	}

	return nil
}

func detectPythonFramework(path string, project *ProjectType, config *LanguageConfig) error {
	reqPath := filepath.Join(path, "requirements.txt")
	data, err := ioutil.ReadFile(reqPath)
	if err != nil {
		return err
	}

	content := string(data)
	for name, framework := range config.Frameworks {
		for _, dep := range framework.Dependencies {
			if strings.Contains(content, dep) {
				project.Framework = name
				if framework.Port != 0 {
					project.Ports = []string{fmt.Sprintf("%d", framework.Port)}
				}
				return nil
			}
		}
	}

	return nil
}

func detectGoFramework(path string, project *ProjectType, config *LanguageConfig) error {
	modPath := filepath.Join(path, "go.mod")
	data, err := ioutil.ReadFile(modPath)
	if err != nil {
		return err
	}

	content := string(data)
	for name, framework := range config.Frameworks {
		for _, dep := range framework.Dependencies {
			if strings.Contains(content, dep) {
				project.Framework = name
				if framework.Port != 0 {
					project.Ports = []string{fmt.Sprintf("%d", framework.Port)}
				}
				return nil
			}
		}
	}

	return nil
}

func detectPHPFramework(path string, project *ProjectType, config *LanguageConfig) error {
	composerPath := filepath.Join(path, "composer.json")
	data, err := ioutil.ReadFile(composerPath)
	if err != nil {
		return err
	}

	var composer struct {
		Require map[string]string `json:"require"`
	}

	if err := yaml.Unmarshal(data, &composer); err != nil {
		return err
	}

	for name, framework := range config.Frameworks {
		for _, dep := range framework.Dependencies {
			if _, ok := composer.Require[dep]; ok {
				project.Framework = name
				if framework.Port != 0 {
					project.Ports = []string{fmt.Sprintf("%d", framework.Port)}
				}
				return nil
			}
		}
	}

	return nil
}

// DetectDependencies attempts to parse dependency files
func DetectDependencies(path string) ([]string, error) {
	// TODO: Implement dependency detection based on project type
	return nil, nil
}
