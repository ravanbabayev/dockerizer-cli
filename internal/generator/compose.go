package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"dockerizer-cli/internal/analyzer"

	"gopkg.in/yaml.v3"
)

// ComposeConfig represents the structure of a docker-compose.yml file
type ComposeConfig struct {
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
	Networks map[string]Network `yaml:"networks,omitempty"`
	Volumes  map[string]Volume  `yaml:"volumes,omitempty"`
}

// Service represents a service in docker-compose.yml
type Service struct {
	Build       *Build        `yaml:"build,omitempty"`
	Image       string        `yaml:"image,omitempty"`
	Ports       []string      `yaml:"ports,omitempty"`
	Environment []string      `yaml:"environment,omitempty"`
	EnvFile     []string      `yaml:"env_file,omitempty"`
	Volumes     []string      `yaml:"volumes,omitempty"`
	DependsOn   []string      `yaml:"depends_on,omitempty"`
	Networks    []string      `yaml:"networks,omitempty"`
	Restart     string        `yaml:"restart,omitempty"`
	HealthCheck *HealthCheck  `yaml:"healthcheck,omitempty"`
	Deploy      *DeployConfig `yaml:"deploy,omitempty"`
}

// Build represents build configuration
type Build struct {
	Context    string            `yaml:"context"`
	Dockerfile string            `yaml:"dockerfile"`
	Args       map[string]string `yaml:"args,omitempty"`
}

// Network represents network configuration
type Network struct {
	Driver string `yaml:"driver,omitempty"`
}

// Volume represents volume configuration
type Volume struct {
	Driver string `yaml:"driver,omitempty"`
}

// HealthCheck represents healthcheck configuration
type HealthCheck struct {
	Test     []string `yaml:"test"`
	Interval string   `yaml:"interval"`
	Timeout  string   `yaml:"timeout"`
	Retries  int      `yaml:"retries"`
}

