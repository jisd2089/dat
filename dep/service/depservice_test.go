package service

import (
	"testing"
	"github.com/valyala/fasthttp"
	_ "drcs/dep/nodelib/crp"
	"runtime"
	"drcs/core"
	"time"
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
		"reqSign": "5f4d604a00df289b6b90b66e4d0e1be9d43cd236fc018197dd27e01a0f7e8a3c",
		"serialNo": "2201611161916567677531846",
		"memId": "0000162",
		"authMode": "00"
	},
	"busiInfo": {
		"fullName": "高尚",
		"identityNumber": "330123197507134199",
		"phoneNumber": "13211109876",
		"timestamp": "1531479822"
	}
}`)

	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetBody(requestBody)

	time.Sleep(time.Duration(time.Second * 1))

	for i:=0; i<b.N; i++ {
		depService.ProcessCrpTrans(ctx)
	}
}