package web

/**
    Author: luzequan
    Created: 2018-01-04 13:42:46
*/
import (
	"testing"
	"encoding/json"

	//_ "drcs/dep/nodelib/demanderonly"
	//_ "drcs/dep/nodelib/dep"
	//_ "drcs/dep/nodelib/batchdistribution"
	//_ "drcs/dep/nodelib/crp"
	_ "drcs/dep/nodelib/fusion"

	"fmt"
	"strconv"
	"reflect"
	"time"
	"drcs/dep/service"
)

//func init() {
//	assetnode.AssetNodeEntity.Init()
//}



func TestDemRun(t *testing.T) {

	RunTest(8095)
}

func TestSupRun(t *testing.T) {

	service.SettingPath = "D:/GoglandProjects/src/drcs/exec/web/properties"

	RunTest(8096)
}


func TestConvert(t *testing.T) {

	/**
{
	"pubReqInfo": {
		"timeStamp": "1469613279966",
		"jobId": "JON20180117000000000",
		"reqSign": "58fcbe63dd9325f93391cd006f0aff272b54ef9b3197941fd31656bda6cdcb8c",
		"serialNo": "2201611161916567677531846",
		"memId": "0000162",
		"authMode": "00"
	},
	"busiInfo": {
		"fullName": "张天",
		"identityNumber": "330123197507134111",
		"cardNumber": "6225768780423333"
	}
}
 */

	var crpReqParams map[string]map[string]interface{}
	crpReqParams1 := make(map[string]map[string]interface{})
	pubReqInfoMap := make(map[string]interface{})
	pubReqInfoMap["timeStamp"] = "1469613279966"
	pubReqInfoMap["jobId"] = "JON20180117000000000"
	pubReqInfoMap["reqSign"] = "58fcbe63dd9325f93391cd006f0aff272b54ef9b3197941fd31656bda6cdcb8c"
	pubReqInfoMap["serialNo"] = "2201611161916567677531846"
	pubReqInfoMap["memId"] = "0000162"
	pubReqInfoMap["authMode"] = "00"

	crpReqParams1["pubReqInfo"] = pubReqInfoMap

	busiInfoMap := make(map[string]interface{})
	busiInfoMap["fullName"] = "张天"
	busiInfoMap["identityNumber"] = "330123197507134111"
	busiInfoMap["cardNumber"] = "6225768780423333"

	crpReqParams1["busiInfo"] = busiInfoMap

	jsonByte, err := json.Marshal(crpReqParams1)
	if err != nil {
		fmt.Println("marshal err")
	}

	err = json.Unmarshal(jsonByte, &crpReqParams)
	if err != nil {
		fmt.Println("Unmarshal err")
	}


	crpReqParamStr := crpReqParamsConvertStr(crpReqParams)
	fmt.Println(crpReqParamStr)

	pubReqInfo := crpReqParamStr["pubReqInfo"]


	fmt.Println(pubReqInfo["timeStamp"])
}

func crpReqParamsConvertStr(crpReqParams map[string]map[string]interface{}) map[string]map[string]string {

	crp := make(map[string]map[string]string)

	for rootKey, rootVal := range crpReqParams {
		pri := make(map[string]string)

		for key, value := range rootVal {
			switch value.(type) {
			case string:
				pri[key] = value.(string)
			case bool:
				pri[key] = reflect.TypeOf(value).String()
			case int:
				pri[key] = strconv.Itoa(value.(int))
			}
		}
		crp[rootKey] = pri
	}

	return crp
}

func TestDate(t *testing.T) {
	tm := time.Now()

	s := tm.Format("20060102")

	loc, _ := time.LoadLocation("Asia/Chongqing")
	locTime, _ := time.ParseInLocation("20060102", s, loc)
	locTimeStr := locTime.Format("20060102")
	fmt.Println(locTimeStr)
}