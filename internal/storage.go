package internal

import (
	"errors"
	"fmt"
	"strconv"
)

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

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gaugeMetrics:   make(map[string]GaugeMetric, 0),
		counterMetrics: make(map[string]CounterMetric, 0),
	}
}

func (storage MemStorage) Add(typ string, name string, value string) error {
	fmt.Println(value)
	switch typ {
	case Gauge:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return ErrImpossibleMetricTypeOrValue
		}
		storage.gaugeMetrics[name] = GaugeMetric{
			name:  name,
			value: val,
		}

		return nil
	case Counter:
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
	case Gauge:
		value, status := storage.gaugeMetrics[name]
		if !status {
			return "", ErrUnknownMetricName
		}
		return value.name + strconv.FormatFloat(value.value, 'g', 2, 64), nil
	case Counter:
		value, status := storage.counterMetrics[name]
		if !status {
			return "", ErrUnknownMetricName
		}
		return value.name + strconv.Itoa(int(value.value)), nil
	default:
		return "", ErrImpossibleMetricTypeOrValue
	}
}
