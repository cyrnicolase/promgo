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
		metricGroup := make(map[*Desc]Metrics)
		for _, metric := range metrics {
			if _, ok := metricGroup[metric.Desc]; !ok {
				metricGroup[metric.Desc] = make([]Metric, 0)
			}
			metricGroup[metric.Desc] = append(metricGroup[metric.Desc], metric)
		}

		descs := make(Descs, 0, len(metricGroup))
		for desc := range metricGroup {
			descs = append(descs, desc)
		}
		sort.Sort(descs)

		lines := []string{}
		for _, desc := range descs {
			lines = append(lines, fmt.Sprintf(`# HELP %s %s`, desc.GetName(), desc.GetHelp()))
			lines = append(lines, fmt.Sprintf(`# TYPE %s %s`, desc.GetName(), desc.GetType()))

			group := metricGroup[desc]
			sort.Sort(group)
			for _, m := range group {
				vv := make([]string, 0, len(m.ConstLabels))
				if len(m.ConstLabels) == 0 {
					lines = append(lines, fmt.Sprintf(`%s %.2f`, m.GetName(), m.GetValue()))
					continue
				}

				for _, l := range m.Desc.Labels {
					vv = append(vv, fmt.Sprintf(`%s="%s"`, l, m.ConstLabels[l]))
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
