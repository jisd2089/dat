package edunwang

/**
    Author: luzequan
    Created: 2018-08-01 10:49:22
*/
type ResponseData struct {
	StatusCode string `json:"statusCode"` // 返回状态码
	Message    string `json:"message"`    // 返回信息
	SeqNo      string `json:"seq_no"`     // 调用流水号
	RspData    string `json:"result"`     // 业务详细信息
}

/**
	Tag标签字段说明
	疑似恶意欺诈 疑似存在欺诈历史
	疑似仿冒包装 疑似用虚假资料包装
	疑似垃圾账户 疑似使用猫池号等工具账户欺诈
	网络恶意行为 在社交、 o2o、社区等疑似有不良的行为
	疑似盗刷 账户被盗刷
	疑似套现
 */
type RspData struct {
	Tag       string `json:"tag"`        // 黑标签
	EvilScore int    `json:"evil_score"` // 恶意等级
}
