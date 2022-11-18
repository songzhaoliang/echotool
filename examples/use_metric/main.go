package main

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/songzhaoliang/echotool"
	"github.com/songzhaoliang/echotool/metric"
)

const (
	MThroughput = "throughput"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	InitMetrics()

	r := echo.New()
	metric.Register(r)

	e := echotool.NewEngine()

	r.POST("/users", e.EchoHandler(CreateUser))

	r.Start(":1323")
}

func CreateUser(c echo.Context, ec *echotool.Context) {
	metric.EmitCounter(MThroughput, 1, NewThroughputLabels(ec.GetHandlerName()))

	user := &User{}
	echotool.New(c, user).JSONBindBody().MustEnd()

	fmt.Printf("user is %+v\n", user)

	ec.Finish(echotool.CodeOKZero, nil)
}

func InitMetrics() {
	c := metric.NewMetricClient(metric.WithNamespace("echotool"))
	c.DefineCounter(MThroughput, &ThroughputLabels{})
	metric.SetMetricClient(c)
}

type ThroughputLabels struct {
	Handler string
}

var _ metric.LabelsParser = (*ThroughputLabels)(nil)

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
