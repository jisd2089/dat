package collector

import (
	"fmt"
)
/**
    Author: luzequan
    Created: 2017-12-28 17:07:15
*/

/************************ posixmq 输出 ***************************/
func init() {
	DataOutput["posixmq"] = func(self *Collector) (err error) {
		defer func() {
			if p := recover(); p != nil {
				err = fmt.Errorf("%v", p)
			}
		}()
		return
	}
}