package prometheus

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestChronicle(t *testing.T) {
	chronicle := NewChronicle()
	prometheus.MustRegister(chronicle)

	chronicle.RegisterRelease(
		"foo.v1",
		time.Date(2018, 7, 8, 10, 20, 30, 0, time.UTC),
		"rollout",
		"foo",
		"1",
		"default",
	)
	chronicle.RegisterRelease(
		"bar.v2",
		time.Date(2018, 1, 2, 3, 4, 5, 0, time.UTC),
		"rollout",
		"bar",
		"2",
		"default",
	)

	// Asserts

	mf, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	var key int
	for k, v := range mf {
		if v.GetName() == "chronologist_releases_total" {
			key = k
			break
		}
	}

	releasesMF := mf[key]

	assert.Equal(t, 2, len(releasesMF.GetMetric()))

	m1 := releasesMF.GetMetric()[0]
	assert.EqualValues(t, m1.GetCounter().GetValue(), 1)

	m1labels := labelMap(m1.GetLabel())
	assert.Equal(t, m1labels["release_name"], "bar")
	assert.Equal(t, m1labels["release_revision"], "2")
	assert.Equal(t, m1labels["release_type"], "rollout")
	assert.Equal(t, m1labels["release_namespace"], "default")

	m2 := releasesMF.GetMetric()[1]
	assert.EqualValues(t, m2.GetCounter().GetValue(), 1)

	m2labels := labelMap(m2.GetLabel())
	assert.Equal(t, m2labels["release_name"], "foo")
	assert.Equal(t, m2labels["release_revision"], "1")
	assert.Equal(t, m2labels["release_type"], "rollout")
	assert.Equal(t, m2labels["release_namespace"], "default")
}

func labelMap(lps []*io_prometheus_client.LabelPair) map[string]string {
	m := make(map[string]string)
	for _, lp := range lps {
		m[lp.GetName()] = lp.GetValue()
	}
	return m
}
