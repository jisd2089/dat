package service

import (
	"testing"
	"github.com/valyala/fasthttp"
	_ "drcs/dep/nodelib/crp"
	"runtime"
	"drcs/core"
)

/**
    Author: luzequan
    Created: 2018-08-14 18:51:55
*/
func init() {
	runtime.GOMAXPROCS(8)

	SettingPath = "/home/deplab/project/drcs/config"

	NewNodeService().Init()

	assetnode.AssetNodeEntity.Init().Run()
}

func BenchmarkProcessCrpTrans(b *testing.B) {

	b.ReportAllocs()

	depService := NewDepService()

	requestBody := []byte(`{
	"pubReqInfo": {
		"timeStamp": "1469613279966",
		"jobId": "JON20180516000000431",
		"reqSign": "58ed911a7b6181e3220b077add2417237a3ce55ba91d3c88bc6d960e43823857",
		"serialNo": "2201611161916567677531846",
		"memId": "0000166",
		"authMode": "00"
	},
	"busiInfo": {
		"fullName": "尚书",
		"phoneNumber": "17316332755",
		"starttime": "1531479822"
	}
}`)

	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetBody(requestBody)

	for i:=0; i<b.N; i++ {
		depService.ProcessCrpTrans(ctx)
	}
}