package entity

/**
    Author: luzequan
    Created: 2018-01-10 18:53:59
*/
type BatchReqestVo struct {
	SerialNo  string `json:"serialNo"`
	ReqType   string `json:"reqType"`
	OrderId   string `json:"orderId"`
	IdType    string `json:"idType"`
	TimeStamp string `json:"timeStamp"`
	BatchNo   string `json:"batchNo"`
	FileNo    string `json:"fileNo"`
	UserId    string `json:"userId"`
	Exid      string `json:"exid"`
	TaskId    string `json:"taskId"`
	MaxDelay  string `json:"maxDelay"`
	Header    string `json:"header"`
	DataBoxId int    `json:"dataBoxId"`
}

const (
	ReqType_Start  = "start"
	ReqType_Normal = "normal"
	ReqType_End    = "end"
)
