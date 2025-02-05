package internal

import (
	"errors"
	"net/http"
)

var storage = MemStorage{
	gaugeMetrics:   make(map[string]GaugeMetric),
	counterMetrics: make(map[string]CounterMetric),
}

func postHandler(writer http.ResponseWriter, request *http.Request) {
	err := storage.Add(
		request.PathValue("metricType"),
		request.PathValue("metricName"),
		request.PathValue("metricValue"),
	)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func getHandler(writer http.ResponseWriter, request *http.Request) {
	value, err := storage.Get(request.PathValue("metricType"), request.PathValue("metricName"))
	if err != nil {
		if errors.Is(err, ErrImpossibleMetricTypeOrValue) {
			writer.WriteHeader(http.StatusBadRequest)
			return
		} else if errors.Is(err, ErrUnknownMetricName) {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
	}
	writer.Write([]byte(value))
}

func RunServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", postHandler)
	mux.HandleFunc("/get/{metricType}/{metricName}", getHandler)
	http.ListenAndServe("localhost:8080", mux)
}
