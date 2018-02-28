package or

// 缓存标志
type CacheFlag int

const (
	// 不缓存
	NO_CACHE CacheFlag = 0
	// 缓存
	NEED_CACHE CacheFlag = 1
)

type KafkaFlag int

const (
	NO_KAFKA   KafkaFlag = 0
	NEED_KAFKA KafkaFlag = 1
)

// 路由方式
type RouteMethod int

const (
	// 单一路由
	ROUTE_SINGLE RouteMethod = 1
	// 静态路由
	ROUTE_STATIC RouteMethod = 2
	// 动态路由
	ROUTE_DYNAMIC RouteMethod = 3
	// 广播路由
	ROUTE_BROADCAST RouteMethod = 4
)

// 服务停止方式
type StopFlag int

const (
	// 成功即止
	CONTORL_STOP_ON_SUCC StopFlag = 1
	// 全部轮询
	CONTORL_FULL_POLLING StopFlag = 2
)

// 服务同步调用标志
type SyncFlag int

const (
	// 同步
	SYNC SyncFlag = 0
	// 异步
	ASYNC SyncFlag = 1
)

const (
	// redis val SPLIT  分隔符
	SPLIT = "_"
)

const (
	GO_BASE_DIR        = "/dep/go/conf"
	ORDER_ROUTE_PREFIX = "order_route_"
	ORDER_ROUTE_SUFFIX = ".xml"
)

const (
	//需方处理错误
	ErrorPanic = "021999"
)
