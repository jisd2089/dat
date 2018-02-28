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
	"time"
	//"bytes"
	"fmt"
	//"strings"
	//"encoding/hex"
	"testing"
)

//func Test_SM2PublicKey(t *testing.T) {
//	var priKey = "81EB26E941BB5AF16DF116495F90695272AE2CD63D6C4AE1678418BE48230029"
//	var res = "03160e12897df4edb61dd812feb96748fbd3ccf4ffe26aa6f6db9540af49c94232"
//	var sm2 = new(SM2)
//	sm2.Init()

//	pubKey := sm2.GetPublicKey(priKey)
//	if res != pubKey {
//		t.Error("SM2PublicKey error")
//	}
//}

//func Test_SM2Encrypt(t *testing.T) {
//	var priKey = "81EB26E941BB5AF16DF116495F90695272AE2CD63D6C4AE1678418BE48230029"
//	var pubKey = "03160e12897df4edb61dd812feb96748fbd3ccf4ffe26aa6f6db9540af49c94232"
//	var msg = "encryption standard"
//	var sm2 = new(SM2)
//	sm2.Init()

//	stext := sm2.EncryptString(pubKey, []byte(msg))
//	//fmt.Println("Encrypt=", stext)
//	mm := sm2.DecryptStringOrigin(priKey, stext)
//	//fmt.Println("Decrypt=", string(mm))
//	if msg != string(mm) {
//		t.Error("SM2Encrypt error")
//	}
//}

//签名性能测试
func Benchmark_SM2Sign(b *testing.B) {
	var priKey = "81EB26E941BB5AF16DF116495F90695272AE2CD63D6C4AE1678418BE48230029"
	var msg = "message digest"
	var sm2 = new(SM2)

	b.ResetTimer()
	sm2.Init()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sm2.Reset()
		input := msg + fmt.Sprintf("%08d", i)
		sm2.Sign(priKey, []byte(input))
	}
}

//验签性能测试
func Benchmark_SM2SignVerify(b *testing.B) {
	var pubKey = "03160e12897df4edb61dd812feb96748fbd3ccf4ffe26aa6f6db9540af49c94232"
	var sign = "99923fa53356725672774292876a96f9c3790fa012dc08d8a76cefea8b6668c10fb476bf616d586607097bd0628bf0d289f38a86bd5fe7ba7715f5b7ed38ddba"
	var msg = "message digest"
	var sm2 = new(SM2)

	b.ResetTimer()
	sm2.Init()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = sm2.SignVerify(pubKey, []byte(msg), sign)

	}
}

//加密性能测试
func Benchmark_SM2Encrypt(b *testing.B) {
	var pubKey = "03160e12897df4edb61dd812feb96748fbd3ccf4ffe26aa6f6db9540af49c94232"
	var msg = "encryption standard"
	var sm2 = new(SM2)
	b.ResetTimer()
	sm2.Init()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = sm2.EncryptString(pubKey, []byte(msg))

	}
}

//解密性能测试
func Benchmark_SM2Decrypt(b *testing.B) {
	var priKey = "81EB26E941BB5AF16DF116495F90695272AE2CD63D6C4AE1678418BE48230029"
	var stext = "03f1bddb0f550e7de0fd07e63a1710695090a62f51d7b8d0c815ebd82fa6400b8d7ac761aba27f39f261c4b7d4ea1cb14e8d6427fdf53bfd76a1147263d827da87"
	var sm2 = new(SM2)

	b.ResetTimer()
	sm2.Init()
	b.StartTimer()
	now := time.Now()
	for i := 0; i < b.N; i++ {
		_ = sm2.DecryptStringOrigin(priKey, stext)

	}
	end_time := time.Now()

	var dur_time time.Duration = end_time.Sub(now)
	var elapsed_min float64 = dur_time.Minutes()
	var elapsed_sec float64 = dur_time.Seconds()
	var elapsed_nano int64 = dur_time.Nanoseconds()
	fmt.Printf("cgo show function elasped %f minutes or \nelapsed %f seconds or \nelapsed %d nanoseconds\n",
		elapsed_min, elapsed_sec, elapsed_nano)
	//fmt.Println("time", end-start)
}

func Test_SM2Encrypt(t *testing.T) {
	var priKey = "C288BFCC033C7D308524286DE082EE87853DE62242F8B3D1FCCE89C3306E9617"
	var pubKey = "03A0D865B3C7FC4EA11A78AD15F92F90B11D251DF5D2418CD8F211109A9CCDC481"
	var msg = "1"
	//var stext = "04cb4e5bca8ad5dbbca6778e9ac42ea1525c7248510d3ca30325dd51ccd35834dc3b364b36a796027b4812bb31a002736e284cd0b6ddbf99867f589b0aab10a959ffc100db1b48679bc8021e424b8cae97"
	var sm2 = new(SM2)
	sm2.Init()

	stext := sm2.EncryptString(pubKey, []byte(msg))
	fmt.Println("Encrypt=", stext)
	mm := sm2.DecryptStringOrigin(priKey, stext)
	fmt.Println("Decrypt=", string(mm))
	var v = 0x98fbaa38
	fmt.Printf("1>>17:%08x", v>>17)
	if msg != string(mm) {
		//t.Error("SM2Encrypt error")
	}

}
