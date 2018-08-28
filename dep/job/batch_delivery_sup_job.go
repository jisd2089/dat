package job

import (
)

type SupFetchFileTask struct{}

func (t SupFetchFileTask) Run() {
	if _SFTP_ENABLE == 1 {
		//logger.Info("start fetching sftp file")
		fetch_sftpFile_task(_SUP_SUFIX)
	} else {
		//logger.Info("start fetching local file")
		//fetch_localFile_task(_SUP_SUFIX)
	}
}

//func (t *FilePushTask) processSup() {
//	var filesToSend = make([]string, 0)
//	logger.Info("sup start process %s", t.path)
//	fileName := path.Base(t.path)
//	fields := strings.Split(fileName, "_")
//
//	orderID := fields[0]
//
//	client, err2 := rediscli.GetRedisClient()
//	if err2 != nil {
//		logger.Error("sup batch job redis get err ", err2)
//		return
//	}
//
//	key := common.GetKeyStr(fileName)
//	if ok, err := client.HExistString(key); !ok || err != nil {
//		if err != nil {
//			logger.Error("error occured ", err)
//			return
//		}
//		if !ok {
//			//TODO 接力配送临时方案--------------------------------------------
//			//logger.Error("invalid file name %s, key:%s does not exists: ", fileName, key)
//			//return
//			logger.Debug("[接力配送临时方案 start][redis-HSetString][%s:%s:%s]", key, common.Field_status, common.Stat_process)
//			tmpKey := common.GetKeyStrTmp(fileName)
//			memId, err := client.HGetString(tmpKey, common.Field_userID)
//			if err != nil  {
//				logger.Error("[接力配送临时方案][redis-GetString][%s:%s] key not existed err: ", tmpKey, common.Field_userID, err)
//				return
//			}
//			seqNo, err := client.HGetString(tmpKey, common.Field_reqNo)
//			if err != nil {
//				logger.Error("[接力配送临时方案][redis-GetString][%s:%s] key not existed err: ", tmpKey, common.Field_reqNo, err)
//				return
//			}
//			taskId, err := client.HGetString(tmpKey, common.Field_taskId)
//			if err != nil{
//				logger.Error("[接力配送临时方案][redis-GetString][%s:%s] key not existed err: ", tmpKey, common.Field_reqNo, err)
//				return
//			}
//			if len(memId) == 0 || len(seqNo) == 0 || len(taskId) == 0 {
//				logger.Error("[接力配送临时方案]memId, seqNo or taskId not existed")
//				return
//			}
//			err = client.HSetString(key, common.Field_status, common.Stat_process)
//			if err != nil {
//				logger.Error("[接力配送临时方案][redis-HSetString][%s:%s:%s] err:%s ", key, common.Field_status, common.Stat_process, err.Error())
//				return
//			}
//			err = client.HSetString(key, common.Field_taskId, taskId)
//			if err != nil {
//				logger.Error("[接力配送临时方案][redis-HSetString][%s:%s:%s] err:%s ", key, common.Field_taskId, taskId, err.Error())
//				return
//			}
//			err = client.HSetString(key, common.Field_userID, memId)
//			if err != nil {
//				logger.Error("[接力配送临时方案][redis-HSetString][%s:%s:%s] err:%s ", key, common.Field_userID, memId, err.Error())
//				return
//			}
//			err = client.HSetString(key, common.Field_reqNo, seqNo)
//			if err != nil {
//				logger.Error("[接力配送临时方案][redis-HSetString][%s:%s:%s] err:%s ", key, common.Field_reqNo, seqNo, err.Error())
//				return
//			}
//			logger.Debug("[接力配送临时方案 end]")
//			//---------------------------------------------------------
//		}
//	}
//
//	status, _ := client.HGetString(key, common.Field_status)
//	if status == common.Stat_finish {
//		logger.Error("file %s already processed ignore....", t.path)
//		return
//	}
//	v, _ := client.HGetString(key, common.Field_recFile)
//	var recFile string
//
//	if v != "" {
//		if !strings.Contains(v, fileName) {
//			recFile = v + "|@|" + fileName
//		} else {
//			recFile = v
//		}
//
//	} else {
//		recFile = fileName
//	}
//	logger.Info("set recfile key:%s recfile:%s", key, recFile)
//	err := client.HSetString(key, common.Field_recFile, recFile)
//	if err != nil {
//		logger.Error("set cache failed key:%s field:%s v:%s", key, common.Field_recFile,
//			recFile)
//	}
//
//	v, _ = client.HGetString(key, common.Field_taskId)
//	if v == "" {
//		logger.Error("redis taskid not found for %s", key)
//		return
//	} else {
//		//taskID2 := strings.Split(v, "|@|")
//		files := strings.Split(recFile, "|@|")
//		demId, _ := client.HGetString(key, common.Field_userID)
//		if demId == "" {
//			logger.Error("error order info demid missing ", key)
//			return
//		}
//		//if !isAllTaskIDRcved(recFile, taskID2) {
//		//	logger.Error("order not ready , recieved taskID:  %s", taskID)
//		//	return
//		//}
//		deadline, _ := client.HGetString(key, common.Field_deadLine)
//		if deadline > getTimeStr() || deadline == "" {
//
//			filesToSend = append(filesToSend, files[:]...)
//			logger.Info("order ready to push file %v", files)
//		}
//	}
//	if len(filesToSend) > 0 {
//		res := supplier.ExecPolicyBatch(filesToSend)
//		if res.ERRCODE != "0000" {
//			logger.Error("sup push file error %s %s ", orderID, res.ERRMSG)
//		} else {
//			client.HSetString(key, common.Field_status, common.Stat_finish)
//			logger.Info("sup push file ok %s %s ", orderID, res.ERRMSG)
//		}
//	} else {
//		logger.Info("sup batch job no file to send....", orderID)
//	}
//
//}
//func isAllTaskIDRcved(files string, dataRngList []string) bool {
//	for _, v := range dataRngList {
//		if !strings.Contains(files, v) {
//			return false
//		}
//	}
//	return true
//}
//func StringSliceEqual(a, b []string) bool {
//	if len(a) != len(b) {
//		return false
//	}
//
//	if (a == nil) != (b == nil) {
//		return false
//	}
//
//	for i, v := range a {
//		if v != b[i] {
//			return false
//		}
//	}
//
//	return true
//}
//func pathExists(path string) (bool, error) {
//	_, err := os.Stat(path)
//	if err == nil {
//		return true, nil
//	}
//	if os.IsNotExist(err) {
//		return false, nil
//	}
//	return false, err
//}
//
//func getTimeStr() string {
//	currentTime := time.Now().Local()
//
//	//print time
//	//fmt.Println(currentTime)
//
//	//format Time, string type
//	newFormat := currentTime.Format("200601021504")
//
//	return newFormat
//}
