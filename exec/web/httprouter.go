package web

/**
    Author: luzequan
    Created: 2018-01-02 19:35:21
*/
import (
	"github.com/buaazp/fasthttprouter"
	. "drcs/dep/handler"
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
	//r.Router.POST("/api/dem/read/", NewDemHandler().ReadFile)
	//r.Router.POST("/api/dem/split/", NewDemHandler().SplitFile)
	//r.Router.POST("/api/dem/send", NewDemHandler().SendDemReqToSup)
	//r.Router.POST("/api/dem/rec", NewDemHandler().RecSupRespAndPushToDem)
	//r.Router.POST("/api/dem/subbox", NewDemHandler().RunParentAndChild)
	//
	//
	//r.Router.POST("/api/sup/rec", NewSupHandler().RecDemReqAndPushToSup)
	//r.Router.POST("/api/sup/send", NewSupHandler().SupRespSendToDem)
	//r.Router.POST("/api/sup/sendfull", NewSupHandler().SupRespWholeSendToDem)
	//r.Router.POST("/api/sup/sendcompress", NewSupHandler().SupCompressFileSendToDem)
	//r.Router.POST("/api/mq/send", NewDemHandler().SendMQ)
	//
	//r.Router.POST("/api/test/rcvfile", NewNodeHandler().RcvData)
	//r.Router.POST("/api/test/rcvalg", NewNodeHandler().RcvAlg)
	//r.Router.POST("/api/fusion/run", NewNodeHandler().RunProcess)
	//r.Router.POST("/api/dis/batch", NewNodeHandler().RunBatchProcess)
	//r.Router.POST("/api/rcv/batch", NewNodeHandler().RunBatchRcv)
	//r.Router.POST("/api/crp/dem", NewNodeHandler().RunCRPProcess)
	r.Router.POST("/api/p/crp", NewNodeHandler().RunCRPProcess)
	r.Router.POST("/api/crp/sup", NewNodeHandler().RunCRPResponse)
	r.Router.POST("/api/p/genKey", NewNodeHandler().GenKeys)
	r.Router.POST("/api/p/initSafeConfig/", NewNodeHandler().InitSecurityConfig)

}