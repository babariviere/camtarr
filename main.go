package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/babariviere/camtarr/internal/tautulli"
	"github.com/caarlos0/env/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var cfg tautulli.Config
	if err := env.ParseWithOptions(&cfg, env.Options{Prefix: "TAUTULLI_"}); err != nil {
		log.Fatalln("failed to parse env:", err)
	}

	fmt.Printf("Parsed env: %+v\n", cfg)

	tautulli := tautulli.NewExporter(cfg)
	prometheus.MustRegister(tautulli)

	fmt.Println("Hello Marcel")
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":"+cfg.ExporterPort, nil)
}
