package service

/**
    Author: luzequan
    Created: 2018-01-02 19:35:21
*/
import (
	"github.com/buaazp/fasthttprouter"
)
type HttpRouter struct {
	Router *fasthttprouter.Router
}

func (r *HttpRouter) Register() {
	r.Router = fasthttprouter.New()
	r.Router.POST("/api/dmp/orderRouteQry/", NewDemService().SendDemToSup)
}