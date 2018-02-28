package errors

import "fmt"

// MeanfulError 带错误码和错误信息的自定义错误类
type MeanfulError struct {
	Code  string
	Msg   string
	Cause error
}

// New 创建MeanfulError实例
func New(code string, msg string, cause error) *MeanfulError {
	// TODO 检查code是否为6位数字
	return &MeanfulError{code, msg, cause}
}

// RawNew 创建MeanfulError实例,无源错误
func RawNew(code string, msg string) *MeanfulError {
	// TODO 检查code是否为6位数字
	return &MeanfulError{Code: code, Msg: msg}
}

// Error 实现error接口
func (err *MeanfulError) Error() string {
	if err.Cause == nil {
		return fmt.Sprintf("%s: %s", err.Code, err.Msg)
	}
	return fmt.Sprintf("%s: %s, cause: [%s]", err.Code, err.Msg, err.Cause.Error())
}
