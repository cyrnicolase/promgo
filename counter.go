package promgo

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

// Counter 计数器
type Counter interface {
	Collector

	Inc(context.Context, ConstLabels)
	IncBy(context.Context, float64, ConstLabels)
}

// CounterOptions ...
type CounterOptions CollectorOptions

// redisCounter ...
type redisCounter struct {
	redisCollector
}

// Incr 自增量，步长为1
func (rc redisCounter) Inc(ctx context.Context, constLables ConstLabels) {
	rc.IncBy(ctx, 1, constLables)
}

// IncrBy 指定增量的增长; 增量 v 必须是一个非负数; 这里没有做校验。。。
func (rc redisCounter) IncBy(ctx context.Context, v float64, constLabels ConstLabels) {
	rc.Rdb.HIncrByFloat(ctx, rc.key(), rc.field(constLabels), v)
}

// Collect 采集数据
func (rc redisCounter) Collect(ch chan<- *MetricErr) {
	values, err := rc.Rdb.HGetAll(context.Background(), rc.key()).Result()
	if err != nil {
		ch <- NewMetricErr(nil, err)
		return
	}

	for field, value := range values {
		v, _ := strconv.ParseFloat(value, 64)
		constlabels := rc.constLabels(field)
		metric := NewMetric(rc.Desc, v, constlabels)
		ch <- NewMetricErr(metric, nil)
	}
}

// Metric
func (rc redisCounter) Describe() *Desc {
	return rc.Desc
}

// NewCounter ...
func NewCounter(rdb redis.Cmdable, opts CounterOptions) Counter {
	desc := &Desc{
		Namespace: opts.Namespace,
		Name:      opts.Name,
		Help:      opts.Help,
		Type:      CounterValue,
		Labels:    opts.Labels,
	}

	rc := redisCollector{
		Rdb:  rdb,
		Desc: desc,
	}

	return &redisCounter{redisCollector: rc}
}
