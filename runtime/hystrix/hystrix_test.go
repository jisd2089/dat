package hystrix

/**
    Author: luzequan
    Created: 2018-02-26 16:15:48
*/

import (
	"testing"
	"net"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
)

func TestHystrix(t *testing.T) {

	hystrix.Go("my_command", func() error {
		// talk to other services
		return nil
	}, nil)

}

func TestHystrixDashboard(t *testing.T) {
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "81"), hystrixStreamHandler)
}