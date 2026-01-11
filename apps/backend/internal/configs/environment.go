package configs

import (
	"fmt"
	"os"
)

const (
	AppStageProd = "AppStageProd"
	AppStageDev  = "AppStageDev"

	envAppPort    = "APP_PORT"
	envAppStage   = "APP_STAGE"
	envDbHost     = "DB_HOST"
	envDbName     = "DB_NAME"
	envDbUsername = "DB_USERNAME"
	envDbPassword = "DB_PASSWORD"
	envDbPort     = "DB_PORT"
)

type database struct {
	Host     string
	Name     string
	Username string
	Password string
	Port     string
	Url      string
}

type Environment struct {
	AppPort  string
	AppStage string
	Database database
}

func NewEnvironment() Environment {
	env := Environment{
		AppPort:  getEnvOrDefault(envAppPort, "3000"),
		AppStage: getEnvOrDefault(envAppStage, AppStageDev),
		Database: database{
			Host:     getEnvOrDefault(envDbHost, "localhost"),
			Name:     getEnvOrDefault(envDbName, "user"),
			Username: getEnvOrDefault(envDbUsername, "user"),
			Password: getEnvOrDefault(envDbPassword, "user"),
			Port:     getEnvOrDefault(envDbPort, "5432"),
		},
	}

	fmt.Println("===============")
	fmt.Println("FLUXIS")
	fmt.Printf("Port: %s\n", env.AppPort)
	fmt.Printf("Stage: %s\n", env.AppStage)
	fmt.Println("===============")

	env.Database.Url = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
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
