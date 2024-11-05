package internal

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type MetricClient struct {
	client    *http.Client
	mstats    *runtime.MemStats
	pollCount *int
}

func getObservableRuntimeMetrics() [27]string {
	return [27]string{
		"Alloc",
		"BuckHashSys",
		"Frees",
		"GCCPUFraction",
		"GCSys",
		"HeapAlloc",
		"HeapIdle",
		"HeapInuse",
		"HeapObjects",
		"HeapReleased",
		"HeapSys",
		"LastGC",
		"Lookups",
		"MCacheInuse",
		"MCacheSys",
		"MSpanInuse",
		"MSpanSys",
		"Mallocs",
		"NextGC",
		"NumForcedGC",
		"NumGC",
		"OtherSys",
		"PauseTotalNs",
		"StackInuse",
		"StackSys",
		"Sys",
		"TotalAlloc",
	}
}

func buildURL(metricType string, metricName string, metricValue string) string {
	baseURL := url.URL{Scheme: "http", Host: "localhost:8080"}
	fullURL := baseURL.JoinPath("update", metricType, metricName, metricValue)
	return fullURL.String()
}

func (mClient MetricClient) UpdateMetrics(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stop update metrics")
			return
		case <-time.After(PollInterval * time.Second):
			runtime.ReadMemStats(mClient.mstats)
			*mClient.pollCount++
			fmt.Println("Successfully update metrics")
		}
	}
}

func (mClient MetricClient) SendMetrics(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stop send metrics")
			return
		case <-time.After(time.Second * ReportInterval):
			for _, metric := range getObservableRuntimeMetrics() {
				r := reflect.ValueOf(mClient.mstats)
				metricValue := reflect.Indirect(r).FieldByName(metric)

				var castedMetricValue float64

				if metricValue.CanUint() {
					castedMetricValue = float64(metricValue.Uint())
				} else if metricValue.CanFloat() {
					castedMetricValue = metricValue.Float()
				} else {
					panic("Unknown metric data type")
				}
				func() {
					runtimeMetricsResponse, err := mClient.client.Post(
						buildURL(Gauge, metric, fmt.Sprintf("%f", castedMetricValue)),
						"text/plain",
						nil,
					)
					if err != nil {
						fmt.Println("Error while send runtime metrics")
						return
					}
					defer runtimeMetricsResponse.Body.Close()
				}()
			}
			func() {
				pollCountMetricResponse, err := mClient.client.Post(
					buildURL(Counter, "PollCount", strconv.Itoa(*mClient.pollCount)),
					"text/plain",
					nil,
				)
				if err != nil {
					fmt.Println("Error while send PollCount metric")
					return
				}
				defer pollCountMetricResponse.Body.Close()
				randomValueMetricResponse, err := mClient.client.Post(
					buildURL(Gauge, "RandomValue", fmt.Sprintf("%f", rand.Float64())),
					"text/plain",
					nil,
				)
				if err != nil {
					fmt.Println("Error while send RandomValue metric")
					return
				}
				defer randomValueMetricResponse.Body.Close()
				fmt.Println("Successfully send metrics")
			}()
		}
	}
}

func NewMetricClient() *MetricClient {
	pollCount := 0
	var mstats runtime.MemStats
	return &MetricClient{
		client:    &http.Client{},
		mstats:    &mstats,
		pollCount: &pollCount,
	}
}
