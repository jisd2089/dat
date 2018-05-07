package formatter

// Auther xiaolie 20170531

import (
	"drcs/settings"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	ESCAPE_CHAR = rune('%')
	SPACE       = rune(' ')
)

var SPACE8 = []rune{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}

// 解析过程的状态
const (
	LITERAL_STATE = iota
	CONVERTER_STATE
	DOT_STATE
	MIN_STATE
	MAX_STATE
)

// NewPatternConverterFunc 用于创建PatternConverter的函数定义,其中参数为选项值
type NewPatternConverterFunc func(string) PatternConverter

// PatternConverter 转换器接口
type PatternConverter interface {
	Convert(fields map[string]interface{}, buffer *RuneBuffer)
}

// NewPatternConverterFunc注册表
var newPatternConverterFuncRegistry = make(map[string]NewPatternConverterFunc)

// RegisterPatternConverter 注册PatternConverter
func RegisterPatternConverter(id string, newFunc NewPatternConverterFunc) {
	newPatternConverterFuncRegistry[id] = newFunc
}

// PatternLayout 使用"%s %m %d"这样的格式字符串来格式化日志，log4j的PatternLayout的阉割版
type PatternLayout struct {
	pattern string
	// charset   string
	converter PatternConverter
}

// NewPatternLayout 创建PatternLayout对象，适配yagrus.NewFormatterFunc
func NewPatternLayout(se settings.Settings) (logrus.Formatter, error) {
	pattern, err := se.GetString("Pattern")
	if err != nil {
		return nil, err
	}

	// charset, err := se.GetString("Charset")
	// if err != nil {
	// return nil, err
	// }

	converter := parse(pattern)
	return &PatternLayout{pattern: pattern, converter: converter}, nil
}

// Format 实现logrus.Formatter接口
func (f *PatternLayout) Format(entry *logrus.Entry) ([]byte, error) {
	fields := map[string]interface{}(entry.Data)
	fields["msg"] = entry.Message
	defer delete(fields, "msg")
	fields["time"] = entry.Time
	defer delete(fields, "time")
	fields["level"] = entry.Level
	defer delete(fields, "level")

	buffer := NewRuneBufferWithCap(1024)
	f.converter.Convert(fields, buffer)
	return buffer.Bytes(), nil
}

// 解析pattern得到一个PatternConverter
func parse(pattern string) PatternConverter {
	var c rune
	var i int // 下标

	converters := make([]PatternConverter, 0, 16)
	state := LITERAL_STATE
	currentLiteral := NewRuneBuffer()
	formattingInfo := NewDefaultFormattingInfo()
	runes := []rune(pattern)
	for i = 0; i < len(runes); {
		c = runes[i]
		i++
		switch state {
		case LITERAL_STATE:
			if i == len(runes) {
				currentLiteral.AppendRune(c)
				continue
			}

			if c == ESCAPE_CHAR {
				// 判断下一字符
				switch runes[i] {
				case ESCAPE_CHAR: // %%表示字面值%
					currentLiteral.AppendRune(c)
					i++
				default:
					if currentLiteral.Length() != 0 {
						converter := &LiteralPatternConverter{currentLiteral.String()}
						converters = append(converters, converter)
					}
					currentLiteral.Reset()
					currentLiteral.AppendRune(c)
					formattingInfo = NewDefaultFormattingInfo()
					state = CONVERTER_STATE
				}
			} else {
				currentLiteral.AppendRune(c)
			}
		case CONVERTER_STATE:
			currentLiteral.AppendRune(c)
			switch c {
			case '-':
				formattingInfo.leftAlign = true
			case '!':
				formattingInfo.rightTruncate = true
			case '.':
				state = DOT_STATE
			default:
				if isDigit(c) {
					formattingInfo.minLength = int(c - '0')
					state = MIN_STATE
				} else {
					nextIndex, converter := finalizeConverter(i, c, runes, currentLiteral, formattingInfo)
					converters = append(converters, converter)
					// reset
					currentLiteral.Reset()
					formattingInfo = NewDefaultFormattingInfo()
					state = LITERAL_STATE
					i = nextIndex
				}
			}
		case MIN_STATE:
			currentLiteral.AppendRune(c)
			if isDigit(c) {
				formattingInfo.minLength = formattingInfo.minLength*10 + int(c-'0')
			} else if c == '.' {
				state = DOT_STATE
			} else {
				nextIndex, converter := finalizeConverter(i, c, runes, currentLiteral, formattingInfo)
				converters = append(converters, converter)
				// reset
				currentLiteral.Reset()
				formattingInfo = NewDefaultFormattingInfo()
				state = LITERAL_STATE
				i = nextIndex
			}
		case DOT_STATE:
			currentLiteral.AppendRune(c)
			if isDigit(c) {
				formattingInfo.maxLength = int(c - '0')
				state = MAX_STATE
			} else {
				// expect digit, 当字面值处理
				state = LITERAL_STATE
			}
		case MAX_STATE:
			currentLiteral.AppendRune(c)
			if isDigit(c) {
				formattingInfo.maxLength = formattingInfo.maxLength*10 + int(c-'0')
			} else {
				nextIndex, converter := finalizeConverter(i, c, runes, currentLiteral, formattingInfo)
				converters = append(converters, converter)
				// reset
				currentLiteral.Reset()
				formattingInfo = NewDefaultFormattingInfo()
				state = LITERAL_STATE
				i = nextIndex
			}
		} // switch

	}
	if currentLiteral.Length() != 0 {
		converter := &LiteralPatternConverter{currentLiteral.String()}
		converters = append(converters, converter)
	}
	return &CompositeConverter{converters}
}

