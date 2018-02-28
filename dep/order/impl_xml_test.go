package order

import "testing"

func initOrderManagerXMLFile() (*OrderManagerXMLFile, error) {
	return NewOrderManagerXMLFile("order_route_ODN20170418000000183.xml")
}

func TestGetOrderInfo(t *testing.T) {
	orderManager, err := initOrderManagerXMLFile()
	if err != nil {
		t.Fatal(err)
	}

	orderInfo := orderManager.GetOrderInfo()
	if orderInfo == nil {
		t.Fatal("not found")
	}
	taskInfo := orderInfo.GetTaskInfo("CTN20161223000010400001090000087")
	if taskInfo == nil {
		t.Fatal("taskinfo not found")
	}

	if taskInfo.SupMemID != "0000138" {
		t.Fatalf("%s != %s", "0000138", taskInfo.SupMemID)
	}
	if taskInfo.ConnObjNo != "DOS20170320000000566" {
		t.Fatalf("%s != %s", "DOS20170320000000566", taskInfo.ConnObjNo)
	}
	if taskInfo.ValuationPrice != 0.01 {
		t.Fatalf("%f != %f", 0.01, taskInfo.ValuationPrice)
	}
	if taskInfo.CacheTime != 36000 {
		t.Fatalf("%d != %d", 36000, taskInfo.CacheTime)
	}
}
