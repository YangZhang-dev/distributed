package registry

// 更新项，包括更新的服务名和URL
type patchEntry struct {
	Name ServiceName
	URL  string
}

// 批量更新
type patch struct {
	Added   []patchEntry
	Removed []patchEntry
}
