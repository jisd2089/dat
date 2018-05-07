package agollo

import "github.com/coocood/freecache"

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
	}
)

func NewAgollo(configFileName string) Agollo {
	return &agollo{
		configFileName: configFileName,
		appConfig:      &AppConfig{},
		repository:     &Repository{currentConnApolloConfig: &ApolloConnConfig{}, apolloConfigCache: freecache.NewCache(apolloConfigCacheSize),},
	}
}

func (a *agollo) Start() {
	//if loggerInterface != nil {
	//	initLogger(loggerInterface)
	//}

	if a.configFileName == "" {
		a.configFileName = appConfigFileName
	}

	a.setConfig(a.configFileName)

	a.initNotify()

	ncc := &NotifyConfigComponent{
		appConfig: a.appConfig,
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

func (a *agollo) setConfig(appConfigFileName string) {

	var err error
	//init config file
	a.appConfig, err = loadJsonConfig(appConfigFileName)

	if err != nil {
		panic(err)
	}

	go func(notifyChan chan *ChangeEvent, appConfig *AppConfig) {
		apolloConfig := &ApolloConfig{}
		apolloConfig.AppId = appConfig.AppId
		apolloConfig.Cluster = appConfig.Cluster
		apolloConfig.NamespaceName = appConfig.NamespaceName

		a.repository.updateApolloConfig(notifyChan, apolloConfig)
	}(a.notifyChan, a.appConfig)
}

func (a *agollo) initNotify() {
	if allNotifications == nil {
		allNotifications = &notificationsMap{
			notifications: make(map[string]int64, 10),
		}
	}

	appConfig := a.appConfig

	allNotifications.setNotify(appConfig.NamespaceName, default_notification_id)
}
