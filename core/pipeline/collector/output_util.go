package collector

import (
	"github.com/henrylee2cn/pholcus/logs"
)

// 主命名空间相对于数据库名，不依赖具体数据内容，可选
func (self *Collector) namespace() string {
	if self.DataBox.Namespace == nil {
		if self.DataBox.GetSubName() == "" {
			return self.DataBox.GetName()
		}
		return self.DataBox.GetName() + "__" + self.DataBox.GetSubName()
	}
	return self.DataBox.Namespace(self.DataBox)
}

// 次命名空间相对于表名，可依赖具体数据内容，可选
func (self *Collector) subNamespace(dataCell map[string]interface{}) string {
	if self.DataBox.SubNamespace == nil {
		return dataCell["RuleName"].(string)
	}
	defer func() {
		if p := recover(); p != nil {
			logs.Log.Error("subNamespace: %v", p)
		}
	}()
	return self.DataBox.SubNamespace(self.DataBox, dataCell)
}

// 下划线连接主次命名空间
func joinNamespaces(namespace, subNamespace string) string {
	if namespace == "" {
		return subNamespace
	} else if subNamespace != "" {
		return namespace + "__" + subNamespace
	}
	return namespace
}
