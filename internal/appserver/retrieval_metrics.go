package appserver

import (
	"sort"
	"sync"
	"time"
)

type RetrievalMetrics struct {
	mu       sync.Mutex
	snapshot RetrievalMetricsSnapshot
}

type RetrievalMetricsSnapshot struct {
	Total             int                         `json:"total"`
	Slow              int                         `json:"slow"`
	Failures          int                         `json:"failures"`
	ByMode            map[string]int              `json:"byMode"`
	FailuresByReason  map[string]int              `json:"failuresByReason"`
	AverageDurationMS int64                       `json:"averageDurationMs"`
	Last              RetrievalObservationSummary `json:"last"`
}

type RetrievalObservationSummary struct {
	Mode          string `json:"mode"`
	DurationMS    int64  `json:"durationMs"`
	Slow          bool   `json:"slow"`
	FailureReason string `json:"failureReason,omitempty"`
	ResultCount   int    `json:"resultCount"`
}

func NewRetrievalMetrics() *RetrievalMetrics {
	return &RetrievalMetrics{
		snapshot: RetrievalMetricsSnapshot{
			ByMode:           map[string]int{},
			FailuresByReason: map[string]int{},
		},
	}
}

func (m *RetrievalMetrics) ObserveRetrieval(observation RetrievalObservation) {
	if m == nil {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.snapshot.ByMode == nil {
		m.snapshot.ByMode = map[string]int{}
	}
	if m.snapshot.FailuresByReason == nil {
		m.snapshot.FailuresByReason = map[string]int{}
	}

	durationMS := durationMilliseconds(observation.Duration)
	previousTotal := m.snapshot.Total
	m.snapshot.Total++
	m.snapshot.ByMode[observation.Mode]++
	m.snapshot.AverageDurationMS = runningAverage(m.snapshot.AverageDurationMS, previousTotal, durationMS)
	if observation.Slow {
		m.snapshot.Slow++
	}
	if observation.FailureReason != "" {
		m.snapshot.Failures++
		m.snapshot.FailuresByReason[observation.FailureReason]++
	}
	m.snapshot.Last = RetrievalObservationSummary{
		Mode:          observation.Mode,
		DurationMS:    durationMS,
		Slow:          observation.Slow,
		FailureReason: observation.FailureReason,
		ResultCount:   observation.ResultCount,
	}
}

func (m *RetrievalMetrics) Snapshot() RetrievalMetricsSnapshot {
	if m == nil {
		return RetrievalMetricsSnapshot{
			ByMode:           map[string]int{},
			FailuresByReason: map[string]int{},
		}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	return RetrievalMetricsSnapshot{
		Total:             m.snapshot.Total,
		Slow:              m.snapshot.Slow,
		Failures:          m.snapshot.Failures,
		ByMode:            cloneSortedMap(m.snapshot.ByMode),
		FailuresByReason:  cloneSortedMap(m.snapshot.FailuresByReason),
		AverageDurationMS: m.snapshot.AverageDurationMS,
		Last:              m.snapshot.Last,
	}
}

func runningAverage(current int64, previousTotal int, next int64) int64 {
	if previousTotal <= 0 {
		return next
	}
	return ((current * int64(previousTotal)) + next) / int64(previousTotal+1)
}

func durationMilliseconds(duration time.Duration) int64 {
	if duration <= 0 {
		return 0
	}
	return duration.Milliseconds()
}

func cloneSortedMap(source map[string]int) map[string]int {
	if len(source) == 0 {
		return map[string]int{}
	}
	keys := make([]string, 0, len(source))
	for key := range source {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	cloned := make(map[string]int, len(source))
	for _, key := range keys {
		cloned[key] = source[key]
	}
	return cloned
}
