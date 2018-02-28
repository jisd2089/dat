//国密算法
//椭圆曲线加密算法算法 SM2
//SM2 Standards: http://www.oscca.gov.cn/News/201012/News_1197.htm

package cncrypt

import (
	"bytes"
	"encoding/binary"
	//	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"time"
)

//定义曲线参数
type SM2Params struct {
	P       *big.Int // the order of the underlying field
	N       *big.Int // the order of the base point
	A       *big.Int // the constant of the curve equation
	B       *big.Int // the constant of the curve equation
	Gx, Gy  *big.Int // (x,y) of the base point
	BitSize int      // the size of the underlying field

	// RInverse contains 1/R mod p - the inverse of the Montgomery constant
	// (2**257).
	RInverse *big.Int
}

//签名上下文
type sm2SignContext struct {
	Init         int      //是否已经初始化,0-否,1-是
	KeyID        string   //密钥ID, 即私钥字符串
	PriKey       *big.Int //私钥
	RData        *big.Int //随机数
	RPubX, RPubY *big.Int //RData *(Gx, Gy)
	TData        *big.Int //中间变量, sjs * (1+da)^(-1)
	ZA           [32]byte //ZA
	StrData      string   //TData的字符串
}

//验签上下文
type sm2SignVerifyContext struct {
	Init       int      //是否已经初始化,0-否,1-是
	KeyID      string   //密钥ID, 即公钥字符串
	PubX, PubY *big.Int //公钥
	ZA         [32]byte //ZA,第1步Hash
}

//加密上下文
type sm2CryptContext struct {
	Init         int           //是否已经初始化,0-否,1-是
	RData        *big.Int      //随机数
	RPubX, RPubY *big.Int      //sjs*G
	C1           [33 + 32]byte //C1, 02/03+32个字节
}

//SM2
type SM2 struct {
	*SM2Params                       //曲线参数
	SignCtx    *sm2SignContext       //签名上下文
	VerifyCtx  *sm2SignVerifyContext //验签上下文
	CryptCtx   *sm2CryptContext      //加密上下文
}

var (
	randSeed, _ = new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFF7203DF6B21C6052B53BBF40939D54123", 16)
	signerID    = "1234567812345678" //签名者ID, 固定
)

func sm2Init(params *SM2Params) {
	//根据国家标准
	params.P, _ = new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFF", 16)
	params.N, _ = new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFF7203DF6B21C6052B53BBF40939D54123", 16)
	params.A, _ = new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFC", 16)
	params.B, _ = new(big.Int).SetString("28E9FA9E9D9F5E344D5A9E4BCF6509A7F39789F515AB8F92DDBCBD414D940E93", 16)
	params.Gx, _ = new(big.Int).SetString("32C4AE2C1F1981195F9904466A39C9948FE30BBFF2660BE1715A4589334C74C7", 16)
	params.Gy, _ = new(big.Int).SetString("BC3736A2F4F6779C59BDCEE36B692153D0A9877CC62A474002DF32E52139F0A0", 16)
	params.BitSize = 256
	//其他计算使用的常量
	//	p256RInverse, _ = new(big.Int).SetString("7fffffff00000001fffffffe8000000100000000ffffffff0000000180000000", 16)
}

