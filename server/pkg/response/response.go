package response

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(data interface{}) Response {
	return Response{Code: 0, Message: "success", Data: data}
}

func SuccessWithCoreError(data interface{}, coreError string) Response {
	if coreError != "" {
		if data == nil {
			data = map[string]interface{}{"core_error": coreError}
		} else if m, ok := data.(map[string]interface{}); ok {
			m["core_error"] = coreError
			data = m
		}
	}
	return Response{Code: 0, Message: "success", Data: data}
}

func SuccessWithCoreResult(data interface{}, msg string) Response {
	if msg == "" {
		msg = "配置已生效"
	}
	return Response{Code: 0, Message: msg, Data: data}
}

// PartialSuccess 表示"数据已持久化但核心未生效"的半成功状态。
//
// 与完全失败的区别在于: 数据已经落库, 用户重试只会造成重复副作用,
// 因此不能谎报成功(否则用户以为已生效), 也不能报完全错误
// (否则用户会重试并造成重复副作用)。
//
// 该状态使用业务码 CodeCoreError(1007), 并在 data 中携带:
//   - core_error:      核心未生效的原因
//   - partial_success: 恒为 true, 客户端据此以警告语义提示
func PartialSuccess(data interface{}, coreError string) Response {
	if data == nil {
		data = map[string]interface{}{}
	}
	m, ok := data.(map[string]interface{})
	if !ok {
		// 非 map 类型的数据包一层, 保留原始数据
		m = map[string]interface{}{"data": data}
	}
	m["core_error"] = coreError
	m["partial_success"] = true
	return Response{Code: CodeCoreError, Message: Msg[CodeCoreError], Data: m}
}

func NeedConfirmRestart(data interface{}, coreError string) Response {
	if data == nil {
		data = map[string]interface{}{}
	}
	if m, ok := data.(map[string]interface{}); ok {
		m["need_confirm"] = true
		m["core_error"] = coreError
		data = m
	}
	return Response{Code: 0, Message: "热重载失败，请确认是否重启内核", Data: data}
}

func Error(code int, message string) Response {
	return Response{Code: code, Message: message}
}

const (
	CodeSuccess       = 0
	CodeBadRequest    = 400
	CodeUnauthorized  = 401
	CodeForbidden     = 403
	CodeNotFound      = 404
	CodeInternalError = 500
)

const (
	CodeInvalidPassword  = 1001
	CodeUserExists       = 1002
	CodeUserNotFound     = 1003
	CodeValidationFailed = 1004
	CodeOperationFailed  = 1005
	CodeConfigError      = 1006
	CodeCoreError        = 1007
	CodeNetworkError     = 1008
)

var Msg = map[int]string{
	CodeSuccess:          "操作成功",
	CodeBadRequest:       "请求参数错误",
	CodeUnauthorized:     "请先登录",
	CodeForbidden:        "没有权限执行此操作",
	CodeNotFound:         "资源不存在",
	CodeInternalError:    "服务器内部错误，请稍后重试",
	CodeInvalidPassword:  "密码错误",
	CodeUserExists:       "用户已存在",
	CodeUserNotFound:     "用户不存在",
	CodeValidationFailed: "数据验证失败",
	CodeOperationFailed:  "操作失败，请稍后重试",
	CodeConfigError:      "配置错误",
	CodeCoreError:        "核心服务错误",
	CodeNetworkError:     "网络连接错误",
}

func MsgWithCode(code int) string {
	if msg, ok := Msg[code]; ok {
		return msg
	}
	return "未知错误"
}

func BadRequest(msg string) Response {
	if msg == "" {
		msg = Msg[CodeBadRequest]
	}
	return Error(CodeBadRequest, msg)
}

func Unauthorized(msg string) Response {
	if msg == "" {
		msg = Msg[CodeUnauthorized]
	}
	return Error(CodeUnauthorized, msg)
}

func NotFound(msg string) Response {
	if msg == "" {
		msg = Msg[CodeNotFound]
	}
	return Error(CodeNotFound, msg)
}

func InternalError(msg string) Response {
	if msg == "" {
		msg = Msg[CodeInternalError]
	}
	return Error(CodeInternalError, msg)
}

func InvalidPassword(msg string) Response {
	if msg == "" {
		msg = Msg[CodeInvalidPassword]
	}
	return Error(CodeInvalidPassword, msg)
}

func OperationFailed(msg string) Response {
	if msg == "" {
		msg = Msg[CodeOperationFailed]
	}
	return Error(CodeOperationFailed, msg)
}

func CoreError(msg string) Response {
	if msg == "" {
		msg = Msg[CodeCoreError]
	}
	return Error(CodeCoreError, msg)
}
