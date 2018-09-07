package smartsail

/**
    Author: luzequan
    Created: 2018-08-03 14:59:59
*/
type ResponseData struct {
	RespCode    int    `json:"code"` // 返回状态码
	RespMessage string `json:"msg"`  // 返回信息
	RespDetail  string `json:"data"` // 业务详细信息
}

type ResponseDecryptData struct {
	RespCode     int           `json:"code"`         // 返回状态码
	RespMessage  string        `json:"msg"`          // 返回信息
	CodeMessage  string        `json:"code_message"` // 返回信息
	ErrorMessage string        `json:"error_msg"`    // 错误信息
	ReqTime      string        `json:"reqTime"`      // 返回信息
	RespDetail   []*RespDetail `json:"var_detail"`   // 返回信息
}

/**
	Tag标签字段说明
	疑似恶意欺诈 疑似存在欺诈历史
	疑似仿冒包装 疑似用虚假资料包装
	疑似垃圾账户 疑似使用猫池号等工具账户欺诈
	网络恶意行为 在社交、 o2o、社区等疑似有不良的行为
	疑似盗刷 账户被盗刷
	疑似套现
 */
type RespDetail struct {
	Earliest_cons_m        int     `json:"earliest_cons_m"`        // 电商最早消费时间距今月份数
	Consvariety            float64 `json:"consvariety"`            // 电商消费多样性
	Cons_type_amtrate      float64 `json:"cons_type_amtrate"`      // 消费金额最大的类别占总消费金额的比值
	Cons_type_cntrate      float64 `json:"cons_type_cntrate"`      // 消费笔数最大的类别占总消费笔数的比值
	Consume_12m_amt        float64 `json:"consume_12m_amt"`        // 最近12个月电商消费金额
	Consume_12m_cnt        int     `json:"consume_12m_cnt"`        // 最近12个月电商消费笔数
	Type_12m_sum           int     `json:"type_12m_sum"`           // 最近12个月的累计购买产品种类数
	Level_12m_consume      int     `json:"level_12m_consume"`      // 最近12个月的消费档次
	Month_num              int     `json:"month_num"`              // 最近12个月有消费记录的月份数
	Mon_max_record         float64 `json:"mon_max_record"`         // 最近12个月单月最大消费金额
	Latest_record_m        int     `json:"latest_record_m"`        // 最近消费月份距今月份数
	Max_interval_month     int     `json:"max_interval_month"`     // 最大无消费间隔月数
	Cos_stab               float64 `json:"cos_stab"`               // 近12个月的消费稳定度
	Label_24_num           int     `json:"label_24_num"`           // 近两年的消费类别个数
	Close_12m_mean_money   float64 `json:"close_12m_mean_money"`   // 近12个月服装类消费的平均订单金额
	Child_6m_money         float64 `json:"child_6m_money"`         // 近6个月孩童类（儿童，婴儿，产妇类）消费金额
	Child_6m_num           int     `json:"child_6m_num"`           // 近6个月孩童类（儿童，婴儿，产妇类）消费笔数
	Car_6m_money           float64 `json:"car_6m_money"`           // 近6个月汽车类消费金额
	Car_6m_num             int     `json:"car_6m_num"`             // 近6个月汽车类消费笔数
	Diamond_6m_money       float64 `json:"diamond_6m_money"`       // 近6个月装饰（手表,饰品，宝石）类消费金额
	Diamond_6m_num         int     `json:"diamond_6m_num"`         // 近6个月装饰（手表,饰品，宝石）类消费笔数
	Sports_6m_money        float64 `json:"sports_6m_money"`        // 近6个月运动类消费金额
	Sports_6m_num          int     `json:"sports_6m_num"`          // 近6个月运动类消费笔数
	Entertainment_6m_money float64 `json:"entertainment_6m_money"` // 近6个月文娱类（书籍，乐器，音乐，影视）消费金额
	Entertainment_6m_num   int     `json:"entertainment_6m_num"`   // 近6个月文娱类（书籍，乐器，音乐，影视）消费笔数
	Digital_6m_money       float64 `json:"digital_6m_money"`       // 近6个月3C类消费金额
	Digital_6m_num         int     `json:"digital_6m_num"`         // 近6个月3C类消费笔数
	Virtual_6m_money       float64 `json:"virtual_6m_money"`       // 近6个月虚拟商品类消费金额
	Virtual_6m_num         int     `json:"virtual_6m_num"`         // 近6个月虚拟商品类消费笔数
	Hosehold_6m_money      float64 `json:"hosehold_6m_money"`      // 近6个月家用电器类商品总消费金额
	Hosehold_6m_num        int     `json:"hosehold_6m_num"`        // 近6个月家用电器类商品总消费笔数
	Phone                  string  `json:"mobile"`                 // 手机号
}
