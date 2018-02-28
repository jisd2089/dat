package formatter

import "testing"

func TestAppendRune(t *testing.T) {
	buffer := NewRuneBuffer()
	buffer.AppendRune('a')
	if buffer.String() != "a" {
		t.Fatalf("expect %s, got %s", "a", buffer.String())
	}
	if buffer.Length() != 1 {
		t.Fatalf("expect length %d, got %d", 1, buffer.Length())
	}

	buffer.AppendRune('b')
	if buffer.String() != "ab" {
		t.Fatalf("expect %s, got %s", "ab", buffer.String())
	}
	if buffer.Length() != 2 {
		t.Fatalf("expect length %d, got %d", 2, buffer.Length())
	}

	buffer.AppendRune('c')
	if buffer.String() != "abc" {
		t.Fatalf("expect %s, got %s", "abc", buffer.String())
	}
	if buffer.Length() != 3 {
		t.Fatalf("expect length %d, got %d", 3, buffer.Length())
	}
}

func TestAppend(t *testing.T) {
	buffer := NewRuneBuffer()
	buffer.Append([]rune{'a', 'b', 'c'})
	if buffer.String() != "abc" {
		t.Fatalf("expect %s, got %s", "abc", buffer.String())
	}
	if buffer.Length() != 3 {
		t.Fatalf("expect length %d, got %d", 3, buffer.Length())
	}
	buffer.Append([]rune{'d', 'e', 'f'})
	if buffer.String() != "abcdef" {
		t.Fatalf("expect %s, got %s", "abcdef", buffer.String())
	}
	if buffer.Length() != 6 {
		t.Fatalf("expect length %d, got %d", 6, buffer.Length())
	}
}

func TestAppendString(t *testing.T) {
	buffer := NewRuneBuffer()
	str := "你好，世界"
	buffer.AppendString(str)
	if buffer.String() != str {
		t.Fatalf("expect %s, got %s", str, buffer.String())
	}
	if buffer.Length() != 5 {
		t.Fatalf("expect length %d, got %d", 5, buffer.Length())
	}
}

func TestInsertRune(t *testing.T) {
	buffer := NewRuneBuffer()
	buffer.InsertRune(0, 'd')
	if buffer.String() != "d" {
		t.Fatalf("expect %s, got %s", "d", buffer.String())
	}
	if buffer.Length() != 1 {
		t.Fatalf("expect length %d, got %d", 1, buffer.Length())
	}

	buffer.InsertRune(1, 'e')
	if buffer.String() != "de" {
		t.Fatalf("expect %s, got %s", "de", buffer.String())
	}
	if buffer.Length() != 2 {
		t.Fatalf("expect length %d, got %d", 2, buffer.Length())
	}
}

func TestInsert(t *testing.T) {
	buffer := NewRuneBuffer()
	buffer.Insert(0, []rune{'d', 'e', 'f'})
	if buffer.String() != "def" {
		t.Fatalf("expect %s, got %s", "def", buffer.String())
	}
	if buffer.Length() != 3 {
		t.Fatalf("expect length %d, got %d", 3, buffer.Length())
	}

	buffer.Insert(1, []rune{'g', 'h', 'i'})
	if buffer.String() != "dghief" {
		t.Fatalf("expect %s, got %s", "dghief", buffer.String())
	}
	if buffer.Length() != 6 {
		t.Fatalf("expect length %d, got %d", 6, buffer.Length())
	}
}

func TestInsertString(t *testing.T) {
	buffer := NewRuneBuffer()
	buffer.InsertString(0, "字符串")
	if buffer.String() != "字符串" {
		t.Fatalf("expect %s, got %s", "字符串", buffer.String())
	}
	if buffer.Length() != 3 {
		t.Fatalf("expect length %d, got %d", 3, buffer.Length())
	}
	// buffer.InsertString(1, "a")
	// if buffer.String() != "字a符串" {
	// t.Fatalf("expect %s, got %s", "字a符串", buffer.String())
	// }
	// if buffer.Length() != 4 {
	// t.Fatalf("expect length %d, got %d", 4, buffer.Length())
	// }
}

func TestReset(t *testing.T) {
	buffer := NewRuneBuffer()
	buffer.Append([]rune{'a', 'b', 'c'})
	if buffer.Length() != 3 {
		t.Fatalf("expect length %d, got %d", 3, buffer.Length())
	}
	buffer.Reset()
	if buffer.Length() != 0 {
		t.Fatalf("expect length %d, got %d", 0, buffer.Length())
	}
}

func TestRunes(t *testing.T) {
	buffer := NewRuneBuffer()
	buffer.AppendString("abc")
	runes := buffer.Runes()

	if len(runes) != 3 {
		t.Fatalf("expect length %d, got %d", 3, len(runes))
	}

	if runes[0] != 'a' {
		t.Fatalf("expect %c, got %c", "a", runes[0])
	}
	if runes[2] != 'c' {
		t.Fatalf("expect %c, got %c", "c", runes[2])
	}
}

func TestBytes(t *testing.T) {
	buffer := NewRuneBuffer()
	buffer.AppendString("abc")
	bytes := buffer.Bytes()

	if len(bytes) != 3 {
		t.Fatalf("expect length %d, got %d", 3, len(bytes))
	}
	if bytes[0] != byte('a') {
		t.Fatalf("expect %c, got %c", "a", bytes[0])
	}
	if bytes[2] != byte('c') {
		t.Fatalf("expect %c, got %c", "c", bytes[2])
	}

}

func TestSlice(t *testing.T) {
	buffer := NewRuneBuffer()
	buffer.Append([]rune{'a', 'b', 'c'})

	sub1 := buffer.Slice(1, 2)
	if string(sub1) != "b" {
		t.Fatalf("expect %s, got %s", "b", string(sub1))
	}

	sub2 := buffer.Slice(1, 3)
	if string(sub2) != "bc" {
		t.Fatalf("expect %s, got %s", "bc", string(sub2))
	}
}

func TestTruncate(t *testing.T) {
	buffer := NewRuneBuffer()
	buffer.Append([]rune{'a', 'b', 'c'})

	buffer.Truncate(2)
	if buffer.String() != "ab" {
		t.Fatalf("expect %s, got %s", "ab", buffer.String())
	}

	if buffer.Length() != 2 {
		t.Fatalf("expect length %d, got %d", 2, buffer.Length())
	}
}

func TestDelete(t *testing.T) {
	buffer := NewRuneBuffer()
	buffer.Append([]rune{'a', 'b', 'c'})

	buffer.Delete(1, 1)
	if buffer.String() != "ac" {
		t.Fatalf("expect %s, got %s", "ac", buffer.String())
	}

	if buffer.Length() != 2 {
		t.Fatalf("expect length %d, got %d", 2, buffer.Length())
	}
}
