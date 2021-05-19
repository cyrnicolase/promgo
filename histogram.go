package promgo

import (
	"context"
	"sort"
	"strconv"

	"github.com/go-redis/redis/v8"
)

const (
	// 小于等于
	leLabel = `le`
)

var (
	// DefaultBuckets 默认
	DefaultBuckets = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
)

// Histogram 直方图
type Histogram interface {
	Collector

	Observe(context.Context, float64, ConstLabels)
	Linear(start, width float64, count int) []float64
	Exponential(start, factor float64, count int) []float64
}

// HistogramOptions 直方图参数
type HistogramOptions CollectorOptions

type redisHistogram struct {
	redisCollector

	Buckets []float64
}

// Observe 观察者
func (rh redisHistogram) Observe(ctx context.Context, v float64, constLables ConstLabels) {
	bucket := rh.findBucket(v)
	constLables[`le`] = strconv.FormatFloat(rh.Buckets[bucket], 'f', -1, 64)

	rh.Rdb.HIncrByFloat(ctx, rh.key(), rh.field(constLables), 1)
}

// Linear 线性buckets
func (rh *redisHistogram) Linear(start, width float64, count int) []float64 {
	if count < 1 {
		panic(`Linear needs a positive count`)
	}

	buckets := make([]float64, count)
	for i := range buckets {
		buckets[i] = start
		start += width
	}
	rh.Buckets = buckets
	return buckets
}

func (rh *redisHistogram) Exponential(start, factor float64, count int) []float64 {
	if count < 1 {
		panic(`Exponential needs a positive count`)
	}
	if start <= 0 {
		panic(`Exponential needs a positive start value`)
	}
	if factor <= 1 {
		panic(`Exponential needs a factor greater than 1`)
	}
	buckets := make([]float64, count)
	for i := range buckets {
		buckets[i] = start
		start *= factor
	}
	rh.Buckets = buckets
	return buckets
}

// 找到对应的bucket
func (rh redisHistogram) findBucket(v float64) int {
	return sort.SearchFloat64s(rh.Buckets, v)
}

// NewHistogram ...
func NewHistogram(rdb redis.Cmdable, opts HistogramOptions, buckets []float64) Histogram {
	desc := &Desc{
		Namespace: opts.Namespace,
		Name:      opts.Name,
		Help:      opts.Help,
		Type:      HistogramValue,
		Labels:    []string{leLabel},
	}

	rc := redisCollector{
		Rdb:  rdb,
		Desc: desc,
	}

	if len(buckets) == 0 {
		buckets = DefaultBuckets
	}
	sort.Float64s(buckets)

	return &redisHistogram{redisCollector: rc, Buckets: buckets}
}
