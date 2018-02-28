package settings

import (
	"testing"
	"fmt"
)

var settings Settings

func init() {
	settings, _ = CreateSettingsFromYAML("settings.yaml")
}

func TestGet(t *testing.T) {
	value, err := settings.Get("a")
	if err != nil {
		t.Fatal(err)
	}

	if value != "Easy" {
		t.Fatalf("%s != %s", "Easy", value)
	}
}

func TestGetString(t *testing.T) {
	value, err := settings.GetString("a")
	if err != nil {
		t.Fatal(err)
	}

	if value != "Easy" {
		t.Fatalf("%s != %s", "Easy", value)
	}
}

func TestGetInt(t *testing.T) {
	value, err := settings.GetInt("b")
	if err != nil {
		t.Fatal(err)
	}

	if value != 1 {
		t.Fatalf("%d != %d", 1, value)
	}

	value, err = settings.GetInt("c")
	if err != nil {
		t.Fatal(err)
	}

	if value != -1 {
		t.Fatalf("%d != %d", -1, value)
	}
}

func TestGetUint(t *testing.T) {
	value, err := settings.GetUint("b")
	if err != nil {
		t.Fatal(err)
	}

	if value != 1 {
		t.Fatalf("%d != %d", 1, value)
	}
}

func TestGetInt8(t *testing.T) {
	value, err := settings.GetInt8("b")
	if err != nil {
		t.Fatal(err)
	}

	if value != 1 {
		t.Fatalf("%d != %d", 1, value)
	}

	value, err = settings.GetInt8("c")
	if err != nil {
		t.Fatal(err)
	}

	if value != -1 {
		t.Fatalf("%d != %d", -1, value)
	}
}

func TestGetUint8(t *testing.T) {
	value, err := settings.GetUint8("b")
	if err != nil {
		t.Fatal(err)
	}

	if value != 1 {
		t.Fatalf("%d != %d", 1, value)
	}
}

func TestGetInt16(t *testing.T) {
	value, err := settings.GetInt16("b")
	if err != nil {
		t.Fatal(err)
	}

	if value != 1 {
		t.Fatalf("%d != %d", 1, value)
	}

	value, err = settings.GetInt16("c")
	if err != nil {
		t.Fatal(err)
	}

	if value != -1 {
		t.Fatalf("%d != %d", -1, value)
	}
}

func TestGetUint16(t *testing.T) {
	value, err := settings.GetUint16("b")
	if err != nil {
		t.Fatal(err)
	}

	if value != 1 {
		t.Fatalf("%d != %d", 1, value)
	}
}

func TestGetInt32(t *testing.T) {
	value, err := settings.GetInt32("b")
	if err != nil {
		t.Fatal(err)
	}

	if value != 1 {
		t.Fatalf("%d != %d", 1, value)
	}

	value, err = settings.GetInt32("c")
	if err != nil {
		t.Fatal(err)
	}

	if value != -1 {
		t.Fatalf("%d != %d", -1, value)
	}
}

func TestGetUint32(t *testing.T) {
	value, err := settings.GetUint32("b")
	if err != nil {
		t.Fatal(err)
	}

	if value != 1 {
		t.Fatalf("%d != %d", 1, value)
	}
}

func TestGetInt64(t *testing.T) {
	value, err := settings.GetInt64("b")
	if err != nil {
		t.Fatal(err)
	}

	if value != 1 {
		t.Fatalf("%d != %d", 1, value)
	}

	value, err = settings.GetInt64("c")
	if err != nil {
		t.Fatal(err)
	}

	if value != -1 {
		t.Fatalf("%d != %d", -1, value)
	}
}

func TestGetUint64(t *testing.T) {
	value, err := settings.GetUint64("b")
	if err != nil {
		t.Fatal(err)
	}

	if value != 1 {
		t.Fatalf("%d != %d", 1, value)
	}
}

func TestGetFloat32(t *testing.T) {
	value, err := settings.GetFloat32("d")
	if err != nil {
		t.Fatal(err)
	}

	if value != 11.1 {
		t.Fatalf("%f != %f", 11.1, value)
	}

	value, err = settings.GetFloat32("f")
	if err != nil {
		t.Fatal(err)
	}

	if value != -11.1 {
		t.Fatalf("%f != %f", -11.1, value)
	}
}

func TestGetFloat64(t *testing.T) {
	value, err := settings.GetFloat64("d")
	if err != nil {
		t.Fatal(err)
	}

	if value != 11.1 {
		t.Fatalf("%f != %f", 11.1, value)
	}

	value, err = settings.GetFloat64("f")
	if err != nil {
		t.Fatal(err)
	}

	if value != -11.1 {
		t.Fatalf("%f != %f", -11.1, value)
	}
}

func TestGetSlice(t *testing.T) {
	value, err := settings.GetSlice("g")
	if err != nil {
		t.Fatal(err)
	}

	// fmt.Println(value)
	if value[0] != 1 {
		t.Fatalf("%d != %d", 1, value[0])
	}
}

func TestGetMap(t *testing.T) {
	value, err := settings.GetMap("h")
	if err != nil {
		t.Fatal(err)
	}

	// fmt.Println(value)
	if value["i"] != 1 {
		t.Fatalf("%d != %d", 1, value["i"])
	}
}

func TestGetSettings(t *testing.T) {
	isettings, err := settings.GetSettings("h")
	if err != nil {
		t.Fatal(err)
	}

	value, err := isettings.GetInt("i")
	if err != nil {
		t.Fatal(err)
	}

	if value != 1 {
		t.Fatalf("%d != %d", 1, value)
	}
}

func TestGetSettings2(t *testing.T)  {
	isettings, err := settings.GetMap("h")
	if err != nil {
		t.Fatal(err)
	}
	kmap := isettings["i"].(map[interface{}]interface{})
	fmt.Println(kmap["k"])

}