//两点相乘, P,A曲线定义, (gx,gy)基点, d乘数
//返回相乘后的结果
func sm2PointsMul(P, A, gx, gy, d *big.Int) (*big.Int, *big.Int, int) {
	var i, nLen int
	var zero = 0

	nLen = d.BitLen()
	//fmt.Println("nLen=", nLen)
	var bitArray = make([]byte, nLen)
	for i = 0; i < nLen; i++ {
		bitArray[i] = (byte)(d.Bit(nLen - i - 1))
	}
	//fmt.Println("Bit Array=", bitArray)

	var x1 = big.NewInt(0)
	var y1 = big.NewInt(0)
	var xx1 = big.NewInt(0)
	var yy1 = big.NewInt(0)
	var x2 = big.NewInt(0)
	var y2 = big.NewInt(0)
	var x3 = big.NewInt(0)
	var y3 = big.NewInt(0)
	x1.Set(gx) //creare X1 & X1=px
	y1.Set(gy)
	x2.Set(gx)
	y2.Set(gy)
	xx1.Set(gx)
	yy1.Set(gy)

	//nLen = 8
	for i = 1; i <= nLen-1; i++ {
		//fmt.Println("")
		//fmt.Println("i=", i, ", Bit=", d.Bit(nLen-i-1))
		//mp_copy(&X2, &X1) //X1=X2
		//mp_copy(&Y2, &Y1)
		x1.Set(x2)
		y1.Set(y2)
		//Two_points_add(&X1, &Y1, &X2, &Y2, &X3, &Y3, &A, zero, &P)
		x3, y3 = sm2PointsAdd(P, A, x1, y1, x2, y2, zero)
		//fmt.Println("i=", i, " ,x3=", x3.Text(16), " ,y3=", y3.Text(16))
		//mp_copy(&X3, &X2)
		//mp_copy(&Y3, &Y2)
		x2.Set(x3)
		y2.Set(y3)
		if bitArray[i] == 1 {
			//mp_copy(&XX1, &X1)
			//mp_copy(&YY1, &Y1)
			x1.Set(xx1)
			y1.Set(yy1)
			//Two_points_add(&X1, &Y1, &X2, &Y2, &X3, &Y3, &A, zero, &P)
			x3, y3 = sm2PointsAdd(P, A, x1, y1, x2, y2, zero)
			//mp_copy(&X3, &X2)
			//mp_copy(&Y3, &Y2)
			x2.Set(x3)
			y2.Set(y3)
		}
	}
	if zero > 0 {
		//cout<<"It is Zero_Unit!";
		return x3, y3, 0 //如果Q为零从新产生D
	}
	return x3, y3, 1
}

// fermatInverse calculates the inverse of k in GF(P) using Fermat's method.
// This has better constant-time properties than Euclid's method (implemented
// in math/big.Int.ModInverse) although math/big itself isn't strictly
// constant-time so it's not perfect.
//invmod
func fermatInverse(k, N *big.Int) *big.Int {
	two := big.NewInt(2)
	nMinus2 := new(big.Int).Sub(N, two)
	return new(big.Int).Exp(k, nMinus2, N)
}

//求差值, v=a-b, if(v<0) v=v+p
func diffValue(a, b, p *big.Int) *big.Int {
	zero := big.NewInt(0)
	value := new(big.Int).Sub(a, b)
	if value.Cmp(zero) == -1 {
		//temp1 := x2x1.Add(x2x1, P)
		//x2x1.Set(temp1)
		value = value.Add(value, p)
	}
	return value
}

//mulmod: return=(a*b) mod p
func mulMod(a, b, p *big.Int) *big.Int {
	value := new(big.Int).Mul(a, b)
	value = value.Mod(value, p)
	return value
}

//sqr: return=a*a
func sqr(a *big.Int) *big.Int {
	value := new(big.Int).Mul(a, a)
	return value
}

//submod: return=(a-b) mod p
func subMod(a, b, p *big.Int) *big.Int {
	value := new(big.Int).Sub(a, b)
	value = value.Mod(value, p)
	return value
}

