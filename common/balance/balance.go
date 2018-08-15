package balance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	logger "drcs/log"
	"bytes"
	"sync"
	"drcs/dep/member"
	"drcs/dep/errors"
)

var getBalanceURL = "http://10.101.12.6:8899/api/getQuota/"
var qryBalanceURL = "http://domain/api/burget/qryQuota"
var BanlanceMutex = map[string]*int64{} //单位：厘
var Key = "b750e1424a2d39ff60b1f57786584640"
var MemIdList = []string{}
var BalanceLock sync.RWMutex

func InitBalanceMutex() {

	memberDetailInfo := member.GetMemberInfoList().MemberDetailList.MemberDetailInfo[0]

	BalanceLock.Lock()
	defer BalanceLock.Unlock()
	var a int64 = 0
	BanlanceMutex[memberDetailInfo.MemberId] = &a
}

func Hasbalance(memId string, unitPrice float64) bool {
	BalanceLock.RLock()
	defer BalanceLock.RUnlock()
	if p, ok := BanlanceMutex[memId]; ok {
		if int64(unitPrice*1000) <= *p {
			return true
		}
	}
	return false
}

func UpdateBalance(accId string, sum float64) error {
	BalanceLock.RLock()
	defer BalanceLock.RUnlock()
	if _, ok := BanlanceMutex[accId]; !ok {
		return fmt.Errorf("account[%s] balance is nil", accId)
	}
	atomic.AddInt64(BanlanceMutex[accId], int64(sum*1000))
	return nil
}

func genrandomfactor() string {
	b := make([]byte, 32)
	io.ReadFull(rand.Reader, b)
	return hex.EncodeToString(b)
}

type ReqBalancePamat struct {
	AccId     string  `json:"accId"`
	UnitPrice float64 `json:"unitPrice"`
	QuotaNum  int64   `json:"quotaNum"`
	Skey      string  `json:"Skey"`
	Vkey      string  `json:"Vkey"`
}

func ApplyBalance(accId string, unitPrice float64, quotaNum int64, balanceUrl string) (float64, *errors.MeanfulError) {

	BalanceLock.RLock()
	if _, ok := BanlanceMutex[accId]; !ok {
		BalanceLock.RUnlock()
		return 0, errors.RawNew("042000", fmt.Sprintf(" memId[%s] balance empty ", accId))
	}
	BalanceLock.RUnlock()

	tep1 := fmt.Sprintf("accId=%s&unitPrice=%v&quotaNum=%d", accId, unitPrice, quotaNum)
	h1 := md5.New()
	h1.Write([]byte(tep1))                      // 需要加密的字符串为 sharejs.com
	tep1Hash := hex.EncodeToString(h1.Sum(nil)) // 输出加密结果
	Skey := genrandomfactor()

	tep2 := fmt.Sprintf("%sSkey=%s", tep1Hash, Skey)
	h2 := md5.New()
	h2.Write([]byte(tep2))                      // 需要加密的字符串为 sharejs.com
	tep2Hash := hex.EncodeToString(h1.Sum(nil)) // 输出加密结果

	tep3 := fmt.Sprintf("%skey=%s", tep2Hash, Key)
	h3 := md5.New()
	h3.Write([]byte(tep3))
	Vkey := hex.EncodeToString(h1.Sum(nil))

	var ReqBalancePamat struct {
		AccId     string  `json:"accId"`
		UnitPrice float64 `json:"unitPrice"`
		QuotaNum  int64   `json:"quotaNum"`
		Skey      string  `json:"Skey"`
		Vkey      string  `json:"Vkey"`
	}
	ReqBalancePamat.AccId = accId
	ReqBalancePamat.UnitPrice = unitPrice
	ReqBalancePamat.QuotaNum = quotaNum
	ReqBalancePamat.Skey = Skey
	ReqBalancePamat.Vkey = Vkey
	reqBalanceParamByte, err := json.Marshal(ReqBalancePamat)
	if err != nil {
		logger.Error("Marshal reqBalanceParam error!")
		return 0, errors.RawNew("042000", "Marshal reqBalanceParam error!")
	}
	client := http.Client{}

	logger.Info("apply balance start:", string(reqBalanceParamByte))
	resp, err := client.Post(balanceUrl, "application/json", bytes.NewReader(reqBalanceParamByte))
	defer resp.Body.Close()
	if err != nil {
		logger.Error(fmt.Sprintf("%s Apply balance error, err : %s", balanceUrl, err.Error()))
		return 0, errors.RawNew("042000", "Apply balance error")
	}

	respData, err := ioutil.ReadAll(resp.Body)
	logger.Info("apply balance end:", string(respData))
	if err != nil {
		return 0, errors.RawNew("042000", "Read resp.Body error")
	}

	var ResPamat struct {
		SerialNo  string `json:"serialNo"`
		ResCode   string `json:"resCode"`
		ResMsg    string `json:"resMsg"`
		QuotaNum  int    `json:"quotaNum"`
		TimeStamp string `json:"timeStamp"`
	}
	err = json.Unmarshal(respData, &ResPamat)
	if err != nil {
		return 0, errors.RawNew("021001", "Unmarshal balanceHttp respData error")
	}

	resCode := ResPamat.ResCode
	balance := ResPamat.QuotaNum
	msg := ResPamat.ResMsg

	// "000000"：全部成功   "000001"：部分成功
	if resCode != "000000" && resCode != "000001" {
		return 0, errors.RawNew("042000", msg)
	}

	applyAmount := (float64)(balance) * unitPrice

	if err := UpdateBalance(accId, applyAmount); err != nil {
		return 0, errors.RawNew("042000", err.Error())
	}

	return applyAmount, nil
}

