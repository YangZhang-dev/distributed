package registry

// Registration 服务注册信息
type Registration struct {
	// 服务名
	ServiceName ServiceName
	// 服务URL
	ServiceURL string
	// 依赖的服务
	RequiredServices []ServiceName
	// 依赖更新时使用的URL
	ServiceUpdateURL string
	// 健康检查URL
	HeartbeatURL string
}
type ServiceName string

const (
	LogService     = ServiceName("Log Service")
	GradingService = ServiceName("Grading Service")
	PortalService  = ServiceName("Portal Service")
)
