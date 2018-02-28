package formatter

import (
	"bytes"
	"testing"
)

func Benchmark_AppendRune(b *testing.B) {
	buffer := NewRuneBuffer()
	for i := 0; i < b.N; i++ {
		buffer.AppendRune('a')
	}
}

func Benchmark_BytesBufferWriteRune(b *testing.B) {
	var buffer bytes.Buffer
	for i := 0; i < b.N; i++ {
		buffer.WriteRune('a')
	}
}

func Benchmark_Append(b *testing.B) {
	buffer := NewRuneBuffer()
	s := "abcdefggggggggggggg"
	runes := []rune(s)
	for i := 0; i < b.N; i++ {
		buffer.Append(runes)
	}
}

func Benchmark_BytesBufferWriteString(b *testing.B) {
	var buffer bytes.Buffer
	s := "abcdefggggggggggggg"
	for i := 0; i < b.N; i++ {
		buffer.WriteString(s)
	}
}

func Benchmark_BytesBufferWrite(b *testing.B) {
	var buffer bytes.Buffer
	s := "abcdefggggggggggggg"
	a := []byte(s)
	for i := 0; i < b.N; i++ {
		buffer.Write(a)
	}
}

func Benchmark_InsertRune(b *testing.B) {
	buffer := NewRuneBuffer()
	s := "abcdefggggggggggggg"
	runes := []rune(s)
	buffer.Append(runes)
	for i := 0; i < b.N; i++ {
		buffer.InsertRune(buffer.Length()-10, 'c')
	}
}

func Benchmark_Insert(b *testing.B) {
	buffer := NewRuneBuffer()
	s := "abcdefggggggggggggg"
	runes := []rune(s)
	buffer.Append(runes)
	for i := 0; i < b.N; i++ {
		buffer.Insert(buffer.Length()-10, runes)
	}
}
