package promgo

import (
	"fmt"
	"strings"
)

// ConstLabels 代表收集指标的标签 name->value 的映射关系
type ConstLabels map[string]string

// MetricErr 指标以及生成错误
type MetricErr struct {
	Metric *Metric
	Err    error
}

// NewMetricErr 指标以及可能的错误
func NewMetricErr(m *Metric, err error) *MetricErr {
	return &MetricErr{Metric: m, Err: err}
}

// Metric 指标
type Metric struct {
	Desc        *Desc
	Value       float64     // 指标值
	ConstLabels ConstLabels // 常量标签值
}

// NewMetric ...
func NewMetric(desc *Desc, v float64, cl ConstLabels) *Metric {
	return &Metric{
		Desc:        desc,
		Value:       v,
		ConstLabels: cl,
	}
}

// ID id
func (m Metric) ID() string {
	return m.Desc.ID()
}

// String ...
func (m Metric) String() string {
	kk := make([]string, 0, len(m.ConstLabels))
	for _, l := range m.Desc.Labels {
		v := m.ConstLabels[l]
		kk = append(kk, fmt.Sprintf(`%s_%s`, l, v))
	}

	return strings.Trim(fmt.Sprintf(`%s_%s_%s`, m.ID(), m.Desc.Type, strings.Join(kk, `_`)), `_`)
}

// GetFQName 获取指标名
func (m Metric) GetFQName() string {
	return m.Desc.ID()
}

// GetHelp 返回指标解释信息
func (m Metric) GetHelp() string {
	return m.Desc.Help
}

// GetType 获取数据类型
func (m Metric) GetType() string {
	return m.Desc.Type.String()
}

// GetValue 获取值
func (m Metric) GetValue() float64 {
	return m.Value
}

// Metrics ...
type Metrics []Metric

// Swap 实现sort.Sort交换两个元素
func (m Metrics) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// Len 实现sort.Sort 获取切片长度
func (m Metrics) Len() int {
	return len(m)
}

// Less 实现sort.Sort比较相邻元素大小
func (m Metrics) Less(i, j int) bool {
	if strings.Compare(m[i].String(), m[j].String()) == 1 {
		return false
	}
	return true
}
