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
	//"drcs/core"
	"drcs/core/interaction/request"
	"time"
)

func init() {

	go initTransConfig("D:/GoglandProjects/src/drcs/dep/service/trans.properties")

	time.Sleep(1*time.Second)
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

	logger.Info("transInfo", transInfo)
	nodeMemberId := transInfo.Trans.MemberId

	dataAddrs := []*request.NodeAddress{}
	algAddrs := []*request.NodeAddress{}

	for _, val := range transInfo.Trans.Dest {
		switch val.Type {
		case "data":
			dataAddrs = append(dataAddrs, &request.NodeAddress{
				MemberId: nodeMemberId,
				IP:       val.DestIp,
				Host:     val.DestHost,
				URL:      "/api/sup/rec",
				Priority: 0,})
		case "algorithm":
			algAddrs = append(algAddrs, &request.NodeAddress{
				MemberId: nodeMemberId,
				IP:       val.DestIp,
				Host:     val.DestHost,
				URL:      "/api/sup/rec",
				Priority: 0,})
		}
	}

	logger.Info("dataAddrs", dataAddrs)
	logger.Info("algAddrs", algAddrs)


	//b := assetnode.AssetNodeEntity.GetDataBoxByName("demsend")
	//if b == nil {
	//	logger.Error("databox is nil!")
	//}
	//b.SetDataFilePath(s.DataPath)
	//
	//addrs := []*request.NodeAddress{}
	//addrs = append(addrs, &request.NodeAddress{
	//	MemberId: transInfo.Trans.MemberId,
	//	IP:       "127.0.0.1",
	//	Host:     "8989",
	//	URL:      "/api/sup/rec",
	//	Priority: 0,})
	//
	//b.SetNodeAddress(addrs)
	//
	//setDataBoxQueue(b)
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
