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
	r.Router.POST("/api/dmp/orderRouteQry/", NewDemService().SendDemToSup)
	r.Router.POST("/api/dem/send", NewDemService().SendDemReqToSup)
	r.Router.POST("/api/sup/rec", NewSupService().RecDemReqAndPushToSup)
}