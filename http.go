package promgo

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// Render ...
func Render() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		m := defaultRegistry.Collect()
		metrics := Metrics(m)
		sort.Sort(metrics)
		metricGroup := make(map[*Desc][]Metric)
		for _, metric := range metrics {
			if _, ok := metricGroup[metric.Desc]; !ok {
				metricGroup[metric.Desc] = make([]Metric, 0)
			}
			metricGroup[metric.Desc] = append(metricGroup[metric.Desc], metric)
		}

		lines := []string{}
		for desc, group := range metricGroup {
			lines = append(lines, fmt.Sprintf(`# HELP %s %s`, desc.GetName(), desc.GetHelp()))
			lines = append(lines, fmt.Sprintf(`# TYPE %s %s`, desc.GetName(), desc.GetType()))

			for _, m := range group {
				vv := make([]string, 0, len(m.ConstLabels))
				if m.ConstLabels == nil {
					lines = append(lines, fmt.Sprintf(`%s %.2f`, m.GetName(), m.GetValue()))
					continue
				}

				for k, v := range m.ConstLabels {
					vv = append(vv, fmt.Sprintf(`%s=%s`, k, v))
				}
				lines = append(lines, fmt.Sprintf(`%s{%s} %.2f`, m.GetName(), strings.Join(vv, `,`), m.GetValue()))
			}
		}
		lines = append(lines, "\n")
		html := strings.Join(lines, "\n")

		rw.Header().Set(`Content-Type`, `text/plain;charset=utf-8`)
		fmt.Fprint(rw, html)
	}
}
