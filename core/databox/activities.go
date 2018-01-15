package databox

import (
	"fmt"
	"dat/common/pinyin"
	"strconv"
)

/**
    Author: luzequan
    Created: 2018-01-12 13:16:10
*/
// DataBox活跃列表
type DataBoxActivites struct {
	list   []*DataBox
	hash   map[string]*DataBox
	sorted bool
}

// 全局DataBox活跃列表实例
var Activites = &DataBoxActivites{
	list: []*DataBox{},
	hash: map[string]*DataBox{},
}

// 向DataBox活跃列表清单添加新种类
func (self *DataBoxActivites) Add(b *DataBox) *DataBox {
	name := b.Name + "_" + strconv.Itoa(b.PairDataBoxId)
	for i := 2; true; i++ {
		if _, ok := self.hash[name]; !ok {
			b.Name = name
			self.hash[b.Name] = b
			break
		}
		name = fmt.Sprintf("%s(%d)", b.Name, i)
	}
	b.Name = name
	self.list = append(self.list, b)
	return b
}

// 从DataBox活跃列表清单移除DataBox
func (self *DataBoxActivites) Remove(b *DataBox) *DataBoxActivites {
	delete(self.hash, b.Name)
	return self
}

// 获取全部DataBox活跃列表
func (self *DataBoxActivites) Get() []*DataBox {
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

func (self *DataBoxActivites) GetByName(name string) *DataBox {
	return self.hash[name]
}