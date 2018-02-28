package handler

/**
    Author: luzequan
    Created: 2018-01-10 15:47:42
*/

import (
	"testing"
	"github.com/valyala/fasthttp"
	"github.com/aarzilli/golua/lua"
)
func TestSendDemReqToSup(t *testing.T) {
	NewDemHandler().SendDemReqToSup(&fasthttp.RequestCtx{})
}

func TestLua(t *testing.T) {
	lua.NewState()
}