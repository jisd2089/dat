package service

import (
	"sync"
	"encoding/xml"

	"drcs/dep/agollo"
	"drcs/dep/member"
	"fmt"
	"path/filepath"
)

/**
    Author: luzequan
    Created: 2018-05-08 17:08:17
*/

func init() {
	//NewMemberService().Init()
}

type MemberService struct {
	lock       sync.RWMutex
}

func NewMemberService() *MemberService {
	return &MemberService{}
}

func (o *MemberService) Init() {

	memberPath := filepath.Join(SETTING_PATH, "member.properties")
	partnersPath := filepath.Join(SETTING_PATH, "partners.properties")
	go initMemberConfig(filepath.Clean(memberPath))
	go initPartnersConfig(filepath.Clean(partnersPath))

}

func initMemberConfig(configDir string) {
	newAgollo := agollo.NewAgollo(configDir)
	go newAgollo.Start()

	event := newAgollo.ListenChangeEvent()
	for {
		changeEvent := <-event

		changesCnt := changeEvent.Changes["content"]
		value := changesCnt.NewValue

		switch changesCnt.ChangeType {
		case 0:
			memberInfoList := &member.MemberInfoList{}
			err := xml.Unmarshal([]byte(value), memberInfoList)
			if err != nil {
			}

			member.SetMemberInfoList(memberInfoList)
		case 1:
			memberInfoList := &member.MemberInfoList{}
			err := xml.Unmarshal([]byte(value), memberInfoList)
			if err != nil {
			}

			member.SetMemberInfoList(memberInfoList)
		}

		fmt.Println(member.GetMemberInfoList())
		//bytes, _ := json.Marshal(changeEvent)
		//fmt.Println("event:", string(bytes))
	}
}

func initPartnersConfig(configDir string) {
	newAgollo := agollo.NewAgollo(configDir)
	go newAgollo.Start()

	event := newAgollo.ListenChangeEvent()
	for {
		changeEvent := <-event

		changesCnt := changeEvent.Changes["content"]
		value := changesCnt.NewValue

		switch changesCnt.ChangeType {
		case 0:
			partnerInfoList := &member.PartnerInfoList{}
			err := xml.Unmarshal([]byte(value), partnerInfoList)
			if err != nil {
			}

			member.SetPartnersInfo(partnerInfoList)
		case 1:
			partnerInfoList := &member.PartnerInfoList{}
			err := xml.Unmarshal([]byte(value), partnerInfoList)
			if err != nil {
			}

			member.SetPartnersInfo(partnerInfoList)
		}

		fmt.Println(member.GetPartnersInfo())
		//bytes, _ := json.Marshal(changeEvent)
		//fmt.Println("event:", string(bytes))
	}
}