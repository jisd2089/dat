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

var commonSettings CommonSettings

type CommonSettings struct {
	DemService DemService
	Redis      Redis
	Conf       Conf
	Routine    Routine
	Kafka      Kafka

	NodeId  string `yaml:"nodeId"`
	MemId   string `yaml:"memId"`
	Version string `yaml:"version"`
	Type    string `yaml:"type"`
	Userkey string `yaml:"userkey"`
	Token   string `yaml:"token"`

	KeysFile           string `yaml:"KeysFile"`
	OrderFile          string `yaml:"OrderFile"`
	OrderFileExpirS    int64  `yaml:"OrderFileExpirS"`
	MemberFile         string `yaml:"MemberFile"`
	MemberFileExpirS   int64  `yaml:"MemberFileExpirS"`
	OrderRouteFilePath string `yaml:"OrderRouteFilePath"`
	DLS_url            string `yaml:"dls_url"`

	SupPort         int    `yaml:"SupPort"`
	LogrusPath      string `yaml:"LogrusPath"`
	BusiLog         BusiLog
	Log             Log
	DMP             DMP
	Sftp            Sftp
	BatchCollison   BatchCollison
	SupLoadDir      string `yaml:"SupLoadDir"`
	FileCleanInterv int    `yaml:"FileCleanInterv"`
	ExIdTimeout     int    `yaml:"ExIdTimeout"`

	BusilogRuleFlow string `yaml:"BusilogRuleFlow"`
}

type Sftp struct {
	Hosts              string `yaml:"sftp.hosts"`
	Port               int    `yaml:"sftp.port"`
	Username           string `yaml:"sftp.username"`
	Password           string `yaml:"sftp.password"`
	DefualtTimeout     int    `yaml:"sftp.defualtTimeout"`
	RemoteDir          string `yaml:"sftp.remoteDir"`
	LocalDir           string `yaml:"sftp.localDir"`
	FetchInterv        int    `yaml:"sftp.fetchInterv"`
	BatchKeyScanInterv int    `yaml:"sftp.batchKeyScanInterv"`
	EnableSftp         int    `yaml:"sftp.enableSftp"`
}

type BatchCollison struct {
	MaxThreadNum   int `yaml:"batchCollison.maxThreadNum"`
	GroupDataCount int `yaml:"batchCollison.groupDataCount"`
}

type DemService struct {
	Protocol   string `yaml:"service.protocol"`
	ListenIP   string `yaml:"service.listenIp"`
	ListenPort int    `yaml:"service.listenPort"`
	ServiceUrl string `yaml:"service.serviceUrl"`
}

//new
type Conf struct {
	CfgDir        string `yaml:"conf.cfgDir"`
	LogDir        string `yaml:"conf.logDir"`
	XmlDir        string `yaml:"conf.xmlDir"`
	XmlReloadTime int    `yaml:"conf.xmlReloadTime"`
}

type Redis struct {
	Addr          string   `yaml:"redis.Addr"`
	Addrs         []string `yaml:"redis.Addrs"`
	DB            int      `yaml:"redis.DB"`
	BatchDB       int      `yaml:"redis.BatchDB"`
	PoolSize      int      `yaml:"redis.PoolSize"`
	BatchPoolSize int      `yaml:"redis.BatchPoolSize"`
	ReadOnly      bool     `yaml:"redis.ReadOnly"`
	BucketSize    int      `yaml:"redis.BucketSize"`

	ReadTimeout  int `yaml:"redis.ReadTimeout"`
	WriteTimeout int `yaml:"redis.WriteTimeout"`
	DialTimeout  int `yaml:"redis.DialTimeout"`

	ClusterOpen     bool     `yaml:"redis.ClusterOpen"`
	ClusterAddrs    []string `yaml:"redis.ClusterAddrs"`
	ClusterPoolSize int      `yaml:"redis.ClusterPoolSize"`
}

type Kafka struct {
	Brokers      []string `yaml:"kafka.brokers"`
	Topic        string   `yaml:"kafka.topic"`
	PartitionNum int      `yaml:"kafka.partitionNum"`
}

type Routine struct {
	MemIds   []string `yaml:"routine.memId"`
	Capacity []int    `yaml:"routine.capacity"`
	MCMap    map[string]int
}

type BusiLog struct {
	MQPath string `yaml:"busilog.MQPath"`
}

type DMP struct {
	URL     string `yaml:"DMP.URL"`
	Timeout int    `yaml:"DMP.Timeout"`
}

type Log struct {
	ConfigPath string `yaml:"Log.ConfigPath"`
}

//func GetCommonSettings() CommonSettings {
//	return commonSettings
//}
func GetCommomSettings() CommonSettings {
	return commonSettings
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
