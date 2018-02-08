package web

/**
    Author: luzequan
    Created: 2018-01-02 19:35:21
*/
import (
	"github.com/buaazp/fasthttprouter"
	."dat/dep/service"
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
	r.Router.POST("/api/dem/read/", NewDemService().ReadFile)
	r.Router.POST("/api/dem/split/", NewDemService().SplitFile)
	r.Router.POST("/api/dem/send", NewDemService().SendDemReqToSup)
	r.Router.POST("/api/dem/rec", NewDemService().RecSupRespAndPushToDem)
	r.Router.POST("/api/dem/subbox", NewDemService().RunParentAndChild)


	r.Router.POST("/api/sup/rec", NewSupService().RecDemReqAndPushToSup)
	r.Router.POST("/api/sup/send", NewSupService().SupRespSendToDem)
	r.Router.POST("/api/sup/sendfull", NewSupService().SupRespWholeSendToDem)
	r.Router.POST("/api/sup/sendcompress", NewSupService().SupCompressFileSendToDem)
}