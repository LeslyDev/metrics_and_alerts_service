package internal

import (
	"errors"
	"fmt"
	"strconv"
)

const gauge = "gauge"
const counter = "counter"

var UnknownMetricType = errors.New("unknown metric type")
var ImpossibleMetricValue = errors.New("impossible metric value")
var UnknownMetricName = errors.New("unknown metric name")

type GaugeMetric struct {
	name  string
	value float64
}

type CounterMetric struct {
	name  string
	value int64
}

type MemStorage struct {
	gaugeMetrics   map[string]GaugeMetric
	counterMetrics map[string]CounterMetric
}

func (storage MemStorage) Add(typ string, name string, value string) error {
	switch typ {
	case gauge:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return ImpossibleMetricValue
		}
		storage.gaugeMetrics[name] = GaugeMetric{
			name:  name,
			value: val,
		}

		return nil
	case counter:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return ImpossibleMetricValue
		}
		value, status := storage.counterMetrics[name]
		if status {
			storage.counterMetrics[name] = CounterMetric{
				name:  name,
				value: value.value + val,
			}
		} else {
			storage.counterMetrics[name] = CounterMetric{
				name:  name,
				value: val,
			}
		}

		return nil
	default:
		return UnknownMetricType
	}
}

func (storage MemStorage) Get(typ string, name string) (string, error) {
	switch typ {
	case gauge:
		value, status := storage.gaugeMetrics[name]
		fmt.Printf("Да это гуага, имя %s, значение %f", value.name, value.value)
		if !status {
			return "", UnknownMetricName
		}
		return value.name + strconv.FormatFloat(value.value, 'g', 2, 64), nil
	case counter:
		value, status := storage.counterMetrics[name]
		fmt.Printf("Да это counter, имя %s, значение %d", value.name, value.value)
		if !status {
			return "", UnknownMetricName
		}
		return value.name + strconv.Itoa(int(value.value)), nil
	default:
		return "", UnknownMetricType
	}
}