//int Two_points_add(mp_int *x1,mp_int *y1,mp_int *x2,mp_int *y2,mp_int *x3,mp_int *y3,mp_int *a,int zero,mp_int *p)
//两点加, x3=x1+x2; y3=y1+y2, P,A由曲线定义
func sm2PointsAdd(P, A, x1, y1, x2, y2 *big.Int, zero int) (*big.Int, *big.Int) {
	var tempzero = big.NewInt(0)
	var temptwo = big.NewInt(2)
	var temp = new(big.Int)
	var x3 = big.NewInt(0)
	var y3 = big.NewInt(0)

	//fmt.Println("add: zero=", zero)
	//fmt.Println("x1=", x1.Text(16))
	//fmt.Println("y1=", y1.Text(16))
	//fmt.Println("x2=", x2.Text(16))
	//fmt.Println("y2=", y2.Text(16))
	if zero > 0 {
		x3.Set(x1)
		y3.Set(y1)
		zero = 0
		return x3, y3
	}
	//计算(x2,x1), (y2,y1)间的差值
	var x2x1 = diffValue(x2, x1, P)
	//fmt.Println("x2x1=", x2x1.Text(16))
	var y2y1 = diffValue(y2, y1, P)
	//fmt.Println("y2y1=", y2y1.Text(16))

	var km = big.NewInt(0)
	if x2x1.Cmp(tempzero) != 0 {
		//mp_invmod(&x2x1,p,&tempk);
		//fmt.Println("!!!x2x1=", x2x1.Text(16))
		tempk := fermatInverse(x2x1, P)
		//fmt.Println("!!!tempk=", tempk.Text(16))
		km = mulMod(y2y1, tempk, P)
		//fmt.Println("!!!km=", km.Text(16))
	} else {
		if y2y1.Cmp(tempzero) == 0 {
			//mp_mulmod(&temp3,y1,p,&tempy);
			tempy := mulMod(temptwo, y1, P)
			//fmt.Println("tempy=", tempy.Text(16))
			//mp_invmod(&tempy,p,&tempk);
			tempk := fermatInverse(tempy, P)
			//fmt.Println("tempk=", tempk.Text(16))
			//mp_sqr(x1, &temp4);
			temp4 := sqr(x1)
			//mp_mul_d(&temp4, 3, &temp5);
			temp5 := temp4.Mul(temp4, big.NewInt(3))
			//mp_add(&temp5, a, &temp6);
			temp6 := temp5.Add(temp5, A)
			//mp_mulmod(&temp6, &tempk, p, &k);
			km = mulMod(temp6, tempk, P)
		} else {
			zero = 1
			return x3, y3
		}
	}

	//fmt.Println("km=", km.Text(16))
	temp7 := sqr(km)
	//fmt.Println("temp7=", temp7.Text(16))
	//mp_sub(&temp7, x1, &temp8);
	temp8 := temp7.Sub(temp7, x1)
	//fmt.Println("x1=", x1.Text(16))
	//fmt.Println("temp8=", temp8.Text(16))
	//mp_submod(&temp8, x2, p, x3);
	x3 = subMod(temp8, x2, P)
	//fmt.Println("x3=", x3.Text(16))
	//mp_sub(x1, x3, &temp9);
	temp9 := temp.Sub(x1, x3)
	//fmt.Println("temp9=", temp9.Text(16))
	//mp_mul(&temp9, &k, &temp10);
	//fmt.Println("km=", km.Text(16))
	temp10 := temp9.Mul(temp9, km)
	//fmt.Println("temp10=", temp10.Text(16))
	//mp_submod(&temp10, y1, p, y3);
	y3 = subMod(temp10, y1, P)
	//fmt.Println("y3=", y3.Text(16))

	return x3, y3
}

//生成大数随机数
func randBigint() *big.Int {
	var rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	return new(big.Int).Rand(rand, randSeed)
}

//bigInt转换成字符串, 需要在前端补0, 长度固定为64
func bigIntStr64(x *big.Int) string {
	var buf bytes.Buffer
	s := x.Text(16)
	for i := 0; i < (64 - len(s)); i++ {
		buf.WriteString("0")
	}
	buf.WriteString(s)
	return buf.String()
}

//bigInt转换成32个字节，签名补0
func bigIntByte32(x *big.Int) []byte {
	arD := x.Bytes()
	var res [32]byte
	start := (32 - len(arD))
	for i := 0; i < 32; i++ {
		if i < start {
			res[i] = 0
		} else {
			res[i] = arD[i-start]
		}
	}
	return res[:]
}

//根据公钥的X获取Y
// y^2=x^3+ax+b
func getPubKeyY(sm2 *SM2Params, sign int, qx *big.Int) *big.Int {
	//x = decode(pub[1:33], 256)
	//beta = pow(int(x*x*x+A*x+B), int((P+1)//4), int(P))
	//y = (P-beta) if ((beta + from_byte_to_int(pub[0])) % 2) else beta
	//if ((beta + bitcoin.from_byte_to_int(pub[0])) % 2):
	//  y = (sm2.P-beta)
	//else:
	//  y = beta

	//fmt.Println("x=", qx.Text(16))
	//计算 y^2
	//mp_mul(&g_curve.A,&QX,&ax);
	//mp_mul(&QX,&QX,&a2);
	//mp_mul(&a2,&QX,&a3);
	//mp_add(&ax,&g_curve.B,&y2);
	//mp_add(&y2,&a3,&y2); //y^2
	ax := new(big.Int).Mul(sm2.A, qx)
	a2 := new(big.Int).Mul(qx, qx)
	a3 := new(big.Int).Mul(a2, qx)
	y2 := ax.Add(ax, sm2.B)
	y2 = y2.Add(y2, a3) //y^2=x^3+ax+b
	//fmt.Println("y2=", y2.Text(16))
	//mp_add_d(&g_curve.P,1,&p);
	p := new(big.Int).Add(sm2.P, big.NewInt(1))
	//fmt.Println("(p+1)=", p.Text(16))
	//mp_div(&p,&four,&p4,NULL);
	p4 := p.Div(p, big.NewInt(4))
	//fmt.Println("(p+1)/4=", p4.Text(16))
	//mp_exptmod (&y2, &p4, &g_curve.P, &beta); //beta
	beta := y2.Exp(y2, p4, sm2.P)
	//fmt.Println("beta=", beta.Text(16))
	//mp_add_d(&beta,nType,&beta2);
	beta2 := new(big.Int).Add(beta, big.NewInt(int64(sign)))
	//fmt.Println("beta2=", beta2.Text(16))
	//mp_mod_d(&beta2, 2, &res);
	res := beta2.Mod(beta2, big.NewInt(2))
	var qy = new(big.Int)
	if res.Cmp(big.NewInt(0)) == 0 {
		//mp_copy(&beta,QY);
		qy = qy.Set(beta)
	} else {
		//mp_sub(&g_curve.P,&beta,QY);
		qy = qy.Sub(sm2.P, beta)
	}
	return qy
}

