package xid

import (
	"encoding/json"
	"time"
	logger "drcs/log"
)

var (
	encrypt_factor = genrandomfactor()
	version        = "1.0.0"
	encrypt_type   = "2"
	sign_type      = "2"
)

type ReqGenXid struct {
	ReqCommon
	App_Id       string `json:"app_id"`
	User_Id_Info string `json:"user_id_info"` //加密信息
}

type ReqConvXid struct {
	ReqCommon
	Id           string `json:"id"`
	App_Id_S     string `json:"app_id_s"`
	Regcode_S    string `json:"xregcode_s"`
	Appxidcode_S string `json:"appxidcode_s"`
	App_Id_D     string `json:"app_id_d"`
	Regcode_D    string `json:"xregcode_d"`
}

//xid返回简明信息
type ResXidSimple struct {
	SecurityFactor struct {
		EncryptFactor string `json:"encrypt_factor"`
	} `json:"security_factor"`
	UserInfo string `json:"user_info"`
}

func (this *ResXidSimple) IsValid() bool {
	if len(this.SecurityFactor.EncryptFactor) == 0 || len(this.UserInfo) == 0 {
		return false
	}
	return true
}

//belows are common info for both genxid and convxid
type ReqCommon struct {
	Version        string    `json:"version"`
	Return_url     string    `json:"return_url,omitempty"`
	BizType        string    `json:"biz_type"`
	BizTime        string    `json:"biz_time"`
	BizSequenceId  string    `json:"biz_sequence_id"`
	SecurityFactor SecFactor `json:"security_factor"`
	EncryptType    string    `json:"encrypt_type,omitempty"`
	SignType       string    `json:"sign_type"`
	Sign           string    `json:"sign,omitempty"`
}

type ResCommon struct {
	Version        string    `json:"version"`
	Result         string    `json:"result"`
	ResultDetail   string    `json:"result_detail"`
	ResultTime     string    `json:"result_time"`
	BizSequenceId  string    `json:"biz_sequence_id"`
	SecurityFactor SecFactor `json:"security_factor"`
	EncryptType    string    `json:"encrypt_type,omitempty"`
	SignType       string    `json:"sign_type"`
	Sign           string    `json:"sign,omitempty"`
}

type SecFactor struct {
	EncryptFactor string `json:"encrypt_factor"`
	SignFactor    string `json:"sign_factor"`
}

func getGenXidParams(app_id, idtype, idnum, appKey string) ([]byte, error) {
	var userid struct {
		Idnum  string `json:"idnum"`
		Idtype string `json:"idtype"`
	}
	userid.Idnum = idnum
	userid.Idtype = idtype
	data, _ := json.Marshal(userid)
	user_id_info := ciphermsg(string(data), appKey)

	reqgenxid := ReqGenXid{
		ReqCommon{
			Version:       version,
			BizType:       "0501001",
			BizTime:       time.Now().Format("20060102150405"),
			BizSequenceId: genuuid(),
			SecurityFactor: SecFactor{
				EncryptFactor: encrypt_factor,
				SignFactor:    genrandomfactor(),
			},
			EncryptType: encrypt_type,
			SignType:    sign_type,
		},
		app_id,
		user_id_info,
	}
	sorted_data := sortReqGenXid(reqgenxid)
	//appkey := settings.GetDemSettings().Appkey
	ks, err := genSignKey(reqgenxid.SecurityFactor.SignFactor, appKey)
	if err != nil {
		return []byte(""), err
	}
	reqgenxid.Sign = sign_SHA256(ks, sorted_data)
	reqJson, _ := json.Marshal(reqgenxid)
	return reqJson, nil
}

func getConvXidParams(app_id_s, regcode_s, app_id_d, regcode_d, appxidcode_s, appKey string) []byte {

	reqCommon := ReqCommon{
		Version:       version,
		BizType:       "0501002",
		BizTime:       time.Now().Format("20060102150405"),
		BizSequenceId: genuuid(),
		SecurityFactor: SecFactor{
			EncryptFactor: encrypt_factor,
			SignFactor:    genrandomfactor(),
		},
		EncryptType: encrypt_type,
		SignType:    sign_type,
	}
	reqconvxid := ReqConvXid{
		ReqCommon:    reqCommon,
		Id:           app_id_s,
		App_Id_S:     app_id_s,
		Regcode_S:    ciphermsg(regcode_s, appKey),
		Appxidcode_S: appxidcode_s,
		App_Id_D:     app_id_d,
		Regcode_D:    ciphermsg(regcode_d, appKey),
	}
	sorted_data := sortReqConvXid(reqconvxid)
	//appkey := settings.GetDemSettings().Appkey
	ks, _ := genSignKey(reqconvxid.SecurityFactor.SignFactor, appKey)
	reqconvxid.Sign = sign_SHA256(ks, sorted_data)
	reqJson, _ := json.Marshal(reqconvxid)
	return reqJson
}

func getXid(xidinfo *ResXidSimple, appKey string) (string, error) {
	//appkey := settings.GetDemSettings().Appkey
	ke, _ := genEncryptKey(xidinfo.SecurityFactor.EncryptFactor, appKey)
	res, err := decryptSM4(ke, xidinfo.UserInfo)
	if err != nil {
		return "", err
	}
	var response struct {
		AppXIDcode string `json:"appxidcode"`
	}
	json.Unmarshal([]byte(res), &response)
	return response.AppXIDcode, nil
}

func ciphermsg(msg, appKey string) string {
	//appkey := settings.GetDemSettings().Appkey
	ke, err := genEncryptKey(encrypt_factor, appKey)
	if err != nil {
		logger.Error("fail to generate encrypt key")
	}
	res, err1 := encryptSM4(ke, msg)
	if err1 != nil {
		logger.Error("fail to generate encrypt key")
	}
	return res
}
