package internal

import (
	"errors"
	"strconv"
)

const gauge = "gauge"
const counter = "counter"

var ErrImpossibleMetricTypeOrValue = errors.New("impossible metric type or value")
var ErrUnknownMetricName = errors.New("unknown metric name")

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
			return ErrImpossibleMetricTypeOrValue
		}
		storage.gaugeMetrics[name] = GaugeMetric{
			name:  name,
			value: val,
		}

		return nil
	case counter:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return ErrImpossibleMetricTypeOrValue
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
		return ErrImpossibleMetricTypeOrValue
	}
}

func (storage MemStorage) Get(typ string, name string) (string, error) {
	switch typ {
	case gauge:
		value, status := storage.gaugeMetrics[name]
		if !status {
			return "", ErrUnknownMetricName
		}
		return value.name + strconv.FormatFloat(value.value, 'g', 2, 64), nil
	case counter:
		value, status := storage.counterMetrics[name]
		if !status {
			return "", ErrUnknownMetricName
		}
		return value.name + strconv.Itoa(int(value.value)), nil
	default:
		return "", ErrImpossibleMetricTypeOrValue
	}
}
