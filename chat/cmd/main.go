package main

import (
	"context"
	"flag"
	"github.com/dimastephen/chatServer/internal/app"
	"log"
)

var configPath string
var level string

func init() {
	flag.StringVar(&configPath, "config", "local.env", "config path file")
	flag.StringVar(&level, "loglvl", "info", "logger level")
}

func main() {
	ctx := context.Background()
	flag.Parse()
	a, err := app.NewApp(ctx, configPath, level)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = a.Run(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer a.ServiceProvider().DBClient(ctx).Close()
}
