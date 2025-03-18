package connection

import (
	"SangXanh/pkg/config"
	"SangXanh/pkg/log"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do/v2"
	"time"
)

// NewSupabaseDatabase connects to the Supabase PostgreSQL database
func NewSupabaseDatabase(di do.Injector) (*pgxpool.Pool, error) {
	ctx := context.Background()
	conf := do.MustInvoke[config.Supabase](di) // Assuming you have a Supabase config struct

	dbpool, err := pgxpool.New(ctx, conf.URL)
	if err != nil {
		log.Errorw("failed to connect to supabase", "error", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := dbpool.Ping(ctx); err != nil {
		log.Errorw("failed to ping supabase", "error", err)
		return nil, err
	}

	log.Infow("Connected to Supabase PostgreSQL")
	return dbpool, nil
}
