package common

/**
    Author: luzequan
    Created: 2018-08-03 09:40:43
*/


type CommonRequest struct {
	DemMemId string
	SupMemId string

	PreTimeOut      int
	AuthMode        string
	Reqsign         string
	JobId           string
	SerialNo        string
	IdType          string
	DmpSerialNo     string
	BusiSerialNo    string
	BusiReqInfo     map[string]interface{}
	BusiInfoStr     string
	BusiReqInfoHash string
	TimeStamp       string
	TaskIdStr       string
	Timeout         int
	//TaskIdInfo      []xml.TaskIdInfo
	ConnObjStr   string
	SvcStartTime int
	TranHash     string
	SignBlock1   string
	ServiceId    string
	UnitPrice    float64
}

type CommonResponse struct {
	SuccFlag bool
	ResCode  string
	ResMsg   string
	ErrCode  string
	ErrMsg   string
	//FlowStatus   busilog.FlowStatus
	SupMemId     string
	TaskId       string
	BusiSerialNo string
	DmpSerialNo  string
	TimeStamp    string
	PubResInfo   map[string]interface{}
	BusiResInfo  map[string]interface{}
	BusiInfoStr  string
	SvcStartTime int
	SignBlock1   string
	SignMemId2   string
	SignBlock2   string
	Chargflag    bool
}

type CommonRequestData struct {
	PubReqInfo PubReqInfo             `json:"pubReqInfo"`
	BusiInfo   map[string]interface{} `json:"busiInfo"`
}

type PubReqInfo struct {
	MemId        string `json:"memId"`
	SerialNo     string `json:"serialNo"`
	JobId        string `json:"jobId"`
	AuthMode     string `json:"authMode"`
	TimeStamp    string `json:"timeStamp"`
	ReqSign      string `json:"reqSign"`
}

type PubAnsInfo struct {
	SerialNo     string `json:"serialNo"`
	BusiSerialNo string `json:"busiSerialNo"`
	ResCode      string `json:"resCode"`
	ResMsg       string `json:"resMsg"`
	TimeStamp    string `json:"timeStamp"`
}

type ResponseInfo struct {
	PubResInfo  *PubResInfo            `json:"PubResInfo"`
	BusiResInfo map[string]interface{} `json:"PubResInfo"`
}

type PubResInfo struct {
	ResCode    string `json:"resCode"`
	ResMsg     string `json:"resMsg"`
	Chargeflag string `json:"chargflag"`
}