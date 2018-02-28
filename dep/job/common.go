package jobs

import (
	logger "drcs/log"
	"drcs/dep/service"
	"drcs/common/sftp"

	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

var (
	LOCAL_DIR    string
	SFTP_DIR     string
	SftpUserName string
	SftpUserPWD  string
	SftpHost     string
	SftpPort     int
	SftpTimeout  int
)

const (
	SourceFolder = "source"
	TargetFolder = "target"

	_SUP_SUFIX = "SOURCE"
	_FILE_SUFIX_TARGET = "TARGET"
)

func getTargetfileList(filelist []string, sufix string) []string {
	sftpClient := service.NewNodeService().SftpClient
	var fileName string
	targetList := make([]string, 0)

	for _, filePath := range filelist {
		fileName = path.Base(filePath)
		fields := strings.Split(fileName, ".")
		dir := path.Dir(filePath)
		if len(fields) < 2 {
			logger.Debug("file name invalid while getting target files  %s", filePath)
			continue
		}
		// fix bug: 添加文件名校验
		if fields[1] == sufix && checkFileName(filePath) {
			flagFile := path.Join(dir, fields[0]+".FLAG")
			if checkFileIsExist(flagFile) {
				logger.Info("find valid file %s ", fileName)
				targetList = append(targetList, filePath)
				if _SFTP_ENABLE == 1 {
					sftpClient.RemoteRM(SFTP_DIR, flagFile)
				} else {

					moveLocalfile(flagFile, path.Join(path.Dir(flagFile), fields[0]+".DONE"))
				}
			}
		}

	}
	return targetList
}

func getResponseFileList(filelist []string, sufix string) []string {
	var fileName string
	targetList := make([]string, 0)

	for _, filePath := range filelist {
		fileName = path.Base(filePath)
		logger.Info("find file %s ", fileName)
		fields := strings.Split(fileName, ".")
		dir := path.Dir(filePath)
		if len(fields) < 2 {
			logger.Debug("file name invalid while getting target files  %s", filePath)
			continue
		}
		if fields[1] == sufix {
			flagFile := path.Join(dir, fields[0]+".FLAG")
			if checkFileIsExist(flagFile) {
				targetList = append(targetList, filePath)
				moveLocalfile(flagFile, path.Join(path.Dir(flagFile), fields[0]+".DONE"))
			}
		}
	}
	logger.Info("getResponseFileList targetLIst: %v ", targetList)
	return targetList
}

func getSftpTargetflieList(filelist []string, sufix string) []string {

	var fileName string
	targetList := make([]string, 0)

	for _, filePath := range filelist {
		fileName = path.Base(filePath)
		fields := strings.Split(fileName, ".")
		dir := path.Dir(filePath)
		if len(fields) < 2 {
			logger.Info("file name invalid wile get sftptarget list %s ", filePath)
			continue
		}
		// fix bug: 添加文件名校验
		if fields[1] == sufix && checkFileName(filePath) {
			flagFile := path.Join(dir, fields[0]+".FLAG")
			for _, flagFile2 := range filelist {
				if flagFile == flagFile2 {
					logger.Info("find valid sftp file %s", fileName)
					targetList = append(targetList, filePath)
					logger.Debug("adding targetfile: %s", filePath)
					rmSftpFile(flagFile)
				}
			}
		}
	}
	logger.Info("targetLIst: %v ", targetList)
	return targetList
}

func getSftpFileList() []string {
	validtList := make([]string, 0)
	sftpClient := service.NewNodeService().SftpClient
	fileInfos, err := sftpClient.RemoteLS(SFTP_DIR)

	if err != nil {
		logger.Error("sftp error while get sftp file list ", err)
	} else {
		for i, fileInfo := range fileInfos {
			logger.Info("No:%d, name: %s, size:%d, dir?:%t \n", i, fileInfo.Name(), fileInfo.Size(), fileInfo.IsDir())
			if fileInfo.IsDir() {
				continue
			}
			ext := path.Ext(fileInfo.Name())
			if ext != ".TARGET" && ext != ".FLAG" && ext != ".SOURCE" {
				logger.Info("file ignored: %s ", fileInfo.Name())
				continue
			}
			validtList = append(validtList, path.Join(SFTP_DIR, fileInfo.Name()))
		}

	}
	logger.Info("valid list****** %v", validtList)
	return validtList
}

func downloadSftpFile(sftpPath, localDir string) string {
	sftpClient := service.NewNodeService().SftpClient
	fcl := &sftp.FileCatalog{
		LocalDir:       localDir,
		LocalFileName:  path.Base(sftpPath),
		RemoteDir:      path.Dir(sftpPath),
		RemoteFileName: path.Base(sftpPath),
	}

	if err := sftpClient.RemoteGet(fcl); err != nil {
		logger.Error("dowloading sftpflie:%s failed", sftpPath, err)
	}
	logger.Info("dowloading sftp file ok: %s", sftpPath)
	return path.Join(localDir, path.Base(sftpPath))
}

func rmSftpFile(sftpPath string) {
	sftpClient := service.NewNodeService().SftpClient
	for i := 0; i < 3; i++ {
		if err := sftpClient.RemoteRM(path.Dir(sftpPath), path.Base(sftpPath)); err != nil {
			logger.Error("rm sftpfile:%s error ", sftpPath, err)
		} else {
			break
		}
	}
}

func fetch_sftpFile_task(fileSufix string) {
	task := FilePushTask{}
	file_list := getSftpFileList()
	targetList := getSftpTargetflieList(file_list, fileSufix)
	for _, targetFile := range targetList {
		if ok := checkFileName(targetFile); ok {
			task.path = downloadSftpFile(targetFile, LOCAL_DIR)
			if task.path == "" {
				continue
			}

			switch fileSufix {
			case "SOURCE":
				service.NewSupService().SendFromSupRespToDem(task.path)
			case "TARGET":
				service.NewDemService().SendFromDemReqToSup(task.path)
			}
		}
	}
}

//func fetch_localFile_task(fileSufix string) {
//	var task FilePushTask
//	file_list, _ := getLocalFileList(LOCAL_DIR)
//	targetList := getTargetfileList(file_list, fileSufix)
//	for _, targetFile := range targetList {
//		if ok := checkFileName(targetFile); ok {
//			task.path = targetFile
//			pushChan <- task
//		}
//	}
//}
//
//func fetchRespFileTask(fileSufix string) {
//	logger.Info("start fetch response file task...")
//	var task CheckBatchUploadTask
//	supRespPath := LOCAL_DIR + string(os.PathSeparator) + "source"
//	file_list, _ := getLocalFileList(supRespPath)
//	targetList := getResponseFileList(file_list, fileSufix)
//	for _, targetFile := range targetList {
//		if ok := checkSupFileName(targetFile); ok {
//			task.path = targetFile
//			pullChan <- task
//		}
//	}
//}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func getLocalFileList(dirPth string) (files []string, err error) {
	files = make([]string, 0, 10)

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)

	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		files = append(files, dirPth+PthSep+fi.Name())
	}

	return files, nil
}

