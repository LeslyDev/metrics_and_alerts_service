package internal

import (
	"errors"
	"net/http"
)

func AddMetricHandler(storage *MemStorage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(writer, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		}
		if request.PathValue("metricName") == "" {
			http.Error(writer, "Metric name is empty", http.StatusNotFound)
			return
		}
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
}

func GetMetricHandler(storage *MemStorage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
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
}
