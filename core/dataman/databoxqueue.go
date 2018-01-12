package dataman

/**
    Author: luzequan
    Created: 2017-12-27 14:36:53
*/
import (
	. "dat/core/databox"
	"dat/common/util"
	"github.com/henrylee2cn/pholcus/logs"
)

type (
	DataBoxQueue interface {
		Reset()                    // 重置清空队列
		Add(box *DataBox)          // 添加一个
		AddChan(box *DataBox)      // 添加一个入Channel
		AddAll([]*DataBox)         // 添加所有
		AddAllChan([]*DataBox)     // 添加所有入Channel
		AddKeyins(string)          // 为队列成员遍历添加Keyin属性，但前提必须是队列成员未被添加过keyin
		GetByIndex(int) *DataBox   // 通过索引查找
		GetByName(string) *DataBox // 通过名称查找
		GetAll() []*DataBox        // 获取所有
		GetOne() *DataBox          // 从Channel取出一个DataBox
		Len() int                  // 返回队列长度
	}
	dbq struct {
		idInc       *util.AutoInc
		dataBoxChan chan *DataBox
		list        []*DataBox
	}
)

func NewDataBoxQueue() DataBoxQueue {
	return &dbq{
		idInc:       util.NewAutoInc(10000, 1),
		dataBoxChan: make(chan *DataBox),
		list:        []*DataBox{},
	}
}

func (q *dbq) Reset() {
	q.dataBoxChan = make(chan *DataBox)
	q.list = []*DataBox{}
}

func (q *dbq) Add(df *DataBox) {
	df.SetId(q.Len())
	q.list = append(q.list, df)
}

func (q *dbq) AddChan(df *DataBox) {
	df.SetId(q.idInc.Id())
	q.dataBoxChan <- df
}

func (q *dbq) AddAll(list []*DataBox) {
	for _, v := range list {
		q.Add(v)
	}
}

func (q *dbq) AddAllChan(list []*DataBox) {
	for _, v := range list {
		q.AddChan(v)
	}
}

func (q *dbq) GetOne() *DataBox {
	return <-q.dataBoxChan
}

// 添加keyin，遍历DataBox队列得到新的队列（已被显式赋值过的DataBox将不再重新分配Keyin）
func (q *dbq) AddKeyins(keyins string) {
	keyinSlice := util.KeyinsParse(keyins)
	if len(keyinSlice) == 0 {
		return
	}

	unit1 := []*DataBox{} // 不可被添加自定义配置的DataBox
	unit2 := []*DataBox{} // 可被添加自定义配置的DataBox
	for _, v := range q.GetAll() {
		if v.GetKeyin() == KEYIN {
			unit2 = append(unit2, v)
			continue
		}
		unit1 = append(unit1, v)
	}

	if len(unit2) == 0 {
		logs.Log.Warning("本批任务无需填写自定义配置！\n")
		return
	}

	q.Reset()

	for _, keyin := range keyinSlice {
		for _, v := range unit2 {
			v.Keyin = keyin
			nv := *v
			q.Add((&nv).Copy())
		}
	}
	if q.Len() == 0 {
		q.AddAll(append(unit1, unit2...))
	}

	q.AddAll(unit1)
}

func (q *dbq) GetByIndex(idx int) *DataBox {
	return q.list[idx]
}

func (q *dbq) GetByName(n string) *DataBox {
	for _, sp := range q.list {
		if sp.GetName() == n {
			return sp
		}
	}
	return nil
}

func (q *dbq) GetAll() []*DataBox {
	return q.list
}

func (q *dbq) Len() int {
	return len(q.list)
}
