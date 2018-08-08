package edunwang

/**
    Author: luzequan
    Created: 2018-08-08 11:18:16
*/

import (
	. "drcs/dep/nodelib/crp/common"
)

const (
	EDUN_SUCC         = "100"     // 成功
	EDUN_PARSE_FAILED = "6000020" // 业务数据解析失败
	EDUN_PARAM_ERR    = "6001001" // 业务参数错误
	EDUN_REQ_TIMEOUT  = "6001002" // 请求超时
	EDUN_SYS_BUSY     = "8888888" // 系统繁忙
	EDUN_UNKNOWN_ERR  = "000000"  // 未知错误

	EDUN_MSG = "null"

)

var edunMapCenterCode = map[string]string{
	EDUN_SUCC:         CenterCodeSucc,
	EDUN_PARSE_FAILED: CenterCodeFormat,
	EDUN_PARAM_ERR:    CenterCodeFormat,
	EDUN_REQ_TIMEOUT:  CenterCodeReqData,
	EDUN_SYS_BUSY:     CenterCodeReqData,
	EDUN_UNKNOWN_ERR:  CenterCodeOther,
}

func GetCenterCodeFromEdun(code string) string {
	str, ok := edunMapCenterCode[code]
	if ok {
		return str
	}
	return CenterCodeOther
}
