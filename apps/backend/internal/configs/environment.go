package configs

import (
	"fmt"
	"os"
)

const (
	AppStageProd = "AppStageProd"
	AppStageDev  = "AppStageDev"

	envAppPort       = "APP_PORT"
	envIsDevEnv      = "IS_DEV_ENV"
	envDbHost        = "DB_HOST"
	envDbName        = "DB_NAME"
	envDbUsername    = "DB_USER"
	envDbPassword    = "DB_PASSWORD"
	envDbPort        = "DB_PORT"
	envAdminUser     = "ADMIN_USERNAME"
	envAdminPassword = "ADMIN_PASSWORD"
)

type database struct {
	Host     string
	Name     string
	Username string
	Password string
	Port     string
	Url      string
}

type admin struct {
	Username, Password string
}

type Environment struct {
	AppPort  string
	AppStage string
	Database database
	Admin    admin
}

func NewEnvironment() Environment {
	isDevStr := os.Getenv(envIsDevEnv)

	var appStage string
	if isDevStr == "true" {
		appStage = AppStageDev
	} else {
		appStage = AppStageProd
	}
	env := Environment{
		AppPort:  getEnvOrDefault(envAppPort, "3000"),
		AppStage: appStage,
		Database: database{
			Host:     getEnvOrDefault(envDbHost, "localhost"),
			Name:     getEnvOrDefault(envDbName, "user"),
			Username: getEnvOrDefault(envDbUsername, "user"),
			Password: getEnvOrDefault(envDbPassword, "user"),
			Port:     getEnvOrDefault(envDbPort, "5432"),
		},
		Admin: admin{
			Username: getEnvOrDefault(envAdminUser, "admin"),
			Password: getEnvOrDefault(envAdminPassword, "password"),
		},
	}

	fmt.Println("===============")
	fmt.Println("Server:")
	fmt.Printf("Port: %s\n", env.AppPort)
	fmt.Printf("Stage: %s\n", env.AppStage)
	fmt.Println("Credential:")
	fmt.Printf("Username: %s\n", env.Admin.Username)
	fmt.Printf("Password: %s\n", env.Admin.Password)
	fmt.Println("===============")

	env.Database.Url = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.Database.Username,
		env.Database.Password,
		env.Database.Host,
		env.Database.Port,
		env.Database.Name,
	)

	return env
}

func getEnvOrDefault(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
