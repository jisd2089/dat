package smartsail

/**
    Author: luzequan
    Created: 2018-08-01 11:37:14
*/
import (
	. "drcs/dep/nodelib/crp/common"
)

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
