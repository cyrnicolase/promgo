package promgo

import (
	"reflect"
	"testing"

	"github.com/go-redis/redis/v8"
)

func newTestRedisCollector() redisCollector {
	desc := &Desc{
		Namespace: `prometheus`,
		Name:      `api_count`,
		Help:      `api stat count`,
		Type:      CounterValue,
		Labels:    []string{`button`, `channel`},
	}

	return redisCollector{
		Rdb:  &redis.Client{},
		Desc: desc,
	}
}

func TestKey(t *testing.T) {
	c := newTestRedisCollector()
	act := c.key()
	exp := `prometheus:counter:prometheus_api_count`

	if act != exp {
		t.Fatalf("redisCollector method key() is not correct\nexp:%s\nact:%s", exp, act)
	}
}

func TestField(t *testing.T) {
	c := newTestRedisCollector()
	act := c.field(ConstLabels{`button`: `pay`, `channel`: `wx`})
	exp := `pay__wx`

	if act != exp {
		t.Fatalf("redisCollector method field() is not correct\nexp:%s\nact:%s", exp, act)
	}
}

func TestConstLabels(t *testing.T) {
	field := `pay__wx`
	c := newTestRedisCollector()

	act := c.constLabels(field)
	exp := ConstLabels{
		`button`:  `pay`,
		`channel`: `wx`,
	}
	if !reflect.DeepEqual(act, exp) {
		t.Fatalf("redisCollector method constLabels() is not correct\nexp:%s\nact:%s", exp, act)
	}
}
