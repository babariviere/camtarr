package tautulli

import (
	"log"
	"strconv"
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
	streamHistory  *prometheus.GaugeVec
}

func NewExporter(config Config) *Exporter {
	namespace := "tautulli"
	return &Exporter{
		config: config,
		currentStreams: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "current_streams",
				Help:      "All current streams.",
			},
			[]string{
				"state",
				"library_name",
				"full_title",
				"title",
				"parent_title",
				"grandparent_title",
				"progress",
				"user",
				"quality_profile",
				"transcode_decision",
				"player",
				"video_full_resolution",
			},
		),
		streamHistory: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "stream_history",
				Help:      "History of past streams.",
			},
			[]string{
				"media_type",
				"full_title",
				"title",
				"parent_title",
				"grandparent_title",
				"user",
				"transcode_decision",
				"player",
				"product",
				"play_duration",
			},
		),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.currentStreams.Describe(ch)
	e.streamHistory.Describe(ch)
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.scrapeActivity()
	e.scrapeHistory()

	e.currentStreams.Collect(ch)
	e.streamHistory.Collect(ch)
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
		e.currentStreams.WithLabelValues(
			session.State,
			session.LibraryName,
			session.FullTitle,
			session.Title,
			session.ParentTitle,
			session.GrandparentTitle,
			session.ProgressPercent,
			session.User,
			session.QualityProfile,
			session.TranscodeDecision,
			session.Player,
			session.VideoFullResolution,
		).Inc()
	}
}

func (e *Exporter) scrapeHistory() {
	resp, err := getHistory(e.config.Uri, e.config.ApiKey)
	if err != nil {
		log.Println("[tautulli] cannot get history:", err)
		return
	}

	// reset
	e.streamHistory.Reset()

	// fill data
	data := resp.Response.Data
	for _, entry := range data.Data {
		e.streamHistory.WithLabelValues(
			entry.MediaType,
			entry.FullTitle,
			entry.Title,
			entry.ParentTitle,
			entry.GrandparentTitle,
			entry.User,
			entry.TranscodeDecision,
			entry.Player,
			entry.Product,
			strconv.Itoa(entry.PlayDuration),
		).Inc()
	}
}
