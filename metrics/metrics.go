package metrics

import (
	"errors"
	"fmt"
	"sync"

	"github.com/labstack/echo"
	"github.com/popeyeio/handy"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/songzhaoliang/echotool/util"
)

type MetricsType int

const (
	MetricsTypeCounter MetricsType = iota
)

var _ fmt.Stringer = (*MetricsType)(nil)

func (t MetricsType) String() string {
	switch t {
	case MetricsTypeCounter:
		return "counter"
	}
	return fmt.Sprintf("unknown metrics type: %d", t)
}

var (
	ErrMetricsExists         = errors.New("metrics exists")
	ErrMetricsNotExists      = errors.New("metrics not exists")
	ErrMetricsTypeNotMatches = errors.New("metrics type not matches")
)

type MetricsClient struct {
	Namespace    string
	GlobalLabels prometheus.Labels
	AllMetrics   sync.Map
}

type MetricsClientOption func(*MetricsClient)

func WithNamespace(namespace string) MetricsClientOption {
	return func(c *MetricsClient) {
		c.Namespace = namespace
	}
}

func WithGlobalLabel(key, value string) MetricsClientOption {
	return func(c *MetricsClient) {
		if !handy.IsEmptyStr(key) && !handy.IsEmptyStr(value) {
			c.GlobalLabels[key] = value
		}
	}
}

func NewMetricsClient(opts ...MetricsClientOption) (c *MetricsClient) {
	c = &MetricsClient{}
	for _, opt := range opts {
		opt(c)
	}
	return
}

func (c *MetricsClient) DefineCounter(name string, parser LabelsParser) error {
	cv := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace:   c.Namespace,
		Name:        name,
		ConstLabels: c.GlobalLabels,
	}, Model(parser))

	if _, exists := c.AllMetrics.LoadOrStore(name, cv); exists {
		return ErrMetricsExists
	}
	return nil
}

func (c *MetricsClient) EmitCounter(name string, value float64, parser LabelsParser) error {
	metrics, exists := c.AllMetrics.Load(name)
	if !exists {
		return ErrMetricsNotExists
	}

	cv, ok := metrics.(*prometheus.CounterVec)
	if !ok {
		return ErrMetricsTypeNotMatches
	}

	counter, err := cv.GetMetricWith(prometheus.Labels(parser.ParseToLabels()))
	if err != nil {
		return err
	}

	counter.Add(value)
	return nil
}

func Model(parser LabelsParser) (keys []string) {
	labels := parser.ParseToLabels()
	keys = make([]string, 0, len(labels))
	for key := range labels {
		keys = append(keys, key)
	}
	return
}

var DefaultMetricsClient = NewMetricsClient()

func SetMetricsClient(c *MetricsClient) {
	if c != nil {
		DefaultMetricsClient = c
	}
}

func EmitCounter(name string, value float64, parser LabelsParser) error {
	return DefaultMetricsClient.EmitCounter(name, value, parser)
}

func Register(r *echo.Echo, prefixes ...string) {
	register(r.Group(util.GetPrefix(prefixes...)))
}

func RouterRegister(g *echo.Group, prefixes ...string) {
	register(g.Group(util.GetPrefix(prefixes...)))
}

func register(g *echo.Group) {
	g.GET("/metrics", util.WrapH(promhttp.Handler()))
}
