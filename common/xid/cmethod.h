/**
 *  功能描述: 打开设备
 *
 *  @param phDev [IN|OUT] 设备句柄
 *
 *  @return  成功返回0，非0错误描述参见ErrorCode.h
 */
unsigned long OpenDevice(unsigned long* phDev);


/**
 *  功能描述: 验证PIN码
 *
 *  @param phDev [IN] 设备句柄
 *  @param pin [IN] 所需验证的PIN码
 *
 *  @return  成功返回0，非0错误描述参见ErrorCode.h
 */
unsigned long VerifyPIN(unsigned long hDev, const char* pin);


/**
 *  功能描述: 获取XIDCode
 *
 *  @param phDev			[IN] 设备句柄
 *  @param appId			[IN] 应用Id，一般为20字节
 *  @param appIdLen			[IN] 应用Id长度
 *  @param idType			[IN] 身份信息类型，一般为8字节
 *  @param idTypeLen		[IN] 身份信息类型长度
 *  @param iD				[IN] iD即id，身份信息,包括个人或机构
 *  @param iDLen			[IN] 身份信息长度
 *  @param appXIDCode		[IN|OUT] 生成的appXIDCode
 *  @param appXIDCodeLen	[IN|OUT] 返回appXIDCode的长度
 *
 *  @return  成功返回0，非0错误描述参见ErrorCode.h
 */
unsigned long GetXIDCode(unsigned long hDev,char* appId, int appIdLen, char* idType, int idTypeLen,char* iD, int iDLen,char* appXIDCode, int *appXIDCodeLen);

/**
 *  功能描述: 转换XIDCode
 *
 *  @param phDev							[IN] 设备句柄
 *  @param appIdOne							[IN] 一个应用Id，一般为20字节
 *  @param appIdOneLen						[IN] 一个应用Id的长度
 *  @param xRegcodeBABase64One 		    	[IN] 一个x注册码，Base64编码
 *  @param xRegcodeBABase64OneLen 			[IN] 一个x注册码的长度
 *  @param appXIDCodeOne 					[IN] 一个appXIDCode
 *  @param appXIDCodeOneLen 				[IN] 一个appXIDCode的长度
 *  @param appIdAnother 					[IN] 另一个应用Id
 *  @param appIdAnotherLen 					[IN] 另一个应用Id的长度
 *  @param xRegcodeBABase64Another 			[IN] 另一个x注册码，Base64编码
 *  @param xRegcodeBABase64AnotherLen		[IN] 另一个x注册码的长度
 *  @param appXIDCodeAnother 				[IN|OUT] 转换得到的另一个appXIDCode
 *  @param appXIDCodeAnotherLen 			[IN|OUT] 转换得到的另一个appXIDCode长度
 *
 *  @return  成功返回0，非0错误描述参见ErrorCode.h
 */

unsigned long ChangeXIDCode(unsigned long hDev,char* appIdOne, int appIdOneLen,char*  xIDregcodeBABase64One, int xIDregcodeBABase64OneLen, char* appXIDCodeOne, int appXIDCodeOneLen,char* appIdAnother, int appIdAnotherLen,char* xIDregcodeBABase64Another, int xIDregcodeBABase64AnotherLen,char* appXIDCodeAnother, int *appXIDCodeAnotherLen);

/**
 *  功能描述: 关闭设备
 *
 *  @param phDev [IN] 设备句柄
 *
 *  @return  成功返回0，非0错误描述参见ErrorCode.h
 */
unsigned long CloseDevice(unsigned long hDev);