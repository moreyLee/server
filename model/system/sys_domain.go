package system

// CreateDomain 创建域名的结构体
type CreateDomain struct {
	Name  string `json:"name"`
	Start string `json:"jump_start"`
}
