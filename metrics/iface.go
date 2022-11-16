package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type LabelsParser interface {
	ParseToLabels() prometheus.Labels
}
