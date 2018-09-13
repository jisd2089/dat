package taifinance

/**
    Author: luzequan
    Created: 2018-09-03 14:35:48
*/

import (
	. "drcs/dep/nodelib/fusion/common"
)

const (
	TAIFIN_SUCC              = 0   // API 成功返回
	TAIFIN_CALL_ERR          = 1   // API 调用异常
	TAIFIN_NOT_EXIST         = 2   // API 不存在
	TAIFIN_SERVER_ERR        = 3   // 服务器异常
	TAIFIN_PATH_NIL          = 100 // 模型路径不能为空
	TAIFIN_FILE_NOTEXIST     = 101 // 模型文件不存在
	TAIFIN_AMOUNT_NIL        = 104 // instancesAmount 参数不能为空
	TAIFIN_AMOUNT_INT        = 105 // instancesAmount 不是整数
	TAIFIN_ARRAY_NIL         = 106 // instancesArray 参数不能为空
	TAIFIN_ARRAY_INT         = 107 // instancesArray 格式不正确，不是整数数组
	TAIFIN_CALCU_ERR         = 108 // instancesArray 参数中数组元素数量不能被instancesAmount 整除
	TAIFIN_POLICY_ERR        = 109 // 模型缺少对应的生成策略
	TAIFIN_POLICY_ERR2       = 115 // 模型缺少对应的审批策略
	TAIFIN_FILE_NOTEXIST2    = 119 // 模型文件不存在
	TAIFIN_INDEX_ERR         = 120 // 列索引值超出范围
	TAIFIN_CALL_OTHERERR     = 121 // 其它模型调用异常
	TAIFIN_UUID_ERR          = 122 // uuid 参数不能为空
	TAIFIN_MODEL_NOEXIST     = 123 // 模型不存在
	TAIFIN_APIURL_NOEXIST    = 124 // APIURL 不存在
	TAIFIN_CASEIDS_NIL       = 125 // caseIDs 参数不能为空
	TAIFIN_EQUAL_ERR         = 126 // caseIDs 数组长度必须与instancesAmount 相等
	TAIFIN_CHECK_ERR         = 127 // 审批策略不存在
	TAIFIN_GEN_ERR           = 128 // 生成策略不存在
	TAIFIN_MODEL_ONLINE      = 129 // 模型必须上线
	TAIFIN_MODEL_DEL         = 130 // 模型已删除
	TAIFIN_MODEL_NOTALLOW    = 131 // 不允许调用此模型
	TAIFIN_POLICY_NOTALLOW   = 132 // 不允许调用此审批策略
	TAIFIN_SCORECARD_NOEXIST = 133 // 此模型评分卡不存在
	TAIFIN_PARAMS_ERR        = 134 // 参数不合法
)

var TAIFINMapCenterCode = map[int]string{
	TAIFIN_SUCC:              CenterCodeSucc,
	TAIFIN_CALL_ERR:          CenterCodeReqFail,
	TAIFIN_NOT_EXIST:         CenterCodeNoService,
	TAIFIN_SERVER_ERR:        CenterCodeReqFail,
	TAIFIN_PATH_NIL:          CenterCodeReqData,
	TAIFIN_FILE_NOTEXIST:     CenterCodeReqData,
	TAIFIN_AMOUNT_NIL:        CenterCodeReqData,
	TAIFIN_AMOUNT_INT:        CenterCodeReqData,
	TAIFIN_ARRAY_NIL:         CenterCodeReqData,
	TAIFIN_ARRAY_INT:         CenterCodeReqData,
	TAIFIN_CALCU_ERR:         CenterCodeReqData,
	TAIFIN_POLICY_ERR:        CenterCodeReqData,
	TAIFIN_POLICY_ERR2:       CenterCodeReqData,
	TAIFIN_FILE_NOTEXIST2:    CenterCodeReqData,
	TAIFIN_INDEX_ERR:         CenterCodeReqData,
	TAIFIN_CALL_OTHERERR:     CenterCodeReqData,
	TAIFIN_UUID_ERR:          CenterCodeReqData,
	TAIFIN_MODEL_NOEXIST:     CenterCodeReqData,
	TAIFIN_APIURL_NOEXIST:    CenterCodeReqData,
	TAIFIN_CASEIDS_NIL:       CenterCodeReqData,
	TAIFIN_EQUAL_ERR:         CenterCodeReqData,
	TAIFIN_CHECK_ERR:         CenterCodeReqData,
	TAIFIN_GEN_ERR:           CenterCodeReqData,
	TAIFIN_MODEL_ONLINE:      CenterCodeReqData,
	TAIFIN_MODEL_DEL:         CenterCodeReqData,
	TAIFIN_MODEL_NOTALLOW:    CenterCodeReqData,
	TAIFIN_POLICY_NOTALLOW:   CenterCodeReqData,
	TAIFIN_SCORECARD_NOEXIST: CenterCodeReqData,
	TAIFIN_PARAMS_ERR:        CenterCodeReqData,
}

func GetCenterCodeFromTAIFIN(code int) string {
	str, ok := TAIFINMapCenterCode[code]
	if ok {
		return str
	}
	return CenterCodeOther
}
