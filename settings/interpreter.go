package settings

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

const (
	INTERPRETER_START = iota
	INTERPRETER_BEGIN_FIELD
	INTERPRETER_BEGIN_ARRAY
	INTERPRETER_END_ARRAY
)

type Parser interface {
	Parse(obj interface{}) (interface{}, error)
}

var (
	supportedMapType   = reflect.TypeOf((map[interface{}]interface{})(nil))
	supportedSliceType = reflect.TypeOf(([]interface{})(nil))
)

type MapParser struct {
	key string
}

func (parser *MapParser) Parse(obj interface{}) (interface{}, error) {
	inType := reflect.TypeOf(obj)
	if inType != supportedMapType {
		return nil, fmt.Errorf("%s: Need %s type, but got %s", parser.key, supportedMapType, inType)
	}

	actual := obj.(map[interface{}]interface{})
	return actual[parser.key], nil
}

func (parser *MapParser) String() string {
	return fmt.Sprintf("MapParser: %s", parser.key)
}

type ArrayParser struct {
	index int
}

func (parser *ArrayParser) Parse(obj interface{}) (interface{}, error) {
	inType := reflect.TypeOf(obj)
	if inType != supportedSliceType {
		return nil, fmt.Errorf("%d: Need %s type, but got %s", parser.index, supportedSliceType, inType)
	}

	actual := obj.([]interface{})
	if parser.index >= len(actual) {
		return nil, fmt.Errorf("%d: index out of boundary", parser.index)
	}

	return actual[parser.index], nil
}

func (parser *ArrayParser) String() string {
	return fmt.Sprintf("ArrayParser: %d", parser.index)
}

type DelegatingParser struct {
	parent    Parser
	delegated Parser
}

func (parser *DelegatingParser) Parse(obj interface{}) (interface{}, error) {
	var err error
	parentOutput := obj
	if parser.parent != nil {
		parentOutput, err = parser.parent.Parse(obj)
		if err != nil {
			return nil, err
		}
	}
	output, err := parser.delegated.Parse(parentOutput)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (parser *DelegatingParser) String() string {
	return fmt.Sprintf("[%s] %s", parser.parent, parser.delegated)
}

func newParserForField(fieldName string, parent Parser) Parser {
	fieldParser := &MapParser{fieldName}
	delegatingParser := &DelegatingParser{parent, fieldParser}
	return delegatingParser
}

func newParserForArray(arrayIndex int, parent Parser) Parser {
	arrayParser := &ArrayParser{arrayIndex}
	delegatingParser := &DelegatingParser{parent, arrayParser}
	return delegatingParser
}

// func (interpreter *Interpreter) Interpret(path string) (Parser, error) {
func Interpret(path string) (Parser, error) {
	var parent Parser
	var text bytes.Buffer
	var index int

	// TODO 使用rune来处理
	input := []byte(path)
	state := INTERPRETER_START
	for index = 0; index < len(input); index++ {
		char := input[index]
		if !isSupportedChar(char) {
			return nil, fmt.Errorf("Invalid character at position %d", index)
		}

		switch state {
		case INTERPRETER_START:
			switch char {
			case '.', '[', ']':
				return nil, fmt.Errorf("Invalid character at position %d", index)
			default:
				text.Reset()
				text.WriteByte(char)
				state = INTERPRETER_BEGIN_FIELD
			}
		case INTERPRETER_BEGIN_FIELD:
			switch char {
			case '.':
				if text.Len() == 0 {
					return nil, fmt.Errorf("Invalid character at position %d", index)
				}
				fieldName := text.String()
				parent = newParserForField(fieldName, parent)
				text.Reset()
				state = INTERPRETER_BEGIN_FIELD
			case '[':
				if text.Len() == 0 {
					return nil, fmt.Errorf("Invalid character at position %d", index)
				}
				fieldName := text.String()
				parent = newParserForField(fieldName, parent)
				text.Reset()
				state = INTERPRETER_BEGIN_ARRAY
			case ']':
				return nil, fmt.Errorf("Invalid character at position %d", index)
			default:
				text.WriteByte(char)
			}
		case INTERPRETER_BEGIN_ARRAY:
			switch char {
			case ']':
				if text.Len() == 0 {
					return nil, fmt.Errorf("Invalid character at position %d", index)
				}
				arryIndex, _ := strconv.Atoi(text.String())
				parent = newParserForArray(arryIndex, parent)
				state = INTERPRETER_END_ARRAY
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				text.WriteByte(char)
			default:
				return nil, fmt.Errorf("Invalid character at position %d", index)
			}
		case INTERPRETER_END_ARRAY:
			switch char {
			case '.':
				text.Reset()
				state = INTERPRETER_BEGIN_FIELD
			case '[':
				text.Reset()
				state = INTERPRETER_BEGIN_ARRAY
			default:
				return nil, fmt.Errorf("Invalid character at position %d", index)
			}
		}
	}

	switch state {
	case INTERPRETER_BEGIN_FIELD:
		if text.Len() == 0 {
			return nil, fmt.Errorf("Invalid character at position %d", index)
		}
		fieldName := text.String()
		return newParserForField(fieldName, parent), nil
	case INTERPRETER_BEGIN_ARRAY:
		return nil, fmt.Errorf("Invalid character at position %d", index)
	case INTERPRETER_END_ARRAY:
		return parent, nil
	}
	// never arrive
	return parent, nil
}

func isSupportedChar(char byte) bool {
	return (char > 47 && char < 58) ||
		(char > 64 && char < 91) ||
		(char > 96 && char < 123) ||
		char == 46 ||
		char == 95 ||
		char == 91 || char == 93
	// char == 123 || char == 125
}
