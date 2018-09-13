package taifinance

/**
    Author: luzequan
    Created: 2018-08-01 11:37:14
*/
import (
	. "drcs/dep/nodelib/fusion/common"
)

type PubResProductMsg struct {
	PubAnsInfo *PubAnsInfo `json:"pubAnsInfo"`
	HitInfo struct {
		HitResult string `json:"hitResult"`
	} `json:"hitInfo"`
	DetailInfo *ResultData  `json:"respData"`
}
