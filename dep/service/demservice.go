package service

import (
	"fmt"
	"drcs/core"
	"drcs/core/databox"
	"drcs/core/interaction/request"
)

/**
    Author: luzequan
    Created: 2018-01-02 19:58:54
*/
type DemService struct {}

func NewDemService() *DemService {
	return &DemService{}
}

/**
 * 批量配送
 * 从需方dmp到供方dmp
 */
func (d *DemService) SendFromDemReqToSup(batchPath string) {
	fmt.Println("SendDemReqToSup^^^^^^^^^^^^^^^^^^^^^^", batchPath)

	// 1.1 匹配相应的DataBox
	b := assetnode.AssetNodeEntity.GetDataBoxByName("demsend")
	if b == nil {
		fmt.Println("databox is nil!")
	}
	b.SetDataFilePath(batchPath)

	addrs := []*request.NodeAddress{}
	//addrs = append(addrs, &request.NodeAddress{MemberId: "000079", IP: "10.101.12.44", Host: "8989", URL: "/api/sup/rec", Priority: 0})
	addrs = append(addrs, &request.NodeAddress{MemberId: "000079", Host: "127.0.0.1", Port: "8989", URL: "/api/sup/rec", Priority: 0})
	//addrs = append(addrs, &request.NodeAddress{MemberId: "000108", IP: "127.0.0.1", Host: "8082", URL: "/api/sup/rec", Priority: 1})
	//addrs = append(addrs, &request.NodeAddress{MemberId: "000109", IP: "127.0.0.1", Host: "8083", URL: "/api/sup/rec", Priority: 2})
	//addrs = append(addrs, &request.NodeAddress{MemberId: "000115", IP: "127.0.0.1", Host: "8084", URL: "/api/sup/rec", Priority: 3})

	//b.SetNodeAddress(addrs)

	// 1.2 setDataBoxQueue
	setDataBoxQueue(b)
}

func setDataBoxQueue(box *databox.DataBox) {
	dataBoxs := []*databox.DataBox{}
	dataBoxs = append(dataBoxs, box)
	assetnode.AssetNodeEntity.PushDataBox(dataBoxs)
}
