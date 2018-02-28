package settings

import (
	"fmt"
	"reflect"
)

type DefaultSettings struct {
	dataProvider map[interface{}]interface{}
}

func (settings *DefaultSettings) Get(path string) (interface{}, error) {
	parser, err := Interpret(path)
	if err != nil {
		return nil, err
	}

	result, err := parser.Parse(settings.dataProvider)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, Nil
	}

	return result, nil
}

func (settings *DefaultSettings) GetString(path string) (string, error) {
	value, err := settings.Get(path)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(value), err
}

func (settings *DefaultSettings) GetInt(path string) (int, error) {
	value, err := settings.GetInt64(path)
	if err != nil {
		return 0, err
	}
	return int(value), nil
}

func (settings *DefaultSettings) GetUint(path string) (uint, error) {
	value, err := settings.GetUint64(path)
	if err != nil {
		return 0, err
	}

	return uint(value), nil
}

func (settings *DefaultSettings) GetInt8(path string) (int8, error) {
	value, err := settings.GetInt64(path)
	if err != nil {
		return 0, err
	}
	return int8(value), nil
}

func (settings *DefaultSettings) GetUint8(path string) (uint8, error) {
	value, err := settings.GetUint64(path)
	if err != nil {
		return 0, err
	}

	return uint8(value), nil
}

func (settings *DefaultSettings) GetInt16(path string) (int16, error) {
	value, err := settings.GetInt64(path)
	if err != nil {
		return 0, err
	}
	return int16(value), nil
}

func (settings *DefaultSettings) GetUint16(path string) (uint16, error) {
	value, err := settings.GetUint64(path)
	if err != nil {
		return 0, err
	}

	return uint16(value), nil
}

func (settings *DefaultSettings) GetInt32(path string) (int32, error) {
	value, err := settings.GetInt64(path)
	if err != nil {
		return 0, err
	}
	return int32(value), nil
}

func (settings *DefaultSettings) GetUint32(path string) (uint32, error) {
	value, err := settings.GetUint64(path)
	if err != nil {
		return 0, err
	}

	return uint32(value), nil
}

func (settings *DefaultSettings) GetInt64(path string) (int64, error) {
	value, err := settings.Get(path)
	if err != nil {
		return 0, err
	}

	switch value.(type) {
	case int:
		intValue := value.(int)
		return int64(intValue), nil
	case uint:
		intValue := value.(uint)
		return int64(intValue), nil
	case int8:
		intValue := value.(int8)
		return int64(intValue), nil
	case uint8:
		intValue := value.(uint8)
		return int64(intValue), nil
	case int16:
		intValue := value.(int16)
		return int64(intValue), nil
	case uint16:
		intValue := value.(uint16)
		return int64(intValue), nil
	case int32:
		intValue := value.(int32)
		return int64(intValue), nil
	case uint32:
		intValue := value.(uint32)
		return int64(intValue), nil
	case int64:
		return value.(int64), nil
	case uint64:
		intValue := value.(uint64)
		return int64(intValue), nil
	case float32:
		floatValue := value.(float32)
		return int64(floatValue), nil
	case float64:
		floatValue := value.(float64)
		return int64(floatValue), nil
	default:
		return 0, fmt.Errorf("type error: %s", reflect.TypeOf(value))
	}
}

func (settings *DefaultSettings) GetUint64(path string) (uint64, error) {
	value, err := settings.Get(path)
	if err != nil {
		return 0, err
	}

	switch value.(type) {
	case int:
		intValue := value.(int)
		return uint64(intValue), nil
	case uint:
		intValue := value.(uint)
		return uint64(intValue), nil
	case int8:
		intValue := value.(int8)
		return uint64(intValue), nil
	case uint8:
		intValue := value.(uint8)
		return uint64(intValue), nil
	case int16:
		intValue := value.(int16)
		return uint64(intValue), nil
	case uint16:
		intValue := value.(uint16)
		return uint64(intValue), nil
	case int32:
		intValue := value.(int32)
		return uint64(intValue), nil
	case uint32:
		intValue := value.(uint32)
		return uint64(intValue), nil
	case int64:
		intValue := value.(int64)
		return uint64(intValue), nil
	case uint64:
		return value.(uint64), nil
	case float32:
		floatValue := value.(float32)
		return uint64(floatValue), nil
	case float64:
		floatValue := value.(float64)
		return uint64(floatValue), nil
	default:
		return 0, fmt.Errorf("type error: %s", reflect.TypeOf(value))
	}
}

func (settings *DefaultSettings) GetFloat32(path string) (float32, error) {
	value, err := settings.GetFloat64(path)
	if err != nil {
		return 0, err
	}

	return float32(value), nil
}

func (settings *DefaultSettings) GetFloat64(path string) (float64, error) {
	value, err := settings.Get(path)
	if err != nil {
		return 0, err
	}

	switch value.(type) {
	case int:
		intValue := value.(int)
		return float64(intValue), nil
	case uint:
		intValue := value.(uint)
		return float64(intValue), nil
	case int8:
		intValue := value.(int8)
		return float64(intValue), nil
	case uint8:
		intValue := value.(uint8)
		return float64(intValue), nil
	case int16:
		intValue := value.(int16)
		return float64(intValue), nil
	case uint16:
		intValue := value.(uint16)
		return float64(intValue), nil
	case int32:
		intValue := value.(int32)
		return float64(intValue), nil
	case uint32:
		intValue := value.(uint32)
		return float64(intValue), nil
	case int64:
		intValue := value.(int64)
		return float64(intValue), nil
	case uint64:
		intValue := value.(uint64)
		return float64(intValue), nil
	case float32:
		floatValue := value.(float32)
		return float64(floatValue), nil
	case float64:
		return value.(float64), nil
	default:
		return 0, fmt.Errorf("type error: %s", reflect.TypeOf(value))
	}
}

func (settings *DefaultSettings) GetBool(path string) (bool, error) {
	value, err := settings.Get(path)
	if err != nil {
		return false, err
	}

	if reflect.TypeOf(value).Kind() == reflect.Bool {
		return value.(bool), nil
	}
	return false, fmt.Errorf("type error: %s", reflect.TypeOf(value))
}

func (settings *DefaultSettings) GetSlice(path string) ([]interface{}, error) {
	value, err := settings.Get(path)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(value) == reflect.TypeOf([]interface{}(nil)) {
		return value.([]interface{}), nil
	}
	return nil, fmt.Errorf("type error: %s", reflect.TypeOf(value))
}

func (settings *DefaultSettings) GetMap(path string) (map[interface{}]interface{}, error) {
	value, err := settings.Get(path)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(value).Kind() == reflect.Map {
		return value.(map[interface{}]interface{}), nil
	}
	return nil, fmt.Errorf("type error: %s", reflect.TypeOf(value))
}

func (settings *DefaultSettings) GetStruct(path string, receiver interface{}) error {
	// TODO shixian
	value, err := settings.Get(path)
	if err != nil {
		return err
	}

	return fmt.Errorf("type error: %s", reflect.TypeOf(value))
}

func (settings *DefaultSettings) GetSettings(path string) (Settings, error) {
	provider, err := settings.GetMap(path)
	if err != nil {
		return nil, err
	}

	return &DefaultSettings{provider}, nil
}
