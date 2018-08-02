package crp

/**
    Author: luzequan
    Created: 2018-08-01 17:30:03
*/
import (
	. "drcs/core/databox"
)

func init() {
	SUPRESPONSE.Register()
}

var SUPRESPONSE = &DataBox{
	Name:        "sup_response",
	Description: "sup_response",
	RuleTree: &RuleTree{
		Root: supResponseRootFunc,

		Trunk: map[string]*Rule{
			"parseparam": {
				ParseFunc: parseReqParamFunc,
			},
			"depauth": {
				ParseFunc: depAuthFunc,
			},
			"getorderinfo": {
				ParseFunc: depAuthFunc,
			},
			"aesencrypt": {
				ParseFunc: aesEncryptParamFunc,
			},
			"base64encode": {
				ParseFunc: base64EncodeFunc,
			},
			"urlencode": {
				ParseFunc: urlEncodeFunc,
			},
			"execquery": {
				ParseFunc: execQueryFunc,
			},
			"aesdecrypt": {
				ParseFunc: aesDecryptFunc,
			},
			"buildresp": {
				ParseFunc: buildResponseFunc,
			},
			"end": {
				ParseFunc: procEndFunc,
			},
		},
	},
}

func supResponseRootFunc(ctx *Context) {

}
