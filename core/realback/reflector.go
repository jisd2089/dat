package realback

import (
	"dat/core/interaction/response"
	"dat/core/interaction/request"
)

/**
    Author: luzequan
    Created: 2018-01-12 22:09:39
*/
type Reflector interface {
	Handle(*request.DataRequest) *response.DataResponse
}