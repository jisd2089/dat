package databox

import (
	"fmt"

	"drcs/common/pinyin"
	"sync"
)

// 数据流产品种类列表
type DataBoxSpecies struct {
	list     []*DataBox
	hash     map[string]*DataBox
	hashPool map[string]*sync.Pool
	sorted   bool
}

// 全局数据流产品种类实例
var Species = &DataBoxSpecies{
	list:     []*DataBox{},
	hash:     map[string]*DataBox{},
	hashPool: make(map[string]*sync.Pool),
}

// 向数据流产品种类清单添加新种类
func (self *DataBoxSpecies) Add(b *DataBox) *DataBox {
	name := b.Name
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

func (self *DataBoxSpecies) AddPool(b DataBox) *DataBox {

	boxPool := &sync.Pool{
		New: func() interface{} {
			return &DataBox{
				status: b.status,
				dnames: make([]string, 0, maxParam),
				dvalues: make([]string, 0, maxParam),
			}
		},
	}
	name := b.Name
	for i := 2; true; i++ {
		if _, ok := self.hashPool[name]; !ok {
			b.Name = name
			self.hashPool[b.Name] = boxPool
			break
		}
		//name = fmt.Sprintf("%s(%d)", b.Name, i)
	}
	return &b
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
		pinyin.SortInitials(initials) // TODO 定制化排序方法
		for i := 0; i < l; i++ {
			self.list[i] = newlist[initials[i]]
		}
		self.sorted = true
	}
	return self.list
}

func (self *DataBoxSpecies) GetByName(name string) *DataBox {
	return self.hash[name]
	//return self.hashPool[name].Get().(*DataBox)
}
