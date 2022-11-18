package metrics

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

type MetricsType int

const (
	MetricsTypeCounter MetricsType = iota
	MetricsTypeGauge
	MetricsTypeHistogram
	MetricsTypeSummary
)

var _ fmt.Stringer = (*MetricsType)(nil)

func (t MetricsType) String() string {
	switch t {
	case MetricsTypeCounter:
		return "counter"
	case MetricsTypeGauge:
		return "gauge"
	case MetricsTypeHistogram:
		return "histogram"
	case MetricsTypeSummary:
		return "summary"
	}
	return fmt.Sprintf("unknown metrics type: %d", t)
}

var (
	DefaultBuckets    = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
	DefaultObjectives = map[float64]float64{.5: .05, .8: .01, .9: .01, .95: .001, .99: .001}
)

var (
	ErrMetricsExists         = errors.New("metrics exists")
	ErrMetricsNotExists      = errors.New("metrics not exists")
	ErrUnknownMetricsType    = errors.New("unknown metrics type")
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

	return c.define(name, cv)
}

func (c *MetricsClient) DefineGauge(name string, parser LabelsParser) error {
	gv := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   c.Namespace,
		Name:        name,
		ConstLabels: c.GlobalLabels,
	}, Model(parser))

	return c.define(name, gv)
}

func (c *MetricsClient) DefineHistogram(name string, parser LabelsParser, buckets []float64) error {
	hv := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   c.Namespace,
		Name:        name,
		ConstLabels: c.GlobalLabels,
		Buckets:     buckets,
	}, Model(parser))

	return c.define(name, hv)
}

func (c *MetricsClient) DefineSummary(name string, parser LabelsParser, objectives map[float64]float64) error {
	sv := promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:   c.Namespace,
		Name:        name,
		ConstLabels: c.GlobalLabels,
		Objectives:  objectives,
	}, Model(parser))

	return c.define(name, sv)
}

func (c *MetricsClient) define(name string, metrics interface{}) error {
	if _, exists := c.AllMetrics.LoadOrStore(name, metrics); exists {
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

	counter, err := cv.GetMetricWith(parser.ParseToLabels())
	if err != nil {
		return err
	}

	counter.Add(value)
	return nil
}

func (c *MetricsClient) EmitGauge(name string, value float64, parser LabelsParser) error {
	metrics, exists := c.AllMetrics.Load(name)
	if !exists {
		return ErrMetricsNotExists
	}

	gv, ok := metrics.(*prometheus.GaugeVec)
	if !ok {
		return ErrMetricsTypeNotMatches
	}

	gauge, err := gv.GetMetricWith(parser.ParseToLabels())
	if err != nil {
		return err
	}

	gauge.Set(value)
	return nil
}

func (c *MetricsClient) EmitHistogram(name string, value float64, parser LabelsParser) error {
	metrics, exists := c.AllMetrics.Load(name)
	if !exists {
		return ErrMetricsNotExists
	}

	hv, ok := metrics.(*prometheus.HistogramVec)
	if !ok {
		return ErrMetricsTypeNotMatches
	}

	histogram, err := hv.GetMetricWith(parser.ParseToLabels())
	if err != nil {
		return err
	}

	histogram.Observe(value)
	return nil
}

func (c *MetricsClient) EmitSummary(name string, value float64, parser LabelsParser) error {
	metrics, exists := c.AllMetrics.Load(name)
	if !exists {
		return ErrMetricsNotExists
	}

	sv, ok := metrics.(*prometheus.SummaryVec)
	if !ok {
		return ErrMetricsTypeNotMatches
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

var DefaultMetricsClient = NewMetricsClient()

func SetMetricsClient(c *MetricsClient) {
	if c != nil {
		DefaultMetricsClient = c
	}
}

func EmitCounter(name string, value float64, parser LabelsParser) error {
	return DefaultMetricsClient.EmitCounter(name, value, parser)
}

func EmitGauge(name string, value float64, parser LabelsParser) error {
	return DefaultMetricsClient.EmitGauge(name, value, parser)
}

func EmitHistogram(name string, value float64, parser LabelsParser) error {
	return DefaultMetricsClient.EmitHistogram(name, value, parser)
}

func EmitHistogramTimer(name string, t time.Time, parser LabelsParser) error {
	return EmitHistogram(name, float64(time.Since(t).Nanoseconds()/1000), parser)
}

func EmitSummary(name string, value float64, parser LabelsParser) error {
	return DefaultMetricsClient.EmitSummary(name, value, parser)
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
