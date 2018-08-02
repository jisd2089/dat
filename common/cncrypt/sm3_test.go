//国密算法
//对称加密算法 SM3 测试模块, 包含单元测试和性能测试
//go test -v
//go test -v -bench=".*"

// Sample 1
// Input:"abc"
// Output:66c7f0f4 62eeedd9 d1f2d46b dc10e4e2 4167c487 5cf2f7a2 297da02b 8f4ba8e0

// Sample 2
// Input:"abcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd"
// Outpuf:debe9ff9 2275b8a1 38604889 c18e5a4d 6fdb70e5 387e5765 293dcba3 9c0c5732

package cncrypt

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func showBytes(input []byte) string {
	var buf bytes.Buffer
	for i := 0; i < len(input); i++ {
		s := fmt.Sprintf("%02x", input[i])
		buf.WriteString(s)
		if ((i + 1) % 4) == 0 {
			buf.WriteString(" ")
		}
	}
	return strings.TrimSpace(buf.String())
}

func Test_SM3Bytes(t *testing.T) {
	var res = "66c7f0f4 62eeedd9 d1f2d46b dc10e4e2 4167c487 5cf2f7a2 297da02b 8f4ba8e0"
	var sm3 = new(SM3)
	var test string = "abc"

	input := []byte(test)
	output := sm3.SM3Bytes(input)
	outstr := showBytes(output[:])
	if res != outstr {
		t.Error("SM3Bytes error")
	}
}

func Test_SM3String(t *testing.T) {
	var res = "debe9ff9 2275b8a1 38604889 c18e5a4d 6fdb70e5 387e5765 293dcba3 9c0c5732"
	var sm3 = new(SM3)

	var test string = "abcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd"
	output := sm3.SM3String(test)
	outstr := showBytes(output[:])
	if res != outstr {
		t.Error("SM3String error")
	}
}

func Test_SM3Debug(t *testing.T) {
	var res = "66c7f0f4 62eeedd9 d1f2d46b dc10e4e2 4167c487 5cf2f7a2 297da02b 8f4ba8e0"
	var sm3 = new(SM3)
	var test string = "abc"
	input := []byte(test)
	sm3.Start()
	//sm3.Debug = true //开启调试
	sm3.Update(input)
	output := sm3.Finish()
	outstr := showBytes(output[:])
	if res != outstr {
		t.Error("SM3Debug error")
	}
}

func Benchmark_SM3(b *testing.B) {
	var test string = "abcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd"

	b.ResetTimer()
	var sm3 = new(SM3)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = sm3.SM3String(test)
	}
}

func Benchmark_SM3Hash(b *testing.B) {
	var test string = "abcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd"

	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		input := []byte(test)
		_ = SM3Hash(input)
	}
}
