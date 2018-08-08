package smartsail

/**
    Author: luzequan
    Created: 2018-08-08 11:18:16
*/

import (
	. "drcs/dep/nodelib/crp/common"
)

const (
	SMARTSAIL_SUCC         = "100"     // 成功
	SMARTSAIL_PARSE_FAILED = "6000020" // 业务数据解析失败
	SMARTSAIL_PARAM_ERR    = "6001001" // 业务参数错误
	SMARTSAIL_REQ_TIMEOUT  = "6001002" // 请求超时
	SMARTSAIL_SYS_BUSY     = "8888888" // 系统繁忙
	SMARTSAIL_UNKNOWN_ERR  = "000000"  // 未知错误

	SMARTSAIL_MSG = "null"

)

var SMARTSAILMapCenterCode = map[string]string{
	SMARTSAIL_SUCC:         CenterCodeSucc,
	SMARTSAIL_PARSE_FAILED: CenterCodeFormat,
	SMARTSAIL_PARAM_ERR:    CenterCodeFormat,
	SMARTSAIL_REQ_TIMEOUT:  CenterCodeReqData,
	SMARTSAIL_SYS_BUSY:     CenterCodeReqData,
	SMARTSAIL_UNKNOWN_ERR:  CenterCodeOther,
}

func GetCenterCodeFromSMARTSAIL(code string) string {
	str, ok := SMARTSAILMapCenterCode[code]
	if ok {
		return str
	}
	return CenterCodeOther
}
