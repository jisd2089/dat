package sftp

import "time"

/**
    Author: luzequan
    Created: 2017-12-29 17:37:35
*/
type FileCatalog struct {
	UserName       string
	Password       string
	Host           string
	Port           int
	TimeOut        time.Duration
	LocalDir       string
	LocalFileName  string
	RemoteDir      string
	RemoteFileName string
}
