package util

import (
	"fmt"
	"math/rand"
	"time"
)

/**
    Author: luzequan
    Created: 2018-06-26 09:47:09
*/
type SeqUtil struct {}

func (s SeqUtil) RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func (s SeqUtil) GenBusiSerialNo(memId string) string {
	now := time.Now()
	//00001092017042801452691075241928197
	serialNo := fmt.Sprintf("%s%04d%02d%02d%02d%02d%02d%09d%05d", memId, now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond(),
		s.RandInt(1000, 99999))

	return serialNo
}