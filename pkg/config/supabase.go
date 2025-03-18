package config

type Supabase struct {
	URL string `envconfig:"DATABASE_URL" required:"true"`
	KEY string `envconfig:"DATABASE_KEY" required:"true"`
}
