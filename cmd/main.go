package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/RenterRus/sausage-auth/internal/app"
	"github.com/labstack/gommon/log"
)

func main() {
	path := flag.String("config", "../config.yaml", "path to config. Example: ../config.yaml")
	flag.Parse()
	if path == nil || len(*path) < 6 {
		log.Fatal("config flag not found")
		os.Exit(1)
	}

	fmt.Println("Path:", *path)

	app, err := app.NewApp(*path)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(); err != nil {
		log.Error(err)
	}
}
