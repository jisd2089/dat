package xid

/*
#include "cmethod.h"
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lGetXIDCode
*/
import "C"
import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os/exec"
	"sync"
	"strconv"
	logger "drcs/log"
)

//monitor usbkey status and open device once globally if usbkey exists
var (
	hasUsb     = false //is usbkey detected on device
	isUsbOpen  = false //is usbkey available
	pin        = C.CString("111111")
	hDev       = C.ulong(0)
	rawxidcode = string(make([]byte, 256)) //placeholder for xid
	mutex      sync.Mutex
)

type XidGenerator struct {
	SrcAppId    string // 源节点appid
	IdType      string // idtype
	IdNo        string // 源id值
	SrcRegCode  string // 源节点regcode
	DesAppId    string // 目的节点appid
	DesXregCode string // 目的节点xregcode
	AppXidCode  string // 源节点xidcode

	XidDealer string // xid生成远程服务标记 “0”：公安三所
	XidIp     string // 公安三所ip地址
	AppKey    string //
}

func init() {
	if hasUsb = detect_usb(); hasUsb {
		isUsbOpen = openUSBKey()
	}
}

func openUSBKey() bool {
	if code := C.OpenDevice(&hDev); code != 0 {
		C.CloseDevice(hDev)
		logger.Error("OpenDevice error, code=%s", strconv.Itoa(int(code)))
		return false
	}
	if code := C.VerifyPIN(hDev, pin); code != 0 {
		C.CloseDevice(hDev)
		logger.Error("VerifyPIN error, code=%s", strconv.Itoa(int(code)))
		return false
	}
	return true
}

func detect_usb() bool {
	cmd := exec.Command("/bin/bash", "-c", `lsusb | grep SCM`)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error("cmd.StdoutPipe error, code=%s", err.Error())
		return false
	}

	if err = cmd.Start(); err != nil {
		logger.Error("cmd.Start error, code=%s", err.Error())
		return false
	}

	mbytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		logger.Error("ioutil.ReadAll error, code=%s", err.Error())
		return false
	}

	if err = cmd.Wait(); err != nil {
		logger.Error("cmd.Wait error, code=%s", err.Error())
		return false
	}
	if len(mbytes) == 0 {
		logger.Error("Info: no device found")
		return false
	}
	return true
}

func (x *XidGenerator) GenXID() (string, error) {

	if hasUsb {
		logger.Info("---will generate xid by usb_key---")
		return genXidByUsb(x.SrcAppId, x.IdType, x.IdNo)
	} else {
		switch x.XidDealer {
		case "0":
			logger.Info("---will generate xid by sansuo---")
			return genXidRemoteSansuo(x.SrcAppId, x.IdType, x.IdNo, x.XidIp, x.AppKey)
		default:
			logger.Info("---will generate xid by chinadep---")
			return genXidRemoteChinadep(x.SrcAppId, x.IdType, x.IdNo, x.XidIp)
		}
	}
}

func (x *XidGenerator) ConvertXID() (string, error) {
	//xid_dealer := settings.GetDemSettings().XidDealer
	if hasUsb {
		logger.Info("---will convert xid by usb_key---")
		return convertXidByUsb(x.SrcAppId, x.SrcRegCode, x.DesAppId, x.DesXregCode, x.AppXidCode)
	} else {
		switch x.XidDealer {
		case "0":
			logger.Info("---will convert xid by sansuo---")
			return convertXidRemoteSansuo(x.SrcAppId, x.SrcRegCode, x.DesAppId, x.DesXregCode, x.AppXidCode, x.XidIp, x.AppKey)
		default:
			logger.Info("---will convert xid by chinadep---")
			return convertXidRemoteChinadep(x.SrcAppId, x.SrcRegCode, x.DesAppId, x.DesXregCode, x.AppXidCode, x.XidIp)
		}
	}
}

func genXidByUsb(appId, idType, idNum string) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()
	app_id := C.CString(appId)
	idtype := C.CString(idType)
	idnum := C.CString(idNum)
	xidcode := C.CString(rawxidcode)
	xidcodeLen := C.int(256)

	if !isUsbOpen {
		if isUsbOpen = openUSBKey(); !isUsbOpen {
			return "", errors.New("Error: failed to open and verify device")
		}
	}
	code := C.GetXIDCode(hDev, app_id, C.int(len(appId)), idtype, C.int(len(idType)), idnum, C.int(len(idNum)),
		xidcode, &xidcodeLen)
	if code != 0 || xidcodeLen == 0 {
		C.CloseDevice(hDev)
		isUsbOpen = false
		logger.Info("Error: failed to getxid from usb key, code=%v", code)
		return "", errors.New("failed to getxid from usb key")
	}

	return C.GoString(xidcode), nil
}

