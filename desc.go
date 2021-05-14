package promgo

import (
	"fmt"
	"strings"
)

// Desc 指标描述信息
type Desc struct {
	Namespace string
	Name      string
	Help      string
	Type      ValueType
	Labels    []string // 标签
}

// GetNamespace ...
func (d Desc) GetNamespace() string {
	return d.Namespace
}

// GetName ...
func (d Desc) GetName() string {
	return d.Name
}

// GetHelp ...
func (d Desc) GetHelp() string {
	return d.Help
}

// GetType ...
func (d Desc) GetType() string {
	return d.Type.String()
}

// ID 唯一标志
func (d Desc) ID() string {
	id := fmt.Sprintf(`%s_%s`, d.Namespace, d.Name)

	return strings.Trim(id, `_`)
}

// Descs ...
type Descs []*Desc

// Less 对比
func (d Descs) Less(i, j int) bool {
	if strings.Compare(d[i].ID(), d[j].ID()) == 1 {
		return false
	}
	return true
}

// Swap 交换
func (d Descs) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Len 长度
func (d Descs) Len() int {
	return len(d)
}
