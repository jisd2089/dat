/*
cncrypt 实现国家密码行业标准算法, 支持的算法如下:
SM2 SM2椭圆曲线公钥密码算法 (GM/T 0003-2012)
SM3 SM3密码杂凑算法 (GM/T 0004-2012)
SM4 SM4分组密码算法 (GM/T 0002-2012)

1、SM2使用示例:

    var sm2 = new(SM2)  //建立SM2实例
	sm2.Init()  //初始化, 必须先调用
	fmt.Println("private key: ", sm2.GetRandomPrivate()) //获取随机私钥

    SM2支持的方法如下:
    //SM2初始化, 必须先调用
    Init() {
    //重置,SM2使用缓存加快计算速度, 调用Reset()可以清除缓存重新计算参数
    Reset()
    //获取随机私钥
    GetRandomPrivate() string
    //根据种子生成私钥，返回私钥字符串，hex编码， seed 种子,hex编码,change 变换值
    GetPrivateKey(seed string, change int) string
    //根据私钥计算公钥，返回公钥字符串，hex编码
    GetPublicKey(privateKey string) string
    //签名，privateKey:私钥,message:要签名的信息, 返回签名hex串
    Sign(privateKey string, message []byte) string
    //签名验证，正确返回1，否则返回0
    SignVerify(publicKey string, message []byte, sign string) int
    //使用公钥加密, publicKey:公钥，message：要加密的信息， 返回字节数组
    Encrypt(publicKey string, message []byte) []byte
    //使用公钥加密, 返回hex字符串
    EncryptString(publicKey string, message []byte) string
    //使用私钥解密, privateKet:私钥,message:密文，返回字节数组
    Decrypt(privateKey string, message []byte) []byte
    //使用私钥解密字符串, message:密文hex串, 返回字节数组
    DecryptString(privateKey string, message string) []byte


2、SM3使用示例：
    var sm3 = new(SM3)
	var test string = "abc"
	input := []byte(test)
	sm3.Start()  //开始
	//sm3.Debug = true //开启调试
	sm3.Update([]byte("abc")) //处理
	output := sm3.Finish() //结束并获取结果

    也可以使用一个方法计算：
    var sm3 = new(SM3)
    output := sm3.SM3String("abc")

	或者直接计算：
    output := SM3Hash([]byte("abc"))

    SM3支持的方法如下:
    //开始SM3
    Start()
    //更新计算
    Update(input []byte)
    //计算结束, 获取结果
    Finish() [32]byte
    //单步骤直接计算SM3 Hash, 返回计算结果
    SM3Bytes(input []byte) [32]byte
    SM3String(input string) [32]byte


3、SM4使用示例：
    var sm4 = new(SM4)
    sm4.EncryptSetKey(key) //设置加密密码
	output := sm4.Encrypt(input) //加密

    SM4支持的方法如下:
    //设置加密密码, 加密前必须先调用
    EncryptSetKey(key []byte)
    //加密
    Encrypt(input []byte) []byte
    //设置解密密码(128-bit), 解密前必须先调用
    DecryptSetKey(key []byte)
    //解密
    Decrypt(input []byte) []byte

4、其他可使用的工具函数（cncrypt中）
   //根据字节数组生成16进制字符串
   HexString(input []byte) string
   //根据16进制字符串获取字节数组
   HexBytes(str string) []byte

*/
package cncrypt
