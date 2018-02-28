//国密算法
//椭圆曲线加密算法 SM2 测试模块, 包含单元测试和性能测试
//go test -v
//go test -v -bench=".*"

// 获取公钥
// 私钥：81EB26E941BB5AF16DF116495F90695272AE2CD63D6C4AE1678418BE48230029
// 公钥：03160e12897df4edb61dd812feb96748fbd3ccf4ffe26aa6f6db9540af49c94232

// Sample 2
// Input:"abcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd"
// Outpuf:debe9ff9 2275b8a1 38604889 c18e5a4d 6fdb70e5 387e5765 293dcba3 9c0c5732

package cncrypt

import (
	//"time"
	//"bytes"
	"fmt"
	//"strings"
	//"encoding/hex"
	"testing"
)

//验签性能测试
//func Benchmark_SM2SignVerify(b *testing.B) {
//	var pubKey = "03160e12897df4edb61dd812feb96748fbd3ccf4ffe26aa6f6db9540af49c94232"
//	var sign = "99923fa53356725672774292876a96f9c3790fa012dc08d8a76cefea8b6668c10fb476bf616d586607097bd0628bf0d289f38a86bd5fe7ba7715f5b7ed38ddba"
//	var msg = "message digest"
//	var sm2 = new(SM2)

//	b.ResetTimer()
//	sm2.Init()
//	b.StartTimer()
//	for i := 0; i < b.N; i++ {
//		_ = sm2.SignVerify(pubKey, []byte(msg), sign)

//	}
//}
//func testRead(i int) {
//	var pubKey = "03160e12897df4edb61dd812feb96748fbd3ccf4ffe26aa6f6db9540af49c94232"
//	for {
//		keys := pubkey_map[pubKey]
//		fmt.Println(keys, "***", i)
//	}

//}
//func testWrite() {
//	var key = [16]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10}
//	//var pubKey = "03160e12897df4edb61dd812feb96748fbd3ccf4ffe26aa6f6db9540af49c94232"
//	for {
//		keys := [2][]byte{key[:], key[:]}
//		//keys[0] = [16]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10}[:]
//		fmt.Println(keys, "***")
//	}

//}

func Test_Envelop_(t *testing.T) {
	var priKey = "81EB26E941BB5AF16DF116495F90695272AE2CD63D6C4AE1678418BE48230029"
	var pubKey = "03160e12897df4edb61dd812feb96748fbd3ccf4ffe26aa6f6db9540af49c94232"
	var msg = "1"
	Init(priKey)

	envelop := EncrpytEnvelop(msg, pubKey)
	output := envelop.DecryptEnvelop(priKey)

	//fmt.Println(key_map, pubkey_map)
	fmt.Println("Decrypt=", string(output))
	//	for i := 0; i < 64; i++ {
	//		go testRead(i)
	//		go testWrite()
	//	}
	//	for {
	//	}

}
func Test_Envelop_Safe(t *testing.T) {
	var priKey = "81EB26E941BB5AF16DF116495F90695272AE2CD63D6C4AE1678418BE48230029"
	var pubKey = "03160e12897df4edb61dd812feb96748fbd3ccf4ffe26aa6f6db9540af49c94232"
	var msg = "1"
	Init(priKey)

	envelop := EncrpytEnvelopSafe(msg, pubKey)
	output := envelop.DecryptEnvelopSafe(priKey)

	fmt.Println(key_map, pubkey_map)
	fmt.Println("Decrypt=", string(output))
	//	for i := 0; i < 64; i++ {
	//		go testRead(i)
	//		go testWrite()
	//	}
	//	for {
	//	}

}
func Test_Sign(t *testing.T) {
	var priKey = "81EB26E941BB5AF16DF116495F90695272AE2CD63D6C4AE1678418BE48230029"
	var pubKey = "03160e12897df4edb61dd812feb96748fbd3ccf4ffe26aa6f6db9540af49c94232"
	var msg = "1"
	Init(priKey)

	sign := Sign(priKey, []byte(msg))
	res := VerifySign(pubKey, []byte(msg), sign)
	fmt.Println("res:::: ", res)
	fmt.Println(key_map, pubkey_map)
	//fmt.Println("Decrypt=", string(output))
	//	for i := 0; i < 64; i++ {
	//		go testRead(i)
	//		go testWrite()
	//	}
	//	for {
	//	}

}

//加密性能测试
func Benchmark_EncryptEnvelop(b *testing.B) {
	var pubKey = "03160e12897df4edb61dd812feb96748fbd3ccf4ffe26aa6f6db9540af49c94232"
	var priKey = "81EB26E941BB5AF16DF116495F90695272AE2CD63D6C4AE1678418BE48230029"
	var msg = "encryption standard"
	b.ResetTimer()
	Init(priKey)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = EncrpytEnvelop(msg, pubKey)

	}
}
func Benchmark_EncryptEnvelopSafe(b *testing.B) {
	var pubKey = "03160e12897df4edb61dd812feb96748fbd3ccf4ffe26aa6f6db9540af49c94232"
	var priKey = "81EB26E941BB5AF16DF116495F90695272AE2CD63D6C4AE1678418BE48230029"

	var msg = "encryption standard"
	b.ResetTimer()
	Init(priKey)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = EncrpytEnvelopSafe(msg, pubKey)

	}
}
