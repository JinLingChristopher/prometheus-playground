package prometheus

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/JinLingChristopher/prometheus-playground/pkg/http"
	"github.com/pingcap/errors"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type TiCDCPromMetricCatcher struct {
	v1.API
}

func NewTiCDCPromMetricCatcher(addr string) (result *TiCDCPromMetricCatcher, err error) {
	client, err := api.NewClient(api.Config{Address: addr})
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &TiCDCPromMetricCatcher{
		API: v1.NewAPI(client),
	}, nil
}

func (c *TiCDCPromMetricCatcher) Query(ctx context.Context, query string) (http.PromResult, http.Warnings, error) {

}

func getFromPromRange(start time.Time, end time.Time, metric string) model.Value {
	client, err := api.NewClient(api.Config{
		Address: "http://localhost:9090",
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}
	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := v1.Range{
		Start: start,
		End:   end,
		Step:  time.Second,
	}
	result, warnings, err := v1api.QueryRange(ctx, metric, r)
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	fmt.Printf("Result:\n%v\n", result)

	return result
}
