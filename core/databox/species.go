package databox

import (
	"fmt"

	"dat/common/pinyin"
)

// 数据流产品种类列表
type DataBoxSpecies struct {
	list   []*DataBox
	hash   map[string]*DataBox
	sorted bool
}

// 全局数据流产品种类实例
var Species = &DataBoxSpecies{
	list: []*DataBox{},
	hash: map[string]*DataBox{},
}

// 向数据流产品种类清单添加新种类
func (self *DataBoxSpecies) Add(sp *DataBox) *DataBox {
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
func (self *DataBoxSpecies) Get() []*DataBox {
	if !self.sorted {
		l := len(self.list)
		initials := make([]string, l)
		newlist := map[string]*DataBox{}
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

func (self *DataBoxSpecies) GetByName(name string) *DataBox {
	return self.hash[name]
}
