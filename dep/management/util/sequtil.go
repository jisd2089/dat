package util

import (
	"fmt"
	"math/rand"
	"time"
	"sync"
)

/**
    Author: luzequan
    Created: 2018-01-15 15:10:04
*/
type SeqUtil struct {
}

var (
	newSeqUtil   *SeqUtil
	once            sync.Once
)

func NewSeqUtil() *SeqUtil {
	return getSeqUtil()
}

func getSeqUtil() *SeqUtil {
	once.Do(func() {
		newSeqUtil = &SeqUtil{}
	})
	return newSeqUtil
}

func (s *SeqUtil) RandInt(min int, max int) int {
	//rand.Seed(time.Now().UTC().UnixNano())
	//fmt.Printf("min : %d, max : %d \n", min, max)
	return min + rand.Intn(max-min)
}

func (s *SeqUtil) GenBusiSerialNo(memId string) string {
	now := time.Now()
	//00001092017042801452691075241928197
	serialNo := fmt.Sprintf("%s%04d%02d%02d%02d%02d%02d%09d%05d", memId, now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond(),
		s.RandInt(1000, 99999))

	return serialNo
}