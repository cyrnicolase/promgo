package promgo

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
)

// Counter 计数器
type Counter interface {
	Collector

	Incr(context.Context, ConstLabels)
	IncrBy(context.Context, float64, ConstLabels)
}

// CounterOptions ...
type CounterOptions struct {
	Namespace string
	Name      string
	Help      string
	Labels    []string
}

// redisCounter ...
type redisCounter struct {
	Rdb  redis.Cmdable
	Desc *Desc
}

// Incr 自增量，步长为1
func (rc redisCounter) Incr(ctx context.Context, constLables ConstLabels) {
	rc.IncrBy(ctx, 1, constLables)
}

// IncrBy 指定增量的增长; 增量 v 必须是一个非负数; 这里没有做校验。。。
func (rc redisCounter) IncrBy(ctx context.Context, v float64, constLabels ConstLabels) {
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
		var constlabels map[string]string
		if field != "" {
			vv := strings.Split(field, `__`)
			constlabels = make(map[string]string)
			for i, l := range rc.Desc.Labels {
				constlabels[l] = vv[i]
			}
		}

		metric := NewMetric(rc.Desc, v, constlabels)
		ch <- NewMetricErr(metric, nil)
	}
}

// Metric
func (rc redisCounter) Describe() *Desc {
	return rc.Desc
}

func (rc redisCounter) key() string {
	return fmt.Sprintf(`prometheus:counter:%s`, rc.Desc.ID())
}

func (rc redisCounter) field(constLables ConstLabels) string {
	vv := make([]string, 0, len(rc.Desc.Labels))
	for _, l := range rc.Desc.Labels {
		if v, ok := constLables[l]; ok {
			vv = append(vv, v)
		}
	}

	return strings.Join(vv, `__`) // 使用双下划线连接标签值
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

	return &redisCounter{
		Rdb:  rdb,
		Desc: desc,
	}
}
