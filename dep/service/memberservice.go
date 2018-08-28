package service

import (
	"sync"
	"encoding/xml"
	"drcs/dep/agollo"
	"drcs/dep/member"
	"fmt"
	"path/filepath"
	"drcs/common/balance"

)

/**
    Author: luzequan
    Created: 2018-05-08 17:08:17
*/

type MemberService struct {
	lock     sync.RWMutex
	memCh    chan bool
}

func NewMemberService() *MemberService {
	return &MemberService{}
}

func (o *MemberService) Init() {
	o.memCh = make(chan bool, 1)

	memberPath := filepath.Join(SettingPath, "member.properties")
	partnersPath := filepath.Join(SettingPath, "partners.properties")
	go o.initMemberConfig(filepath.Clean(memberPath))
	go initPartnersConfig(filepath.Clean(partnersPath))

	select {
	case ret := <-o.memCh:
		fmt.Println("balance init", ret)
		balance.InitBalanceMutex()

		wg.Done()
		break
	}
}

func (o *MemberService) initMemberConfig(configDir string) {
	newAgollo := agollo.NewAgollo(configDir)
	go newAgollo.Start()

	event := newAgollo.ListenChangeEvent()
	for {
		changeEvent := <-event

		fmt.Println("initMemberConfig")

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

		o.memCh <- true
	}
}

func initPartnersConfig(configDir string) {
	newAgollo := agollo.NewAgollo(configDir)
	go newAgollo.Start()

	event := newAgollo.ListenChangeEvent()
	for {
		changeEvent := <-event

		fmt.Println("initPartnersConfig")

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


