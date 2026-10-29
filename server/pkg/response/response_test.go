package response

import (
	"encoding/json"
	"testing"
)

// TestPartialSuccessCode 验证半成功状态使用业务码 1007 而非 0。
func TestPartialSuccessCode(t *testing.T) {
	res := PartialSuccess(nil, "核心重载失败: YAML 语法错误")
	if res.Code != CodeCoreError {
		t.Errorf("期望 code=%d, 实际 %d", CodeCoreError, res.Code)
	}
	if res.Code == CodeSuccess {
		t.Errorf("半成功状态不得返回 code 0, 否则客户端无法区分未生效与成功")
	}
}

// TestPartialSuccessFields 验证 data 中携带 core_error 与 partial_success 标记。
func TestPartialSuccessFields(t *testing.T) {
	res := PartialSuccess(map[string]interface{}{"subscription": "sub-1"}, "核心重载失败")
	data, ok := res.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("期望 data 为 map, 实际 %T", res.Data)
	}
	if data["partial_success"] != true {
		t.Errorf("期望 partial_success=true, 实际 %v", data["partial_success"])
	}
	if data["core_error"] != "核心重载失败" {
		t.Errorf("期望 core_error 被保留, 实际 %v", data["core_error"])
	}
	// 原有数据不得丢失
	if data["subscription"] != "sub-1" {
		t.Errorf("期望原数据被保留, 实际 %v", data["subscription"])
	}
}

// TestPartialSuccessEmptyError 验证 coreError 为空时仍返回 1007 而非降级为 0。
func TestPartialSuccessEmptyError(t *testing.T) {
	res := PartialSuccess(nil, "")
	if res.Code != CodeCoreError {
		t.Errorf("coreError 为空时也应返回 %d, 实际 %d", CodeCoreError, res.Code)
	}
}

// TestPartialSuccessJSONShape 验证序列化后的响应体符合 spec 约定的结构。
func TestPartialSuccessJSONShape(t *testing.T) {
	res := PartialSuccess(map[string]interface{}{"id": 1}, "重载失败")
	raw, err := json.Marshal(res)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}
	if got["code"] != float64(CodeCoreError) {
		t.Errorf("期望 code=%d, 实际 %v", CodeCoreError, got["code"])
	}
	data := got["data"].(map[string]interface{})
	if data["partial_success"] != true {
		t.Errorf("期望 partial_success=true, 实际 %v", data["partial_success"])
	}
	if _, ok := data["core_error"]; !ok {
		t.Errorf("响应体缺少 core_error 字段")
	}
}
