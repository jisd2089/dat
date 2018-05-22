package settings

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	//defaultCommonSettingYaml = "/dep/go/conf/settings.yaml"
	defaultDemSettingYaml = "settings.yaml"
	defaultTag            = "yaml"
)

var commonSettings *CommonSettings

type CommonSettings struct {
	Node          Node          `yaml:"node"`
	ConfigFile    ConfigFile    `yaml:"configfile"`
	NodeService   NodeService   `yaml:"service"`
	Xid           Xid           `yaml:"xid"`
	Redis         Redis         `yaml:"redis"`
	Conf          Conf          `yaml:"conf"`
	Routine       Routine       `yaml:"routine"`
	Kafka         Kafka         `yaml:"kafka"`
	BusiLog       BusiLog       `yaml:"busilog"`
	DMP           DMP           `yaml:"dmp"`
	Log           Log           `yaml:"log"`
	Sftp          Sftp          `yaml:"sftp"`
	Hdfs          Hdfs          `yaml:"hdfs"`
	BatchCollison BatchCollison `yaml:"batchCollison"`
	Other         Other         `yaml:"other"`
}

type Node struct {
	NodeId   string `yaml:"nodeId"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	MemberId string `yaml:"memberId"`
	Role     string `yaml:"role"`
	Version  string `yaml:"version"`
	Userkey  string `yaml:"userkey"`
	Token    string `yaml:"token"`
	DlsUrl   string `yaml:"dlsUrl"`
}

type ConfigFile struct {
	KeysFile             string `yaml:"keysFile"`
	OrderFile            string `yaml:"orderFile"`
	OrderFileExpireTime  int64  `yaml:"orderFileExpireTime"`
	MemberFile           string `yaml:"memberFile"`
	MemberFileExpireTime int64  `yaml:"memberFileExpireTime"`
	OrderRouteFilePath   string `yaml:"orderRouteFilePath"`
}

type Xid struct {
	XidIp     string `yaml:"xidIp"`
	XidDealer string `yaml:"xidDealer"`
}

type NodeService struct {
	Protocol   string `yaml:"protocol"`
	ListenIP   string `yaml:"listenIp"`
	ListenPort int    `yaml:"listenPort"`
	ServiceUrl string `yaml:"serviceUrl"`
}

type Conf struct {
	CfgDir        string `yaml:"cfgDir"`
	LogDir        string `yaml:"logDir"`
	XmlDir        string `yaml:"xmlDir"`
	XmlReloadTime int    `yaml:"xmlReloadTime"`
}

type Redis struct {
	Addr       string `yaml:"addr"`
	DB         int    `yaml:"db"`
	PoolSize   int    `yaml:"poolSize"`
	ReadOnly   bool   `yaml:"readOnly"`
	BucketSize int    `yaml:"bucketSize"`

	ReadTimeout  int `yaml:"readTimeout"`
	WriteTimeout int `yaml:"writeTimeout"`
	DialTimeout  int `yaml:"dialTimeout"`

	Mode string `yaml:"mode"`
}

type Routine struct {
	MemberIds []string `yaml:"memberId"`
	Capacity  []int    `yaml:"capacity"`
	MCMap     map[string]int
}

type Kafka struct {
	Brokers      []string `yaml:"brokers"`
	Topic        string   `yaml:"topic"`
	PartitionNum int      `yaml:"partitionNum"`
}

type BusiLog struct {
	MQPath string `yaml:"mqPath"`
}

type DMP struct {
	URL     string `yaml:"url"`
	Timeout int    `yaml:"timeout"`
}

type Log struct {
	ConfigPath string `yaml:"configPath"`
	LogrusPath string `yaml:"logrusPath"`
}

type Sftp struct {
	Hosts              string `yaml:"hosts"`
	Port               int    `yaml:"port"`
	Username           string `yaml:"username"`
	Password           string `yaml:"password"`
	DefualtTimeout     int    `yaml:"defualtTimeout"`
	RemoteDir          string `yaml:"remoteDir"`
	LocalDir           string `yaml:"localDir"`
	EnableSftp         int    `yaml:"enableSftp"`
	FetchInterv        int    `yaml:"fetchInterv"`
	BatchKeyScanInterv int    `yaml:"batchKeyScanInterv"`
}

type Hdfs struct {
	InputDir  string `yaml:"inputDir"`
	OutputDir string `yaml:"outputDir"`
}

type BatchCollison struct {
	MaxThreadNum   int `yaml:"maxThreadNum"`
	GroupDataCount int `yaml:"groupDataCount"`
}

type Other struct {
	SupLoadDir      string `yaml:"supLoadDir"`
	FileCleanInterv int    `yaml:"fileCleanInterv"`
	ExidTimeout     string `yaml:"exidTimeout"`
	GuardFlag       string `yaml:"guardFlag"`
	Crp             Crp    `yaml:"crp"`
}

type Crp struct {
	ReqTimeout int `yaml:"reqTimeout"`
}

//func GetCommonSettings() CommonSettings {
//	return commonSettings
//}
func GetCommonSettings() *CommonSettings {
	return commonSettings
}

func SetCommonSettings(settings *CommonSettings) {
	commonSettings = settings
}

func LoadCommonSettings(yamlPath string) error {
	var err error
	var setting Settings
	fmt.Println("yamlPath :", yamlPath)
	if yamlPath == "" {
		setting, err = CreateSettingsFromYAML(defaultDemSettingYaml)
		if err != nil {
			return fmt.Errorf("%s, %s", "loadCommonSettings err", err.Error())
		}
	} else {
		setting, err = CreateSettingsFromYAML(yamlPath)
		if err != nil {
			return fmt.Errorf("%s, %s", "loadCommonSettings err", err.Error())
		}
	}

	unmarshal(&commonSettings, setting)
	fmt.Printf("%+v", commonSettings)
	//yagrusLog.Info("dem setting loading: %+v", CommonSettings)

	return nil
}

func unmarshal(v interface{}, s Settings) {
	unmarshal1(reflect.ValueOf(v).Elem(), s)
}

func unmarshal1(val reflect.Value, s Settings) {
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		fv := val.Field(i)
		ft := typ.Field(i)
		var data []byte
		tag := ft.Tag.Get(defaultTag)
		switch fv.Kind() {
		default:
		case reflect.String:
			val, err := s.GetString(tag)
			if err != nil {
				fmt.Println(tag, ":", err)
			}
			data = []byte(val)
		case reflect.Int:
			val, err := s.GetInt(tag)
			if err != nil {
				fmt.Println(tag, ":", err)
			}
			data = []byte(strconv.Itoa(val))
		case reflect.Slice:
			val, err := s.GetSlice(tag)
			if err != nil {
				fmt.Println(tag, ":", err)
			}
			elemSlice := reflect.MakeSlice(ft.Type, 0, len(val))
			for _, v := range val {
				elemSlice = appendSlice(fv, v, elemSlice)
			}
			fv.Set(elemSlice)
			continue
		case reflect.Bool:
			val, err := s.GetBool(tag)
			if err != nil {
				fmt.Println(tag, ":", err)
			}
			data = []byte(strconv.FormatBool(val))
		case reflect.Struct:
			unmarshal1(fv, s)
			continue
		}
		if tag != "" && len(data) > 0 {
			copyValue(fv, data)
		}
	}
}

func appendSlice(v reflect.Value, val interface{}, elemSlice reflect.Value) reflect.Value {
	switch v.Type().Elem().Kind() {
	case reflect.String:
		elemSlice = reflect.Append(elemSlice, reflect.ValueOf(val.(string)))
	case reflect.Int:
		elemSlice = reflect.Append(elemSlice, reflect.ValueOf(val.(int)))
	}
	return elemSlice
}

func copyValue(dst reflect.Value, src []byte) error {
	dst0 := dst
	if dst.Kind() == reflect.Ptr {
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type()).Elem())
		}
		dst = dst.Elem()
	}

	switch dst.Kind() {
	default:
		return fmt.Errorf("%s", "copyValue error, cannot unmarshal into"+dst0.Type().String())
	case reflect.Invalid:
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		itmp, err := strconv.ParseInt(string(src), 10, dst.Type().Bits())
		if err != nil {
			return err
		}
		dst.SetInt(itmp)
	case reflect.Float32, reflect.Float64:
		ftmp, err := strconv.ParseFloat(string(src), dst.Type().Bits())
		if err != nil {
			return err
		}
		dst.SetFloat(ftmp)
	case reflect.Bool:
		value, err := strconv.ParseBool(strings.TrimSpace(string(src)))
		if err != nil {
			return err
		}
		dst.SetBool(value)
	case reflect.String:
		dst.SetString(string(src))
	case reflect.Slice:
		if len(src) == 0 {
			// non-nil to flag presence
			src = []byte{}
		}
		dst.SetBytes(src)
	}
	return nil
}
