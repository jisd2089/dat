package dataman

/**
    Author: luzequan
    Created: 2017-12-27 14:36:53
*/
import (
	. "dat/core/dataflow"
	"dat/common/util"
	"github.com/henrylee2cn/pholcus/logs"
)

type (
	DataFlowQueue interface {
		Reset()                     //重置清空队列
		Add(flow *DataFlow)         // 添加一个
		AddAll([]*DataFlow)         // 添加所有
		AddKeyins(string)           //为队列成员遍历添加Keyin属性，但前提必须是队列成员未被添加过keyin
		GetByIndex(int) *DataFlow   //通过索引查找
		GetByName(string) *DataFlow //通过名称查找
		GetAll() []*DataFlow        //获取所有
		Len() int                   // 返回队列长度
	}
	dfq struct {
		list []*DataFlow
	}
)

func NewDataFlowQueue() DataFlowQueue {
	return &dfq{
		list: []*DataFlow{},
	}
}

func (dfq *dfq) Reset() {
	dfq.list = []*DataFlow{}
}

func (dfq *dfq) Add(df *DataFlow) {
	df.SetId(dfq.Len())
	dfq.list = append(dfq.list, df)
}

func (dfq *dfq) AddAll(list []*DataFlow) {
	for _, v := range list {
		dfq.Add(v)
	}
}

// 添加keyin，遍历DataFlow队列得到新的队列（已被显式赋值过的DataFlow将不再重新分配Keyin）
func (dfq *dfq) AddKeyins(keyins string) {
	keyinSlice := util.KeyinsParse(keyins)
	if len(keyinSlice) == 0 {
		return
	}

	unit1 := []*DataFlow{} // 不可被添加自定义配置的DataFlow
	unit2 := []*DataFlow{} // 可被添加自定义配置的DataFlow
	for _, v := range dfq.GetAll() {
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

	dfq.Reset()

	for _, keyin := range keyinSlice {
		for _, v := range unit2 {
			v.Keyin = keyin
			nv := *v
			dfq.Add((&nv).Copy())
		}
	}
	if dfq.Len() == 0 {
		dfq.AddAll(append(unit1, unit2...))
	}

	dfq.AddAll(unit1)
}

func (dfq *dfq) GetByIndex(idx int) *DataFlow {
	return dfq.list[idx]
}

func (dfq *dfq) GetByName(n string) *DataFlow {
	for _, sp := range dfq.list {
		if sp.GetName() == n {
			return sp
		}
	}
	return nil
}

func (dfq *dfq) GetAll() []*DataFlow {
	return dfq.list
}

func (dfq *dfq) Len() int {
	return len(dfq.list)
}
