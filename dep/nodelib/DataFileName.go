package nodelib

import (
	"strings"
	"errors"
	"fmt"
)

/**
    Author: luzequan
    Created: 2017-12-06 18:55:57
*/

type DataFileName struct {
	MemId	string // 会员编号

	JobId   string // 工单号
	IdType  string // id类型
	BatchNo string // 批次号
	FileNo  string // 文件排序号
	Suffix  string // 后缀名
}

func (d *DataFileName) ParseAndValidFileName(fileName string) error {
	fullName := strings.Split(fileName, ".")
	if len(fullName) != 2 {
		return errors.New("The pattern of fileName is error.")
	}
	d.Suffix = fullName[1]

	prefixName := fullName[0]
	names := strings.Split(prefixName, "_")


	switch len(names) {
	case 4:
		d.JobId = names[0]
		d.IdType = names[1]
		d.BatchNo = names[2]
		d.FileNo = names[3]
	case 5:
		d.MemId = names[0]
		d.JobId = names[1]
		d.IdType = names[2]
		d.BatchNo = names[3]
		d.FileNo = names[4]
	default:
		return errors.New("The pattern of prefixName is error.")
	}

	return nil
}

func (d *DataFileName) GetFullName() string {
	if d.MemId != "" {
		return fmt.Sprintf("%s_%s_%s_%s_%s.%s", d.MemId, d.JobId, d.IdType, d.BatchNo, d.FileNo, d.Suffix)
	}
	return fmt.Sprintf("%s_%s_%s_%s.%s", d.JobId, d.IdType, d.BatchNo, d.FileNo, d.Suffix)
}

func (d *DataFileName) GetPrefixName() string {
	if d.MemId != "" {
		return fmt.Sprintf("%s_%s_%s_%s_%s", d.MemId, d.JobId, d.IdType, d.BatchNo, d.FileNo)
	}
	return fmt.Sprintf("%s_%s_%s_%s", d.JobId, d.IdType, d.BatchNo, d.FileNo)
}

func (d *DataFileName) GetCacheKey() string {
	return fmt.Sprintf("%s_%s_%s_%s", d.JobId, d.IdType, d.BatchNo, d.FileNo)
}

func (d *DataFileName) GetBatchSendName() string {
	return fmt.Sprintf("%s_%s_%s_batchsend", d.JobId, d.IdType, d.BatchNo)
}

func (d *DataFileName) GetBatchRecName() string {
	return fmt.Sprintf("%s_%s_%s_batchrec", d.JobId, d.IdType, d.BatchNo)
}

func (d *DataFileName) GetOrderId_BatchNo() string{
	return fmt.Sprintf("%s_%s", d.JobId, d.BatchNo)
}

type HeadRecord struct {
	DataType   string //批量文件中数据类型
	StatPeriod string //统计周期
	FileCount  int    //批量文件数
	ConnObjIds []string
}
