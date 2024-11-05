package internal

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddMetricHandler(t *testing.T) {
	storage := NewMemStorage()

	type want struct {
		statusCode  int
		contentType string
		body        string
	}
	type request struct {
		metricType  string
		metricName  string
		metricValue string
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "happy_gauge_path",
			request: request{
				metricType:  "gauge",
				metricName:  "kek",
				metricValue: "1.1",
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "",
				body:        "",
			},
		},
		{
			name: "happy_counter_path",
			request: request{
				metricType:  "counter",
				metricName:  "kek",
				metricValue: "3",
			},
			want: want{
				statusCode:  200,
				contentType: "",
				body:        "",
			},
		},
		{
			name: "unknown_metric_type",
			request: request{
				metricType:  "gaug",
				metricName:  "kek",
				metricValue: "1.1",
			},
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "",
				body:        "",
			},
		},
		{
			name: "empty_metric_name",
			request: request{
				metricType:  "gauge",
				metricValue: "1.1",
			},
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				body:        "Metric name is empty\n",
			},
		},
		{
			name: "bad_counter_value",
			request: request{
				metricType:  "counter",
				metricName:  "kek",
				metricValue: "1.1",
			},
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "",
				body:        "",
			},
		},
		{
			name: "bad_gauge_value",
			request: request{
				metricType:  "gauge",
				metricName:  "kek",
				metricValue: "lol",
			},
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "",
				body:        "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", nil)
			request.SetPathValue("metricType", test.request.metricType)
			request.SetPathValue("metricName", test.request.metricName)
			request.SetPathValue("metricValue", test.request.metricValue)
			w := httptest.NewRecorder()
			AddMetricHandler(storage)(w, request)

			res := w.Result()

			assert.Equal(t, test.want.statusCode, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, test.want.body, string(resBody))
		})
	}
}
