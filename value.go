package promgo

const (
	// CounterValue 计数器数字
	CounterValue ValueType = `counter`
	// GaugeValue 面版
	GaugeValue ValueType = `gauge`
)

// ValueType ...
type ValueType string

// String 格式转换
func (vt ValueType) String() string {
	return string(vt)
}
