package configs

import (
	"github.com/danielgtaylor/huma/v2"
)

func GetOpenapiConfig(env Environment) huma.Config {
	config := huma.DefaultConfig("Fluxis", "1.0.0")
	url := "http://localhost:" + env.AppPort
	desc := "Development server"
	if env.AppStage == AppStageProd {
		url = "/api"
		desc = "Proxied server"
	}
	config.Servers = []*huma.Server{{URL: url, Description: desc}}

	config.CreateHooks = []func(huma.Config) huma.Config{}
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearer": {Type: "http", Scheme: "bearer", BearerFormat: "JWT"},
	}

	return config
}
