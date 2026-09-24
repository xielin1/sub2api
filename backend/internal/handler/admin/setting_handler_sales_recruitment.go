package admin

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
)

// validateSalesRecruitment 只在管理员输入进入系统时规范化和校验招募配置。
func validateSalesRecruitment(value *dto.SalesRecruitment) error {
	// 1. 规范化普通文案；所有内容以纯文本展示，不接收可执行 HTML。
	texts := []*string{&value.Title, &value.Subtitle, &value.Badge, &value.ButtonText, &value.Intro, &value.WechatID, &value.ContactNote, &value.Footer}
	for _, text := range texts {
		*text = strings.TrimSpace(*text)
		if utf8.RuneCountInString(*text) > 500 {
			return fmt.Errorf("销售招募文案不能超过 500 字")
		}
	}
	value.Rules = strings.TrimSpace(value.Rules)
	if utf8.RuneCountInString(value.Rules) > 20000 {
		return fmt.Errorf("合作规则不能超过 20000 字")
	}
	// 2. 启用时必须有可用联系渠道和规则，避免发布无入口或无约定的招募页。
	value.HeroImage = strings.TrimSpace(value.HeroImage)
	value.WechatQRCode = strings.TrimSpace(value.WechatQRCode)
	if value.Enabled && (value.Title == "" || value.ButtonText == "" || value.Rules == "" || (value.WechatID == "" && value.WechatQRCode == "")) {
		return fmt.Errorf("启用销售招募前，请填写标题、按钮文字、合作规则，以及微信号或二维码")
	}
	// 3. 限制展示条数和文案长度；佣金是展示文字，不引入第二套结算规则。
	if value.Tiers == nil {
		value.Tiers = []dto.SalesRecruitmentTier{}
	}
	if value.Benefits == nil {
		value.Benefits = []dto.SalesRecruitmentBenefit{}
	}
	if len(value.Tiers) > 6 || len(value.Benefits) > 6 {
		return fmt.Errorf("佣金档位和合作权益各最多 6 条")
	}
	for i := range value.Tiers {
		tier := &value.Tiers[i]
		tier.Label = strings.TrimSpace(tier.Label)
		tier.Rate = strings.TrimSpace(tier.Rate)
		if tier.Label == "" || tier.Rate == "" || utf8.RuneCountInString(tier.Label) > 100 || utf8.RuneCountInString(tier.Rate) > 30 {
			return fmt.Errorf("佣金档位须填写业绩条件（最多 100 字）和佣金说明（最多 30 字）")
		}
	}
	for i := range value.Benefits {
		benefit := &value.Benefits[i]
		benefit.Title = strings.TrimSpace(benefit.Title)
		benefit.Description = strings.TrimSpace(benefit.Description)
		if benefit.Title == "" || utf8.RuneCountInString(benefit.Title) > 100 || utf8.RuneCountInString(benefit.Description) > 500 {
			return fmt.Errorf("合作权益标题必填且最多 100 字，说明最多 500 字")
		}
	}
	// 4. 上传图片只接收真实的 PNG/JPEG/WebP/GIF，拒绝伪造 MIME、SVG 脚本和非 HTTP 链接。
	for _, image := range []string{value.HeroImage, value.WechatQRCode} {
		if image == "" {
			continue
		}
		if len(image) > 700*1024 {
			return fmt.Errorf("销售招募图片不能超过 500KB，链接不能超过 2048 字符")
		}
		if strings.HasPrefix(image, "data:") {
			header, encoded, ok := strings.Cut(image, ",")
			if !ok || (header != "data:image/png;base64" && header != "data:image/jpeg;base64" && header != "data:image/webp;base64" && header != "data:image/gif;base64") {
				return fmt.Errorf("销售招募图片仅支持 PNG、JPEG、WebP 或 GIF")
			}
			decoded, err := base64.StdEncoding.DecodeString(encoded)
			if err != nil || len(decoded) > 500*1024 || "data:"+http.DetectContentType(decoded)+";base64" != header {
				return fmt.Errorf("销售招募图片内容无效或超过 500KB")
			}
		} else if len(image) > 2048 || config.ValidateAbsoluteHTTPURL(image) != nil {
			return fmt.Errorf("销售招募图片须上传文件或填写绝对 http(s) 链接")
		}
	}
	return nil
}
