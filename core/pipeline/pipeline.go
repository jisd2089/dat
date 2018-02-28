// 数据收集
package pipeline

import (
	"drcs/core/pipeline/collector/data"
	"drcs/core/pipeline/collector"
	"drcs/core/databox"
)

// 数据拆包/核验管道
type Pipeline interface {
	Start()                          //启动
	Stop()                           //停止
	CollectData(data.DataCell) error //收集数据单元
	CollectFile(data.FileCell) error //收集文件
}

func New(df *databox.DataBox) Pipeline {
	return collector.NewCollector(df)
}
