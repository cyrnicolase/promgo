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
	// Collect 收集指标值
	Collect(ch chan<- *MetricErr)

	// Describe 指标描述
	Describe() *Desc
}

// CollectorOptions 参数管理
type CollectorOptions struct {
	Namespace string
	Name      string
	Help      string
	Labels    []string
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

// 在缓存中存储的key值
func (rc redisCollector) key() string {
	return fmt.Sprintf(`%s:%s:%s`, CollectorPrefix, rc.Desc.GetType(), rc.Desc.ID())
}

// Hash类型中的field值
func (rc redisCollector) field(constLabels ConstLabels) string {
	vv := make([]string, 0, len(rc.Desc.Labels))
	for _, l := range rc.Desc.Labels {
		if v, ok := constLabels[l]; ok {
			vv = append(vv, v)
		}
	}

	return strings.Join(vv, FieldSeperator)
}

// 解析Hash类型中的field值，并与指标label进行映射
func (rc redisCollector) constLabels(field string) ConstLabels {
	vv := strings.Split(field, FieldSeperator)
	cl := make(ConstLabels)
	for i, l := range rc.Desc.Labels {
		cl[l] = vv[i]
	}
	return cl
}
