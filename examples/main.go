package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/cyrnicolase/promgo"
	"github.com/go-redis/redis/v8"
)

var (
	// RequestTotal ...
	RequestTotal promgo.Counter
	// AllRequest ...
	AllRequest promgo.Counter

	cpuDegree promgo.Gauge
)

func init() {
	rand.Seed(time.Now().UnixNano())
	rdb := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})

	RequestTotal = promgo.NewCounter(rdb, promgo.CounterOptions{
		Name:   `request_total`,
		Help:   `接口请求总数`,
		Labels: []string{`method`, `endpoint`},
	})
	AllRequest = promgo.NewCounter(rdb, promgo.CounterOptions{
		Name: `all_request`,
		Help: `接口请求总数`,
	})
	cpuDegree = promgo.NewGauge(rdb, promgo.GaugeOptions{
		Name: `cpu_gauge`,
		Help: `cpu 压力`,
	})

	promgo.GetDefaultRegistry().MustRegister(RequestTotal)
	promgo.GetDefaultRegistry().MustRegister(AllRequest)
	promgo.GetDefaultRegistry().MustRegister(cpuDegree)

}

func main() {
	go func() {
		for {
			time.Sleep(time.Second)
			v := rand.Float64()
			cpuDegree.Set(context.Background(), v, nil)
		}
	}()
	http.HandleFunc(`/hello`, func(rw http.ResponseWriter, r *http.Request) {
		AllRequest.Inc(r.Context(), nil)
		RequestTotal.Inc(r.Context(), promgo.ConstLabels{
			`method`:   r.Method,
			`endpoint`: r.URL.Path,
		})

		fmt.Fprintf(rw, `hello`)
	})
	http.HandleFunc(`/index`, func(rw http.ResponseWriter, r *http.Request) {
		AllRequest.Inc(r.Context(), nil)
		RequestTotal.Inc(r.Context(), promgo.ConstLabels{
			`method`:   r.Method,
			`endpoint`: r.URL.Path,
		})

		fmt.Fprint(rw, `index`)
	})
	http.HandleFunc(`/metrics`, promgo.Render())
	http.ListenAndServe(`:8774`, nil)
}
