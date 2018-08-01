package interaction

/**
    Author: luzequan
    Created: 2017-12-28 14:20:45
*/
import (
	"drcs/core/databox"
	"drcs/core/interaction/request"
	"drcs/core/interaction/transfer"
	"drcs/core/interaction/response"
	"sync"
)

type Cross struct {
	httpTsf     transfer.Transfer
	fastHttpTsf transfer.Transfer
	sftpTsf     transfer.Transfer
	noneTsf     transfer.Transfer
	fileTsf     transfer.Transfer
	shellTsf    transfer.Transfer
	sshTsf      transfer.Transfer
	redisTsf    transfer.Transfer
	encryptTsf  transfer.Transfer
	encodeTsf   transfer.Transfer
	depauthTsf  transfer.Transfer
	sync.RWMutex
}

func NewCross() Carrier {
	return &Cross{
		httpTsf:     transfer.NewHttpTransfer(),
		fastHttpTsf: transfer.NewFastTransfer(),
		sftpTsf:     transfer.NewSftpTransfer(),
		noneTsf:     transfer.NewNoneTransfer(),
		fileTsf:     transfer.NewFileTransfer(),
		shellTsf:    transfer.NewShellTransfer(),
		sshTsf:      transfer.NewSshTransfer(),
		redisTsf:    transfer.NewRedisTransfer(),
		encryptTsf:  transfer.NewEncryptTransfer(),
		encodeTsf:   transfer.NewEncodeTransfer(),
		depauthTsf:  transfer.NewDepAuthTransfer(),
	}
}

var CrossHandler = &Cross{
	httpTsf:     transfer.NewHttpTransfer(),
	fastHttpTsf: transfer.NewFastTransfer(),
	sftpTsf:     transfer.NewSftpTransfer(),
	noneTsf:     transfer.NewNoneTransfer(),
	fileTsf:     transfer.NewFileTransfer(),
	shellTsf:    transfer.NewShellTransfer(),
	sshTsf:      transfer.NewSshTransfer(),
	redisTsf:    transfer.NewRedisTransfer(),
}

func (c *Cross) Handle(b *databox.DataBox, cReq *request.DataRequest) *databox.Context {
	c.RLock()
	defer c.RUnlock()

	ctx := databox.GetContext(b, cReq)

	var resp *response.DataResponse
	var err error

	switch cReq.GetTransferType() {
	case request.HTTP:
		resp = c.httpTsf.ExecuteMethod(cReq).(*response.DataResponse)
	case request.FASTHTTP:
		resp = c.fastHttpTsf.ExecuteMethod(cReq).(*response.DataResponse)
	case request.SFTP:
		resp = c.sftpTsf.ExecuteMethod(cReq).(*response.DataResponse)
	case request.FILETYPE:
		resp = c.fileTsf.ExecuteMethod(cReq).(*response.DataResponse)
	case request.NONETYPE:
		resp = c.noneTsf.ExecuteMethod(cReq).(*response.DataResponse)
	case request.SHELLTYPE:
		resp = c.shellTsf.ExecuteMethod(cReq).(*response.DataResponse)
	case request.SSH:
		resp = c.sshTsf.ExecuteMethod(cReq).(*response.DataResponse)
	case request.REDIS:
		resp = c.redisTsf.ExecuteMethod(cReq).(*response.DataResponse)
	case request.ENCRYPT:
		resp = c.encryptTsf.ExecuteMethod(cReq).(*response.DataResponse)
	case request.ENCODE:
		resp = c.encodeTsf.ExecuteMethod(cReq).(*response.DataResponse)
	case request.DEPAUTH:
		resp = c.depauthTsf.ExecuteMethod(cReq).(*response.DataResponse)
	default:
		resp = c.noneTsf.ExecuteMethod(cReq).(*response.DataResponse)
	}

	if resp.GetStatusCode() >= 400 {
		//err = errors.New("响应状态 " + resp.Status)
	}

	ctx.SetResponse(resp).SetError(err)

	return ctx
}

func (c *Cross) Close() {
	c.sftpTsf.Close()
}

func (c *Cross) Process(b *databox.DataBox, req *request.DataRequest) {
	var ctx = c.Handle(b, req)
	//var ctx = self.Downloader.Download(sp, req) // download page

	if err := ctx.GetError(); err != nil {
		// 返回是否作为新的失败请求被添加至队列尾部
		//if b.DoHistory(req, false) {
		//	// 统计失败数
		//	cache.PageFailCount()
		//}
		// 提示错误
		//logs.Log.Error(" *     Fail  [download][%v]: %v\n", downUrl, err)
		return
	}

	// 过程处理，提炼数据
	ctx.Parse(req.GetRuleName())

	// 该条请求文件结果存入pipeline
	//for _, f := range ctx.PullFiles() {
	//	if m.Pipeline.CollectFile(f) != nil {
	//		break
	//	}
	//}
	//// 该条请求文本结果存入pipeline
	//for _, item := range ctx.PullItems() {
	//	if m.Pipeline.CollectData(item) != nil {
	//		break
	//	}
	//}

	// 处理成功请求记录
	//b.DoHistory(req, true)

	// 统计成功页数
	//cache.PageSuccCount()

	// 提示抓取成功
	//logs.Log.Informational(" *     Success: %v\n", downUrl)

	// 释放ctx准备复用
	databox.PutContext(ctx)
}
