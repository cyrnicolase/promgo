package promgo

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// Gauge dashboard
type Gauge interface {
	Collector

	Inc(context.Context, ConstLabels)
	IncBy(context.Context, float64, ConstLabels)
	Dec(context.Context, ConstLabels)
	DecBy(context.Context, float64, ConstLabels)
	Set(context.Context, float64, ConstLabels)
}

// GaugeOptions ...
type GaugeOptions CollectorOptions

// redisGauge ...
type redisGauge struct {
	redisCollector
}

// Inc 自增
func (rg redisGauge) Inc(ctx context.Context, cl ConstLabels) {
	rg.IncBy(ctx, 1, cl)
}

func (rg redisGauge) IncBy(ctx context.Context, v float64, cl ConstLabels) {
	rg.Rdb.HIncrByFloat(ctx, rg.key(), rg.field(cl), v)
}

func (rg redisGauge) Dec(ctx context.Context, cl ConstLabels) {
	rg.DecBy(ctx, 1, cl)
}

func (rg redisGauge) DecBy(ctx context.Context, v float64, cl ConstLabels) {
	rg.IncBy(ctx, -1*v, cl)
}

func (rg redisGauge) Set(ctx context.Context, v float64, cl ConstLabels) {
	rg.Rdb.HSet(ctx, rg.key(), rg.field(cl), v)
}

// NewGauge ...
func NewGauge(rdb redis.Cmdable, opts GaugeOptions) Gauge {
	desc := &Desc{
		Namespace: opts.Namespace,
		Name:      opts.Name,
		Help:      opts.Help,
		Type:      GaugeValue,
		Labels:    opts.Labels,
	}

	rc := redisCollector{Rdb: rdb, Desc: desc}

	return redisGauge{redisCollector: rc}
}
