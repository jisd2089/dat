//SM 算法公用函数

package cncrypt

import (
	"bytes"
	"fmt"
)

// 32-bit integer manipulation macros (big endian)
func GetUlongByte(b []byte, i int) uint32 {
	var b1, b2, b3, b4, n uint32
	b1 = ((uint32)(b[i]) << 24)
	b2 = ((uint32)(b[i+1]) << 16)
	b3 = ((uint32)(b[i+2]) << 8)
	b4 = (uint32)(b[i+3])
	n = b1 | b2 | b3 | b4
	return n
}

func PutUlongByte(b []byte, n uint32, i int) {
	b[i] = (byte)((n) >> 24)
	b[i+1] = (byte)((n) >> 16)
	b[i+2] = (byte)((n) >> 8)
	b[i+3] = (byte)((n))
}

// rotate shift left marco definition
func SHL(x uint32, n uint) uint32 {
	//#define  SHL(x,n) (((x) & 0xFFFFFFFF) << n)
	var rs uint32
	if n > 32 {
		rs = (x & 0xFFFFFFFF) << (n - 32)
	} else {
		rs = (x & 0xFFFFFFFF) << n
	}
	return rs
}

func ROTL(x uint32, n uint) uint32 {
	//#define ROTL(x,n) (SHL((x),n) | ((x) >> (32 - n)))
	var rs, rs2 uint32
	rs1 := SHL(x, n)
	if n > 32 {
		n = (0x7F + 33 - n) % 32
		rs2 = x >> n
	} else {
		rs2 = x >> (32 - n)
	}
	rs = rs1 | rs2
	return rs
}

//根据字节数组生成16进制字符串
func HexString(input []byte) string {
	var buf bytes.Buffer
	for i := 0; i < len(input); i++ {
		s := fmt.Sprintf("%02x", input[i])
		buf.WriteString(s)
	}
	return buf.String()
}

func char2Int(ch byte) int {
	var d1 int
	switch {
	case '0' <= ch && ch <= '9':
		d1 = int(ch - '0')
	case 'a' <= ch && ch <= 'z':
		d1 = int(ch - 'a' + 10)
	case 'A' <= ch && ch <= 'Z':
		d1 = int(ch - 'A' + 10)
	default:
		d1 = 0
	}
	return d1
}

//根据16进制字符串获取字节数组
func HexBytes(str string) []byte {
	var buf bytes.Buffer
	for i := 0; i < len(str)/2; i++ {
		c1 := str[i*2+0]
		c2 := str[i*2+1]
		bt := char2Int(c1)*16 + char2Int(c2)
		buf.WriteByte((byte)(bt))
	}
	return buf.Bytes()
}
