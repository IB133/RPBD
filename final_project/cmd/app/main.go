package main

import (
	"context"
	"log"

	"github.com/IB133/RPBD/final_project/internal/app"
	"github.com/IB133/RPBD/final_project/internal/config"
)

func main() {
	ctx := context.Background()
	conf := config.NewConfig()
	if err := app.Start(conf, ctx); err != nil {
		log.Fatal(err)
	}
}
