package formatter

// 一个rune缓冲区实现
// Auther xiaolie 20170531

const (
	defalutCapacity = 16
)

// RuneBuffer rune缓冲区
type RuneBuffer struct {
	runeSlice []rune
}

// NewRuneBuffer 创建RuneBuffer实例, 默认初始容量
func NewRuneBuffer() *RuneBuffer {
	runeSlice := make([]rune, 0, defalutCapacity)
	return &RuneBuffer{runeSlice}
}

// NewRuneBufferWithCap 创建RuneBuffer实例，指定的初始容量
func NewRuneBufferWithCap(capacity int) *RuneBuffer {
	runeSlice := make([]rune, 0, capacity)
	return &RuneBuffer{runeSlice}
}

// AppendRune 追加单个rune
func (r *RuneBuffer) AppendRune(c rune) {
	r.runeSlice = append(r.runeSlice, c)
}

// Append 追加rune切片
func (r *RuneBuffer) Append(runes []rune) {
	r.runeSlice = append(r.runeSlice, runes...)
}

// AppendString 追加字符串
func (r *RuneBuffer) AppendString(str string) {
	r.Append([]rune(str))
}

// InsertRune 在指定位置插入单个rune
func (r *RuneBuffer) InsertRune(index int, c rune) {
	r.AppendRune(c)
	for i := len(r.runeSlice) - 1; i > index; i-- {
		r.runeSlice[i] = r.runeSlice[i-1]
	}
	r.runeSlice[index] = c
}

// Insert 在指定位置插入rune切片
func (r *RuneBuffer) Insert(index int, runes []rune) {
	r.Append(runes)
	length := len(r.runeSlice)
	runesLen := len(runes)
	for i := length - 1; i > index+runesLen-1; i-- {
		r.runeSlice[i] = r.runeSlice[i-runesLen]
	}

	for i, c := range runes {
		r.runeSlice[index+i] = c
	}
}

// InsertString 在指定位置插入字符串
func (r *RuneBuffer) InsertString(index int, str string) {
	r.Insert(index, []rune(str))
}

// Reset 重置，谨慎调用，当心内存泄露
func (r *RuneBuffer) Reset() {
	r.runeSlice = r.runeSlice[:0]
}

// Length 返回当前缓冲区的长度
func (r *RuneBuffer) Length() int {
	return len(r.runeSlice)
}

// Runes 以rune切片形式返回缓冲区内容
func (r *RuneBuffer) Runes() []rune {
	return r.runeSlice
}

// Bytes 以字节切片形式返回缓冲区内容
func (r *RuneBuffer) Bytes() []byte {
	return []byte(r.String())
}

// String 以字符串形式返回缓冲区内容
func (r *RuneBuffer) String() string {
	return string(r.runeSlice)
}

// Slice 返回底层切片起止位置的切片
func (r *RuneBuffer) Slice(start, end int) []rune {
	return r.runeSlice[start:end]
}

// Truncate 截去多余数据，保留指定长度的数据
func (r *RuneBuffer) Truncate(remain int) {
	r.runeSlice = r.runeSlice[:remain]
}

// Delete 从index位置删除长度len的数据
func (r *RuneBuffer) Delete(index, delLen int) {
	sub := r.Slice(index+delLen, r.Length())
	r.Truncate(index)
	r.Append(sub)
}
