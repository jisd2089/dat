package collector

import (
	"fmt"
	"os"
	"path"

	//"github.com/henrylee2cn/pholcus/logs"
	"dat/runtime/output"
)

/**
    Author: luzequan
    Created: 2018-01-15 10:29:21
*/
/************************ FILE 输出 ***************************/
func init() {
	DataOutput["file"] = func(self *Collector) (err error) {
		defer func() {
			if p := recover(); p != nil {
				err = fmt.Errorf("%v", p)
			}
		}()
		var (
			//namespace = util.FileNameReplace(self.namespace())
			outputfile *os.File
			//err error
		)

		for _, datacell := range self.dataDocker {

			dataMap := datacell["Data"].(map[string]interface{})
			fileName := dataMap["FileName"].(string)
			localDir := dataMap["LocalDir"].(string)
			targetFolder := dataMap["TargetFolder"].(string)
			writeType := dataMap["WriteType"]
			content := dataMap["Content"].(string)

			switch writeType.(int) {
			case output.CTW:
				outputfile, err = os.OpenFile(path.Join(localDir, targetFolder, fileName), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
				defer outputfile.Close()
				if err != nil {

				}
			case output.WA:
				outputfile, err = os.OpenFile(path.Join(localDir, targetFolder, fileName), os.O_APPEND|os.O_WRONLY, 0644)
				defer outputfile.Close()
				if err != nil {

				}
			}
			outputfile.WriteString(content)
		}
		return
	}
}