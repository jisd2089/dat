package smartsail

/**
    Author: luzequan
    Created: 2018-08-08 11:18:16
*/

import (
	. "drcs/dep/nodelib/crp/common"
)

const (
	SMARTSAIL_SUCC           = 0   // 系统正常
	SMARTSAIL_RESP_SUCC      = 200 // 查询成功
	SMARTSAIL_DECRYPT_FAILED = 1   // 解密失败
	SMARTSAIL_RESP_TIMEOUT   = 2   // 响应过时
	SMARTSAIL_VALID_ERR      = 3   // 无校验数据
	SMARTSAIL_JSON_ERR       = 5   // json格式错误
	SMARTSAIL_URLTYPE_ERR    = 6   // 非法接口类型
	SMARTSAIL_NO_BALANCE     = 7   // 没有余额
	SMARTSAIL_CLI_NOEXIST    = 8   // 客户不存在或已经失效
	SMARTSAIL_PHONE_ERR      = 9   // 手机号格式错误
	SMARTSAIL_SYS_ERR        = 11  // 系统异常
	SMARTSAIL_NO_PHONE       = -1  // 无手机号信息

)

var SMARTSAILMapCenterCode = map[int]string{
	SMARTSAIL_SUCC:           CenterCodeSucc,
	SMARTSAIL_RESP_SUCC:      CenterCodeSucc,
	SMARTSAIL_DECRYPT_FAILED: CenterCodeFormat,
	SMARTSAIL_VALID_ERR:      CenterCodeFormat,
	SMARTSAIL_JSON_ERR:       CenterCodeFormat,
	SMARTSAIL_RESP_TIMEOUT:   CenterCodeReqData,
	SMARTSAIL_URLTYPE_ERR:    CenterCodeNoAccess,
	SMARTSAIL_NO_BALANCE:     CenterCodeNoMoney,
	SMARTSAIL_CLI_NOEXIST:    CenterCodeNoAccess,
	SMARTSAIL_PHONE_ERR:      CenterCodeFormat,
	SMARTSAIL_SYS_ERR:        CenterCodeOther,
	SMARTSAIL_NO_PHONE:       CenterCodeReqData,
}

func GetCenterCodeFromSMARTSAIL(code int) string {
	str, ok := SMARTSAILMapCenterCode[code]
	if ok {
		return str
	}
	return CenterCodeOther
}