//----------------------------- 签名操作函数 -----------------------------
//初始化签名上下文
func signInitContext(sm2 *SM2Params, ctx *sm2SignContext, privateKey string) {
	if ctx.Init == 1 && ctx.KeyID == privateKey {
		return //已经初始化
	}
	ctx.KeyID = privateKey
	//获取私钥, 公钥
	ctx.PriKey, _ = new(big.Int).SetString(privateKey, 16)
	qx, qy, _ := sm2PointsMul(sm2.P, sm2.A, sm2.Gx, sm2.Gy, ctx.PriKey)
	//fmt.Println("QX=", qx.Text(16))
	//fmt.Println("QY=", qy.Text(16))
	//预处理
	ctx.ZA = signPreHash(sm2, qx, qy, []byte(signerID))
	//随机数
	ctx.RData = randBigint()
	ctx.RPubX, ctx.RPubY, _ = sm2PointsMul(sm2.P, sm2.A, sm2.Gx, sm2.Gy, ctx.RData)
	//TData , 中间变量, sjs * (1+da)^(-1)
	tmp1 := new(big.Int).Add(ctx.PriKey, big.NewInt(1)) //1+d
	//mp_invmod(&unit,&g_curve.N,&g_ctxSign.s); //1/(1+dA)
	ctx.TData = fermatInverse(tmp1, sm2.N) //1/(1+dA)
	//mp_mul(&g_ctxSign.sjs,&g_ctxSign.s,&g_ctxSign.s); //k/(1+dA)
	//mp_mod(&g_ctxSign.s,&g_curve.N,&g_ctxSign.s);
	ctx.TData = mulMod(ctx.RData, ctx.TData, sm2.N) //sjs /(1+d)
	ctx.StrData = bigIntStr64(ctx.TData)
	//初始化完成
	ctx.Init = 1
}

//签名预处理--只生成ZA
func signPreHash(sm2 *SM2Params, qx *big.Int, qy *big.Int, signerID []byte) [32]byte {
	//var curveA = "FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFC"
	//var curveB = "28E9FA9E9D9F5E344D5A9E4BCF6509A7F39789F515AB8F92DDBCBD414D940E93"
	//var curveGx = "32C4AE2C1F1981195F9904466A39C9948FE30BBFF2660BE1715A4589334C74C7"
	//var curveGy = "BC3736A2F4F6779C59BDCEE36B692153D0A9877CC62A474002DF32E52139F0A0"

	// 预处理1
	idBitLen := len(signerID) * 8 //位长度
	// 判断当前环境是 big-endian 还是 little-endian。
	// 国密规范中要求把 ENTL(用 2 个字节表示的 ID 的比特长度)
	// 以 big-endian 方式作为预处理 1 输入的前两个字节
	var idLen [2]byte
	binary.BigEndian.PutUint16(idLen[:], (uint16)(idBitLen))
	var buf bytes.Buffer
	buf.Write(idLen[:])
	buf.Write(signerID)
	buf.Write(bigIntByte32(sm2.A)) //[]byte(curveA))
	buf.Write(bigIntByte32(sm2.B))
	buf.Write(bigIntByte32(sm2.Gx))
	buf.Write(bigIntByte32(sm2.Gy))
	buf.Write(bigIntByte32(qx))
	buf.Write(bigIntByte32(qy))
	//input := buf.Bytes()
	//fmt.Println("input len=", len(input))
	//for i := 0; i < len(input); i++ {
	//	fmt.Printf("%02x", input[i])
	//}
	//fmt.Println()
	return SM3Hash(buf.Bytes())
}

