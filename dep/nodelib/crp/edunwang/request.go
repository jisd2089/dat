package edunwang

/**
    Author: luzequan
    Created: 2018-08-01 10:49:08
*/
type URIData struct {
	Appid     string `json:"appid"`      // 用户标识
	SecretId  string `json:"secret_id"`  // 云api标识
	ProductId int    `json:"product_id"` //产品ID
	SeqNo     string `json:"seq_no"`     // 调用流水号
	ReqData   string `json:"req_data"`   // 加密信息
}

type RequestData struct {
	TimeStamp       int    `json:"timestamp"`        // 时间戳
	Name            string `json:"name"`             // 姓名
	IdNum           string `json:"id_num"`           // 身份证号
	PhoneNum        string `json:"phone_num"`        // 手机号码
	//Imei            string `json:"imei"`             // 手机IMEI
	//CreditNo        string `json:"credit_no"`        // 统一社会信用代码
	//CompanyName     string `json:"company_name"`     // 公司名称
	//RegisterAddress string `json:"register_address"` // 注册地址
	//ExpiryDate      string `json:"expiry_date"`      // 注册地址
}
