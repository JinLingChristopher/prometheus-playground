package http

import "context"

const (
	PromApiPrefix = "/api/v1"
	epQuery       = PromApiPrefix + "/query"
	epQueryRange  = PromApiPrefix + "/query_range"

	// epSeries      = PromApiPrefix + "/series"
)

type Warnings []string

type PromResult interface {
}

type PromMetricsFetcher interface {
	Query(ctx context.Context, query string) (PromResult, Warnings, error)
}
