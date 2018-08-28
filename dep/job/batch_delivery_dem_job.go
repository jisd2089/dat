package job

import (
)

// 需方前置机扫描文件
type FetchFileTask struct{}

func (t FetchFileTask) Run() {

	// 需方前置机扫描查询请求文件
	if _SFTP_ENABLE == 1 {
		//logger.Info("start fetch sftp file......")
		fetch_sftpFile_task(_FILE_SUFIX_TARGET)
	} else {
		//logger.Info("start fetch local file......")
		//fetch_localFile_task(_FILE_SUFIX_TARGET)
	}
}

// 需方前置机发送文件至供方
type FilePushTask struct {
	path string

	//redisCli      rediscli.RedisClient
	//redisCacheKey string
	//dataFileName  *common.DataFileName
}


// 需方前置机接口供方返回文件
//type CheckBatchUploadTask struct {
//	path string
//
//	redisCli      rediscli.RedisClient
//	redisCacheKey string
//	dataFileName  *common.DataFileName
//}
//
//func (c CheckBatchUploadTask) Run() {
//	bc := &demander.BatchCollision{
//		MaxThreadNum: settings.GetCommomSettings().BatchCollison.MaxThreadNum,
//	}
//	bc.ScanBatchKeys()
//}
//
//func (c CheckBatchUploadTask) process() {
//	logger.Info("CheckBatchUploadTask start process..., filePath: %v", c.path)
//
//	fileName := path.Base(c.path)
//
//	dataFileName := &common.DataFileName{}
//	if err := dataFileName.ParseAndValidFileName(fileName); err != nil {
//		logger.Error("Parse and valid fileName: [%s] error: %s", fileName, err)
//		return
//	}
//
//	c.dataFileName = dataFileName
//	c.redisCacheKey = dataFileName.GetCacheKey()
//
//	logger.Info("CheckBatchUploadTask GetCacheKey: %v", dataFileName.GetCacheKey())
//
//	client, err2 := rediscli.GetBatchRedisClient()
//	if err2 != nil {
//		logger.Error("dem batch job redis get err: %s", err2)
//		return
//	}
//
//	c.redisCli = client
//
//	logger.Info("CheckBatchUploadTask execPullFileTask pending...")
//
//	c.execPullFileTask()
//}
