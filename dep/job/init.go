package job

import (
	"sync"

	"github.com/robfig/cron"
)

var (
	scheduler *cron.Cron
	mutex     sync.Mutex
)

func Init() {
	if scheduler != nil {
		scheduler.Start()
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	if scheduler == nil {
		scheduler = cron.New()
		scheduler.Start()
		scheduler.Run()
	}
}


//var (
//	_FETCH_Interv          int
//	_SFTP_ENABLE           int
//	_Batch_Key_Scan_Interv int
//)
//
//const _WORK_NUM = 1
//
//func Init() {}
//
//func InitSup() {
//	logger.Info("供方定时任务初始化开始。。。。。。")
//	commInit()
//	scheduler.Schedule(time.Minute*time.Duration(_FETCH_Interv), SupFetchFileTask{})
//
//	// TODO 扫描时间间隔
//	//scheduler.Schedule(time.Minute*time.Duration(_FETCH_Interv), SupLoadDataTask{})
//
//}
//
//func InitDem() {
//	logger.Info("需方定时任务初始化开始。。。。。。")
//	commInit()
//	logger.Info("需方定时任务初始化 扫描key时间间隔 : %d 分钟 ; 本地路径: %s \n", _Batch_Key_Scan_Interv, LOCAL_DIR)
//	scheduler.Schedule(time.Minute*time.Duration(_FETCH_Interv), FetchFileTask{})
//	//scheduler.Schedule(time.Minute*time.Duration(_Batch_Key_Scan_Interv), CheckBatchUploadTask{})
//
//}
//
//func commInit() {
//	scheduler.Init()
//
//	settings := settings.GetCommomSettings()
//	SFTP_DIR = settings.Sftp.RemoteDir
//	LOCAL_DIR = settings.Sftp.LocalDir
//	_Batch_Key_Scan_Interv = settings.Sftp.BatchKeyScanInterv
//	if _Batch_Key_Scan_Interv == 0 {
//		_Batch_Key_Scan_Interv = 10
//		logger.Info("min vaule for sftp._Batch_Key_Scan_Interv is 10 minute, ignore..")
//	}
//	_FETCH_Interv = settings.Sftp.FetchInterv
//	if _FETCH_Interv == 0 {
//		_FETCH_Interv = 1
//		logger.Info("min vaule for sftp.fetchInterv is 1, ignore..")
//	}
//
//	_SFTP_ENABLE = settings.Sftp.EnableSftp
//
//	SftpUserName = settings.Sftp.Username
//	SftpUserPWD = settings.Sftp.Password
//	SftpHost = settings.Sftp.Hosts
//	SftpPort = settings.Sftp.Port
//	SftpTimeout = settings.Sftp.DefualtTimeout
//	// fileScan init
//	fileScanInit()
//
//	// 文件清理job per 12H
//	scheduler.Schedule(time.Hour*time.Duration(12), FileScanCleanTask{})
//}
