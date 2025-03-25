package connection

import (
	"SangXanh/pkg/config"
	"SangXanh/pkg/log"
	"fmt"
	"github.com/samber/do/v2"
	"github.com/supabase-community/supabase-go"
)

// NewSupabaseDatabase connects to the Supabase PostgreSQL database
func NewSupabaseDatabase(di do.Injector) (*supabase.Client, error) {
	conf := do.MustInvoke[config.Supabase](di) // Assuming you have a Supabase config struct

	client, err := supabase.NewClient(conf.URL, conf.KEY, &supabase.ClientOptions{})
	if err != nil {
		fmt.Println("cannot initalize client", err)
		return nil, err
	}

	log.Infow("Connected to Supabase PostgreSQL")
	return client, nil
}
