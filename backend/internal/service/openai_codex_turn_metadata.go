package service

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"
	"unicode/utf16"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// rewriteCodexTurnMetadataJSON 仅改写顶层身份字段，保留键序、数值精度和已有转义。
// 原始 Unicode 字符转成 JSON 转义，确保同一份 metadata 可安全用于 HTTP 请求头。
func rewriteCodexTurnMetadataJSON(raw string, rebuildInvalid bool, updates func(map[string]any) map[string]any) string {
	// 1. 校验外部 metadata；按调用方约定保留非法输入或重建最小对象。
	original := raw
	var metadata map[string]any
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.UseNumber()
	if !json.Valid([]byte(raw)) || decoder.Decode(&metadata) != nil || metadata == nil {
		if !rebuildInvalid {
			return original
		}
		raw, metadata = "{}", map[string]any{}
	}
	// 2. 只编码需要改写的字段，不重新序列化整个对象。
	fields := updates(metadata)
	encoded := make(map[string]string, len(fields))
	for name, value := range fields {
		next, err := marshalCodexTurnMetadataValue(value)
		if err != nil {
			return original
		}
		encoded[name] = next
	}
	out := make([]byte, 0, len(raw))
	offset := 0
	seen := make(map[string]bool, len(fields))
	gjson.Parse(raw).ForEach(func(key, value gjson.Result) bool {
		name := key.String()
		next, ok := encoded[name]
		if !ok {
			return true
		}
		seen[name] = true
		// 3. 已相等的字符串保留原有转义写法。
		if text, ok := fields[name].(string); ok && value.Type == gjson.String && value.Str == text {
			return true
		}
		// 4. 覆盖全部同名字段，身份派生仍以最后一个字段值为准。
		out = append(out, raw[offset:value.Index]...)
		out = append(out, next...)
		offset = value.Index + len(value.Raw)
		return true
	})
	out = append(out, raw[offset:]...)
	next := string(out)
	// 5. 仅追加原对象缺少的身份字段，不移动已有字段。
	for _, name := range slices.Sorted(maps.Keys(encoded)) {
		if seen[name] {
			continue
		}
		var err error
		next, err = sjson.SetRaw(next, name, encoded[name])
		if err != nil {
			return original
		}
	}
	// 6. 转义 HTTP 请求头不能携带的字符，保留其他字节及未知字段。
	return escapeCodexTurnMetadataUnicode(next)
}

// marshalCodexTurnMetadataValue 将新值编码成 ASCII JSON，并保留 HTML 字符。
func marshalCodexTurnMetadataValue(value any) (string, error) {
	// 1. 复用出站编码器，再与完整 metadata 共用字符转义规则。
	raw, err := marshalOpenAIUpstreamJSON(value)
	if err != nil {
		return "", err
	}
	return escapeCodexTurnMetadataUnicode(string(raw)), nil
}

func escapeCodexTurnMetadataUnicode(raw string) string {
	// 1. 只替换原始 Unicode 和 DEL；已有 JSON 转义、空白和数值写法均保持不变。
	out := make([]byte, 0, len(raw))
	for _, r := range raw {
		switch {
		case r < 0x7f:
			out = append(out, byte(r))
		case r <= 0xffff:
			out = fmt.Appendf(out, `\u%04x`, r)
		default:
			high, low := utf16.EncodeRune(r)
			out = fmt.Appendf(out, `\u%04x\u%04x`, high, low)
		}
	}
	return string(out)
}
