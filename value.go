package promgo

const (
	// CounterValue 计数器数字
	CounterValue ValueType = `counter`
	// GaugeValue 面版
	GaugeValue ValueType = `gauge`
	// HistogramValue 直方图
	HistogramValue ValueType = `histogram`
)

// ValueType ...
type ValueType string

// String 格式转换
func (vt ValueType) String() string {
	return string(vt)
}
