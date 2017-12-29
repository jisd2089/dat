package dataflow

import (
	"fmt"

	"dat/common/pinyin"
)

// 数据流产品种类列表
type DataFlowSpecies struct {
	list   []*DataFlow
	hash   map[string]*DataFlow
	sorted bool
}

// 全局数据流产品种类实例
var Species = &DataFlowSpecies{
	list: []*DataFlow{},
	hash: map[string]*DataFlow{},
}

// 向数据流产品种类清单添加新种类
func (self *DataFlowSpecies) Add(sp *DataFlow) *DataFlow {
	name := sp.Name
	for i := 2; true; i++ {
		if _, ok := self.hash[name]; !ok {
			sp.Name = name
			self.hash[sp.Name] = sp
			break
		}
		name = fmt.Sprintf("%s(%d)", sp.Name, i)
	}
	sp.Name = name
	self.list = append(self.list, sp)
	return sp
}

// 获取全部数据流产品种类
func (self *DataFlowSpecies) Get() []*DataFlow {
	if !self.sorted {
		l := len(self.list)
		initials := make([]string, l)
		newlist := map[string]*DataFlow{}
		for i := 0; i < l; i++ {
			initials[i] = self.list[i].GetName()
			newlist[initials[i]] = self.list[i]
		}
		pinyin.SortInitials(initials)  // TODO 定制化排序方法
		for i := 0; i < l; i++ {
			self.list[i] = newlist[initials[i]]
		}
		self.sorted = true
	}
	return self.list
}

func (self *DataFlowSpecies) GetByName(name string) *DataFlow {
	return self.hash[name]
}
