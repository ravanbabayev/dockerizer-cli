package analyzer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

// NodePackageJSON represents package.json structure
type NodePackageJSON struct {
	Engines struct {
		Node string `json:"node"`
	} `json:"engines"`
}

// ComposerJSON represents composer.json structure
type ComposerJSON struct {
	Require    map[string]string `json:"require"`
	RequireDev map[string]string `json:"require-dev"`
	Config     struct {
		Platform struct {
			PHP string `json:"php"`
		} `json:"platform"`
	} `json:"config"`
}

// PythonConfig represents Python version configuration
type PythonConfig struct {
	PythonVersion string `json:"python_version"`
}

func detectNodeVersion(path string) (string, error) {
	packagePath := filepath.Join(path, "package.json")
	data, err := ioutil.ReadFile(packagePath)
	if err != nil {
		return "", err
	}

	var pkg NodePackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return "", err
	}

	if pkg.Engines.Node != "" {
		version := parseVersionConstraint(pkg.Engines.Node)
		if version != "" {
			return version, nil
		}
	}

	// If no version specified, use latest LTS
	return "18", nil
}

func detectPHPVersion(path string) (string, error) {
	composerPath := filepath.Join(path, "composer.json")
	data, err := ioutil.ReadFile(composerPath)
	if err != nil {
		return "", err
	}

	var composer ComposerJSON
	if err := json.Unmarshal(data, &composer); err != nil {
		return "", err
	}

	// Check PHP version requirement in require section
	if phpVersion, ok := composer.Require["php"]; ok {
		version := parseVersionConstraint(phpVersion)
		if version != "" {
			return version, nil
		}
	}

	// Check platform config
	if composer.Config.Platform.PHP != "" {
		version := parseVersionConstraint(composer.Config.Platform.PHP)
		if version != "" {
			return version, nil
		}
	}

	// If no specific version found, use latest stable version
	return getLatestPHPVersion(), nil
}

func detectPythonVersion(path string) (string, error) {
	// Check pyproject.toml
	if data, err := ioutil.ReadFile(filepath.Join(path, "pyproject.toml")); err == nil {
		re := regexp.MustCompile(`python\s*=\s*["'](\d+\.\d+)`)
		if matches := re.FindStringSubmatch(string(data)); len(matches) > 1 {
			return matches[1], nil
		}
	}

	// Check Pipfile
	if data, err := ioutil.ReadFile(filepath.Join(path, "Pipfile")); err == nil {
		re := regexp.MustCompile(`python_version\s*=\s*["'](\d+\.\d+)`)
		if matches := re.FindStringSubmatch(string(data)); len(matches) > 1 {
			return matches[1], nil
		}
	}

	// Check runtime.txt (common in Django projects)
	if data, err := ioutil.ReadFile(filepath.Join(path, "runtime.txt")); err == nil {
		re := regexp.MustCompile(`python-(\d+\.\d+)`)
		if matches := re.FindStringSubmatch(string(data)); len(matches) > 1 {
			return matches[1], nil
		}
	}

	// Default to latest stable
	return "3.9", nil
}

func detectGoVersion(path string) (string, error) {
	modPath := filepath.Join(path, "go.mod")
	data, err := ioutil.ReadFile(modPath)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`go\s+(\d+\.\d+)`)
	if matches := re.FindStringSubmatch(string(data)); len(matches) > 1 {
		return matches[1], nil
	}

	// Default to latest stable
	return "1.21", nil
}

func detectJavaVersion(path string) (string, error) {
	// Check pom.xml
	if data, err := ioutil.ReadFile(filepath.Join(path, "pom.xml")); err == nil {
		re := regexp.MustCompile(`<java.version>(\d+)</java.version>`)
		if matches := re.FindStringSubmatch(string(data)); len(matches) > 1 {
			return matches[1], nil
		}
	}

	// Check build.gradle
	if data, err := ioutil.ReadFile(filepath.Join(path, "build.gradle")); err == nil {
		re := regexp.MustCompile(`sourceCompatibility\s*=\s*['"](\d+)['"]`)
		if matches := re.FindStringSubmatch(string(data)); len(matches) > 1 {
			return matches[1], nil
		}
	}

	// Default to latest LTS
	return "17", nil
}

func parseVersionConstraint(constraint string) string {
	// Remove version operators and spaces
	constraint = strings.TrimSpace(constraint)
	operators := []string{">=", "^", "~", "<", ">", "="}
	for _, op := range operators {
		constraint = strings.TrimPrefix(constraint, op)
	}
	constraint = strings.TrimSpace(constraint)

	// Get major.minor version
	parts := strings.Split(constraint, ".")
	if len(parts) >= 2 {
		return fmt.Sprintf("%s.%s", parts[0], parts[1])
	}
	return ""
}

func getLatestPHPVersion() string {
	// Bu fonksiyon Dockerhub API'sini kullanarak
	// mevcut en son PHP sürümünü alabilir
	// Şimdilik sabit bir değer döndürelim
	return "8.3"
}

// UpdateBaseImage updates the base image according to detected version
func UpdateBaseImage(project *ProjectType) error {
	var version string
	var err error

	switch project.Language {
	case "Node.js":
		version, err = detectNodeVersion(".")
		if err == nil && version != "" {
			project.BaseImage = fmt.Sprintf("node:%s-alpine", version)
		}

	case "Python":
		version, err = detectPythonVersion(".")
		if err == nil && version != "" {
			project.BaseImage = fmt.Sprintf("python:%s-slim", version)
		}

	case "Go":
		version, err = detectGoVersion(".")
		if err == nil && version != "" {
			project.BaseImage = fmt.Sprintf("golang:%s-alpine", version)
		}

	case "Java":
		version, err = detectJavaVersion(".")
		if err == nil && version != "" {
			project.BaseImage = fmt.Sprintf("openjdk:%s-slim", version)
		}

	case "PHP":
		version, err = detectPHPVersion(".")
		if err == nil && version != "" {
			project.BaseImage = fmt.Sprintf("php:%s-fpm", version)
		}
	}

	return nil
}
