package agollo

//start apollo
func Start(appConfigPath string) {
	if appConfigPath == "" {
		appConfigPath = appConfigFileName
	}

	StartWithLogger(nil, appConfigPath)
}

func StartWithLogger(loggerInterface LoggerInterface, appConfigFileName string) {
	if loggerInterface != nil {
		initLogger(loggerInterface)
	}

	setConfig(appConfigFileName)

	initNotify()

	//notifyChan := make(chan *ChangeEvent, 1)

	//first sync
	//notifySyncConfigServices(notifyChan, &AppConfig{})

	//start auto refresh config
	go StartRefreshConfig(&AutoRefreshConfigComponent{})

	//start long poll sync config
	go StartRefreshConfig(&NotifyConfigComponent{})
}
