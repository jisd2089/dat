package appender

import (
	"dds/settings"
	"io"
	"os"
	"strconv"
	"time"
	"path"
	"github.com/pkg/errors"
)

const (
	flag     = os.O_WRONLY | os.O_APPEND | os.O_CREATE
	filemode = 0644
)

func NewNormalFileAppender(se settings.Settings) (io.Writer, error) {
	path, err := se.GetString("Path")
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, flag, filemode)
	if err != nil {
		return nil, err
	}
	return file, nil
}

type RotateLogWriter struct {
	file   *os.File
	name   string
	dir    string
	/*
	* h, per hour
	* d, per day
	* m, per month
	 */
	rotate string
	// current opened file's flag
	curFileTime string
}

type RollingLogWriter struct {
	file *os.File
	name string
	dir  string
	size int64
}

// 按日期记录日志
func NewRotateFileAppender(se settings.Settings) (io.Writer, error) {
	name, err := se.GetString("Name")
	if err != nil {
		return nil, err
	}

	dir, err := se.GetString("Dir")
	if err != nil {
		return nil, err
	}
	/*
	* h, per hour
	* d, per day
	* m, per month
	 */
	rotate, err := se.GetString("Rotate")
	if err != nil {
		return nil, err
	}
	return &RotateLogWriter{nil, name, dir, rotate, ""}, nil
}

// 按大小记录日志
func NewRollingFileAppender(se settings.Settings) (io.Writer, error) {
	name, err := se.GetString("Name")
	if err != nil {
		return nil, err
	}

	dir, err := se.GetString("Dir")
	if err != nil {
		return nil, err
	}

	// size：KB
	size, err := se.GetString("Size")
	if err != nil {
		return nil, err
	}
	fileSize, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		return nil, err
	}
	return &RollingLogWriter{nil, name, dir, fileSize}, nil
}

func (logWriter *RotateLogWriter) Write(p []byte) (int, error) {
	fileName := ""
	name := logWriter.name
	dir := logWriter.dir
	/*
	* h, per hour
	* d, per day
	* m, per month
	 */
	rotate := logWriter.rotate
	curDate := time.Now().Unix()
	var timeAppendName = ""
	if rotate == "h" {
		timeAppendName = time.Unix(curDate-3600,0).Format("2006010215")
		//fileName = path.Join(dir, name + "." + yearMonthDayHour)

	} else if rotate == "d" {
		timeAppendName = time.Unix(curDate-86400,0).Format("20060102")
		//fileName = path.Join(dir, name + "." + yearMonthDay)

	} else if rotate == "m" {
		timeAppendName = time.Unix(curDate-2592000,0).Format("200601")
		//fileName = path.Join(dir, name + "." + yearMonth)
	} else {
		return -1, errors.New("rotate type is no included")
	}
	oldFileName := path.Join(dir, name)
	//
	if logWriter.file == nil {
		file, err := os.OpenFile(oldFileName, flag, filemode)
		if err != nil {
			return -1, err
		}

		logWriter.file = file
		logWriter.curFileTime = timeAppendName
	} else if timeAppendName != logWriter.curFileTime {
		logWriter.file.Close()
		fileName = path.Join(dir, name+"."+timeAppendName)
		file, err := openRotateFileWithRename(oldFileName, fileName)
		if err != nil {
			return -1, err
		}
		logWriter.file = file
		logWriter.curFileTime = timeAppendName
	}
	return logWriter.file.Write(p)
}

func openRotateFileWithRename(oldFileName string, newFileName string) (*os.File, error) {

	_, err := os.Stat(newFileName)
	if err != nil {
		os.Rename(oldFileName, newFileName)
	}

	file, err := os.OpenFile(oldFileName, flag, filemode)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (logWriter *RollingLogWriter) Write(p []byte) (int, error) {

	name := logWriter.name

	dir := logWriter.dir

	// size：KB
	fileSize := logWriter.size
	oldFileName := path.Join(dir, name)
	if logWriter.file == nil {
		file, err := os.OpenFile(oldFileName, flag, filemode)
		if err != nil {
			return -1, err
		}

		logWriter.file = file
	} else {
		fileStat, err := logWriter.file.Stat()
		if err != nil {
			return -1, err
		}

		if (fileStat.Size() + 512)/1024 >= fileSize {
			logWriter.file.Close()
			file, err := openFileWithRename(dir, name)
			if err != nil {
				return -1, err
			}
			logWriter.file = file
		}
	}
	return logWriter.file.Write(p)
}

func openFileWithRename(dir string, name string) (*os.File, error) {

	for i := 1; ; i++ {
		_, err := os.Stat(path.Join(dir, name+"."+strconv.Itoa(i)))
		if err != nil {
			os.Rename(path.Join(dir, name), path.Join(dir, name+"."+strconv.Itoa(i)))
			break
		}
	}

	file, err := os.OpenFile(path.Join(dir, name), flag, filemode)
	if err != nil {
		return nil, err
	}
	return file, nil
}
