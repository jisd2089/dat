package agollo

import (
	"github.com/coocood/freecache"
	"sync"
)

/**
    Author: luzequan
    Created: 2018-05-07 15:01:35
*/

type (
	Agollo interface {
		Start()
		ListenChangeEvent() <-chan *ChangeEvent
	}
	agollo struct {
		configFileName string
		appConfig      *AppConfig
		notifyChan     chan *ChangeEvent
		repository     *Repository
		sync.RWMutex
	}
)

func NewAgollo(configFileName string) Agollo {
	return &agollo{
		configFileName: configFileName,
		repository: &Repository{
			currentConnApolloConfig: &ApolloConnConfig{},
			apolloConfigCache:       freecache.NewCache(apolloConfigCacheSize),
		},
	}
}

func (a *agollo) Start() {

	//if loggerInterface != nil {
	//	initLogger(loggerInterface)
	//}

	if a.appConfig != nil {
		go a.updateConfig()
	} else {
		a.setConfig()
	}

	a.initNotify()

	ncc := &NotifyConfigComponent{
		appConfig:  a.appConfig,
		notifyChan: a.notifyChan,
		repository: a.repository,
	}

	//first sync
	ncc.notifySyncConfigServices()

	//start auto refresh config
	go StartRefreshConfig(&AutoRefreshConfigComponent{a.appConfig, a.notifyChan, a.repository})

	//start long poll sync config
	go StartRefreshConfig(ncc)
}

func (a *agollo) ListenChangeEvent() <-chan *ChangeEvent {
	if a.notifyChan == nil {
		a.notifyChan = make(chan *ChangeEvent, 1)
	}
	return a.notifyChan
}

func (a *agollo) setConfig() {

	if a.configFileName == "" {
		a.configFileName = appConfigFileName
	}

	var err error
	//init config file
	a.appConfig, err = loadJsonConfig(a.configFileName)

	if err != nil {
		panic(err)
	}

	go a.updateConfig()
}

func (a *agollo) updateConfig() {
	apolloConfig := &ApolloConfig{}
	apolloConfig.AppId = a.appConfig.AppId
	apolloConfig.Cluster = a.appConfig.Cluster
	apolloConfig.NamespaceName = a.appConfig.NamespaceName

	a.repository.updateApolloConfig(a.notifyChan, apolloConfig)
}

func (a *agollo) initNotify() {
	a.RLock()
	defer a.RUnlock()
	if allNotifications == nil {
		allNotifications = &notificationsMap{
			notifications: make(map[string]int64, 50),
		}
	}

	appConfig := a.appConfig

	allNotifications.setNotify(appConfig.NamespaceName, default_notification_id)
}
