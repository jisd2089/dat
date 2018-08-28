package job

import (
	logger "drcs/log"
	"drcs/settings"
	"drcs/dep/service"

	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

var (
	ToScanDir     []FileScanInfo
	ToScanSftpDir []FileScanInfo
)

// 文件拓展名
const (
	BatchSufix  = "batch"
	TargetSufix = "target"
	SourceSufix = "source"
	FlagSufix   = "flag"
	DoneSufix   = "done"
	ErrSufix    = "err"
	SuccSufix   = "succ"
	ResultSufix = "result"
	SortSufix   = "sort"

	SupType = "sup"
	DemType = "dem"
	// 临时 需方 目录
	ResultTmpFolder = "/tmp"
)

/*
path 文件路径
sufix 拓展名或者该拓展名的备份文件
*/
type FileScanInfo struct {
	Path  string
	Sufix []string
}
type FileScanCleanTask struct{}

var suffixList = []string{TargetSufix, SourceSufix, FlagSufix, DoneSufix, ErrSufix, SuccSufix, ResultSufix}
var fileCleanInterv int // 全局文件过期时间

func fileScanInit() {
	se := settings.GetCommomSettings()
	// 本地batch 文件清理
	ToScanDir = append(ToScanDir, FileScanInfo{se.Sftp.LocalDir, suffixList})
	// 远程sftp服务器 文件清理
	ToScanSftpDir = append(ToScanSftpDir, FileScanInfo{se.Sftp.RemoteDir, suffixList})

	if se.Type == DemType {
		ToScanDir = append(ToScanDir, FileScanInfo{path.Join(se.Sftp.LocalDir, SourceFolder), suffixList})

		ToScanSftpDir = append(ToScanSftpDir, FileScanInfo{path.Join(se.Sftp.RemoteDir, SourceFolder), suffixList})
		// 需方 删除 临时 文件
		tmpSufixList := append(suffixList, ResultSufix, SortSufix)
		ToScanDir = append(ToScanDir, FileScanInfo{path.Join(se.Sftp.LocalDir, ResultTmpFolder), tmpSufixList})

	} else if se.Type == SupType {
		ToScanDir = append(ToScanDir, FileScanInfo{se.SupLoadDir, []string{BatchSufix}})

		ToScanDir = append(ToScanDir, FileScanInfo{path.Join(se.Sftp.LocalDir, TargetFolder), suffixList})

		ToScanSftpDir = append(ToScanSftpDir, FileScanInfo{path.Join(se.Sftp.RemoteDir, TargetFolder), suffixList})
	}

	fileCleanInterv = settings.GetCommomSettings().FileCleanInterv
	if fileCleanInterv <= 0 {
		fileCleanInterv = 7
		logger.Error("fileScanClean job fileScanInterv(0) is reseted to default value 7 ")
	}
}

func (fsct FileScanCleanTask) Run() {
	for _, fileScanInfo := range ToScanDir {
		scanPathWithSufix(fileScanInfo)
	}
	for _, fileSftpScanInfo := range ToScanSftpDir {
		scanSftpPathWithSufix(fileSftpScanInfo)
	}
}

func scanPathWithSufix(fileScanInfo FileScanInfo) {
	fileInfoList, err := ioutil.ReadDir(fileScanInfo.Path)
	if err != nil {
		logger.Error("file scan and clean error : ", err)
		return
	}
	for _, item := range fileInfoList {
		if item.IsDir() {
			continue
		}
		for _, suffix := range fileScanInfo.Sufix {
			if isEndWithSufix(strings.ToLower(item.Name()), suffix) {
				delFile(fileScanInfo.Path, item)
			} else {
				fields := strings.Split(item.Name(), ".")
				if len(fields) > 1 && isEndWithSufix(suffix, "."+strings.ToLower(fields[len(fields)-2])) {
					delFile(fileScanInfo.Path, item)
				}
			}
		}
	}
}

func scanSftpPathWithSufix(fileScanInfo FileScanInfo) {
	sftpClient := service.NewNodeService().SftpClient

	sftpFile, err := sftpClient.RemoteLS(fileScanInfo.Path)
	if err != nil {
		logger.Error("get sftp file list error : ", err)
		return
	}

	for _, item := range sftpFile {
		if item.IsDir() {
			continue
		}
		for _, suffix := range fileScanInfo.Sufix {
			if isEndWithSufix(strings.ToLower(item.Name()), suffix) {
				delSftpFile(fileScanInfo.Path, item)
			} else {
				fields := strings.Split(item.Name(), ".")
				if len(fields) > 1 && isEndWithSufix(suffix, "."+strings.ToLower(fields[len(fields)-2])) {
					delFile(fileScanInfo.Path, item)
				}
			}
		}
	}
}

func isEndWithSufix(fileName string, sufix string) bool {
	sufix = "." + sufix
	return strings.Contains(fileName, sufix)
}

func delFile(filePath string, file os.FileInfo) {
	curTime := time.Now()
	durationTime := curTime.Sub(file.ModTime())

	timeInterv := time.Duration(fileCleanInterv) * 24
	if durationTime/time.Hour >= timeInterv {
		os.Remove(path.Join(filePath, file.Name()))
		logger.Info("%s is removed %d days ago", file.Name(), durationTime/time.Hour/24)
	}
}

func delSftpFile(filePath string, file os.FileInfo) {
	curTime := time.Now()
	durationTime := curTime.Sub(file.ModTime())

	timeInterv := time.Duration(fileCleanInterv) * 24
	if durationTime/time.Hour >= timeInterv {
		sftpClient := service.NewNodeService().SftpClient

		if err := sftpClient.RemoteRM(filePath, file.Name()); err != nil {
			logger.Error("rm sftp file error : ", err)
			return
		}
		logger.Info("%s is removed %d days ago", file.Name(), durationTime/time.Hour/24)
	}
}