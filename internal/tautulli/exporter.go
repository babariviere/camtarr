package tautulli

import (
	"log"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type Config struct {
	ApiKey       string `env:"API_KEY"`
	Uri          string `env:"URI"`
	ExporterPort string `env:"EXPORTER_PORT" envDefault:"9301"`
}

type Exporter struct {
	config Config
	mu     sync.Mutex

	currentStreams *prometheus.GaugeVec
}

func NewExporter(config Config) *Exporter {
	namespace := "tautulli"
	return &Exporter{
		config: config,
		currentStreams: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Name: "current_streams", Help: "All current streams."},
			[]string{"state", "library_name", "title", "parent_title", "grandparent_title", "progress", "user"},
		),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.currentStreams.Describe(ch)
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.scrapeActivity()

	e.currentStreams.Collect(ch)
}

func (e *Exporter) scrapeActivity() {
	resp, err := getActivity(e.config.Uri, e.config.ApiKey)
	if err != nil {
		log.Println("[tautulli] cannot get activity:", err)
		return
	}

	// reset
	e.currentStreams.Reset()

	// fill data
	data := resp.Response.Data
	for _, session := range data.Sessions {
		e.currentStreams.WithLabelValues(session.State, session.LibraryName, session.Title, session.ParentTitle, session.GrandparentTitle, session.ProgressPercent, session.User).Inc()
	}
}
