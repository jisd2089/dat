package security

import "bytes"

// StepStatus 步骤状态
type StepStatus string

const (
	// StepStatusFail 失败状态
	StepStatusFail StepStatus = "0"
	// StepStatusSucc 成功状态
	StepStatusSucc StepStatus = "1"
)

// StepInfo 步骤信息
type StepInfo struct {
	No        int        // 步骤号
	MemID     string     // 处理会员
	Status    StepStatus // 状态
	Signature string     // 签名信息
}

// GenerateStepInfo 生成步骤信息的方法
func GenerateStepInfo(prevSignature string, currentNo int, currentMemID string, currentStatus StepStatus) (*StepInfo, error) {
	var signatureSrc bytes.Buffer
	signatureSrc.WriteString(string(currentNo))
	signatureSrc.WriteString(currentMemID)
	signatureSrc.WriteString(string(currentStatus))
	signatureSrc.WriteString(prevSignature)
	signature, err := Signature(signatureSrc.Bytes())
	if err != nil {
		return nil, err
	}
	return &StepInfo{currentNo, currentMemID, currentStatus, signature}, nil
}
