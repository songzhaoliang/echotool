package metric

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/labstack/echo"
	"github.com/popeyeio/handy"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/songzhaoliang/echotool/util"
)

type MetricType int

const (
	MetricTypeCounter MetricType = iota
	MetricTypeGauge
	MetricTypeHistogram
	MetricTypeSummary
)

var _ fmt.Stringer = (*MetricType)(nil)

func (t MetricType) String() string {
	switch t {
	case MetricTypeCounter:
		return "counter"
	case MetricTypeGauge:
		return "gauge"
	case MetricTypeHistogram:
		return "histogram"
	case MetricTypeSummary:
		return "summary"
	}
	return fmt.Sprintf("unknown metric type: %d", t)
}

var (
	DefaultBuckets    = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
	DefaultObjectives = map[float64]float64{.5: .05, .8: .01, .9: .01, .95: .001, .99: .001}
)

var (
	ErrMetricExists         = errors.New("metric exists")
	ErrMetricNotExists      = errors.New("metric not exists")
	ErrUnknownMetricType    = errors.New("unknown metric type")
	ErrMetricTypeNotMatches = errors.New("metric type not matches")
)

type MetricClient struct {
	Namespace    string
	GlobalLabels prometheus.Labels
	AllMetrics   sync.Map
}

type MetricClientOption func(*MetricClient)

func WithNamespace(namespace string) MetricClientOption {
	return func(c *MetricClient) {
		c.Namespace = namespace
	}
}

func WithGlobalLabel(key, value string) MetricClientOption {
	return func(c *MetricClient) {
		if !handy.IsEmptyStr(key) && !handy.IsEmptyStr(value) {
			c.GlobalLabels[key] = value
		}
	}
}

func NewMetricClient(opts ...MetricClientOption) (c *MetricClient) {
	c = &MetricClient{}
	for _, opt := range opts {
		opt(c)
	}
	return
}

func (c *MetricClient) DefineCounter(name string, parser LabelsParser) error {
	cv := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace:   c.Namespace,
		Name:        name,
		ConstLabels: c.GlobalLabels,
	}, Model(parser))

	return c.define(name, cv)
}

func (c *MetricClient) DefineGauge(name string, parser LabelsParser) error {
	gv := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   c.Namespace,
		Name:        name,
		ConstLabels: c.GlobalLabels,
	}, Model(parser))

	return c.define(name, gv)
}

func (c *MetricClient) DefineHistogram(name string, parser LabelsParser, buckets []float64) error {
	hv := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   c.Namespace,
		Name:        name,
		ConstLabels: c.GlobalLabels,
		Buckets:     buckets,
	}, Model(parser))

	return c.define(name, hv)
}

func (c *MetricClient) DefineSummary(name string, parser LabelsParser, objectives map[float64]float64) error {
	sv := promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:   c.Namespace,
		Name:        name,
		ConstLabels: c.GlobalLabels,
		Objectives:  objectives,
	}, Model(parser))

	return c.define(name, sv)
}

func (c *MetricClient) define(name string, metric interface{}) error {
	if _, exists := c.AllMetrics.LoadOrStore(name, metric); exists {
		return ErrMetricExists
	}
	return nil
}

func (c *MetricClient) EmitCounter(name string, value float64, parser LabelsParser) error {
	metric, exists := c.AllMetrics.Load(name)
	if !exists {
		return ErrMetricNotExists
	}

	cv, ok := metric.(*prometheus.CounterVec)
	if !ok {
		return ErrMetricTypeNotMatches
	}

	counter, err := cv.GetMetricWith(parser.ParseToLabels())
	if err != nil {
		return err
	}

	counter.Add(value)
	return nil
}

func (c *MetricClient) EmitGauge(name string, value float64, parser LabelsParser) error {
	metric, exists := c.AllMetrics.Load(name)
	if !exists {
		return ErrMetricNotExists
	}

	gv, ok := metric.(*prometheus.GaugeVec)
	if !ok {
		return ErrMetricTypeNotMatches
	}

	gauge, err := gv.GetMetricWith(parser.ParseToLabels())
	if err != nil {
		return err
	}

	gauge.Set(value)
	return nil
}

func (c *MetricClient) EmitHistogram(name string, value float64, parser LabelsParser) error {
	metric, exists := c.AllMetrics.Load(name)
	if !exists {
		return ErrMetricNotExists
	}

	hv, ok := metric.(*prometheus.HistogramVec)
	if !ok {
		return ErrMetricTypeNotMatches
	}

	histogram, err := hv.GetMetricWith(parser.ParseToLabels())
	if err != nil {
		return err
	}

	histogram.Observe(value)
	return nil
}

func (c *MetricClient) EmitSummary(name string, value float64, parser LabelsParser) error {
	metric, exists := c.AllMetrics.Load(name)
	if !exists {
		return ErrMetricNotExists
	}

	sv, ok := metric.(*prometheus.SummaryVec)
	if !ok {
		return ErrMetricTypeNotMatches
	}

	summary, err := sv.GetMetricWith(parser.ParseToLabels())
	if err != nil {
		return err
	}

	summary.Observe(value)
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

var DefaultMetricClient = NewMetricClient()

func SetMetricClient(c *MetricClient) {
	if c != nil {
		DefaultMetricClient = c
	}
}

func EmitCounter(name string, value float64, parser LabelsParser) error {
	return DefaultMetricClient.EmitCounter(name, value, parser)
}

func EmitGauge(name string, value float64, parser LabelsParser) error {
	return DefaultMetricClient.EmitGauge(name, value, parser)
}

func EmitHistogram(name string, value float64, parser LabelsParser) error {
	return DefaultMetricClient.EmitHistogram(name, value, parser)
}

func EmitHistogramTimer(name string, t time.Time, parser LabelsParser) error {
	return EmitHistogram(name, time.Since(t).Seconds(), parser)
}

func EmitSummary(name string, value float64, parser LabelsParser) error {
	return DefaultMetricClient.EmitSummary(name, value, parser)
}

func EmitSummaryTimer(name string, t time.Time, parser LabelsParser) error {
	return EmitSummary(name, time.Since(t).Seconds(), parser)
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
