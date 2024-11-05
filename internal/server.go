package internal

import (
	"net/http"
)

func RunServer() {
	mux := http.NewServeMux()

	storage := NewMemStorage()

	mux.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", AddMetricHandler(storage))
	mux.HandleFunc("/get/{metricType}/{metricName}", GetMetricHandler(storage))
	http.ListenAndServe("localhost:8080", mux)
}
