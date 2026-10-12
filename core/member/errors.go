package member

import "errors"

var (
	ErrInvalidConfig    = errors.New("初始化配置无效")
	ErrEmptyCredentials = errors.New("用户名或密码不能为空")
	ErrBadCredentials   = errors.New("用户名或密码错误")
	ErrDisabled         = errors.New("账户已禁用")
	ErrNotFound         = errors.New("管理员不存在")
	ErrInternal         = errors.New("内部错误")
)
