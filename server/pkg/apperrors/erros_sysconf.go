package apperrors

const (
	TypeSystemConfigNotFound     = "SYSTEM_CONFIG_NOT_FOUND"
	TypeSystemConfigListNotFound = "SYSTEM_CONFIG_LIST_NOT_FOUND"
)

var (
	ErrSystemConfigNotFound     = New(TypeSystemConfigNotFound, "系统配置不存在")
	ErrSystemConfigListNotFound = New(TypeSystemConfigListNotFound, "系统配置列表不存在")
)