func signFast(sm2 *SM2Params, ctx *sm2SignContext, message []byte) string {
	//生成e
	var buf bytes.Buffer
	buf.Write(ctx.ZA[:])
	buf.Write(message)
	hash := SM3Hash(buf.Bytes())
	var e = new(big.Int).SetBytes(hash[:])
	//生成r
	var r = new(big.Int).Add(e, ctx.RPubX)
	r = r.Mod(r, sm2.N)
	//mp_add(&e,&g_ctxSign.X1,&r);
	//mp_mod(&r,&g_curve.N,&r);
	//fmt.Println("e=", e.Text(16))
	//fmt.Println("r=", r.Text(16))
	//返回结果i
	szR := bigIntStr64(r) + ctx.StrData
	return szR
}

//初始化验签上下文
func verifyInitContext(sm2 *SM2Params, ctx *sm2SignVerifyContext, publicKey string) {
	if ctx.Init == 1 && ctx.KeyID == publicKey {
		return //已经初始化
	}
	ctx.KeyID = publicKey
	//还原公钥
	ctx.PubX, _ = new(big.Int).SetString(publicKey[2:], 16)
	sign, _ := strconv.Atoi(publicKey[0:2])
	ctx.PubY = getPubKeyY(sm2, sign, ctx.PubX)
	//fmt.Println("QY=", ctx.PubY.Text(16))
	//计算步骤1 Hash
	ctx.ZA = signPreHash(sm2, ctx.PubX, ctx.PubY, []byte(signerID))
	ctx.Init = 1
}

//协程方式执行两点加，提高效率
func coEccPointMul(P, A, gx, gy, d *big.Int) chan *big.Int {
	out := make(chan *big.Int, 2)
	go func() {
		xt, yt, _ := sm2PointsMul(P, A, gx, gy, d)
		out <- xt
		out <- yt
	}()
	return out
}

//签名验证。正确返回1，错误返回0
func signVerify(sm2 *SM2Params, ctx *sm2SignVerifyContext, message []byte, signValue string) int {
	//读入签名中的r',s'
	var r = new(big.Int)
	var s = new(big.Int)
	r, _ = r.SetString(signValue[0:64], 16)
	s, _ = s.SetString(signValue[64:], 16)
	//fmt.Println("r=", r.Text(16))
	//fmt.Println("s=", s.Text(16))
	//Ecc_points_mul(&x0,&y0,&g_curve.GX,&g_curve.GY,&s,&g_curve.A,&g_curve.P);
	//x0, y0, _ := sm2PointsMul(sm2.P, sm2.A, sm2.Gx, sm2.Gy, s)
	//Ecc_points_mul 比较耗时，两次调用并行执行
	out := coEccPointMul(sm2.P, sm2.A, sm2.Gx, sm2.Gy, s)
	//Ecc_points_mul(&x00,&y00,&QX,&QY,&s,&g_curve.A,&g_curve.P);
	x00, y00, _ := sm2PointsMul(sm2.P, sm2.A, ctx.PubX, ctx.PubY, s)
	//Two_points_add(&x0,&y0,&x00,&y00,&x1,&y1,&g_curve.A,0,&g_curve.P);
	x0 := <-out
	y0 := <-out
	x1, _ := sm2PointsAdd(sm2.P, sm2.A, x0, y0, x00, y00, 0)
	//fmt.Println("x0=", x0.Text(16))
	//fmt.Println("y0=", y0.Text(16))
	//fmt.Println("x00=", x00.Text(16))
	//fmt.Println("y00=", y00.Text(16))
	//fmt.Println("x1=", x1.Text(16))
	//fmt.Println("y1=", y1.Text(16))
	//签名处理
	//步骤2
	var buf bytes.Buffer
	buf.Write(ctx.ZA[:])
	buf.Write(message)
	hash := SM3Hash(buf.Bytes())
	//计算e
	var e = new(big.Int)
	e = e.SetBytes(hash[:])
	//mp_add(&e,&x1,&R);
	r2 := e.Add(e, x1)
	//mp_mod(&R,&g_curve.N,&R);
	r2 = r2.Mod(r2, sm2.N)
	//fmt.Println("r2=", r2.Text(16))
	if r.Cmp(r2) == 0 {
		return (1)
	}
	return (0)
}

