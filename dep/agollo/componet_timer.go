package agollo

import (
	"time"
)

type AutoRefreshConfigComponent struct {
	appConfig *AppConfig
	notifyChan chan *ChangeEvent
	repository *Repository
}

func (this *AutoRefreshConfigComponent) Start() {
	t2 := time.NewTimer(refresh_interval)
	for {
		select {
		case <-t2.C:
			this.notifySyncConfigServices()
			t2.Reset(refresh_interval)
		}
	}
}

func (this *AutoRefreshConfigComponent) notifySyncConfigServices() error {

	remoteConfigs, err := this.getRemoteConfig()

	if err != nil || len(remoteConfigs) == 0 {
		return err
	}

	updateAllNotifications(remoteConfigs)

	//sync all config
	this.SyncConfig()

	return nil
}

func (this *AutoRefreshConfigComponent) getRemoteConfig() ([]*apolloNotify, error) {
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

func (this *AutoRefreshConfigComponent) SyncConfig() error {
	return autoSyncConfigServices(this.notifyChan, this.appConfig, this.repository)
}

func autoSyncConfigServicesSuccessCallBack(notifyChan chan *ChangeEvent, repository *Repository, responseBody []byte) (o interface{}, err error) {
	apolloConfig, err := createApolloConfigWithJson(responseBody)

	if err != nil {
		logger.Error("Unmarshal Msg Fail,Error:", err)
		return nil, err
	}

	repository.updateApolloConfig(notifyChan, apolloConfig)

	return nil, nil
}






