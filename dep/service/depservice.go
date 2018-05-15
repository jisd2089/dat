package service

/**
    Author: luzequan
    Created: 2018-05-10 14:59:03
*/
import (
	"sync"
	"gopkg.in/yaml.v2"
	logger "drcs/log"

	"drcs/dep/agollo"
	"drcs/core"
	"drcs/core/interaction/request"
	"time"
	"drcs/settings"
)

func init() {

	go initTransConfig("D:/GoglandProjects/src/drcs/dep/service/trans.properties")

	time.Sleep(1 * time.Second)
	NewDepService().Process()
}

type DepService struct {
	DataPath  string
	JobId     string
	PartnerId string
	lock      sync.RWMutex
}

func NewDepService() *DepService {
	return &DepService{}
}

func (s *DepService) Process() {

	transInfo := GetTransInfo()

	common := settings.GetCommonSettings()

	logger.Info("transInfo", transInfo)
	logger.Info("common setting", common)

	nodeMemberId := transInfo.Trans.MemberId

	fsAddress := &request.FileServerAddress{
		Host:      common.Sftp.Hosts,
		Port:      common.Sftp.Port,
		UserName:  common.Sftp.Username,
		Password:  common.Sftp.Password,
		TimeOut:   common.Sftp.DefualtTimeout,
		LocalDir:  common.Sftp.LocalDir,
		RemoteDir: common.Sftp.RemoteDir,
	}

	dataAddrs := []*Dest{}
	algAddrs := []*Dest{}

	for _, val := range transInfo.Trans.Dest {
		switch val.Type {
		case "data":
			dataAddrs = append(dataAddrs, val)
		case "algorithm":
			algAddrs = append(algAddrs, val)
		}
	}

	logger.Info("dataAddrs", dataAddrs)
	logger.Info("algAddrs", algAddrs)

	runDataBox(dataAddrs, "datasend", nodeMemberId, fsAddress)

	runDataBox(algAddrs, "algorithmsend", nodeMemberId, fsAddress)

}

func initTransConfig(configDir string) {
	newAgollo := agollo.NewAgollo(configDir)
	go newAgollo.Start()

	event := newAgollo.ListenChangeEvent()
	for {
		changeEvent := <-event

		changesCnt := changeEvent.Changes["content"]
		value := changesCnt.NewValue

		transInfo := &TransmissionInfo{}
		err := yaml.Unmarshal([]byte(value), transInfo)
		if err != nil {
		}

		SetTransInfo(transInfo)
	}
}

func runDataBox(addrs []*Dest, boxName string, nodeMemberId string, fsAddress *request.FileServerAddress) {
	for _, v := range addrs {

		b := assetnode.AssetNodeEntity.GetDataBoxByName(boxName)
		if b == nil {
			logger.Error("databox is nil!")
		}
		b.SetDataFilePath(v.DataPath)

		addrs := []*request.NodeAddress{}
		addrs = append(addrs, &request.NodeAddress{
			MemberId: nodeMemberId,
			IP:       v.DestIp,
			Host:     v.DestHost,
			URL:      v.Api,
			Priority: 0,})

		b.SetNodeAddress(addrs)
		b.FileServerAddress = fsAddress

		setDataBoxQueue(b)
	}
}