func finalizeConverter(nextIndex int, lastChar rune, pattern []rune,
	currentLiteral *RuneBuffer, formattingInfo *FormattingInfo) (int, PatternConverter) {
	index := nextIndex
	index, converterID := extractConverterID(index, lastChar, pattern, currentLiteral)
	if converterID == "" {
		return nextIndex, &LiteralPatternConverter{currentLiteral.String()}
	}

	index, options := extractOptions(index, pattern, currentLiteral)
	converter := createConverter(converterID, options)
	if converter == nil {
		return index, &LiteralPatternConverter{currentLiteral.String()}
	}
	return index, &FormattableConverter{converter, formattingInfo}
}

func extractConverterID(nextIndex int, lastChar rune, pattern []rune, currentLiteral *RuneBuffer) (int, string) {
	if !isLetter(lastChar) {
		return nextIndex, ""
	}

	buffer := NewRuneBuffer()
	buffer.AppendRune(lastChar)
	j := nextIndex
	for ; j < len(pattern) && isLetter(pattern[j]); j++ {
		buffer.AppendRune(pattern[j])
		currentLiteral.AppendRune(pattern[j])
	}

	return j, buffer.String()
}

func isLetter(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func extractOptions(nextIndex int, pattern []rune, currentLiteral *RuneBuffer) (int, string) {
	if nextIndex == len(pattern) {
		return nextIndex, ""
	}

	if pattern[nextIndex] != '{' {
		return nextIndex, ""
	}

	buffer := NewRuneBuffer()
	for j := nextIndex + 1; j < len(pattern); j++ {
		if pattern[j] == '}' {
			currentLiteral.AppendRune(pattern[nextIndex]) // {
			currentLiteral.Append(buffer.Runes())
			currentLiteral.AppendRune(pattern[j]) // }
			return j + 1, buffer.String()
		}

		buffer.AppendRune(pattern[j])
	}

	return nextIndex, ""
}

// 从注册表获取newFunc创建PatternConverter实例
func createConverter(converterID string, options string) PatternConverter {
	newFunc := newPatternConverterFuncRegistry[converterID]
	if newFunc == nil {
		return nil
	}
	return newFunc(options)
}

// FormattingInfo 格式化信息
type FormattingInfo struct {
	leftAlign     bool // 是否左对齐
	minLength     int  // 最小长度
	maxLength     int  // 最大长度
	rightTruncate bool // 是否右截取
}

// NewDefaultFormattingInfo 新建一个含有默认值的FormattingInfo实例
func NewDefaultFormattingInfo() *FormattingInfo {
	return &FormattingInfo{
		leftAlign:     false,
		minLength:     0,
		maxLength:     int(math.MaxInt32), // int至少为32位
		rightTruncate: false,
	}
}

// Format 格式化,offset是原始内容在buffer中的偏移量
func (f *FormattingInfo) Format(offset int, buffer *RuneBuffer) {
	rawLength := buffer.Length() - offset
	if rawLength > f.maxLength {
		if f.rightTruncate {
			buffer.Truncate(offset + f.maxLength)
		} else {
			buffer.Delete(offset, rawLength-f.maxLength)
		}
	} else if rawLength < f.minLength {
		if f.leftAlign {
			for i := 0; i < f.minLength-rawLength; i++ {
				buffer.AppendRune(SPACE)
			}
		} else {
			for i := 0; i < (f.minLength-rawLength)/len(SPACE8); i++ {
				buffer.Insert(offset, SPACE8)
			}
			for i := 0; i < (f.minLength-rawLength)%len(SPACE8); i++ {
				buffer.InsertRune(offset, SPACE)
			}
		}
	}
}

// CompositeConverter 转换器组合装饰器
type CompositeConverter struct {
	converters []PatternConverter
}

// Convert 实现PatternConverter接口
func (c *CompositeConverter) Convert(fields map[string]interface{}, buffer *RuneBuffer) {
	for _, converter := range c.converters {
		converter.Convert(fields, buffer)
	}
}

// FormattableConverter 带FormattingInfo的转换器装饰器
type FormattableConverter struct {
	converter      PatternConverter
	formattingInfo *FormattingInfo
}

// Convert 实现PatternConverter接口
func (c *FormattableConverter) Convert(fields map[string]interface{}, buffer *RuneBuffer) {
	offset := buffer.Length()
	c.converter.Convert(fields, buffer)
	c.formattingInfo.Format(offset, buffer)
}

// LiteralPatternConverter 字面值转换器
type LiteralPatternConverter struct {
	literal string
}

// Convert 实现PatternConverter接口
func (c *LiteralPatternConverter) Convert(fields map[string]interface{}, buffer *RuneBuffer) {
	buffer.AppendString(c.literal)
}

// MsgPatternConverter 错误信息 %m %msg
type MsgPatternConverter struct {
}

// 单例
var defaultMsgPatternConverter = &MsgPatternConverter{}

// Convert 实现PatternConverter接口
func (c *MsgPatternConverter) Convert(fields map[string]interface{}, buffer *RuneBuffer) {
	msg := fields["msg"].(string)
	buffer.AppendString(msg)
}

// NewMsgPatternConverter 适配 NewPatternConverterFunc
func NewMsgPatternConverter(options string) PatternConverter {
	return defaultMsgPatternConverter
}

// TimePatternConverter 时间 %d %date
type TimePatternConverter struct {
	layout string // Time.Format 参数
}

// Convert 实现PatternConverter接口
func (c *TimePatternConverter) Convert(fields map[string]interface{}, buffer *RuneBuffer) {
	logTime := fields["time"].(time.Time)
	if c.layout == "" {
		buffer.AppendString(logTime.String())
	} else {
		buffer.AppendString(logTime.Format(c.layout))
	}
}

// NewTimePatternConverter 适配 NewPatternConverterFunc
func NewTimePatternConverter(options string) PatternConverter {
	return &TimePatternConverter{layout: options}
}

// LevelPatternConverter 错误级别 %p %level
type LevelPatternConverter struct {
}

// Convert 实现PatternConverter接口
func (c *LevelPatternConverter) Convert(fields map[string]interface{}, buffer *RuneBuffer) {
	level := fields["level"].(logrus.Level)
	buffer.AppendString(strings.ToUpper(level.String()))
}

var defaultLevelPatternConverter = &LevelPatternConverter{}

// NewLevelPatternConverter 适配 NewPatternConverterFunc
func NewLevelPatternConverter(options string) PatternConverter {
	return defaultLevelPatternConverter
}

// LinePatternConverter 行号 %l %line
type LinePatternConverter struct {
}

// Convert 实现PatternConverter接口
func (c *LinePatternConverter) Convert(fields map[string]interface{}, buffer *RuneBuffer) {
	line := fields["line"].(int)
	buffer.AppendString(strconv.Itoa(line))
}

var defaultLinePatternConverter = &LinePatternConverter{}

// NewLinePatternConverter 适配 NewPatternConverterFunc
func NewLinePatternConverter(options string) PatternConverter {
	return defaultLinePatternConverter
}

// FilePatternConverter 文件名 %F %file
type FilePatternConverter struct {
}

// Convert 实现PatternConverter接口
func (c *FilePatternConverter) Convert(fields map[string]interface{}, buffer *RuneBuffer) {
	file := fields["file"].(string)
	buffer.AppendString(file)
}

var defaultFilePatternConverter = &FilePatternConverter{}

// NewFilePatternConverter 适配 NewPatternConverterFunc
func NewFilePatternConverter(options string) PatternConverter {
	return defaultFilePatternConverter
}

// MethodPatternConverter 文件名 %M %method
type MethodPatternConverter struct {
}

// Convert 实现PatternConverter接口
func (c *MethodPatternConverter) Convert(fields map[string]interface{}, buffer *RuneBuffer) {
	method := fields["func"].(string)
	buffer.AppendString(method)
}

var defaultMethodPatternConverter = &MethodPatternConverter{}

// NewMethodPatternConverter 适配 NewPatternConverterFunc
func NewMethodPatternConverter(options string) PatternConverter {
	return defaultMethodPatternConverter
}

// ErrorPatternConverter 错误 %e %error
type ErrorPatternConverter struct {
}

// Convert 实现PatternConverter接口
func (c *ErrorPatternConverter) Convert(fields map[string]interface{}, buffer *RuneBuffer) {
	err := fields["error"]
	if err == nil {
	    return
	}

	buffer.AppendString(err.(error).Error())
}

var defaultErrorPatternConverter = &ErrorPatternConverter{}

// NewErrorPatternConverter 适配 NewPatternConverterFunc
func NewErrorPatternConverter(options string) PatternConverter {
	return defaultErrorPatternConverter
}

// 初始化
func init() {
	// 注册PatternConverter
	RegisterPatternConverter("m", NewMsgPatternConverter)
	RegisterPatternConverter("msg", NewMsgPatternConverter)

	RegisterPatternConverter("d", NewTimePatternConverter)
	RegisterPatternConverter("date", NewTimePatternConverter)

	RegisterPatternConverter("p", NewLevelPatternConverter)
	RegisterPatternConverter("level", NewLevelPatternConverter)

	RegisterPatternConverter("l", NewLinePatternConverter)
	RegisterPatternConverter("line", NewLinePatternConverter)

	RegisterPatternConverter("F", NewFilePatternConverter)
	RegisterPatternConverter("file", NewFilePatternConverter)

	RegisterPatternConverter("M", NewMethodPatternConverter)
	RegisterPatternConverter("method", NewMethodPatternConverter)

	RegisterPatternConverter("e", NewErrorPatternConverter)
	RegisterPatternConverter("error", NewErrorPatternConverter)
}
