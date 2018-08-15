package handler

import (
	"github.com/valyala/fasthttp"
	"fmt"
	"drcs/core"
	"drcs/core/databox"
	"drcs/core/interaction/request"
	"os"
	"path"
	"io"

	"bytes"
	"crypto/md5"
	"bufio"
)

/**
    Author: luzequan
    Created: 2018-01-02 19:58:54
*/
type DemHandler struct {}

func NewDemHandler() *DemHandler {
	return &DemHandler{}
}

/**********************************************************************************************
  * 以一对一批量碰撞为例
  * 1. 需方exid文件，发往供方前置机
  */
func (d *DemHandler) SendDemReqToSup(ctx *fasthttp.RequestCtx) {

	fmt.Println("SendDemReqToSup^^^^^^^^^^^^^^^^^^^^^^")

	filePath := string(ctx.FormValue("filePath"))
	fmt.Println("filePath:" + filePath)
	//filePath := "D:/dds_send/JON20171102000000276_ID010201_20171213175701_0002.TARGET"

	// 1.1 匹配相应的DataBox
	b := assetnode.AssetNodeEntity.GetDataBoxByName("demsend")
	if b == nil {
		fmt.Println("databox is nil!")
	}
	b.SetDataFilePath(filePath)

	addrs := []*request.NodeAddress{}
	//addrs = append(addrs, &request.NodeAddress{MemberId: "000079", IP: "10.101.12.44", Host: "8989", URL: "/api/sup/rec", Priority: 0})
	addrs = append(addrs, &request.NodeAddress{MemberId: "000079", Host: "127.0.0.1", Port: "8989", URL: "/api/sup/rec", Priority: 0})
	//addrs = append(addrs, &request.NodeAddress{MemberId: "000108", IP: "127.0.0.1", Host: "8082", URL: "/api/sup/rec", Priority: 1})
	//addrs = append(addrs, &request.NodeAddress{MemberId: "000109", IP: "127.0.0.1", Host: "8083", URL: "/api/sup/rec", Priority: 2})
	//addrs = append(addrs, &request.NodeAddress{MemberId: "000115", IP: "127.0.0.1", Host: "8084", URL: "/api/sup/rec", Priority: 3})

	//b.SetNodeAddress(addrs)

	// 1.2 setDataBoxQueue
	setDataBoxQueue(b)

	// 1.3 执行，单条http执行碰撞请求
}

func setDataBoxQueue(box *databox.DataBox) {
	dataBoxs := []*databox.DataBox{}
	dataBoxs = append(dataBoxs, box)
	assetnode.AssetNodeEntity.PushDataBox(dataBoxs)
}

/**********************************************************************************************
  * 以一对一批量碰撞为例
  * 4. 需方前置机接收到供方返回文件，推送给需方
  */
func (d *DemHandler) RecSupRespAndPushToDem(ctx *fasthttp.RequestCtx) {

	fmt.Println("seqNo: ", string(ctx.FormValue("seqNo")))
	fmt.Println(string(ctx.Request.Header.Peek("seqNo")))
	fmt.Println(string(ctx.Request.Header.Peek("dataFile")))
	fmt.Println(string(ctx.Request.Header.Peek("taskId")))
	fmt.Println(string(ctx.Request.Header.Peek("md5Str")))

	//dataFile, err := ctx.FormFile("file")
	//if err != nil {
	//	fmt.Println("filePath err:", err)
	//}
	//fmt.Println("filePath***********: ", dataFile.Filename)

	targetFileDir := "D:/dds_receive/tmp"
	//targetFileDir := "/home/ddsdev/data/test/dem/rec"
	targetFilePath := path.Join(targetFileDir, "JON20171102000000276_ID010201_20180627175701_00011.TARGET")

	targetFile, err := os.OpenFile(targetFilePath, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
	}

	//dataFileContent, err := dataFile.Open()
	//defer dataFileContent.Close()
	//if err != nil {
	//}

	io.Copy(targetFile, bytes.NewReader(ctx.Request.Body()))
	targetFile.Close()

	dataFile, err := os.Open(targetFilePath)
	//dataFile, err := os.Open("D:/dds_send/rcv/ODN20161222000000071_ID010201_20171213175701_0003.TARGET")
	defer dataFile.Close()
	if err != nil {
		return
	}

	buf := bufio.NewReader(dataFile)

	md5Hash := md5.New()
	lineCnt := 300
	cntBuf := &bytes.Buffer{}
	c := lineCnt
	for {
		c--
		line, _, err := buf.ReadLine()
		if (err != nil || err == io.EOF) && cntBuf.Len() == 0 {
			//fmt.Println("file end ###############################")
			break
		}

		if (err != nil || err == io.EOF) && cntBuf.Len() > 0 {
			md5Hash.Write(cntBuf.Bytes())
			break
		}

		cntBuf.Write(line)
		cntBuf.WriteByte('\n')

		if c == 0 {

			md5Hash.Write(cntBuf.Bytes())

			c = lineCnt
			cntBuf.Reset()
		}
	}

	md5Str := fmt.Sprintf("%x", md5Hash.Sum(nil))
	fmt.Println("rcv md5: ", md5Str)


	//io.Copy(targetFile, dataFileContent)

	// 1.1 匹配相应的DataBox
	//b := assetnode.AssetNodeEntity.GetDataBoxByName("demrec")
	//b := assetnode.AssetNodeEntity.GetDataBoxByName("demrecbig")
	//
	//if b == nil {
	//	fmt.Println("databox is nil!")
	//}
	//b.DataFilePath = targetFilePath
	//
	//// 1.2 setDataBoxQueue
	//setDataBoxQueue(b)
	// 1.3 执行DataBox，sftp推送文件，核验
}


