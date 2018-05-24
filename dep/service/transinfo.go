package service

/**
    Author: luzequan
    Created: 2018-05-10 16:42:16
*/

var (
	transInfo *TransmissionInfo
)

type TransmissionInfo struct {
	Trans Trans `yaml:"transmission"`
}

type Trans struct {
	MemberId string  `yaml:"memberId"`
	Dest     []*Dest `yaml:"dest"`
}

type Dest struct {
	Type          string `yaml:"type"`
	JobId         string `yaml:"jobId"`
	DataPath      string `yaml:"dataPath"`
	DestinationId string `yaml:"destinationId"`
	DestHost      string `yaml:"destHost"`
	DestPort      string `yaml:"destPort"`
	Api           string `yaml:"api"`
	BoxName       string `yaml:"boxName"`
}

func SetTransInfo(transmissionInfo *TransmissionInfo) {
	transInfo = transmissionInfo
}

func GetTransInfo() *TransmissionInfo {
	return transInfo
}