// DeployConfig represents deployment configuration
type DeployConfig struct {
	Replicas int `yaml:"replicas,omitempty"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Type     string
	Version  string
	Port     string
	Username string
	Password string
	Database string
}

// GenerateCompose creates a docker-compose.yml file
func GenerateCompose(project *analyzer.ProjectType, outputPath string) error {
	compose := &ComposeConfig{
		Version:  "3.8",
		Services: make(map[string]Service),
		Networks: map[string]Network{
			"app-network": {Driver: "bridge"},
		},
		Volumes: make(map[string]Volume),
	}

	// Add main application service
	appService := Service{
		Build: &Build{
			Context:    ".",
			Dockerfile: "Dockerfile",
		},
		Networks: []string{"app-network"},
		Restart:  "unless-stopped",
		EnvFile:  []string{".env"},
	}

	// Special handling for Laravel
	if project.Framework == "laravel" {
		appService.Volumes = []string{
			".:/var/www/html",
		}
		// Add nginx service for Laravel
		compose.Services["nginx"] = Service{
			Image: "nginx:alpine",
			Ports: []string{"80:80"},
			Volumes: []string{
				".:/var/www/html",
				"./docker/nginx/conf.d:/etc/nginx/conf.d",
			},
			Networks:  []string{"app-network"},
			DependsOn: []string{"app"},
		}

		// Create nginx config directory and configuration
		nginxConfigDir := filepath.Join(outputPath, "docker", "nginx", "conf.d")
		if err := os.MkdirAll(nginxConfigDir, 0755); err != nil {
			return fmt.Errorf("failed to create nginx config directory: %w", err)
		}

		nginxConfig := `server {
    listen 80;
    index index.php index.html;
    server_name localhost;
    error_log  /var/log/nginx/error.log;
    access_log /var/log/nginx/access.log;
    root /var/www/html/public;

    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    location ~ \.php$ {
        try_files $uri =404;
        fastcgi_split_path_info ^(.+\.php)(/.+)$;
        fastcgi_pass app:9000;
        fastcgi_index index.php;
        include fastcgi_params;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        fastcgi_param PATH_INFO $fastcgi_path_info;
    }
}`

		if err := os.WriteFile(filepath.Join(nginxConfigDir, "default.conf"), []byte(nginxConfig), 0644); err != nil {
			return fmt.Errorf("failed to create nginx configuration: %w", err)
		}
	} else if len(project.Ports) > 0 {
		appService.Ports = project.Ports
	}

	compose.Services["app"] = appService

	// Add database service if needed
	if project.Database != "" {
		dbConfig := getDefaultDBConfig(project)
		if dbConfig != nil {
			dbService := createDatabaseService(dbConfig)
			compose.Services[dbConfig.Type] = dbService
			compose.Volumes[fmt.Sprintf("%s-data", dbConfig.Type)] = Volume{Driver: "local"}

			// Update app service to depend on database
			appService := compose.Services["app"]
			appService.DependsOn = append(appService.DependsOn, dbConfig.Type)
			compose.Services["app"] = appService
		}
	}

	// Add cache service if needed
	if needsCache(project) {
		compose.Services["redis"] = createRedisService()
		appService := compose.Services["app"]
		appService.DependsOn = append(appService.DependsOn, "redis")
		compose.Services["app"] = appService
	}

	data, err := yaml.Marshal(compose)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(outputPath, "docker-compose.yml"), data, 0644)
}

func getDefaultDBConfig(project *analyzer.ProjectType) *DatabaseConfig {
	switch project.Framework {
	case "django", "flask", "fastapi":
		return &DatabaseConfig{
			Type:     "postgres",
			Version:  "13-alpine",
			Port:     "5432",
			Username: "postgres",
			Password: "postgres",
			Database: "app",
		}
	case "rails":
		return &DatabaseConfig{
			Type:     "postgres",
			Version:  "13-alpine",
			Port:     "5432",
			Username: "postgres",
			Password: "postgres",
			Database: "app",
		}
	case "laravel", "symfony":
		return &DatabaseConfig{
			Type:     "mysql",
			Version:  "8.0",
			Port:     "3306",
			Username: "root",
			Password: "root",
			Database: "app",
		}
	case "express", "nestjs":
		return &DatabaseConfig{
			Type:     "mongodb",
			Version:  "4.4",
			Port:     "27017",
			Username: "root",
			Password: "root",
			Database: "app",
		}
	default:
		return nil
	}
}

func createDatabaseService(config *DatabaseConfig) Service {
	switch config.Type {
	case "postgres":
		return Service{
			Image: fmt.Sprintf("postgres:%s", config.Version),
			Environment: []string{
				fmt.Sprintf("POSTGRES_USER=%s", config.Username),
				fmt.Sprintf("POSTGRES_PASSWORD=%s", config.Password),
				fmt.Sprintf("POSTGRES_DB=%s", config.Database),
			},
			Ports: []string{fmt.Sprintf("%s:5432", config.Port)},
			Volumes: []string{
				"postgres-data:/var/lib/postgresql/data",
			},
			Networks: []string{"app-network"},
			HealthCheck: &HealthCheck{
				Test:     []string{"CMD-SHELL", "pg_isready -U postgres"},
				Interval: "10s",
				Timeout:  "5s",
				Retries:  5,
			},
			Restart: "unless-stopped",
		}
	case "mysql":
		return Service{
			Image: fmt.Sprintf("mysql:%s", config.Version),
			Environment: []string{
				fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", config.Password),
				fmt.Sprintf("MYSQL_DATABASE=%s", config.Database),
			},
			Ports: []string{fmt.Sprintf("%s:3306", config.Port)},
			Volumes: []string{
				"mysql-data:/var/lib/mysql",
			},
			Networks: []string{"app-network"},
			HealthCheck: &HealthCheck{
				Test:     []string{"CMD", "mysqladmin", "ping", "-h", "localhost"},
				Interval: "10s",
				Timeout:  "5s",
				Retries:  5,
			},
			Restart: "unless-stopped",
		}
	case "mongodb":
		return Service{
			Image: fmt.Sprintf("mongo:%s", config.Version),
			Environment: []string{
				fmt.Sprintf("MONGO_INITDB_ROOT_USERNAME=%s", config.Username),
				fmt.Sprintf("MONGO_INITDB_ROOT_PASSWORD=%s", config.Password),
				fmt.Sprintf("MONGO_INITDB_DATABASE=%s", config.Database),
			},
			Ports: []string{fmt.Sprintf("%s:27017", config.Port)},
			Volumes: []string{
				"mongodb-data:/data/db",
			},
			Networks: []string{"app-network"},
			HealthCheck: &HealthCheck{
				Test:     []string{"CMD", "mongo", "--eval", "db.adminCommand('ping')"},
				Interval: "10s",
				Timeout:  "5s",
				Retries:  5,
			},
			Restart: "unless-stopped",
		}
	default:
		return Service{}
	}
}

func createRedisService() Service {
	return Service{
		Image: "redis:alpine",
		Ports: []string{"6379:6379"},
		Volumes: []string{
			"redis-data:/data",
		},
		Networks: []string{"app-network"},
		HealthCheck: &HealthCheck{
			Test:     []string{"CMD", "redis-cli", "ping"},
			Interval: "10s",
			Timeout:  "5s",
			Retries:  5,
		},
		Restart: "unless-stopped",
	}
}

func needsCache(project *analyzer.ProjectType) bool {
	return project.Framework == "laravel" ||
		project.Framework == "rails" ||
		project.Framework == "django" ||
		project.Framework == "nestjs"
}
