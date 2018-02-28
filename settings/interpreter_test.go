package settings

import "testing"

func TestInterpreter1(t *testing.T) {
	// interpreter := &Interpreter{}
	parser, err := Interpret("a")
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Println(err)
	// fmt.Println(parser)
	provider := make(map[interface{}]interface{})

	key := "a"
	value := "b"
	provider[key] = value

	output, err := parser.Parse(provider)
	if err != nil {
		t.Fatal(err)
	}

	if output != value {
		t.Fatal("Not equal", output, value)
	}
}

func TestInterpreter2(t *testing.T) {
	// interpreter := &Interpreter{}
	parser, err := Interpret("a.b")
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Println(err)
	// fmt.Println(parser)
	provider := make(map[interface{}]interface{})

	key := "a"
	value := make(map[interface{}]interface{})
	value["b"] = "abcdefg"
	provider[key] = value

	output, err := parser.Parse(provider)
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Println(output)
	if output != value["b"] {
		t.Fatal("Not equal", output, value["b"])
	}
}

func TestInterpreter3(t *testing.T) {
	// interpreter := &Interpreter{}
	parser, err := Interpret("a[0]")
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Println(err)
	// fmt.Println(parser)
	provider := make(map[interface{}]interface{})

	key := "a"
	value := make([]interface{}, 1)
	value[0] = "abcdefg"
	provider[key] = value

	output, err := parser.Parse(provider)
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Println(output)
	if output != value[0] {
		t.Fatal("Not equal", output, value[0])
	}
}
