package app

import (
	"fmt"
	"net/http"

	"github.com/IB133/RPBD/final_project/internal/config"
	"github.com/IB133/RPBD/final_project/internal/db"
)

func Start(cfg *config.Config) error {
	store, err := db.NewConnect(fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBNname))
	if err != nil {
		return err
	}

	srv := newServer(store)
	return http.ListenAndServe(cfg.HTTPPort, srv)
}
