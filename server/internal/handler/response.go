package handler

import (
	"clash-server/pkg/response"
)

// responseWithCore 根据核心应用结果返回对应的响应体。
//
// 核心应用成功时返回普通成功响应 (code 0)。
// 核心应用失败时返回半成功响应 (业务码 1007 + partial_success 标记),
// 明确告知调用方"数据已持久化但核心未生效", 避免调用方误以为操作完全成功。
//
// 之所以不用错误语义: 数据已经落库, 用户重试只会造成重复副作用;
// 之所以不用成功语义: 那会掩盖核心未生效的事实。
func responseWithCore(data interface{}, coreErr string) response.Response {
	if coreErr != "" {
		return response.PartialSuccess(data, coreErr)
	}
	return response.Success(data)
}
