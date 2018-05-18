package web

/**
    Author: luzequan
    Created: 2018-01-02 19:35:21
*/
import (
	"github.com/buaazp/fasthttprouter"
	."drcs/dep/handler"
)

type HttpRouter struct {
	Router *fasthttprouter.Router
}

func NewHttpRouter() *HttpRouter {
	return &HttpRouter{
		Router: fasthttprouter.New(),
	}
}

func (r *HttpRouter) Register() {
	r.Router.POST("/api/dem/read/", NewDemHandler().ReadFile)
	r.Router.POST("/api/dem/split/", NewDemHandler().SplitFile)
	r.Router.POST("/api/dem/send", NewDemHandler().SendDemReqToSup)
	r.Router.POST("/api/dem/rec", NewDemHandler().RecSupRespAndPushToDem)
	r.Router.POST("/api/dem/subbox", NewDemHandler().RunParentAndChild)


	r.Router.POST("/api/sup/rec", NewSupHandler().RecDemReqAndPushToSup)
	r.Router.POST("/api/sup/send", NewSupHandler().SupRespSendToDem)
	r.Router.POST("/api/sup/sendfull", NewSupHandler().SupRespWholeSendToDem)
	r.Router.POST("/api/sup/sendcompress", NewSupHandler().SupCompressFileSendToDem)

	r.Router.POST("/api/test/rcvfile", NewNodeHandler().RcvData)
	r.Router.POST("/api/test/rcvalg", NewNodeHandler().RcvAlg)
	r.Router.POST("/api/fusion/run", NewNodeHandler().RunProcess)

	// 生产业务流程
	//r.Router.POST("/api/dmp/orderRouteQry/", demander.HTTPService{}.DoService)
	//r.Router.POST("/api/p/pushFile/", common.CommonSvc.PlatPushFile)
	//r.Router.POST("/api/p/genKey", common.CommonSvc.GenKeys)
	//r.Router.POST("/api/p/initSafeConfig/", common.CommonSvc.InitSafeConfig)
	//r.Router.POST("/api/p/batchfile/", common.CommonSvc.AcceptBatchFile)
	//
	//r.Router.POST("/api/d/qryData/", singleQuery)
	//r.Router.POST("/api/p/pushFile/", common.CommonSvc.PlatPushFile)
	//r.Router.POST("/api/p/genKey", common.CommonSvc.GenKeys)
	//r.Router.POST("/api/p/initSafeConfig/", common.CommonSvc.InitSafeConfig)
	//r.Router.POST("/api/p/batchfile/", common.CommonSvc.AcceptBatchFile)
	//r.Router.POST("/api/p/batchmode/", supplier.BatchSvc.Serve)
}