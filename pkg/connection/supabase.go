package connection

import (
	"SangXanh/pkg/config"
	"SangXanh/pkg/log"
	"github.com/nedpals/supabase-go"
	"github.com/samber/do/v2"
)

// NewSupabaseDatabase connects to the Supabase PostgreSQL database
func NewSupabaseDatabase(di do.Injector) (*supabase.Client, error) {
	conf := do.MustInvoke[config.Supabase](di) // Assuming you have a Supabase config struct

	client := supabase.CreateClient(conf.URL, conf.KEY)

	log.Infow("Connected to Supabase PostgreSQL")
	return client, nil
}
