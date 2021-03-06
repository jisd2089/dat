package interaction

/**
    Author: luzequan
    Created: 2017-12-28 14:15:55
*/

import (
	"drcs/core/databox"
	"drcs/core/interaction/request"
)

// The Handler interface.
// You can implement the interface by implement function Handler.
// Function Handler need to return http response from Request.
type Carrier interface {
	Handle(*databox.DataBox, *request.DataRequest) *databox.Context
	Close()
	Process(*databox.DataBox, *request.DataRequest)
}