package settings

// Nil 当从Settings缺失设置时返回此error
const Nil = SettingsError("Settings: nil")

type SettingsError string

func (e SettingsError) Error() string { return string(e) }

type Settings interface {
	Get(path string) (interface{}, error)

	GetString(path string) (string, error)

	GetInt(path string) (int, error)
	GetUint(path string) (uint, error)

	GetInt8(path string) (int8, error)
	GetUint8(path string) (uint8, error)

	GetInt16(path string) (int16, error)
	GetUint16(path string) (uint16, error)

	GetInt32(path string) (int32, error)
	GetUint32(path string) (uint32, error)

	GetInt64(path string) (int64, error)
	GetUint64(path string) (uint64, error)

	GetFloat32(path string) (float32, error)
	GetFloat64(path string) (float64, error)

	GetBool(path string) (bool, error)

	GetSlice(path string) ([]interface{}, error)
	GetMap(path string) (map[interface{}]interface{}, error)

	GetStruct(path string, receiver interface{}) error

	GetSettings(path string) (Settings, error)
}

var defaultSettings = &DefaultSettings{make(map[interface{}]interface{})}

// GetSettings 获取设置类实例
func GetSettings() Settings {
	return defaultSettings
}
