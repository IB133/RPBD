package main

import (
	"log"

	"github.com/IB133/RPBD/final_project/internal/app"
	"github.com/IB133/RPBD/final_project/internal/config"
)

func main() {
	conf := config.NewConfig()
	if err := app.Start(conf); err != nil {
		log.Fatal(err)
	}
}
