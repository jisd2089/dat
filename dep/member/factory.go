package member

import (
	logger "drcs/log"
	"drcs/settings"
	"fmt"
	"sync"
)

const (
	settings_xpath_memid    = "MemberID"
	settings_xpath_filepath = "MemberFile"
	settings_xpath_expirs   = "MemberFileExpirS"
)

var (
	mutex                sync.Mutex
	defaultMemberManager MemberManager
	localMemberID        string
)

func createMemberManager() (MemberManager, error) {
	settings := settings.GetCommonSettings()
	filePath := settings.ConfigFile.MemberFile
	if filePath == "" {
		logger.Error("%s missing in setting ", settings_xpath_filepath)
		return nil, fmt.Errorf("配置缺失:%s", settings_xpath_filepath)
	}

	expirS := settings.ConfigFile.MemberFileExpireTime

	memberManager, err := NewMemberManagerXMLFile(filePath)
	if err != nil {
		logger.Error("getting new member manager from xml error ", err)
		return nil, err
	}

	expirManager := NewExpireMemberManager(expirS, memberManager)
	return expirManager, nil
}

// GetMemberManager MemberManager的工厂方法
func GetMemberManager() (MemberManager, error) {
	if defaultMemberManager != nil {
		return defaultMemberManager, nil
	}
	mutex.Lock()
	defer mutex.Unlock()
	if defaultMemberManager == nil {
		memberManager, err := createMemberManager()
		if err != nil {
			return nil, err
		}
		defaultMemberManager = memberManager
	}
	return defaultMemberManager, nil
}

func initLocalMemberID() (string, error) {
	settings := settings.GetCommonSettings()
	memID := settings.Node.MemberId
	if memID == "" {
		logger.Error(" %s missing in setting", settings_xpath_memid)
		fmt.Errorf("配置缺失:%s", settings_xpath_memid)
	}
	return memID, nil
}

// GetLocalMemberID 获取本地会员编号
func GetLocalMemberID() (string, error) {
	if localMemberID != "" {
		return localMemberID, nil
	}
	mutex.Lock()
	defer mutex.Unlock()
	if localMemberID == "" {
		memID, err := initLocalMemberID()
		if err != nil {
			return "", err
		}
		localMemberID = memID
	}
	return localMemberID, nil
}
