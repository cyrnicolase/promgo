package promgo

// Collector ...
type Collector interface {
	// Collect 收集器
	Collect(ch chan<- *MetricErr)

	// Describe 指标
	Describe() *Desc
}
