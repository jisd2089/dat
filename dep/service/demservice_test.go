package service

/**
    Author: luzequan
    Created: 2018-01-10 15:47:42
*/

import (
	"testing"
	"github.com/valyala/fasthttp"
)
func TestSendDemReqToSup(t *testing.T) {
	NewDemService().SendDemReqToSup(&fasthttp.RequestCtx{})
}