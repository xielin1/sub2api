package service

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"log/slog"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// accountErrorNotifyInterval 账号异常扫描周期。
	accountErrorNotifyInterval = time.Minute
	// accountErrorNotifyNotifiedKey 保存各类异常中已通知且仍未恢复的账号 ID，重启后不重复通知。
	accountErrorNotifyNotifiedKey = "account_error_notify_notified_ids"
	// accountErrorNotifyMaxMessageLen 邮件中单个账号异常信息的最大长度。
	accountErrorNotifyMaxMessageLen = 300
	// accountErrorNotify5xxWindow 统计上游 5xx 的时间窗口。
	accountErrorNotify5xxWindow = 5 * time.Minute
	// accountErrorNotify5xxThreshold 窗口内同一账号上游 5xx 达到该次数才通知。
	accountErrorNotify5xxThreshold = 5
	// accountErrorNotify5xxPageSize 单次拉取上游错误日志的最大条数。
	accountErrorNotify5xxPageSize = 500
)

// 异常类型：写入已通知集合的键，同时决定邮件里的类型标签。
const (
	accountErrorNotifyKindError       = "error"
	accountErrorNotifyKindUnavailable = "unavailable"
	accountErrorNotifyKindUpstream5xx = "upstream_5xx"
)

// accountErrorNotifyKindLabels 邮件中各异常类型的标签。
var accountErrorNotifyKindLabels = map[string]string{
	accountErrorNotifyKindError:       "错误停用 / Error",
	accountErrorNotifyKindUnavailable: "暂时不可用 / Temporarily unavailable",
	accountErrorNotifyKindUpstream5xx: "上游 5xx 频繁 / Frequent upstream 5xx",
}

// accountErrorNotifyItem 一个异常账号及其说明。
type accountErrorNotifyItem struct {
	Kind     string
	ID       int64
	Name     string
	Platform string
	Detail   string
}

// AccountErrorNotifyService 定时扫描异常的上游账号，并邮件通知配置的邮箱。
// 覆盖三类异常：错误停用（status=error）、暂时不可用（限流/过载/临时停调度）、上游 5xx 频繁。
type AccountErrorNotifyService struct {
	accountRepo              AccountRepository
	opsRepo                  OpsRepository
	settingRepo              SettingRepository
	emailService             *EmailService
	notificationEmailService *NotificationEmailService
	stopCh                   chan struct{}
	stopOnce                 sync.Once
	wg                       sync.WaitGroup
}

// NewAccountErrorNotifyService 创建账号异常通知服务。
func NewAccountErrorNotifyService(accountRepo AccountRepository, opsRepo OpsRepository, settingRepo SettingRepository, emailService *EmailService, notificationEmailService *NotificationEmailService) *AccountErrorNotifyService {
	return &AccountErrorNotifyService{
		accountRepo:              accountRepo,
		opsRepo:                  opsRepo,
		settingRepo:              settingRepo,
		emailService:             emailService,
		notificationEmailService: notificationEmailService,
		stopCh:                   make(chan struct{}),
	}
}

// Start 启动后台扫描。
func (s *AccountErrorNotifyService) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		// 1. 启动后立即扫描一次，之后按固定周期扫描。
		ticker := time.NewTicker(accountErrorNotifyInterval)
		defer ticker.Stop()
		s.runOnce()
		for {
			select {
			case <-ticker.C:
				s.runOnce()
			case <-s.stopCh:
				return
			}
		}
	}()
}

// Stop 停止后台扫描并等待当前一轮结束。
func (s *AccountErrorNotifyService) Stop() {
	s.stopOnce.Do(func() {
		close(s.stopCh)
	})
	s.wg.Wait()
}

