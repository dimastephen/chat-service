package main

import (
	"context"
	"flag"
	"github.com/dimastephen/auth/internal/app"
	"log"
)

var configPath string
var level string

func init() {
	flag.StringVar(&configPath, "config", "local.env", "Path to env config file")
	flag.StringVar(&level, "loglvl", "info", "logger level")
}

func main() {
	ctx := context.Background()
	flag.Parse()
	a, err := app.NewApp(ctx, configPath, level)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = a.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}
