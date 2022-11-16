package main

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/songzhaoliang/echotool"
	"github.com/songzhaoliang/echotool/metrics"
)

var ThroughputCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "echotool",
	Name:      "throughput",
}, metrics.Model(&ThroughputTag{}))

type ThroughputTag struct {
	Handler string
}

var _ metrics.LabelsParser = (*ThroughputTag)(nil)

func NewThroughputTag(handler string) *ThroughputTag {
	return &ThroughputTag{
		handler,
	}
}

func (t *ThroughputTag) ParseToLabels() prometheus.Labels {
	return prometheus.Labels{
		"handler": t.Handler,
	}
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	r := echo.New()
	metrics.Register(r)

	e := echotool.NewEngine()

	r.POST("/users", e.EchoHandler(CreateUser))

	r.Start(":1323")
}

func CreateUser(c echo.Context, ec *echotool.Context) {
	ThroughputCounter.With(NewThroughputTag("CreateUser").ParseToLabels()).Add(1)

	user := &User{}
	echotool.New(c, user).JSONBindBody().MustEnd()

	fmt.Printf("user is %+v\n", user)

	ec.Finish(echotool.CodeOKZero, nil)
}
