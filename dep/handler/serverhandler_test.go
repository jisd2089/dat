package handler

import (
	"testing"
	"runtime"
	"drcs/core"
	"drcs/dep/service"
	"github.com/valyala/fasthttp"
)

func init() {
	runtime.GOMAXPROCS(8)

	service.SettingPath = "D:/gopath/src/drcs/settings/properties"

	service.NewNodeService().Init()

	assetnode.AssetNodeEntity.Init().Run()
}

func TestAcceptCSVfile(t *testing.T) {

	requestBody := []byte(`{
	"pubReqInfo": {
		"timeStamp": "1469613279966",
		"jobId": "JON20180816000000631",
		"reqSign": "dd4239bbbaca226924a4cf6babd002b9d5f02d33d03025589e937b4ce1b3b3dc",
		"serialNo": "2201611161916567677531846",
		"memId": "0000162",
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

	handler :=NewCustomerHandler()

	handler.UploadCSVfile(ctx)
}

func TestCustomerPredictCreditScore(t *testing.T) {

}

func BenchmarkCustomerPredictCreditScore(b *testing.B) {

}