package server

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
)

var (
	// Metrics metrics
	Metrics               Prometheus
	responseStatusCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "rrmine",
			Name:      "response_status_counter",
			Help:      "http status counter",
		},
		[]string{"status", "uri"},
	)
	requestURICounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "rrmine",
			Name:      "request_uri_counter",
			Help:      "http uri counter",
		},
		[]string{"uri"},
	)
	panicCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "rrmine",
			Name:      "request_panic_counter",
			Help:      "http panic counter",
		},
		[]string{"uri"},
	)
)

func init() {
	prometheus.MustRegister(panicCounter)
	prometheus.MustRegister(requestURICounter)
	prometheus.MustRegister(responseStatusCounter)
}

func restrain(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	fn()
}

// Prometheus 服务指标统计
type Prometheus struct{}

// HTTPResponseStatusCounter http 响应状态码统计
func (Prometheus) HTTPResponseStatusCounter(uri string, status int) {
	var label = prometheus.Labels{"uri": uri, "status": strconv.Itoa(status)}
	if status == http.StatusNotFound {
		label = prometheus.Labels{"uri": "", "status": "404"}
	}
	restrain(responseStatusCounter.With(label).Inc)
}

// HTTPRequestURICounter http 请求资源计数
func (Prometheus) HTTPRequestURICounter(uri string) {
	restrain(requestURICounter.With(prometheus.Labels{"uri": uri}).Inc)
}

// RequestPanicCounter http 发生崩溃的统计
func (Prometheus) RequestPanicCounter(uri string) {
	restrain(panicCounter.With(prometheus.Labels{"uri": uri}).Inc)
}
