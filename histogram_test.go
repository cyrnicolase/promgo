package promgo

import (
	"reflect"
	"testing"

	"github.com/go-redis/redis/v8"
)

func newTestHistogram() Histogram {
	rdb := &redis.Client{}
	return NewHistogram(rdb, HistogramOptions{}, nil)
}

func TestLinear(t *testing.T) {
	h := newTestHistogram()
	buckets := h.Linear(10, 5, 10)
	exp := []float64{10, 15, 20, 25, 30, 35, 40, 45, 50, 55}
	if !reflect.DeepEqual(buckets, exp) {
		t.Fatalf("Histogram linear not correct;\nexp=%v\nact=%v", exp, buckets)
	}
}

func TestExponential(t *testing.T) {
	h := newTestHistogram()
	buckets := h.Exponential(1.5, 2, 5)
	exp := []float64{1.5, 3, 6, 12, 24}

	if !reflect.DeepEqual(buckets, exp) {
		t.Fatalf("Histogram exponential not correct;\nexp=%v\nact=%v", exp, buckets)
	}
}
