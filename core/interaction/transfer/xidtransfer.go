package transfer

import (
	"drcs/common/xid"
	"drcs/core/interaction/response"
	"reflect"
	"strconv"
)

type XidTransfer struct{}

func NewXidTransfer() Transfer {
	return &XidTransfer{}
}

// 封装xid服务
func (xt *XidTransfer) ExecuteMethod(req Request) Response {

	var (
		returnCode = "000000"
		retMap     map[string]string
	)

	xidParams := ConvertItoStr(req.GetXidParams())

	switch req.GetMethod() {
	case "GEN":
		for k, v := range xidParams {
			code, err := genXidCode(req, v)
			if err != nil {
				returnCode = "000008"
				break
			}
			retMap[k] = code
		}

	case "CVT":
		for k, v := range xidParams {
			code, err := convertXidCode(req, v)
			if err != nil {
				returnCode = "000008"
				break
			}
			retMap[k] = code
		}
	}

	return &response.DataResponse{
		StatusCode: 200,
		ReturnCode: returnCode,
		Bobject:    retMap,
	}
}

func genXidCode(req Request, idNo string) (string, error) {

	xidGenerator := &xid.XidGenerator{
		SrcAppId: req.Param("srcAppId"),
		IdType:   req.Param("idType"),
		IdNo:     idNo,
		XidIp:    req.Param("xidIp"),
		AppKey:   req.Param("appKey"),
	}

	return xidGenerator.GenXID()
}

func convertXidCode(req Request, appXidCode string) (string, error) {

	xidGenerator := &xid.XidGenerator{
		SrcAppId:    req.Param("srcAppId"),
		XidIp:       req.Param("xidIp"),
		AppKey:      req.Param("appKey"),
		SrcRegCode:  req.Param("srcRegCode"),
		DesAppId:    req.Param("desAppId"),
		DesXregCode: req.Param("desXregCode"),
		AppXidCode:  appXidCode,
	}

	return xidGenerator.ConvertXID()
}

func ConvertItoStr(pre map[string]interface{}) map[string]string {

	after := make(map[string]string)

	for key, value := range pre {
		switch value.(type) {
		case string:
			after[key] = value.(string)
		case bool:
			after[key] = reflect.TypeOf(value).String()
		case int:
			after[key] = strconv.Itoa(value.(int))
		}
	}

	return after
}


func (xt *XidTransfer) Close() {}
