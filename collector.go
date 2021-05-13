package promgo

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
)

// Collector ...
type Collector interface {
	// Collect 收集器
	Collect(ch chan<- *MetricErr)

	// Describe 指标
	Describe() *Desc
}

const (
	// CollectorPrefix ...
	CollectorPrefix = `prometheus`
	// FieldSeperator 域分隔符
	FieldSeperator = `__`
)

type redisCollector struct {
	Collector

	Rdb  redis.Cmdable
	Desc *Desc
}

func (rc redisCollector) Describe() *Desc {
	return rc.Desc
}

func (rc redisCollector) Collect(ch chan<- *MetricErr) {
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

func (rc redisCollector) key() string {
	return fmt.Sprintf(`%s:%s:%s`, CollectorPrefix, rc.Desc.GetType(), rc.Desc.ID())
}

func (rc redisCollector) field(constLables ConstLabels) string {
	vv := make([]string, 0, len(rc.Desc.Labels))
	for _, l := range rc.Desc.Labels {
		if v, ok := constLables[l]; ok {
			vv = append(vv, v)
		}
	}

	return strings.Join(vv, FieldSeperator)
}

func (rc redisCollector) constLabels(field string) map[string]string {
	vv := strings.Split(field, FieldSeperator)
	cl := make(map[string]string)
	for i, l := range rc.Desc.Labels {
		cl[l] = vv[i]
	}
	return cl
}
