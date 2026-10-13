package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/tidwall/sjson"
)

// 智商雷达（check_mode = iq）：向渠道发送需要推理的糖果题，回答中出现独立的 21 判为通过。
// 满血模型通常思考一段时间后答 21；降智渠道常在很短时间内答出 29 或 34。

// monitorIQPrompt 糖果题原文，与 codex-candy-eval 保持一致，便于和外部检测结果对照。
const monitorIQPrompt = `不使用任何外部工具回答以下问题：

在一个黑色的袋子里放有三种口味的糖果，每种糖果有两种不同的形状（圆形和五角星形，不同的形状靠手感可以分辨）。现已知不同口味的糖和不同形状的数量统计如下表。参赛者需要在活动前决定摸出的糖果数目，那么，最少取出多少个糖果才能保证手中同时拥有不同形状的苹果味和桃子味的糖？（同时手中有圆形苹果味匹配五角星桃子味糖果，或者有圆形桃子味匹配五角星苹果味糖果都满足要求）

        苹果味  桃子味  西瓜味
圆形       7      9      8
五角星形   7      6      4
`

// monitorIQAnswerRegex 回答中出现前后都不是数字的 21 即判为正确。
var monitorIQAnswerRegex = regexp.MustCompile(`(^|[^0-9])21([^0-9]|$)`)

// monitorIQHTTPClient 智商雷达专用 client：推理耗时远长于算术探活，单独放宽总超时。
var monitorIQHTTPClient = newSSRFSafeHTTPClient(monitorIQRequestTimeout)

// runIQCheckForModel 对单个模型做一次智商雷达检测。
// 不返回 error：请求失败记 error，答错记 failed，答对记 operational。
func runIQCheckForModel(ctx context.Context, provider, endpoint, apiKey, model string, opts *CheckOptions) *CheckResult {
	res := &CheckResult{
		Model:     model,
		Status:    MonitorStatusError,
		CheckedAt: time.Now(),
	}

	// 1. 发送糖果题；耗时记录整次请求时长，反映模型的思考时间。
	start := time.Now()
	respText, rawBody, statusCode, _, err := callProvider(ctx, provider, endpoint, apiKey, model, monitorIQPrompt, opts)
	latencyMs := int(time.Since(start) / time.Millisecond)
	res.LatencyMs = &latencyMs

	// 2. 网络错误与非 2xx 响应记为请求失败，保留上游错误信息。
	if err != nil {
		res.Message = truncateMessage(sanitizeErrorMessage(err.Error()))
		return res
	}
	if statusCode < 200 || statusCode >= 300 {
		res.Message = truncateMessage(sanitizeErrorMessage(fmt.Sprintf("upstream HTTP %d: %s", statusCode, truncateForErrorBody(rawBody))))
		return res
	}

	// 3. 回答中没有独立的 21 判为未通过，记录回答结尾便于核对模型给出的答案。
	if !monitorIQAnswerRegex.MatchString(respText) {
		res.Status = MonitorStatusFailed
		res.Message = truncateMessage(sanitizeErrorMessage("answer is not 21: " + monitorIQAnswerTail(respText)))
		return res
	}

	// 4. 答对即通过，不按耗时降级：满血模型本来就需要较长思考时间。
	res.Status = MonitorStatusOperational
	return res
}

// monitorIQAnswerTail 取回答末尾片段：结论通常在最后，截断时优先保留结尾。
func monitorIQAnswerTail(text string) string {
	text = strings.Join(strings.Fields(text), " ")
	if text == "" {
		return "(empty)"
	}
	const maxRunes = 150
	runes := []rune(text)
	if len(runes) <= maxRunes {
		return text
	}
	return "..." + string(runes[len(runes)-maxRunes:])
}

// monitorIQInstructions Responses 协议必填的系统指令；不能沿用探活的「简短回答算术题」，否则会诱导模型不推理。
const monitorIQInstructions = "You are a helpful assistant."

// setIQMaxTokens 把默认请求体中的输出上限放宽到 monitorIQMaxTokens，容纳推理过程与完整解答。
func setIQMaxTokens(body []byte, provider, apiMode string) ([]byte, error) {
	// 1. 各协议的输出上限字段名不同，按 adapter 默认 body 的字段改写。
	path := "max_tokens"
	responses := provider == MonitorProviderOpenAI && apiMode == MonitorAPIModeResponses
	switch {
	case provider == MonitorProviderGemini:
		path = "generationConfig.maxOutputTokens"
	case responses:
		path = "max_output_tokens"
	}
	// 2. 改写失败说明默认 body 不是合法 JSON，直接返回错误。
	out, err := sjson.SetBytes(body, path, monitorIQMaxTokens)
	if err != nil {
		return nil, fmt.Errorf("set iq max tokens: %w", err)
	}
	// 3. Responses 协议替换掉探活用的算术题系统指令。
	if responses {
		if out, err = sjson.SetBytes(out, "instructions", monitorIQInstructions); err != nil {
			return nil, fmt.Errorf("set iq instructions: %w", err)
		}
	}
	return out, nil
}