//--------------------------- 加解密操作函数 ------------------------------
//初始化加密上下文
func cryptInitContext(sm2 *SM2Params, ctx *sm2CryptContext) {
	if ctx.Init == 1 {
		return
	}
	//随机数
	ctx.RData = randBigint()
	//Ecc_points_mul(&g_ctxEncrypt.X1,&g_ctxEncrypt.Y1,&g_curve.GX,&g_curve.GY,&g_ctxEncrypt.sjs,&g_curve.A,&g_curve.P);
	ctx.RPubX, ctx.RPubY, _ = sm2PointsMul(sm2.P, sm2.A, sm2.Gx, sm2.Gy, ctx.RData)
	//生成C1
	//判断奇偶
	res := new(big.Int).Mod(ctx.RPubY, big.NewInt(2))
	var buf bytes.Buffer
	if res.Cmp(big.NewInt(0)) == 0 {
		buf.WriteByte(0x02)
	} else {
		buf.WriteByte(0x03)
	}
	buf.Write(bigIntByte32(ctx.RPubX))
	//origin
	buf.Write(bigIntByte32(ctx.RPubY))
	copy(ctx.C1[:], buf.Bytes())
	//初始化完成
	ctx.Init = 1
}

//密码派生函数计算
func sm2KDF(str []byte, klen int) []byte {
	//#define HASH_BYTE_LENGTH 32
	//#define HASH_BIT_LENGTH 256
	var hashBitLength = 256
	var ct = 0x00000001
	group_number := ((klen + (hashBitLength - 1)) / hashBitLength)
	var sm3 = new(SM3)
	var bc [4]byte
	var hBuf bytes.Buffer
	for i := 0; i < group_number; i++ {
		//ct复制到字符串最后，big-endian
		bc[0] = (byte)((ct >> 24) & 0xFF)
		bc[1] = (byte)((ct >> 16) & 0xFF)
		bc[2] = (byte)((ct >> 8) & 0xFF)
		bc[3] = (byte)((ct >> 0) & 0xFF)
		sm3.Start()
		sm3.Update(str)
		sm3.Update(bc[:])
		//sm3_finish( &ctx, (BYTE *)&H[i * HASH_BYTE_LENGTH]);
		//memset( &ctx, 0, sizeof( sm3_context ) );
		hash := sm3.Finish()
		hBuf.Write(hash[:])
		ct = ct + 1
	}
	return hBuf.Bytes()
}

//公钥加密，返回密文
func cryptPublicEncrypt(sm2 *SM2Params, ctx *sm2CryptContext, publicKey string, message []byte) []byte {
	//还原公钥
	qx, _ := new(big.Int).SetString(publicKey[2:], 16)
	sign, _ := strconv.Atoi(publicKey[0:2])
	qy := getPubKeyY(sm2, sign, qx)

	//相乘
	//Ecc_points_mul(&x2,&y2,&QX,&QY,&g_ctxEncrypt.sjs,&g_curve.A,&g_curve.P);
	x2, y2, _ := sm2PointsMul(sm2.P, sm2.A, qx, qy, ctx.RData)
	//fmt.Println("x2=", x2.Text(16))
	//fmt.Println("y2=", y2.Text(16))
	//fmt.Println("x1=", ctx.RPubX.Text(16))
	//fmt.Println("y1=", ctx.RPubY.Text(16))
	//计算KDF
	mlen := len(message)
	klen := mlen * 8
	//printf("klen=%d\n",klen);
	var buf bytes.Buffer
	buf.Write(bigIntByte32(x2))
	buf.Write(bigIntByte32(y2))
	kdf := sm2KDF(buf.Bytes(), klen)
	//c2
	for i := 0; i < mlen; i++ {
		kdf[i] = kdf[i] ^ message[i]
	}
	//ShowByte("c2",kdf,mlen);
	//输出密文 C1|C3|C2
	var buf2 bytes.Buffer
	buf2.Write(ctx.C1[:])
	buf2.Write(kdf[0:mlen])
	return buf2.Bytes()
}

