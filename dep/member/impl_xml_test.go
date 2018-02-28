package member

import "testing"

func initMemberManagerXMLFile() (*MemberManagerXMLFile, error) {
	return NewMemberManagerXMLFile("memInfo.xml")
}

func TestGetMemberInfo(t *testing.T) {
	memberManager, err := initMemberManagerXMLFile()
	if err != nil {
		t.Fatal(err)
	}

	memberInfo := memberManager.GetMemberInfo("0000109")
	if memberInfo == nil {
		t.Fatal("not found")
	}
	if memberInfo.MemID != "0000109" {
		t.Fatalf("%s != %s", "0000109", memberInfo.MemID)
	}
	if memberInfo.Status != "01" {
		t.Fatalf("%s != %s", "01", memberInfo.Status)
	}
	if memberInfo.TotLmt != 1000.00 {
		t.Fatalf("%f != %f", 1000.00, memberInfo.TotLmt)
	}
}
