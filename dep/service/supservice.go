package service

import (
	"github.com/valyala/fasthttp"
	"fmt"
	"drcs/core"
	"sync"
	//"drcs/core/databox"
	"drcs/dep/management/entity"
	"drcs/dep/management/constant"
	"strconv"
	"encoding/json"
	"github.com/micro/misc/lib/ctx"
)

/**
    Author: luzequan
    Created: 2018-01-10 13:54:30
*/

type SupService struct {
	lock sync.RWMutex
}

func NewSupService() *SupService {
	return &SupService{}
}

/**
 * 批量配送
 * 从供方dmp到需方dmp
 */
func (d *SupService) SendFromSupRespToDem(batchPath string) {

	// 3.1 匹配相应的DataBox
	b := assetnode.AssetNodeEntity.GetDataBoxByName("supsend")
	if b == nil {
		fmt.Println("databox is nil!")
	}
	b.DataFilePath = batchPath

	// 3.2 setDataBoxQueue
	setDataBoxQueue(b)
	// 3.3 执行DataBox，通过Sftp传输
}