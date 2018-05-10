package agollo

import (
	"time"
	"encoding/json"
	"sync"
)

const (
	default_notification_id = -1
)

var (
	allNotifications *notificationsMap
	lock sync.RWMutex
)

type NotifyConfigComponent struct {
	appConfig *AppConfig
	notifyChan chan *ChangeEvent
	repository *Repository
}

type apolloNotify struct {
	NotificationId int64  `json:"notificationId"`
	NamespaceName  string `json:"namespaceName"`
}

func (this *NotifyConfigComponent) Start() {
	t2 := time.NewTimer(long_poll_interval)
	//long poll for sync
	for {
		select {
		case <-t2.C:
			this.notifySyncConfigServices()
			t2.Reset(long_poll_interval)
		}
	}
}

func toApolloConfig(resBody []byte) ([]*apolloNotify, error) {
	remoteConfig := make([]*apolloNotify, 0)

	err := json.Unmarshal(resBody, &remoteConfig)

	if err != nil {
		logger.Error("Unmarshal Msg Fail,Error:", err)
		return nil, err
	}
	return remoteConfig, nil
}

func getRemoteConfigSuccessCallBack(notifyChan chan *ChangeEvent, repository *Repository, responseBody []byte) (o interface{}, err error) {
	return toApolloConfig(responseBody)
}

func (this *NotifyConfigComponent) getRemoteConfig() ([]*apolloNotify, error) {
	appConfig := this.appConfig
	if appConfig == nil {
		panic("can not find apollo config!please confirm!")
	}
	notifyChan := this.notifyChan

	repository := this.repository

	urlSuffix := getNotifyUrlSuffix(allNotifications.getNotifies(), appConfig)

	//seelog.Debugf("allNotifications.getNotifies():%s",allNotifications.getNotifies())

	notifies, err := requestRecovery(notifyChan, appConfig, repository, &ConnectConfig{
		Uri:     urlSuffix,
		Timeout: nofity_connect_timeout,
	}, &CallBack{
		SuccessCallBack: getRemoteConfigSuccessCallBack,
	})

	if notifies == nil {
		return nil, err
	}

	return notifies.([]*apolloNotify), err
}

func (this *NotifyConfigComponent) notifySyncConfigServices() error {

	remoteConfigs, err := this.getRemoteConfig()
	if err != nil || len(remoteConfigs) == 0 {
		return err
	}

	updateAllNotifications(remoteConfigs)

	//sync all config
	this.SyncConfig()

	return nil
}

func updateAllNotifications(remoteConfigs []*apolloNotify) {
	for _, remoteConfig := range remoteConfigs {
		if remoteConfig.NamespaceName == "" {
			continue
		}

		allNotifications.setNotify(remoteConfig.NamespaceName, remoteConfig.NotificationId)
	}
}

//func init()  {
//	allNotifications=&notificationsMap{
//		notifications:make(map[string]int64,1),
//	}
//	appConfig:=GetAppConfig()
//
//	allNotifications.setNotify(appConfig.NamespaceName,default_notification_id)
//}

func initNotify() {
	if allNotifications == nil {
		allNotifications = &notificationsMap{
			notifications: make(map[string]int64, 10),
		}
	}

	appConfig := GetAppConfig()

	allNotifications.setNotify(appConfig.NamespaceName, default_notification_id)
}

type notification struct {
	NamespaceName  string `json:"namespaceName"`
	NotificationId int64  `json:"notificationId"`
}

type notificationsMap struct {
	notifications map[string]int64
	sync.RWMutex
}

func (this *notificationsMap) setNotify(namespaceName string, notificationId int64) {
	this.Lock()
	defer this.Unlock()
	this.notifications[namespaceName] = notificationId
}

func (this *notificationsMap) getNotifies() string {
	this.RLock()
	defer this.RUnlock()

	notificationArr := make([]*notification, 0)
	for namespaceName, notificationId := range this.notifications {
		notificationArr = append(notificationArr,
			&notification{
				NamespaceName:  namespaceName,
				NotificationId: notificationId,
			})
	}

	j, err := json.Marshal(notificationArr)

	if err != nil {
		return ""
	}

	return string(j)
}

func (this *NotifyConfigComponent) SyncConfig() error {
	return autoSyncConfigServices(this.notifyChan, this.appConfig, this.repository)
}


func autoSyncConfigServices(notifyChan chan *ChangeEvent, appConfig *AppConfig, repository *Repository) error {
	//appConfig := GetAppConfig()
	//if appConfig == nil {
	//	panic("can not find apollo config!please confirm!")
	//}
	lock.RLock()
	defer lock.RUnlock()

	urlSuffix := getConfigUrlSuffix(appConfig)

	_, err := requestRecovery(notifyChan, appConfig, repository, &ConnectConfig{
		Uri: urlSuffix,
	}, &CallBack{
		SuccessCallBack:   autoSyncConfigServicesSuccessCallBack,
		NotModifyCallBack: repository.touchApolloConfigCache,
	})

	return err
}