//国密算法
//摘要算法 SM3

// Sample 1
// Input:"abc"
// Output:66c7f0f4 62eeedd9 d1f2d46b dc10e4e2 4167c487 5cf2f7a2 297da02b 8f4ba8e0

// Sample 2
// Input:"abcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd"
// Outpuf:debe9ff9 2275b8a1 38604889 c18e5a4d 6fdb70e5 387e5765 293dcba3 9c0c5732

package cncrypt

import (
	"fmt"
)

type SM3Context struct {
	Total  [2]uint32 // 需要处理的字节数
	State  [8]uint32 // intermediate digest state
	Buffer [64]byte  // data block being processed
	IPad   [64]byte  // HMAC: inner padding
	OPad   [64]byte  //HMAC: outer padding
}

type SM3 struct {
	*SM3Context
	Debug bool //是否输出调试信息
}

var sm3Padding = [64]byte{
	0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

//默认不输出调试信息
var sm3Debug = false

//初始化 context
func sm3Starts(ctx *SM3Context) {
	ctx.Total[0] = 0
	ctx.Total[1] = 0

	ctx.State[0] = 0x7380166F
	ctx.State[1] = 0x4914B2B9
	ctx.State[2] = 0x172442D7
	ctx.State[3] = 0xDA8A0600
	ctx.State[4] = 0xA96F30BC
	ctx.State[5] = 0x163138AA
	ctx.State[6] = 0xE38DEE4D
	ctx.State[7] = 0xB0FB0E4E
	for i := 0; i < 64; i++ {
		ctx.Buffer[i] = 0
		ctx.IPad[i] = 0
		ctx.OPad[i] = 0
	}
}

func debugOutput(title string, data []uint32) {
	if sm3Debug == false {
		return
	}
	var i int
	fmt.Printf("%s", title)
	for i = 0; i < len(data); i++ {
		fmt.Printf("%08x ", data[i])
		if ((i + 1) % 8) == 0 {
			fmt.Printf("\n")
		}
	}
	if i > 0 {
		fmt.Printf("\n")
	}
}

func debugOutputTitle() {
	if sm3Debug == false {
		return
	}
	fmt.Printf("j     A       B        C         D         E        F        G       H\n")
}

func debugOutputLine(j int, A uint32, B uint32, C uint32, D uint32, E uint32, F uint32, G uint32, H uint32) {
	if sm3Debug == false {
		return
	}
	if j < 0 {
		fmt.Printf("   %08x %08x %08x %08x %08x %08x %08x %08x\n", A, B, C, D, E, F, G, H)
	} else {
		fmt.Printf("%02d %08x %08x %08x %08x %08x %08x %08x %08x\n", j, A, B, C, D, E, F, G, H)
	}
}

//#define FF0(x,y,z) ( (x) ^ (y) ^ (z))
func FF0(x uint32, y uint32, z uint32) uint32 {
	return x ^ y ^ z
}

//#define FF1(x,y,z) (((x) & (y)) | ( (x) & (z)) | ( (y) & (z)))
func FF1(x uint32, y uint32, z uint32) uint32 {
	return ((x & y) | (x & z) | (y & z))
}

//#define GG0(x,y,z) ( (x) ^ (y) ^ (z))
func GG0(x uint32, y uint32, z uint32) uint32 {
	return x ^ y ^ z
}

//#define GG1(x,y,z) (((x) & (y)) | ( (~(x)) & (z)) )
func GG1(x uint32, y uint32, z uint32) uint32 {
	return ((x & y) | ((^x) & z))
}

//#define P0(x) ((x) ^  ROTL((x),9) ^ ROTL((x),17))
func P0(x uint32) uint32 {
	return (x ^ ROTL(x, 9) ^ ROTL(x, 17))
}

//#define P1(x) ((x) ^  ROTL((x),15) ^ ROTL((x),23))
func P1(x uint32) uint32 {
	return (x ^ ROTL(x, 15) ^ ROTL(x, 23))
}

func sm3Process(ctx *SM3Context, data []byte) {
	var SS1, SS2, TT1, TT2 uint32
	var W [68]uint32
	var W1 [64]uint32
	var A, B, C, D, E, F, G, H uint32
	var T [64]uint32
	var Temp1, Temp2, Temp3, Temp4, Temp5 uint32
	var j int

	for j = 0; j < 16; j++ {
		T[j] = 0x79CC4519
	}
	for j = 16; j < 64; j++ {
		T[j] = 0x7A879D8A
	}
	for j = 0; j < 16; j++ {
		W[j] = GetUlongByte(data, j*4)
	}
	debugOutput("Message with padding:\n", W[0:16])

	for j = 16; j < 68; j++ {
		//W[j] = P1( W[j-16] ^ W[j-9] ^ ROTL(W[j-3],15)) ^ ROTL(W[j - 13],7 ) ^ W[j-6];
		//Why thd release's result is different with the debug's ?
		//Below is okay. Interesting, Perhaps VC6 has a bug of Optimizaiton.
		Temp1 = W[j-16] ^ W[j-9]
		Temp2 = ROTL(W[j-3], 15)
		Temp3 = Temp1 ^ Temp2
		Temp4 = P1(Temp3)
		Temp5 = ROTL(W[j-13], 7) ^ W[j-6]
		W[j] = Temp4 ^ Temp5
	}
	debugOutput("Expanding message W0-67:\n", W[0:68])

	for j = 0; j < 64; j++ {
		W1[j] = W[j] ^ W[j+4]
	}
	debugOutput("Expanding message W'0-63:\n", W1[0:64])

	A = ctx.State[0]
	B = ctx.State[1]
	C = ctx.State[2]
	D = ctx.State[3]
	E = ctx.State[4]
	F = ctx.State[5]
	G = ctx.State[6]
	H = ctx.State[7]
	debugOutputTitle() //输出标题
	debugOutputLine(-1, A, B, C, D, E, F, G, H)

	for j = 0; j < 16; j++ {
		SS1 = ROTL((ROTL(A, 12) + E + ROTL(T[j], (uint)(j))), 7)
		SS2 = SS1 ^ ROTL(A, 12)
		TT1 = FF0(A, B, C) + D + SS2 + W1[j]
		TT2 = GG0(E, F, G) + H + SS1 + W[j]
		D = C
		C = ROTL(B, 9)
		B = A
		A = TT1
		H = G
		G = ROTL(F, 19)
		F = E
		E = P0(TT2)
		debugOutputLine(j, A, B, C, D, E, F, G, H)
	}

	for j = 16; j < 64; j++ {
		SS1 = ROTL((ROTL(A, 12) + E + ROTL(T[j], (uint)(j))), 7)
		SS2 = SS1 ^ ROTL(A, 12)
		TT1 = FF1(A, B, C) + D + SS2 + W1[j]
		TT2 = GG1(E, F, G) + H + SS1 + W[j]
		D = C
		C = ROTL(B, 9)
		B = A
		A = TT1
		H = G
		G = ROTL(F, 19)
		F = E
		E = P0(TT2)
		debugOutputLine(j, A, B, C, D, E, F, G, H)
	}

	ctx.State[0] ^= A
	ctx.State[1] ^= B
	ctx.State[2] ^= C
	ctx.State[3] ^= D
	ctx.State[4] ^= E
	ctx.State[5] ^= F
	ctx.State[6] ^= G
	ctx.State[7] ^= H
	debugOutputLine(-1, ctx.State[0], ctx.State[1], ctx.State[2],
		ctx.State[3], ctx.State[4], ctx.State[5], ctx.State[6], ctx.State[7])
}

//SM3 process buffer
func sm3Update(ctx *SM3Context, input []byte, ilen int) {
	var fill int
	var left uint32
	var inputStart int = 0
	var bufStart, bufEnd int

	if ilen <= 0 {
		return
	}

	left = ctx.Total[0] & 0x3F
	fill = (int)(64 - left)

	ctx.Total[0] += (uint32)(ilen)
	ctx.Total[0] &= 0xFFFFFFFF

	if ctx.Total[0] < (uint32)(ilen) {
		ctx.Total[1]++
	}

	if left > 0 && ilen >= fill {
		bufStart = (int)(left)
		bufEnd = bufStart + fill
		copy(ctx.Buffer[bufStart:bufEnd], input[inputStart:inputStart+fill])
		//		memcpy( (void *) (ctx->buffer + left),
		//                (void *) input, fill );
		sm3Process(ctx, ctx.Buffer[:])
		inputStart += fill
		ilen -= fill
		left = 0
	}

	for ilen >= 64 {
		sm3Process(ctx, input[inputStart:inputStart+64])
		inputStart += 64
		ilen -= 64
	}

	if ilen > 0 {
		//memcpy( (void *) (ctx->buffer + left),
		//        (void *) input, ilen );
		bufStart = (int)(left)
		bufEnd = bufStart + ilen
		copy(ctx.Buffer[bufStart:bufEnd], input[inputStart:inputStart+ilen])
	}
}

// SM3 final digest
func sm3Finish(ctx *SM3Context, output []byte) {
	var last, padn uint32
	var high, low uint32
	var msglen [8]byte

	high = (ctx.Total[0] >> 29) | (ctx.Total[1] << 3)
	low = (ctx.Total[0] << 3)

	PutUlongByte(msglen[:], high, 0)
	PutUlongByte(msglen[:], low, 4)

	last = ctx.Total[0] & 0x3F
	if last < 56 {
		padn = 56 - last
	} else {
		padn = 120 - last
	}

	sm3Update(ctx, sm3Padding[:], (int)(padn))
	sm3Update(ctx, msglen[:], 8)

	PutUlongByte(output, ctx.State[0], 0)
	PutUlongByte(output, ctx.State[1], 4)
	PutUlongByte(output, ctx.State[2], 8)
	PutUlongByte(output, ctx.State[3], 12)
	PutUlongByte(output, ctx.State[4], 16)
	PutUlongByte(output, ctx.State[5], 20)
	PutUlongByte(output, ctx.State[6], 24)
	PutUlongByte(output, ctx.State[7], 28)
}

//不使用结构, 直接计算
func SM3Hash(input []byte) [32]byte {
	var output [32]byte
	var ctx *SM3Context = &SM3Context{}

	ilen := len(input)
	sm3Starts(ctx)
	sm3Update(ctx, input, ilen)
	sm3Finish(ctx, output[:])
	return output
}

//--------------------------- SM3方法 ------------------------------------
//构造方法
func (sm3 *SM3) New() *SM3 {
	if sm3 == nil {
		sm3 = &SM3{}
	}
	if sm3.SM3Context == nil {
		sm3.SM3Context = &SM3Context{}
	}
	sm3.Debug = false
	return sm3
}

//开始SM3
func (sm3 *SM3) Start() {
	if sm3.SM3Context == nil {
		sm3.SM3Context = &SM3Context{}
		sm3.Debug = false
	}
	sm3Starts(sm3.SM3Context)
}

//更新计算
func (sm3 *SM3) Update(input []byte) {
	if sm3.SM3Context == nil {
		sm3.SM3Context = &SM3Context{}
		sm3Starts(sm3.SM3Context)
	}
	ilen := len(input)
	sm3Debug = sm3.Debug
	sm3Update(sm3.SM3Context, input, ilen)
}

//计算结束, 获取结果
func (sm3 *SM3) Finish() [32]byte {
	var output [32]byte
	sm3Finish(sm3.SM3Context, output[:])
	return output
}

//单步骤直接计算SM3 Hash, 返回计算结果
func (sm3 *SM3) SM3Bytes(input []byte) [32]byte {
	sm3.Start()
	sm3.Update(input)
	return sm3.Finish()
}

func (sm3 *SM3) SM3String(input string) [32]byte {
	inputBytes := []byte(input)
	sm3.Start()
	sm3.Update(inputBytes)
	return sm3.Finish()
}
