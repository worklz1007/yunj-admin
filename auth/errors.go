package auth

import "errors"

var (
	ErrInvalidConfig   = errors.New("初始化配置无效")
	ErrUnauthenticated = errors.New("未登录")
	ErrTokenRevoked    = errors.New("登录已失效")
	ErrInternal        = errors.New("内部错误")
)