func (s *AccountErrorNotifyService) runOnce() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// 1. 未开启或没有可用收件人时不扫描。
	enabled, err := s.settingRepo.GetValue(ctx, SettingKeyAccountErrorNotifyEnabled)
	if err != nil || enabled != "true" {
		return
	}
	rawEmails, err := s.settingRepo.GetValue(ctx, SettingKeyAccountErrorNotifyEmails)
	if err != nil {
		return
	}
	recipients := filterVerifiedEmails(ParseNotifyEmails(rawEmails))
	if len(recipients) == 0 {
		return
	}

	// 2. 收集三类当前异常的账号；任一类查询失败则本轮放弃，下一轮重试。
	now := time.Now()
	errorItems, err := s.collectErrorAccounts(ctx)
	if err != nil {
		slog.Error("account error notify: list error accounts failed", "error", err)
		return
	}
	unavailableItems, err := s.collectUnavailableAccounts(ctx, now)
	if err != nil {
		slog.Error("account error notify: list active accounts failed", "error", err)
		return
	}
	upstreamItems, err := s.collectUpstream5xxAccounts(ctx, now)
	if err != nil {
		slog.Error("account error notify: list upstream errors failed", "error", err)
		return
	}
	current := map[string][]accountErrorNotifyItem{
		accountErrorNotifyKindError:       errorItems,
		accountErrorNotifyKindUnavailable: unavailableItems,
		accountErrorNotifyKindUpstream5xx: upstreamItems,
	}

	// 3. 与已通知集合比较，每类只挑出新进入该异常的账号。
	notified := s.loadNotified(ctx)
	newItems := make([]accountErrorNotifyItem, 0)
	nextNotified := make(map[string][]int64, len(current))
	for _, kind := range []string{accountErrorNotifyKindError, accountErrorNotifyKindUnavailable, accountErrorNotifyKindUpstream5xx} {
		seen := map[int64]bool{}
		for _, id := range notified[kind] {
			seen[id] = true
		}
		ids := make([]int64, 0, len(current[kind]))
		for _, item := range current[kind] {
			ids = append(ids, item.ID)
			if !seen[item.ID] {
				newItems = append(newItems, item)
			}
		}
		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
		nextNotified[kind] = ids
	}

	// 4. 有新异常时发送一封汇总邮件；全部发送失败则不更新已通知集合，下一轮重试。
	if len(newItems) > 0 && !s.sendAlert(ctx, recipients, newItems, now) {
		return
	}

	// 5. 用当前异常账号覆盖已通知集合：已恢复的账号被移除，再次异常会重新通知。
	payload, _ := json.Marshal(nextNotified)
	if err := s.settingRepo.Set(ctx, accountErrorNotifyNotifiedKey, string(payload)); err != nil {
		slog.Error("account error notify: save notified ids failed", "error", err)
	}
}

// collectErrorAccounts 查询状态为 error（已停止调度，需要人工处理）的账号。
func (s *AccountErrorNotifyService) collectErrorAccounts(ctx context.Context) ([]accountErrorNotifyItem, error) {
	accounts, err := s.accountRepo.ListAllWithFilters(ctx, "", "", StatusError, "", 0, "")
	if err != nil {
		return nil, err
	}
	items := make([]accountErrorNotifyItem, 0, len(accounts))
	for _, account := range accounts {
		items = append(items, accountErrorNotifyItem{
			Kind:     accountErrorNotifyKindError,
			ID:       account.ID,
			Name:     account.Name,
			Platform: account.Platform,
			Detail:   account.ErrorMessage,
		})
	}
	return items, nil
}

// collectUnavailableAccounts 查询状态正常但暂时不可调度的账号：
// 429 限流、529 过载冷却，以及 OAuth 401 / 402 / OpenAI 403 等触发的临时停调度。
func (s *AccountErrorNotifyService) collectUnavailableAccounts(ctx context.Context, now time.Time) ([]accountErrorNotifyItem, error) {
	accounts, err := s.accountRepo.ListActive(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]accountErrorNotifyItem, 0)
	for _, account := range accounts {
		// 1. 每个账号把命中的不可用原因拼在一起。
		reasons := make([]string, 0, 3)
		if account.RateLimitResetAt != nil && now.Before(*account.RateLimitResetAt) {
			reasons = append(reasons, "限流(429)至 / rate limited until "+formatAccountErrorNotifyTime(*account.RateLimitResetAt))
		}
		if account.OverloadUntil != nil && now.Before(*account.OverloadUntil) {
			reasons = append(reasons, "过载(529)冷却至 / overloaded until "+formatAccountErrorNotifyTime(*account.OverloadUntil))
		}
		if account.TempUnschedulableUntil != nil && now.Before(*account.TempUnschedulableUntil) {
			reasons = append(reasons, "临时停调度至 / temporarily unschedulable until "+formatAccountErrorNotifyTime(*account.TempUnschedulableUntil)+": "+account.TempUnschedulableReason)
		}
		if len(reasons) == 0 {
			continue
		}
		items = append(items, accountErrorNotifyItem{
			Kind:     accountErrorNotifyKindUnavailable,
			ID:       account.ID,
			Name:     account.Name,
			Platform: account.Platform,
			Detail:   strings.Join(reasons, "\n"),
		})
	}
	return items, nil
}

