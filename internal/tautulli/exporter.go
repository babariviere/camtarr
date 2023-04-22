package tautulli

import (
	"log"
	"strconv"
	"sync"
	"time"

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

	streams       *prometheus.GaugeVec
	streamHistory *prometheus.GaugeVec

	totalBandwidth, lanBandwidth, wanBandwidth                                        prometheus.Gauge
	streamCount, streamCountDirectPlay, streamCountDirectStream, streamCountTranscode prometheus.Gauge
}

func NewExporter(config Config) *Exporter {
	namespace := "tautulli"
	return &Exporter{
		config: config,
		streams: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "streams",
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
				"session_id",
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
				"date",
			},
		),
		totalBandwidth: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "total_bandwidth",
			Help:      "Total bandwidth usage for streams.",
		}),
		wanBandwidth: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "wan_bandwidth",
			Help:      "WAN bandwidth usage for streams.",
		}),
		lanBandwidth: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "lan_bandwidth",
			Help:      "LAN bandwidth usage for streams.",
		}),
		streamCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "stream_count",
			Help:      "Total number of streams.",
		}),
		streamCountDirectPlay: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "stream_count_direct_play",
			Help:      "Total number of streams using direct play",
		}),
		streamCountDirectStream: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "stream_count_direct_stream",
			Help:      "Total number of streams using direct stream",
		}),
		streamCountTranscode: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "stream_count_transcode",
			Help:      "Total number of streams using transcode",
		}),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.streams.Describe(ch)
	e.streamHistory.Describe(ch)

	ch <- e.totalBandwidth.Desc()
	ch <- e.wanBandwidth.Desc()
	ch <- e.lanBandwidth.Desc()

	ch <- e.streamCount.Desc()
	ch <- e.streamCountDirectPlay.Desc()
	ch <- e.streamCountDirectStream.Desc()
	ch <- e.streamCountTranscode.Desc()
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.scrapeActivity()
	e.scrapeHistory()

	e.streams.Collect(ch)
	e.streamHistory.Collect(ch)

	ch <- e.totalBandwidth
	ch <- e.wanBandwidth
	ch <- e.lanBandwidth

	ch <- e.streamCount
	ch <- e.streamCountDirectPlay
	ch <- e.streamCountDirectStream
	ch <- e.streamCountTranscode

}

func (e *Exporter) scrapeActivity() {
	resp, err := getActivity(e.config.Uri, e.config.ApiKey)
	if err != nil {
		log.Println("[tautulli] cannot get activity:", err)
		return
	}

	// reset
	e.streams.Reset()

	// fill data
	data := resp.Response.Data
	for _, session := range data.Sessions {
		e.streams.WithLabelValues(
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
			session.SessionID,
		).Inc()
	}

	e.totalBandwidth.Set(float64(data.TotalBandwidth))
	e.wanBandwidth.Set(float64(data.WanBandwidth))
	e.lanBandwidth.Set(float64(data.LanBandwidth))

	streamCount, _ := strconv.Atoi(data.StreamCount)
	e.streamCount.Set(float64(streamCount))
	e.streamCountDirectPlay.Set(float64(data.StreamCountDirectPlay))
	e.streamCountDirectStream.Set(float64(data.StreamCountDirectStream))
	e.streamCountTranscode.Set(float64(data.StreamCountTranscode))
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
			time.Unix(int64(entry.Date), 0).UTC().Format("2006-01-02 15:04:05"),
		).Inc()
	}
}
