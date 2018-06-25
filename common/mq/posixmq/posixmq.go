package posixmq

import (
	"mime/multipart"
	"strconv"
	"bytes"
	"github.com/oxtoacart/bpool"
)

/**
    Author: luzequan
    Created: 2018-05-28 18:07:36
*/
type (
	RecordType string // RecordType 记录类型

	FlowStatus string // FlowStatus 流程状态

	StepStatus string // StepStatus 步骤状态
)

var _bufferPool = bpool.NewSizedBufferPool(1000, 1024)

const (
	// DELI 分隔符
	DELI = "|@|"
	DOT  = "."
)

// Record 业务日志记录结构
type Record struct {
	recordSeq string      // 日志序号
	rdate     string      // 记录日期 YYYYMMDD
	rtime     string      // 记录时间 HHMMSS
	DemMemID  string      // 需方会员
	SupMemID  string      // 供方会员
	TaskID    string      // 任务编号
	SeqNo     string      // 业务流水号
	ExID      string      // ExID
	DmpSeqNo  string      // DMP流水号
	Type      string      // 记录类型
	SuccCount string      // 成功计数
	Status    string      // 流程状态
	UsedTime  int         // 处理耗时
	ErrCode   string      // 错误码
	stepCount int         // 步骤数
	StepInfos []*StepInfo // 步骤信息

	FileAddr *multipart.File
}

// StepInfo 步骤信息
type StepInfo struct {
	No        int    // 步骤号
	MemID     string // 处理会员
	Status    string // 状态
	Signature string // 签名信息
}

func flatRecordIntoBuffer(record *Record, buffer *bytes.Buffer) {
	buffer.WriteString(string(record.Status))
	buffer.WriteString(DELI)
	buffer.WriteString(record.rdate)
	buffer.WriteString(DELI)
	buffer.WriteString(record.rtime)
	buffer.WriteString(DELI)
	buffer.WriteString(record.DemMemID)
	buffer.WriteString(DELI)
	buffer.WriteString(record.SupMemID)
	buffer.WriteString(DELI)
	buffer.WriteString(record.TaskID)
	buffer.WriteString(DELI)
	buffer.WriteString(record.SeqNo)
	buffer.WriteString(DELI)
	buffer.WriteString(record.ExID)
	buffer.WriteString(DELI)
	buffer.WriteString(record.DmpSeqNo)
	buffer.WriteString(DELI)
	buffer.WriteString(string(record.Type))
	buffer.WriteString(DELI)
	buffer.WriteString(record.SuccCount)
	buffer.WriteString(DELI)
	buffer.WriteString(strconv.Itoa(record.UsedTime))
	buffer.WriteString(DELI)
	buffer.WriteString(record.ErrCode)
	buffer.WriteString(DELI)
	buffer.WriteString(strconv.Itoa(record.stepCount))

	for _, stepInfo := range record.StepInfos {
		buffer.WriteString(DELI)
		flatStepInfoIntoBuffer(stepInfo, buffer)
	}
}

func flatStepInfoIntoBuffer(stepInfo *StepInfo, buffer *bytes.Buffer) {
	buffer.WriteString(strconv.Itoa(stepInfo.No))
	buffer.WriteString(DELI)
	buffer.WriteString(stepInfo.MemID)
	buffer.WriteString(DELI)
	buffer.WriteString(string(stepInfo.Status))
	buffer.WriteString(DELI)
	buffer.WriteString(stepInfo.Signature)
}
