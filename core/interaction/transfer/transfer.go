package transfer

/**
    Author: luzequan
    Created: 2017-12-28 14:39:09
*/

type Transfer interface {
	ExecuteMethod(Request) Response
	Close()
}