//私钥解密
func cryptDecryptByte(sm2 *SM2Params, priKey string, message []byte) []byte {
	//获取私钥
	pk, _ := new(big.Int).SetString(priKey, 16)
	//读入密文中的C1
	x1 := new(big.Int).SetBytes(message[1:33])
	sign := (int)(message[0])
	//还原C1中的y1
	//sign, _ := strconv.Atoi(publicKey[0:2])
	y1 := getPubKeyY(sm2, sign, x1)
	//y1 := new(big.Int).SetBytes(message[33:65])
	//fmt.Println("x1=", x1.Text(16))
	//fmt.Println("y1=", y1.Text(16))

	//解密
	mlen := len(message) - 33
	klen := mlen * 8
	var mw = make([]byte, mlen)
	copy(mw, message[33:])

	//Ecc_points_mul(&x21,&y21,&X1,&Y1,&K,&g_curve.A,&g_curve.P);
	x2, y2, _ := sm2PointsMul(sm2.P, sm2.A, x1, y1, pk)
	//fmt.Println("x2=", x2.Text(16))
	//fmt.Println("y2=", y2.Text(16))

	var buf bytes.Buffer
	buf.Write(bigIntByte32(x2))
	buf.Write(bigIntByte32(y2))
	kdf := sm2KDF(buf.Bytes(), klen)
	//ShowByte("KDF",kdf,mlen);
	//c2
	for i := 0; i < mlen; i++ {
		kdf[i] = kdf[i] ^ mw[i]
	}
	return (kdf[0:mlen])
}

//私钥解密
func cryptDecryptByteOrigin(sm2 *SM2Params, priKey string, message []byte) []byte {
	//获取私钥
	pk, _ := new(big.Int).SetString(priKey, 16)
	//读入密文中的C1
	//fmt.Println(message[1:33])
	x1 := new(big.Int).SetBytes(message[1:33])
	//fmt.Println("x1:" + x1.Text(16))
	y1 := new(big.Int).SetBytes(message[33:65])
	//fmt.Println("y1:" + y1.Text(16))
	//sign := (int)(message[0])
	//还原C1中的y1
	//sign, _ := strconv.Atoi(publicKey[0:2])
	//y1 := getPubKeyY(sm2, sign, x1)
	//y1 := new(big.Int).SetBytes(message[33:65])
	//fmt.Println("x1=", x1.Text(16))
	//fmt.Println("y1=", y1.Text(16))

	//解密
	mlen := len(message) - 33 - 32
	klen := mlen * 8
	var mw = make([]byte, mlen)
	copy(mw, message[33+32:])

	//Ecc_points_mul(&x21,&y21,&X1,&Y1,&K,&g_curve.A,&g_curve.P);
	x2, y2, _ := sm2PointsMul(sm2.P, sm2.A, x1, y1, pk)
	//fmt.Println("x2=", x2.Text(16))
	//fmt.Println("y2=", y2.Text(16))

	var buf bytes.Buffer
	buf.Write(bigIntByte32(x2))
	buf.Write(bigIntByte32(y2))
	kdf := sm2KDF(buf.Bytes(), klen)
	//fmt.Println("kdf:" + HexString(kdf))

	//ShowByte("KDF",kdf,mlen);
	//c2
	for i := 0; i < mlen; i++ {
		kdf[i] = kdf[i] ^ mw[i]
	}
	//fmt.Println(HexString(kdf[0:mlen]))
	return (kdf[0:mlen])
}

//--------------------------- SM2方法 ------------------------------------
//构造方法
func (sm2 *SM2) New() *SM2 {
	if sm2 == nil {
		sm2 = &SM2{}
	}
	sm2.Init()
	return sm2
}

//SM2初始化
func (sm2 *SM2) Init() {
	if sm2.SM2Params == nil {
		sm2.SM2Params = &SM2Params{}
		sm2Init(sm2.SM2Params)
	}
}

//SM2重置
func (sm2 *SM2) Reset() {
	sm2.Init()
	if sm2.SignCtx != nil {
		sm2.SignCtx.Init = 0
	}
	if sm2.VerifyCtx != nil {
		sm2.VerifyCtx.Init = 0
	}
	if sm2.CryptCtx != nil {
		sm2.CryptCtx.Init = 0
	}
}

//获取随机私钥
func (sm2 *SM2) GetRandomPrivate() string {
	return randBigint().Text(16)
}

