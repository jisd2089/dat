package collector

import (
	"drcs/common/mq/posixmq"
	"fmt"
	"drcs/common/mq"
)

/**
    Author: luzequan
    Created: 2017-12-28 17:07:15
*/

/************************ posixmq 输出 ***************************/
func init() {
	DataOutput["posixmq"] = func(self *Collector) (err error) {
		defer func() {
			if p := recover(); p != nil {
				err = fmt.Errorf("%v", p)
			}
		}()

		recordor, err := mq.GetRecordor()
		if err != nil {
			err = fmt.Errorf("get mq recordor err: ", err.Error())
			return
		}

		for _, dataCell := range self.dataDocker {
			//fmt.Println("dataCell: ", dataCell)

			data := dataCell["Data"].(map[string]interface{})

			record := &posixmq.Record{}
			if exID, ok := data["exID"]; ok {
				record.ExID = exID.(string) // exid
			}
			if demMemID, ok := data["demMemID"]; ok {
				record.DemMemID = demMemID.(string) // 需方会员
			}
			if supMemID, ok := data["supMemID"]; ok {
				record.SupMemID = supMemID.(string) // 供方会员
			}
			if taskID, ok := data["taskID"]; ok {
				record.TaskID = taskID.(string) // 任务编号
			}
			if seqNo, ok := data["seqNo"]; ok {
				record.SeqNo = seqNo.(string) // DMP流水号
			}
			if dmpSeqNo, ok := data["dmpSeqNo"]; ok {
				record.DmpSeqNo = dmpSeqNo.(string) // DMP流水号
			}
			if recordType, ok := data["recordType"]; ok {
				record.Type = recordType.(string) // 记录类型
			}
			if succCount, ok := data["succCount"]; ok {
				record.SuccCount = succCount.(string) // 成功计数
			}
			if flowStatus, ok := data["flowStatus"]; ok {
				record.Status = flowStatus.(string) // 流程状态
			}
			if usedTime, ok := data["usedTime"]; ok {
				record.UsedTime = usedTime.(int) // 处理耗时
			}
			if errCode, ok := data["errCode"]; ok {
				record.ErrCode = errCode.(string) // 错误码
			}
			if stepInfoM, ok := data["stepInfoM"]; ok {

				stepInfos := []*posixmq.StepInfo{}

				stepInfoMap := stepInfoM.([]map[string]interface{})
				for _, stepInfo := range stepInfoMap {
					s := &posixmq.StepInfo{}
					if no, ok := stepInfo["no"]; ok {
						s.No = no.(int) // 步骤号
					}
					if memID, ok := stepInfo["memID"]; ok {
						s.MemID = memID.(string) // 处理会员
					}
					if status, ok := stepInfo["stepStatus"]; ok {
						s.Status = status.(string) // 状态
					}
					if sign, ok := stepInfo["signature"]; ok {
						s.Signature = sign.(string) // 签名信息
					}
					stepInfos = append(stepInfos, s)
				}

				record.StepInfos = stepInfos // 步骤信息
			}

			recordor.Record(record)
		}

		return
	}
}
