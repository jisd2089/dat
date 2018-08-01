package edunwang

/**
    Author: luzequan
    Created: 2018-08-01 11:37:14
*/
import (
	. "drcs/dep/nodelib/crp"
)

//报文1-003-006：手机实名认证（详版）
type PubReqProductMsg_1_003_006 struct {
	OrderId        string `json:"orderId"`
	FullName       string `json:"fullName"`
	IdentityNumber string `json:"identityNumber"`
	Mobile         string `json:"mobile"`
}

//报文1-003-006：手机实名认证（详版）
type PubResProductMsg_1_003_006 struct {
	PubAnsInfo PubAnsInfo `json:"pubAnsInfo"`
	HitInfo struct {
		HitResult string `json:"hitResult"`
	} `json:"hitInfo"`
	DetailInfo struct {
		CheckResult string `json:"checkResult"`
		Message     string `json:"message"`
		Province    string `json:"province"`
		City        string `json:"city"`
		Isp         string `json:"isp"`
	} `json:"detailInfo"`
}

type PubResProductMsg_0_000_000 struct {
	PubAnsInfo PubAnsInfo `json:"pubAnsInfo"`
	HitInfo struct {
		HitResult string `json:"hitResult"`
	} `json:"hitInfo"`
	DetailInfo struct {
		Tag       string `json:"tag"`
		EvilScore int `json:"evilScore"`
	} `json:"respData"`
}
