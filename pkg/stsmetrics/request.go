package stsmetrics

import (
	"time"
)

//STSRequest requests generation of timeseries data for a particular stream
type STSRequest struct {
	Stream int32 `json:"stream"`
}

//MetricTimeSeries the data structure
type MetricTimeSeries struct {
	Stream int32     `json:"stream"`
	Count  int       `json:"count"`
	Time   time.Time `json:"time"`
}
