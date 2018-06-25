package mq

import "drcs/common/mq/posixmq"

/**
    Author: luzequan
    Created: 2017-12-28 16:30:25
*/



// 记录数据资产接口
type Recordor interface {
	Record(record *posixmq.Record)
}

// 数据资产清单
type AssetList struct {
}
