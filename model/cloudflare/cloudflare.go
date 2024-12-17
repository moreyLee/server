package system

// CreateDomain 创建域名的结构体
type CreateDomain struct {
	Name  string `json:"name,omitempty"` // 允许json 字段在缺失时自动忽略
	Start string `json:"jump_start,omitempty"`
}
