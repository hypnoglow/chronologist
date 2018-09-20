package prometheus

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func NewChronicle() *Chronicle {
	c := &Chronicle{
		releases: prometheus.NewDesc(
			"chronologist_releases",
			"Chronologist releases",
			[]string{
				"release_type",
				"release_name",
				"release_revision",
				"release_namespace",
			},
			nil,
		),
		metrics: make(map[string]prometheus.Metric),
	}
	return c
}

type Chronicle struct {
	releases *prometheus.Desc

	metrics map[string]prometheus.Metric
	mtx     sync.RWMutex // Protects metrics.
}

// Describe implements Describe method of prometheus.Chronicle interface.
func (c *Chronicle) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.releases
}

// Collect implements Collect method of prometheus.Chronicle interface.
func (c *Chronicle) Collect(ch chan<- prometheus.Metric) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	for _, metric := range c.metrics {
		ch <- metric
	}
}

func (c *Chronicle) RegisterRelease(id string, t time.Time, kind, name, revision, namespace string) {
	m := prometheus.MustNewConstMetric(
		c.releases,
		prometheus.CounterValue,
		1,
		kind,
		name,
		revision,
		namespace,
	)
	tm := prometheus.NewMetricWithTimestamp(t, m)

	// If release metric already exists, it will be overwritten.
	c.metrics[id] = tm
}

func (c *Chronicle) UnregisterRelease(id string) {
	delete(c.metrics, id)
}
