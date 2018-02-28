package cncrypt

import (
	"fmt"
	"testing"
)

func Test_Exid_to_Etid(t *testing.T) {
	var exid = "01234567890123456789012345678912"
	etid, err := EXID_to_ETID(exid)
	if err != nil {
		fmt.Println(err)
		t.Error("not passed")

	}
	fmt.Println(etid)
	exid_res, err := ETID_to_EXID(etid)
	if err != nil {
		fmt.Println(err)
		t.Error("not passed")
	}
	fmt.Println(exid_res)
	if exid != exid_res {
		t.Error("value wrong!!")
	}

}
