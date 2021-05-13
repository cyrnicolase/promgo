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
	id := fmt.Sprintf(`%s%s%s`,
		d.Namespace,
		FieldSeperator,
		d.Name,
	)

	return strings.Trim(id, FieldSeperator)
}