func convertXidByUsb(appIdS, regcodeS, appIdD, regcodeD, appxidcodeS string) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()
	app_id_s := C.CString(appIdS)
	regcode_s := C.CString(regcodeS)
	app_id_d := C.CString(appIdD)
	regcode_d := C.CString(regcodeD)
	appxidcode_s := C.CString(appxidcodeS)
	xidcode := C.CString(rawxidcode)
	xidcodeLen := C.int(256)

	if !isUsbOpen {
		if isUsbOpen = openUSBKey(); !isUsbOpen {
			return "", errors.New("Error: failed to open and verify device")
		}
	}
	code := C.ChangeXIDCode(hDev, app_id_s, C.int(len(appIdS)), regcode_s, C.int(len(regcodeS)), appxidcode_s,
		C.int(len(appxidcodeS)), app_id_d, C.int(len(appIdD)), regcode_d, C.int(len(regcodeD)), xidcode, &xidcodeLen)
	if code != 0 || xidcodeLen == 0 {
		C.CloseDevice(hDev)
		isUsbOpen = false
		return "", errors.New("failed to getxid from usb key")
	}

	return C.GoString(xidcode), nil
}

func genXidRemoteSansuo(app_id, idtype, idnum, xidIp, appKey string) (string, error) {
	url := "http://" + xidIp + "/xidserver/rest/coding/appxidcode/create/sync/" + app_id
	content_type := "application/json;charset=UTF-8;idsp-protocol-version=1.0.0"
	data, err := getGenXidParams(app_id, idtype, idnum, appKey)
	if err != nil {
		return "", err
	}
	res, err := http.Post(url, content_type, bytes.NewBuffer(data))
	defer res.Body.Close()
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", errors.New("post error while visit sansuo to generate xid")
	}

	body, _ := ioutil.ReadAll(res.Body)

	xidinfo := &ResXidSimple{}
	json.Unmarshal(body, xidinfo)
	if !xidinfo.IsValid() {
		return "", errors.New("failed to generate xid, wrong params")
	}
	return getXid(xidinfo, appKey)
}

func genXidRemoteChinadep(app_id, idtype, idnum, xidIp string) (string, error) {
	//ip := settings.GetDemSettings().XidIp
	//if len(ip) == 0 {
	//	logger.Error("no ip found for chinadep")
	//	return "", errors.New("no ip found for chinadep")
	//}
	url := "http://" + xidIp + ":8080/api/p/genxid/"
	content_type := "application/json;charset=UTF-8"
	params := map[string]string{"app_id": app_id, "idtype": idtype, "idnum": idnum}
	data, _ := json.Marshal(params)
	res, err := http.Post(url, content_type, bytes.NewBuffer([]byte(data)))
	defer res.Body.Close()
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", errors.New("post error while visit chinadep to generate xid")
	}

	body, _ := ioutil.ReadAll(res.Body)
	return string(body), nil
}

func convertXidRemoteSansuo(app_id_s, regcode_s, app_id_d, regcode_d, appxidcode_s, xidIp, appKey string) (string, error) {
	//ip := settings.GetDemSettings().XidIp
	//if len(ip) == 0 {
	//	logger.Error("no ip found for sansuo")
	//	return "", errors.New("no ip found for sansuo")
	//}
	url := "http://" + xidIp + "/xidserver/rest/coding/appxidcode/exchange/sync/" + app_id_s
	content_type := "application/json;charset=UTF-8;idsp-protocol-version=1.0.0"
	data := getConvXidParams(app_id_s, regcode_s, app_id_d, regcode_d, appxidcode_s, appKey)
	res, err := http.Post(url, content_type, bytes.NewBuffer(data))
	defer res.Body.Close()
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", errors.New("post error while visit sansuo to convert xid")
	}

	body, _ := ioutil.ReadAll(res.Body)
	xidinfo := &ResXidSimple{}
	json.Unmarshal(body, xidinfo)
	if !xidinfo.IsValid() {
		return "", errors.New("failed to convert xid, wrong params")
	}
	return getXid(xidinfo, appKey)
}

func convertXidRemoteChinadep(app_id_s, regcode_s, app_id_d, regcode_d, appxidcode_s, xidIp string) (string, error) {
	//ip := settings.GetDemSettings().XidIp
	//if len(ip) == 0 {
	//	logger.Error("no ip found for chinadep")
	//	return "", errors.New("no ip found for chinadep")
	//}
	url := "http://" + xidIp + ":8080/api/p/convertxid/"
	content_type := "application/json;charset=UTF-8"
	params := map[string]string{"app_id_s": app_id_s, "regcode_s": regcode_s, "app_id_d": app_id_d, "regcode_d": regcode_d, "appxidcode_s": appxidcode_s}
	data, _ := json.Marshal(params)
	res, err := http.Post(url, content_type, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", errors.New("post error while visit chinadep to convert xid")
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return string(body), nil
}
