package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/IB133/RPBD/final_project/internal/config"
	"github.com/IB133/RPBD/final_project/internal/db"
)

func Start(cfg *config.Config, ctx context.Context) error {
	store, err := db.NewConnect(fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))
	if err != nil {
		return err
	}

	srv := newServer(store, ctx)
	return http.ListenAndServe(cfg.HTTPPort, srv)
}
