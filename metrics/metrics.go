package metrics

import (
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/songzhaoliang/echotool/util"
)

func Register(r *echo.Echo, prefixes ...string) {
	register(r.Group(util.GetPrefix(prefixes...)))
}

func RouterRegister(g *echo.Group, prefixes ...string) {
	register(g.Group(util.GetPrefix(prefixes...)))
}

func register(g *echo.Group) {
	g.GET("/metrics", util.WrapH(promhttp.Handler()))
}

func Model(p LabelsParser) (keys []string) {
	labels := p.ParseToLabels()
	keys = make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	return
}

type MetricsClient struct {
	Namespace string
	Labels    prometheus.Labels
}
