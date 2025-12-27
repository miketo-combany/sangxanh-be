package config

type App struct {
	Environment string `envconfig:"APP_ENV" default:"prod"`
}

// IsDevelopment checks if the app is running in development mode
func (a *App) IsDevelopment() bool {
	return a.Environment == "development" || a.Environment == "dev"
}

// IsProduction checks if the app is running in production mode
func (a *App) IsProduction() bool {
	return a.Environment == "production" || a.Environment == "prod"
}
