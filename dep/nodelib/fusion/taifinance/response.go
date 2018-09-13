package taifinance

/**
    Author: luzequan
    Created: 2018-09-03 14:35:00
*/
type ResponseData struct {
	RespCode    int         `json:"resultCode"` // 返回状态码
	RespMessage string      `json:"resultMsg"`  // 返回信息
	RespDetail  *ResultData `json:"resultData"` // 业务详细信息
}

type ResultData struct {
	DefaultProbability []float64 `json:"defaultProbability"` // 风险概率浮点数组
	CreditScore        []float64 `json:"creditScore"`        // 信用评分结果浮点数组
}