/**********************************************************************************************
	test
 */

func (d *DemHandler) SplitFile(ctx *fasthttp.RequestCtx) {

	fmt.Println("hello data")
	filePath := string(ctx.FormValue("filePath"))
	fmt.Println("filePath:" + filePath)

	b := assetnode.AssetNodeEntity.GetDataBoxByName("demsplit")
	if b == nil {
		fmt.Println("databox is nil!")
	}
	b.SetDataFilePath(filePath)

	setDataBoxQueue(b)
}

func (d *DemHandler) ReadFile(ctx *fasthttp.RequestCtx) {

	fmt.Println("hello data")
	filePath := string(ctx.FormValue("filePath"))
	fmt.Println("filePath:" + filePath)

	b := assetnode.AssetNodeEntity.GetDataBoxByName("fileread")
	if b == nil {
		fmt.Println("databox is nil!")
	}
	b.SetDataFilePath(filePath)

	setDataBoxQueue(b)
}


func (d *DemHandler) RunParentAndChild(ctx *fasthttp.RequestCtx) {

	filePath := string(ctx.FormValue("filePath"))

	b := assetnode.AssetNodeEntity.GetDataBoxByName("demsendsub")
	if b == nil {
		fmt.Println("databox is nil!")
	}

	b.SetDataFilePath(filePath)

	addrs := []*request.NodeAddress{}
	addrs = append(addrs, &request.NodeAddress{MemberId: "000079", Host: "127.0.0.1", Port: "8989", URL: "/api/sup/rec", Priority: 0})

	//b.SetNodeAddress(addrs)

	setDataBoxQueue(b)
}



func (d *DemHandler) RecSupRespUncompressAndPushToDem(ctx *fasthttp.RequestCtx) {

	//fmt.Println(string(ctx.Request.Body()))

	dataFile, err := ctx.FormFile("file")
	if err != nil {
		fmt.Println("filePath err:", err)
	}
	fmt.Println("filePath***********: ", dataFile.Filename)

	targetFileDir := "D:/dds_receive/tmp"
	//targetFileDir := "/home/ddsdev/data/test/dem/rec"
	targetFilePath := path.Join(targetFileDir, "JON20171102000000276_ID010201_20171213175701_00011.TARGET")
	//
	targetFile, err := os.OpenFile(targetFilePath, os.O_WRONLY|os.O_CREATE, 0644)
	defer targetFile.Close()
	if err != nil {

	}

	dataFileContent, err := dataFile.Open()
	defer dataFileContent.Close()
	if err != nil {

	}

	//io.Copy(targetFile, bytes.NewReader(ctx.Request.Body()))

	io.Copy(targetFile, dataFileContent)

	// 1.1 匹配相应的DataBox
	//b := assetnode.AssetNodeEntity.GetDataBoxByName("demrec")
	b := assetnode.AssetNodeEntity.GetDataBoxByName("demrecbig")

	if b == nil {
		fmt.Println("databox is nil!")
	}
	b.DataFilePath = targetFilePath

	// 1.2 setDataBoxQueue
	setDataBoxQueue(b)
	// 1.3 执行DataBox，sftp推送文件，核验
}


func (d *DemHandler) SendMQ(ctx *fasthttp.RequestCtx) {

	fmt.Println("SendMQ")

	b := assetnode.AssetNodeEntity.GetDataBoxByName("datamq")
	if b == nil {
		fmt.Println("databox is nil!")
		return
	}

	setDataBoxQueue(b)
}