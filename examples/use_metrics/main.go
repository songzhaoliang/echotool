package main

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/songzhaoliang/echotool"
	"github.com/songzhaoliang/echotool/metrics"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	InitMetrics()

	r := echo.New()
	metrics.Register(r)

	e := echotool.NewEngine()

	r.POST("/users", e.EchoHandler(CreateUser))

	r.Start(":1323")
}

func CreateUser(c echo.Context, ec *echotool.Context) {
	metrics.EmitCounter("throughput", 1, NewThroughputLabels(ec.GetHandlerName()))

	user := &User{}
	echotool.New(c, user).JSONBindBody().MustEnd()

	fmt.Printf("user is %+v\n", user)

	ec.Finish(echotool.CodeOKZero, nil)
}

func InitMetrics() {
	c := metrics.NewMetricsClient(metrics.WithNamespace("echotool"))
	c.DefineCounter("throughput", &ThroughputLabels{})
	metrics.SetMetricsClient(c)
}

type ThroughputLabels struct {
	Handler string
}

var _ metrics.LabelsParser = (*ThroughputLabels)(nil)

func NewThroughputLabels(handler string) *ThroughputLabels {
	return &ThroughputLabels{
		Handler: handler,
	}
}

func (ls *ThroughputLabels) ParseToLabels() map[string]string {
	return map[string]string{
		"handler": ls.Handler,
	}
}
