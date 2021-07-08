package prometheus

import "encoding/json"

var RespStatusInfo = map[int]string{
	400: "Bad Request",
	422: "Unprocessable Entity",
	503: "Service Unavailable",
}

type ErrorType string

const (
	// Possible values for ErrorType.
	ErrBadData     ErrorType = "bad_data"
	ErrTimeout     ErrorType = "timeout"
	ErrCanceled    ErrorType = "canceled"
	ErrExec        ErrorType = "execution"
	ErrBadResponse ErrorType = "bad_response"
	ErrServer      ErrorType = "server_error"
	ErrClient      ErrorType = "client_error"
)

type PrometheusResp struct {
	Status    string          `json:"status"`
	Data      json.RawMessage `json:"data"`
	ErrorType string          `json:"errorType"`
	Error     string          `json:"error"`
	Warnings  []string        `json:"warnings,omitempty"`
}

type Result interface {
}

// MetricsFetch should be implemented by each metrics if need to be automatically monitored
type MetricsFetcher interface {
	// Get return the final
	Get() (Result, error)
}
