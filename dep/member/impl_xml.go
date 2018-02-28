package member

import (
	logger "drcs/log"
	"encoding/xml"
	"io/ioutil"
	"sync/atomic"
)

// MemberManagerXMLFile 基于XML文件的 MemberManager实现
type MemberManagerXMLFile struct {
	filePath string
	value    *atomic.Value
}

func (manager *MemberManagerXMLFile) GetMemberInfo(memID string) *MemberInfo {
	memberInfoMap := manager.value.Load().(map[string]*MemberInfo)
	return memberInfoMap[memID]
}

// Update 重新读取XML文件，更新成员信息
func (manager *MemberManagerXMLFile) Update() error {
	memberInfoMap, err := parseXMLFile(manager.filePath)
	if err != nil {
		return err
	}
	manager.value.Store(memberInfoMap)
	return nil
}

// NewMemberManagerXMLFile 创建MemberManagerXMLFile实例
func NewMemberManagerXMLFile(filePath string) (*MemberManagerXMLFile, error) {
	memberInfoMap, err := parseXMLFile(filePath)
	if err != nil {
		return nil, err
	}

	manager := &MemberManagerXMLFile{}
	manager.filePath = filePath
	manager.value = &atomic.Value{}
	manager.value.Store(memberInfoMap)
	return manager, nil
}

func parseXMLFile(filePath string) (map[string]*MemberInfo, error) {
	text, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Error("Reading member info %s err", filePath, err)
		return nil, err
	}

	var conf memberInfoConf
	err = xml.Unmarshal(text, &conf)
	if err != nil {
		logger.Error("unmarshal member info %s err", filePath, err)
		return nil, err
	}

	return memberInfoConf2MemberInfoMap(&conf), nil
}

func memberInfoConf2MemberInfoMap(conf *memberInfoConf) map[string]*MemberInfo {
	memberInfoMap := make(map[string]*MemberInfo, len(conf.MemberInfoList.MemberInfoXMLs))
	for _, memberInfoXML := range conf.MemberInfoList.MemberInfoXMLs {
		memberInfo := &MemberInfo{
			MemID:      memberInfoXML.MemId,
			PubKey:     memberInfoXML.PubKey,
			SvrURL:     memberInfoXML.SvrURL,
			Status:     memberInfoXML.Status,
			TotLmt:     memberInfoXML.TotLmt,
			SettFlag:   memberInfoXML.SettFlag,
			SettPoint:  memberInfoXML.SettPoint,
			Threashold: memberInfoXML.Threshold,
		}
		memberInfoMap[memberInfo.MemID] = memberInfo
		logger.Info("init mem info: %+v", memberInfo)
	}
	return memberInfoMap
}

type memberInfoXML struct {
	MemId     string  `xml:"memId"`
	PubKey    string  `xml:"pubKey"`
	SvrURL    string  `xml:"svrURL"`
	Status    string  `xml:"status"`
	TotLmt    float64 `xml:"totLmt"`
	SettFlag  string  `xml:"settFlag"`
	SettPoint string  `xml:"settPoint"`
	Threshold string  `xml:"threshold"`
}

type memberInfoList struct {
	MemberInfoXMLs []memberInfoXML `xml:"mem_dtl_info"`
}

type memberInfoConf struct {
	MemberInfoList memberInfoList `xml:"member_dtl_list"`
}
