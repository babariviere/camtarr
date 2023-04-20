package main

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v8"
)

type tautulliConfig struct {
	ApiKey       string `env:"API_KEY"`
	Uri          string `env:"URI"`
	ExporterPort string `env:"EXPORTER_PORT" envDefault:"9301"`
}

type config struct {
	Tautulli tautulliConfig `envPrefix:"TAUTULLI_"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln("failed to parse env:", err)
	}

	fmt.Printf("Parsed env: %+v\n", cfg)

	fmt.Println("Hello Marcel")
}
