package xid

import (
	"testing"
	"fmt"
)

var (
	appkey1 = "FAE11C71715C232ADBE2FE422FA67209"
	idnum1  = "320123197507134111"
	idtype1 = "ID010104"
	//
	app_id_s1  = "IDCORRECTAPP00000000"
	regcode_s1 = "5H1Xd0nlVaG/64bnKTH9/BEjnLVRfaBEioTwEsrBpfPEmeTY2ptzP8RAilj+6LgN"
	app_id1    = app_id_s1
	//163
	app_id_d1  = "APPXIDCODED000000000"
	regcode_d1 = "stXue5Ihsw2bCn7W8ir8yo726OU+fb6De+c2//zanuXvtjki5Fl3sk7B9zvqPH0z"
	//167
	//app_id_d1  = "01JR1502021030342817"
	//regcode_d1 = "JwTZRiGseMv74e45YFLyvqIxGk+kGZdnXORkt1SMwGJTxFcPFMsvgI2onGdbg8Yx"
	//app_id_d1  = "01JR1502021030342818"
	//regcode_d1 = "JwTZRiGseMv74e45YFLyvo5ko/+lEuMGDKM69H4DdiweFtjKxV61njP3ox+9bp+t"
	//app_id_d1  = "01JR1502021030342819"
	//regcode_d1 = "JwTZRiGseMv74e45YFLyvnhZlN1SRVkDwB08TVD/ABjSa/hHKb4+4XOKAkmCo2Jc"

	appxidcode_s1 = "ID010104zfWhg6UiSHNZZz24VBxd1u7y+iyqhVY6aoCEyJyPMs4LhCvp"
)

func Test_GenXID(t *testing.T) {
	xid, err := GenXID(app_id1, idtype1, idnum1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(xid)
}

func Test_ConvertID(t *testing.T) {
	xid, err := ConvertXID(app_id_s1, regcode_s1, app_id_d1, regcode_d1, appxidcode_s1)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	fmt.Println(xid)
}
