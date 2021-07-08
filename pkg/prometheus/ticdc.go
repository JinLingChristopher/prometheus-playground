package prometheus

import (
	"context"
	"fmt"
	"time"

	"github.com/pingcap/errors"
	"github.com/pingcap/log"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

type TiCDCUptime struct {
	v1.API
	query string
}

var serverPromQLFormat = map[string]string{
	"uptime":                       "time() - process_start_time_seconds{}",
	"cpu_usage":                    "rate(process_cpu_seconds_total{}[1m])",
	"process_mem_usage":            "process_resident_memory_bytes{}",
	"go_memstats_heap_alloc_bytes": "go_memstats_heap_alloc_bytes{}",
	"go_goroutines":                "go_goroutines{}",
}

var changefeedPromQLFormat = map[string]string{
	"processor-resolved-ts":     "max(ticdc_processor_resolved_ts{}) by (capture)",
	"changefeed-checkpoint":     "max(ticdc_processor_checkpoint_ts{}) by (capture) > 0",
	"table-resolved-ts":         "max(ticdc_processor_table_resolved_ts{}) by (capture, table)",
	"pd-etcd-requests":          "sum(rate(ticdc_etcd_request_count{}[1m])) by (capture, type)",
	"changefeed-checkpoint-lag": "max(ticdc_owner_checkpoint_ts_lag{}) by (changefeed)",
	"processor-resolved-ts-lag": "sum(ticdc_processor_resolved_ts_lag{}) by (capture)",

	"sink-write-duration-p95":  "histogram_quantile(0.95, sum(rate(ticdc_sink_txn_exec_duration_bucket{}[1m])) by (le,instance))",
	"sink-write-duration-p99":  "histogram_quantile(0.99, sum(rate(ticdc_sink_txn_exec_duration_bucket{}[1m])) by (le,instance))",
	"sink-write-duration-p999": "histogram_quantile(0.999, sum(rate(ticdc_sink_txn_exec_duration_bucket{}[1m])) by (le,instance))",

	"sink-write-row-counts-by-capture":    "sum (rate(ticdc_sink_txn_batch_size_sum{}[1m])) by (capture)",
	"sink-write-row-counts-by-changefeed": "sum (rate(ticdc_sink_txn_batch_size_sum{}[1m])) by (changefeed)",
}

var eventsPromQLFormat = map[string]string{
	"kv-client-receive-events-per-second": "sum(rate(ticdc_kvclient_pull_event_count{}[1m])) by (instance, type)",
}

var unifiedSorterPromQLFormat = map[string]string{
	"unified-sorter-resolved-ts": "min(ticdc_sorter_resolved_ts_gauge{}) by (capture)",
}

var tikvCDCComponentPromQLFormat = map[string]string{
	"cdc-endpoint-cpu":                    "sum(rate(tikv_thread_cpu_seconds_total{name=~\"cdc.*\"}[1m])) by (instance)",
	"cdc-worker-cpu-worker":               "sum(rate(tikv_thread_cpu_seconds_total{name=~\"cdcwkr.*\"}[1m])) by (instance)",
	"cdc-worker-cpu-tso":                  "sum(rate(tikv_thread_cpu_seconds_total{name=~\"tso\"}[1m])) by (instance)",
	"resolved-ts-lag-duration-percentile": "histogram_quantile(0.99999, sum(rate(tikv_cdc_resolved_ts_gap_seconds_bucket{}[1m])) by (le, instance))",
	"memory-without-block-cache":          "(avg(process_resident_memory_bytes{job=~\"tikv.*\"}) by (instance)) - (avg(tikv_engine_block_cache_size_bytes{db=\"kv\"}) by(instance))",
	"cdc-pending-bytes-in-memory":         "",
	"captured-region-counts":              "avg(tikv_cdc_captured_region_total{}) by (instance)",
}

var metrics = []map[string]string{serverPromQLFormat, changefeedPromQLFormat, eventsPromQLFormat, unifiedSorterPromQLFormat, tikvCDCComponentPromQLFormat}

func NewTiCDCUptime(addr string) (result *TiCDCUptime, err error) {
	client, err := api.NewClient(api.Config{Address: addr})
	if err != nil {
		return nil, errors.Trace(err)
	}
	query := fmt.Sprintf(serverPromQLFormat["uptime"], "172.16.6.191:8302", "ticdc")

	return &TiCDCUptime{
		API:   v1.NewAPI(client),
		query: query,
	}, nil
}

func (c *TiCDCUptime) Get() (Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, warnings, err := c.Query(ctx, c.query, time.Now())
	if len(warnings) > 0 {
		log.Warn("get ticdc uptime warnings")
	}
	if err != nil {
		return nil, errors.Trace(err)
	}

	return resp, errors.Trace(err)
}

type TiCDCMonitor struct {
}
