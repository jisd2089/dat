package jobs

import (
	logger "drcs/log"
	//"dds/rediscli"
	"drcs/settings"

	"os"
	"strings"
)

const LOADSUFIX = "BATCH"
const LOADROUTESIZE = 50
const LOADREDISPIPESIZE = 10000

type SupLoadDataTask struct{}

func (sldt SupLoadDataTask) Run() {
	// 扫描 路径下 所有文件，不做判断是否合法
	fileNameList, err := getLocalFileList(settings.GetCommomSettings().SupLoadDir)
	if err != nil {
		logger.Error("sup batch load job scan dir err ", err)
	}
	exidTimeout := settings.GetCommomSettings().ExIdTimeout
	if exidTimeout <= 0 {
		exidTimeout = 30
		logger.Error("sup batch load job exidTimeout(0) is reseted to default value 30 ")
	}
	for _, fileName := range fileNameList {
		// name judge
		if !judgeName(fileName) {
			continue
		}
		newFileName := fileName + ".done"
		err := os.Rename(fileName, newFileName)
		if err != nil {
			logger.Error("sup batch load job rename file error : ", err)
			return
		}
		//for i := 0; i < 3; i++ {
		//	err = loadRedis(newFileName)
		//	if err == nil {
		//		break
		//	}
		//}
	}

}
func judgeName(fileName string) bool {
	fileNameComp := strings.Split(fileName, ".")
	if len(fileNameComp) < 2 {
		return false
	}
	if LOADSUFIX == fileNameComp[len(fileNameComp)-1] {
		idTypeComp := strings.Split(fileNameComp[len(fileNameComp)-2], "_")
		if len(idTypeComp) < 1 || len(idTypeComp[len(idTypeComp)-1]) < 8 || idTypeComp[len(idTypeComp)-1][:2] != "ID" {
			return false
		}
		return true
	}
	return false
}

// 文件名命名 ： *_IDTYPE.BATCH
//func loadRedis(targetFileName string) error {
//	client, err := rediscli.GetRedisClient()
//	if err != nil {
//		logger.Error("sup batch load job redis get err ", err)
//		return err
//	}
//	fs, err := os.Open(targetFileName)
//	defer fs.Close()
//	if err != nil {
//		logger.Error("sup batch load job open file err ", err)
//		return err
//	}
//	fileNameComponent := strings.Split(targetFileName, ".")
//	forIdtypeComp := strings.Split(fileNameComponent[0], "_")
//	idType := forIdtypeComp[len(forIdtypeComp)-1]
//	content, err := ioutil.ReadAll(fs)
//	if err != nil {
//		logger.Error("sup batch load job read file err ", err)
//		return err
//	}
//	// 超时时间 30天
//	var EXPIRED = time.Duration(settings.GetCommomSettings().ExIdTimeout*24) * time.Hour
//	if settings.GetCommomSettings().ExIdTimeout <= 0 {
//		logger.Error("sup batch load job ExIdTimeout is <= 0, reset to default value 30 days ")
//		EXPIRED = time.Duration(30*24) * time.Hour
//	}
//	exidList := strings.Split(string(content), "\n")
//
//	countChan := make(chan int64, 1)
//	count := int64(0)
//
//	realRouteSize := LOADROUTESIZE
//	if len(exidList) < LOADREDISPIPESIZE {
//		realRouteSize = 1
//	}
//	allocToRouteNum := (len(exidList) + realRouteSize) / realRouteSize
//	for i := 0; i < realRouteSize; i++ {
//		start := i * allocToRouteNum
//		end := (i + 1) * allocToRouteNum
//		if end >= len(exidList) {
//			end = len(exidList) - 1
//		}
//		go pipeLineSetExid(client, exidList[start:end], idType, EXPIRED, countChan)
//	}
//
//	for i := 0; i < realRouteSize; i++ {
//		select {
//		case succCount := <-countChan:
//			if succCount == -1 {
//				return errors.New("panic in loading redis using goroute")
//			}
//			count += succCount
//		}
//	}
//
//	logger.Info("%v records have been imported from %v", count, targetFileName)
//	return nil
//}

//func pipeLineSetExid(client rediscli.RedisClient, exids []string, idType string, EXPIRED time.Duration, countChan chan int64) {
	//count := int64(0)
	//defer func() {
	//	if p := recover(); p != nil {
	//		logger.Error("load exid panic occurs from %s to %s", exids[0], exids[len(exids)-1])
	//		// -1 presents panic
	//		countChan <- -1
	//	}
	//}()
	//
	//r, _ := regexp.Compile("^[0-9A-Za-z]{32}$")
	//tmpExidList := []rediscli.PipeKeyValue{}
	//for i, exid := range exids {
	//	// exid has 32 characters
	//	exid = strings.Trim(exid, " ")
	//	if !r.Match([]byte(exid)) {
	//		continue
	//	}
	//	tmpExidList = append(tmpExidList, rediscli.PipeKeyValue{Key: idType + "_" + exid, Value: "1", Expiration: EXPIRED})
	//	count++
	//	// call redis pipeline for LOADROUTESIZE and the last tmpExidList
	//	if count%LOADREDISPIPESIZE == 0 || i == len(exids)-1 {
	//		err := client.PipeLineSetString(tmpExidList)
	//		if err == nil {
	//			tmpExidList = []rediscli.PipeKeyValue{}
	//		} else {
	//			logger.Error("error occurs : %s", err.Error())
	//			if count%LOADREDISPIPESIZE == 0 {
	//				count -= LOADREDISPIPESIZE
	//			} else {
	//				count -= int64((i + 1) % LOADREDISPIPESIZE)
	//			}
	//		}
	//	}
	//
	//}
	//// count of successful loaded exids
	//countChan <- count
//}
