package appserver

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestBytesToVectorLiteral(t *testing.T) {
	got := bytesToVectorLiteral([]byte{1, 2, 255})
	if got != "[1,2,255]" {
		t.Fatalf("vector literal = %q", got)
	}
}

func TestRankedChunksPGVectorSQLBuildsTopKQuery(t *testing.T) {
	query, args, err := rankedChunksPGVectorSQL("opsmate", []string{"doc-1", "doc-2"}, "[1,2]", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(args) != 5 {
		t.Fatalf("args len = %d, want 5", len(args))
	}
	if args[0] != "opsmate" || args[1] != "[1,2]" || args[4] != 5 {
		t.Fatalf("args = %+v", args)
	}
	for _, want := range []string{
		"e.embedding <=> $2::vector",
		"d.id IN ($3, $4)",
		"LIMIT $5",
	} {
		if !contains(query, want) {
			t.Fatalf("query = %q, want %q", query, want)
		}
	}
}

func TestRankedChunksPGVectorSQLValidatesInput(t *testing.T) {
	if _, _, err := rankedChunksPGVectorSQL("", []string{"doc"}, "[1]", 1); err == nil {
		t.Fatal("expected namespace error")
	}
	if _, _, err := rankedChunksPGVectorSQL("opsmate", nil, "[1]", 1); err == nil {
		t.Fatal("expected document ids error")
	}
	if _, _, err := rankedChunksPGVectorSQL("opsmate", []string{"doc"}, "", 1); err == nil {
		t.Fatal("expected vector error")
	}
}

func TestPostgresRetrieverObservesNoopRetrieval(t *testing.T) {
	observer := &recordingRetrievalObserver{}
	_, err := PostgresRetriever{
		Repository: nil,
		Observer:   observer,
	}.Retrieve(context.Background(), RetrievalRequest{Message: "status"})
	if err != nil {
		t.Fatal(err)
	}
	if len(observer.items) != 1 {
		t.Fatalf("observations len = %d, want 1", len(observer.items))
	}
	if observer.items[0].Mode != "bytea" {
		t.Fatalf("mode = %q, want bytea", observer.items[0].Mode)
	}
}

func TestRetrievalMetricsAggregatesObservations(t *testing.T) {
	metrics := NewRetrievalMetrics()
	metrics.ObserveRetrieval(RetrievalObservation{
		Mode:        "bytea",
		Duration:    10 * time.Millisecond,
		ResultCount: 1,
	})
	metrics.ObserveRetrieval(RetrievalObservation{
		Mode:          "pgvector",
		Duration:      30 * time.Millisecond,
		Slow:          true,
		FailureReason: "pgvector_query_failed",
	})

	snapshot := metrics.Snapshot()

	if snapshot.Total != 2 {
		t.Fatalf("total = %d, want 2", snapshot.Total)
	}
	if snapshot.Slow != 1 {
		t.Fatalf("slow = %d, want 1", snapshot.Slow)
	}
	if snapshot.Failures != 1 {
		t.Fatalf("failures = %d, want 1", snapshot.Failures)
	}
	if snapshot.ByMode["bytea"] != 1 || snapshot.ByMode["pgvector"] != 1 {
		t.Fatalf("by mode = %+v, want bytea and pgvector counts", snapshot.ByMode)
	}
	if snapshot.FailuresByReason["pgvector_query_failed"] != 1 {
		t.Fatalf("failures by reason = %+v, want pgvector failure", snapshot.FailuresByReason)
	}
	if snapshot.AverageDurationMS != 20 {
		t.Fatalf("average duration = %d, want 20", snapshot.AverageDurationMS)
	}
	if snapshot.Last.FailureReason != "pgvector_query_failed" {
		t.Fatalf("last failure reason = %q, want pgvector_query_failed", snapshot.Last.FailureReason)
	}
}

type recordingRetrievalObserver struct {
	items []RetrievalObservation
}

func (o *recordingRetrievalObserver) ObserveRetrieval(observation RetrievalObservation) {
	o.items = append(o.items, observation)
}

func contains(value string, substr string) bool {
	return strings.Contains(value, substr)
}