//根据种子生成私钥，返回私钥字符串，hex编码， seed 种子,hex编码,change 变换值
func (sm2 *SM2) GetPrivateKey(seed string, change int) string {
	//计算种子的SM3
	inputBytes := []byte(seed)
	hash := SM3Hash(inputBytes)
	//变换
	var buf bytes.Buffer
	for i := 0; i < len(hash); i++ {
		s := fmt.Sprintf("%02x", hash[i])
		buf.WriteString(s)
	}
	strn := fmt.Sprintf("%s:%d:%s", seed, change, buf.String())
	//再次Hash
	input2 := []byte(strn)
	hash2 := SM3Hash(input2)
	//生成密码并返回
	var buf2 bytes.Buffer
	for i := 0; i < len(hash2); i++ {
		s := fmt.Sprintf("%02x", hash2[i])
		buf2.WriteString(s)
	}
	return buf2.String()
}

//根据私钥计算公钥，返回公钥字符串，hex编码
func (sm2 *SM2) GetPublicKey(privateKey string) string {
	//私钥
	k := big.NewInt(0)
	k.SetString(privateKey, 16)
	//ShowMPInt("私钥 K 是",K);
	//fmt.Println("私钥 K 是: ", k.Text(16))
	//计算公钥
	qx, qy, _ := sm2PointsMul(sm2.P, sm2.A, sm2.Gx, sm2.Gy, k)

	//返回结果
	strQX := qx.Text(16)
	//判断奇偶
	res := qy.Mod(qy, big.NewInt(2))
	var result string
	if res.Cmp(big.NewInt(0)) == 0 {
		result = "02" + strQX
	} else {
		result = "03" + strQX
	}
	return result
}

//签名，返回签名hex串
func (sm2 *SM2) Sign(privateKey string, message []byte) string {
	//初始化签名上下文
	if sm2.SignCtx == nil {
		sm2.SignCtx = &sm2SignContext{}
		sm2.SignCtx.Init = 0
	}
	signInitContext(sm2.SM2Params, sm2.SignCtx, privateKey)
	//fmt.Println("step1: ", bytesHex(sm2.SignCtx.ZA[:]))
	return signFast(sm2.SM2Params, sm2.SignCtx, message)
}

//签名验证，正确返回1，否则返回0
func (sm2 *SM2) SignVerify(publicKey string, message []byte, sign string) int {
	//初始化验签上下文
	if sm2.VerifyCtx == nil {
		sm2.VerifyCtx = &sm2SignVerifyContext{}
		sm2.VerifyCtx.Init = 0
	}
	verifyInitContext(sm2.SM2Params, sm2.VerifyCtx, publicKey)
	return signVerify(sm2.SM2Params, sm2.VerifyCtx, message, sign)
}

//加密,使用公钥, 返回字节数组
func (sm2 *SM2) Encrypt(publicKey string, message []byte) []byte {
	//初始化加密上下文
	if sm2.CryptCtx == nil {
		sm2.CryptCtx = &sm2CryptContext{}
		sm2.CryptCtx.Init = 0
	}
	cryptInitContext(sm2.SM2Params, sm2.CryptCtx)
	return cryptPublicEncrypt(sm2.SM2Params, sm2.CryptCtx, publicKey, message)
}

//加密,使用公钥, 返回hex字符串
func (sm2 *SM2) EncryptString(publicKey string, message []byte) string {
	//初始化加密上下文
	if sm2.CryptCtx == nil {
		sm2.CryptCtx = &sm2CryptContext{}
		sm2.CryptCtx.Init = 0
	}
	cryptInitContext(sm2.SM2Params, sm2.CryptCtx)
	cr := cryptPublicEncrypt(sm2.SM2Params, sm2.CryptCtx, publicKey, message)
	return HexString(cr)
}

//解密,使用私钥, 返回字节数组
func (sm2 *SM2) Decrypt(privateKey string, message []byte) []byte {
	return cryptDecryptByte(sm2.SM2Params, privateKey, message)
}

//解密,使用私钥, 返回字节数组
func (sm2 *SM2) DecryptString(privateKey string, message string) []byte {
	return cryptDecryptByte(sm2.SM2Params, privateKey, HexBytes(message))
}

//解密,使用私钥, 返回字节数组
func (sm2 *SM2) DecryptStringOrigin(privateKey string, message string) []byte {
	return cryptDecryptByteOrigin(sm2.SM2Params, privateKey, HexBytes(message))
}
