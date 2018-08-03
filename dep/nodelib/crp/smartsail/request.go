package smartsail

/**
    Author: luzequan
    Created: 2018-08-03 14:57:26
*/

type RequestMsg struct {
	CliKey      string `json:"cliKey"`
	RequestData string `json:"data"`
}

type RequestData struct {
	Phone     string `json:"phone"`     // 手机号
	Name      string `json:"name"`      // 姓名
	StartTime string `json:"starttime"` // 开始时间
}
