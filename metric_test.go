package promgo

import "testing"

func newTestMetric() *Metric {
	return &Metric{
		Desc: &Desc{
			Namespace: `app`,
			Name:      `api_request_total`,
			Help:      `api request count`,
			Type:      CounterValue,
			Labels:    []string{`method`, `endpoint`},
		},
		Value: 3,
		ConstLabels: ConstLabels{
			`method`:   `GET`,
			`endpoint`: `/index`,
		},
	}
}

func TestString(t *testing.T) {
	m := newTestMetric()
	s := m.String()
	exp := `app_api_request_total_counter_method_GET_endpoint_/index`

	if s != exp {
		t.Fatalf("metric string method not correct\nexp:%s\nact:%s", exp, s)
	}
}

func TestGetFQName(t *testing.T) {
	m := newTestMetric()
	exp := `app_api_request_total`

	if m.GetFQName() != exp {
		t.Fatalf("metric GetFQName method not correct\nexp:%s\nact:%s", exp, m.GetFQName())
	}
}