// collectUpstream5xxAccounts 统计最近窗口内上游 5xx 次数达到阈值的账号。
// 5xx 不会改变账号状态，只能从运维错误日志统计。
func (s *AccountErrorNotifyService) collectUpstream5xxAccounts(ctx context.Context, now time.Time) ([]accountErrorNotifyItem, error) {
	// 1. 拉取窗口内的上游错误，包含已被重试恢复的请求。
	start := now.Add(-accountErrorNotify5xxWindow)
	list, err := s.opsRepo.ListErrorLogs(ctx, &OpsErrorLogFilter{
		StartTime:                &start,
		EndTime:                  &now,
		Owner:                    "provider",
		ErrorPhasesAny:           []string{"upstream"},
		IncludeRecoveredUpstream: true,
		View:                     "all",
		Page:                     1,
		PageSize:                 accountErrorNotify5xxPageSize,
	})
	if err != nil {
		return nil, err
	}

	// 2. 按账号统计 5xx 次数，记录状态码分布和最近一条错误信息。
	type stat struct {
		item    accountErrorNotifyItem
		count   int
		codes   map[int]int
		message string
	}
	stats := map[int64]*stat{}
	for _, e := range list.Errors {
		if e.AccountID == nil || e.StatusCode < 500 {
			continue
		}
		st := stats[*e.AccountID]
		if st == nil {
			st = &stat{
				item:  accountErrorNotifyItem{Kind: accountErrorNotifyKindUpstream5xx, ID: *e.AccountID, Name: e.AccountName, Platform: e.Platform},
				codes: map[int]int{},
				// 列表按时间倒序，首条即最近一条。
				message: e.Message,
			}
			stats[*e.AccountID] = st
		}
		st.count++
		st.codes[e.StatusCode]++
	}

	// 3. 只保留达到阈值的账号。
	items := make([]accountErrorNotifyItem, 0)
	for _, st := range stats {
		if st.count < accountErrorNotify5xxThreshold {
			continue
		}
		codes := make([]int, 0, len(st.codes))
		for code := range st.codes {
			codes = append(codes, code)
		}
		sort.Ints(codes)
		parts := make([]string, 0, len(codes))
		for _, code := range codes {
			parts = append(parts, fmt.Sprintf("%d×%d", code, st.codes[code]))
		}
		st.item.Detail = fmt.Sprintf("最近 %d 分钟 %d 次 / %d in last %d min: %s\n%s",
			int(accountErrorNotify5xxWindow.Minutes()), st.count, st.count, int(accountErrorNotify5xxWindow.Minutes()),
			strings.Join(parts, ", "), st.message)
		items = append(items, st.item)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items, nil
}

// loadNotified 读取各类异常已通知的账号 ID。旧格式或解析失败时视为空。
func (s *AccountErrorNotifyService) loadNotified(ctx context.Context) map[string][]int64 {
	result := map[string][]int64{}
	raw, err := s.settingRepo.GetValue(ctx, accountErrorNotifyNotifiedKey)
	if err != nil || strings.TrimSpace(raw) == "" {
		return result
	}
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return map[string][]int64{}
	}
	return result
}

// formatAccountErrorNotifyTime 按服务器本地时区格式化时间。
func formatAccountErrorNotifyTime(t time.Time) string {
	return t.Local().Format("2006-01-02 15:04:05 MST")
}

// sendAlert 给每个收件人发送异常汇总邮件，至少一封成功返回 true。
func (s *AccountErrorNotifyService) sendAlert(ctx context.Context, recipients []string, items []accountErrorNotifyItem, now time.Time) bool {
	// 1. 组装账号列表文本，每个账号一段：类型、ID、名称、平台和异常说明。
	lines := make([]string, 0, len(items))
	for _, item := range items {
		detail := strings.TrimSpace(item.Detail)
		if len([]rune(detail)) > accountErrorNotifyMaxMessageLen {
			detail = string([]rune(detail)[:accountErrorNotifyMaxMessageLen]) + "..."
		}
		lines = append(lines, fmt.Sprintf("[%s] #%d %s (%s)\n%s", accountErrorNotifyKindLabels[item.Kind], item.ID, item.Name, item.Platform, detail))
	}
	accountList := strings.Join(lines, "\n\n")
	triggeredAt := formatAccountErrorNotifyTime(now)
	count := strconv.Itoa(len(items))

	anySent := false
	for _, to := range recipients {
		// 2. 优先使用可编辑的通知模板发送。
		sendCtx, cancel := context.WithTimeout(ctx, emailSendTimeout)
		err := s.notificationEmailService.Send(sendCtx, NotificationEmailSendInput{
			Event:          NotificationEmailEventAccountErrorAlert,
			RecipientEmail: to,
			RecipientName:  emailRecipientName(to),
			Variables: map[string]string{
				"account_count": count,
				"account_list":  accountList,
				"triggered_at":  triggeredAt,
			},
		})
		cancel()
		if err == nil {
			anySent = true
			continue
		}
		if !shouldFallbackNotificationEmail(err) {
			slog.Warn("account error notify: send failed", "to", to, "error", err)
			continue
		}

		// 3. 模板或配置不可用时回退到内置正文。
		subject := fmt.Sprintf("账号异常告警 / Account Alert - %s", count)
		body := fmt.Sprintf(`<p>%s 个上游账号出现异常 / %s upstream account(s) need attention.</p><p>%s</p><pre style="white-space:pre-wrap;">%s</pre>`,
			count, count, triggeredAt, html.EscapeString(accountList))
		sendCtx, cancel = context.WithTimeout(ctx, emailSendTimeout)
		err = s.emailService.SendEmail(sendCtx, to, subject, body)
		cancel()
		if err != nil {
			slog.Warn("account error notify: fallback send failed", "to", to, "error", err)
			continue
		}
		anySent = true
	}
	if anySent {
		slog.Info("account error notify: alert sent", "accounts", len(items), "recipients", len(recipients))
	}
	return anySent
}
