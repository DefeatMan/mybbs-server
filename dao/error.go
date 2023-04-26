package dao

import "errors"

var (
	EQueryFailed     = errors.New("查询失败")
	ECreateFailed    = errors.New("新建失败")
	EUpdateFailed    = errors.New("更新失败")
	EDeleteFailed    = errors.New("删除失败")
	EHadExist        = errors.New("数据已存在")
	ENotExist        = errors.New("数据不存在")
	EPasswordInvalid = errors.New("密码无效或错误")
)
