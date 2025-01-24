package analyzer

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// ProjectType represents the type of project detected
type ProjectType struct {
	Language     string
	Framework    string
	BaseImage    string
	Dependencies []string
	Ports        []string
}

// PackageJSON represents a Node.js package.json file
type PackageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Scripts         map[string]string `json:"scripts"`
}

// AnalyzeProject analyzes the given directory and returns project information
func AnalyzeProject(path string) (*ProjectType, error) {
	project := &ProjectType{}

	// Check for common project files
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fileName := file.Name()
		switch {
		case fileName == "package.json":
			project.Language = "nodejs"
			project.BaseImage = "node:18-alpine"
			if err := detectNodeFramework(path, project); err != nil {
				return nil, err
			}

		case fileName == "requirements.txt":
			project.Language = "python"
			project.BaseImage = "python:3.9-slim"
			if err := detectPythonFramework(path, project); err != nil {
				return nil, err
			}

		case fileName == "go.mod":
			project.Language = "go"
			project.BaseImage = "golang:1.21-alpine"
			if err := detectGoFramework(path, project); err != nil {
				return nil, err
			}

		case fileName == "pom.xml":
			project.Language = "java"
			project.BaseImage = "openjdk:17-slim"
			project.Framework = "spring-boot" // Assuming Spring Boot for now

		case fileName == "composer.json":
			project.Language = "php"
			project.BaseImage = "php:8.2-fpm"
			if err := detectPHPFramework(path, project); err != nil {
				return nil, err
			}

		case fileName == "Gemfile":
			project.Language = "ruby"
			project.BaseImage = "ruby:3.2-alpine"
			if err := detectRubyFramework(path, project); err != nil {
				return nil, err
			}
		}
	}

	return project, nil
}

func detectNodeFramework(path string, project *ProjectType) error {
	packageJSONPath := filepath.Join(path, "package.json")
	data, err := ioutil.ReadFile(packageJSONPath)
	if err != nil {
		return err
	}

	var packageJSON PackageJSON
	if err := json.Unmarshal(data, &packageJSON); err != nil {
		return err
	}

	// Check for common frameworks in dependencies
	deps := make(map[string]string)
	for k, v := range packageJSON.Dependencies {
		deps[k] = v
	}
	for k, v := range packageJSON.DevDependencies {
		deps[k] = v
	}

	switch {
	case deps["next"] != "":
		project.Framework = "nextjs"
		project.Ports = []string{"3000"}
	case deps["react"] != "":
		project.Framework = "react"
		project.Ports = []string{"3000"}
	case deps["@angular/core"] != "":
		project.Framework = "angular"
		project.Ports = []string{"4200"}
	case deps["express"] != "":
		project.Framework = "express"
		project.Ports = []string{"3000"}
	case deps["@nestjs/core"] != "":
		project.Framework = "nestjs"
		project.Ports = []string{"3000"}
	}

	project.Dependencies = make([]string, 0, len(deps))
	for dep := range deps {
		project.Dependencies = append(project.Dependencies, dep)
	}

	return nil
}

func detectPythonFramework(path string, project *ProjectType) error {
	reqPath := filepath.Join(path, "requirements.txt")
	data, err := ioutil.ReadFile(reqPath)
	if err != nil {
		return err
	}

	requirements := strings.Split(string(data), "\n")
	project.Dependencies = requirements

	for _, req := range requirements {
		req = strings.ToLower(strings.TrimSpace(req))
		switch {
		case strings.HasPrefix(req, "django"):
			project.Framework = "django"
			project.Ports = []string{"8000"}
		case strings.HasPrefix(req, "flask"):
			project.Framework = "flask"
			project.Ports = []string{"5000"}
		case strings.HasPrefix(req, "fastapi"):
			project.Framework = "fastapi"
			project.Ports = []string{"8000"}
		}
	}

	return nil
}

func detectGoFramework(path string, project *ProjectType) error {
	modPath := filepath.Join(path, "go.mod")
	data, err := ioutil.ReadFile(modPath)
	if err != nil {
		return err
	}

	content := string(data)
	switch {
	case strings.Contains(content, "github.com/gin-gonic/gin"):
		project.Framework = "gin"
		project.Ports = []string{"8080"}
	case strings.Contains(content, "github.com/gofiber/fiber"):
		project.Framework = "fiber"
		project.Ports = []string{"3000"}
	case strings.Contains(content, "github.com/labstack/echo"):
		project.Framework = "echo"
		project.Ports = []string{"1323"}
	}

	return nil
}

func detectPHPFramework(path string, project *ProjectType) error {
	composerPath := filepath.Join(path, "composer.json")
	data, err := ioutil.ReadFile(composerPath)
	if err != nil {
		return err
	}

	var composer struct {
		Require map[string]string `json:"require"`
	}
	if err := json.Unmarshal(data, &composer); err != nil {
		return err
	}

	for dep := range composer.Require {
		switch {
		case strings.Contains(dep, "laravel"):
			project.Framework = "laravel"
			project.Ports = []string{"8000"}
		case strings.Contains(dep, "symfony"):
			project.Framework = "symfony"
			project.Ports = []string{"8000"}
		}
	}

	return nil
}

func detectRubyFramework(path string, project *ProjectType) error {
	gemfilePath := filepath.Join(path, "Gemfile")
	data, err := ioutil.ReadFile(gemfilePath)
	if err != nil {
		return err
	}

	content := string(data)
	if strings.Contains(content, "rails") {
		project.Framework = "rails"
		project.Ports = []string{"3000"}
	}

	return nil
}

// DetectDependencies attempts to parse dependency files
func DetectDependencies(path string) ([]string, error) {
	// TODO: Implement dependency detection based on project type
	return nil, nil
}