func moveLocalfile(oldPath, newPath string) {

	os.Rename(oldPath, newPath)
}

func checkFileName(filepath string) bool {
	fileName := path.Base(filepath)
	isDem := strings.HasSuffix(fileName, _FILE_SUFIX_TARGET)
	isSup := strings.HasSuffix(fileName, _SUP_SUFIX)
	fields := strings.Split(fileName, "_")
	if isDem && len(fields) != 4 {
		logger.Error("dem batch job invalid file name: ", fileName)
		return false
	}
	if isSup && len(fields) != 4 {
		logger.Error("sup batch job invalid file name: ", fileName)
		return false
	}
	if ok, _ := regexp.MatchString("^(?:ID)[0-9]{6}$", fields[1]); !ok {
		logger.Error("invalid Idtype format ", fields[1])
		return false
	}
	if isDem {
		if ok, _ := regexp.MatchString("[0-9]{14}$", fields[2]); !ok {
			logger.Error("invalid SerNum format ", fields[2])
			return false
		}
		file_num := strings.Split(fields[3], ".")[0]
		if ok, _ := regexp.MatchString("[0-9]{4}$", file_num); !ok {
			logger.Error("invalid file_num format ", file_num)
			return false
		}
	}
	if isSup {
		if ok, _ := regexp.MatchString("[0-9]{14}$", fields[2]); !ok {
			logger.Error("invalid SerNum format ", fields[2])
			return false
		}
		file_num := strings.Split(fields[3], ".")[0]
		if ok, _ := regexp.MatchString("[0-9]{4}$", file_num); !ok {
			logger.Error("invalid file_num format ", file_num)
			return false
		}
	}
	return true
}

func checkSupFileName(filepath string) bool {
	fileName := path.Base(filepath)
	isSup := strings.HasSuffix(fileName, _SUP_SUFIX)
	fields := strings.Split(fileName, "_")
	if isSup && len(fields) != 5 {
		logger.Error("sup batch job invalid file name: ", fileName)
		return false
	}
	if ok, _ := regexp.MatchString("^(?:ID)[0-9]{6}$", fields[2]); !ok {
		logger.Error("invalid Idtype format ", fields[2])
		return false
	}
	if isSup {
		if ok, _ := regexp.MatchString("[0-9]{14}$", fields[3]); !ok {
			logger.Error("invalid SerNum format ", fields[3])
			return false
		}
		file_num := strings.Split(fields[4], ".")[0]
		if ok, _ := regexp.MatchString("[0-9]{4}$", file_num); !ok {
			logger.Error("invalid file_num format ", file_num)
			return false
		}
	}
	return true
}
