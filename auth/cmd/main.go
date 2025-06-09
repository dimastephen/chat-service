package main

import (
	"context"
	"flag"
	"github.com/dimastephen/auth/internal/app"
	"log"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "local.env", "Path to env config file")
}

func main() {
	ctx := context.Background()
	a, err := app.NewApp(ctx, configPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = a.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}
