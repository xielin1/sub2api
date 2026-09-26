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
	// accountErrorNotifyInterval 账号故障扫描周期。
	accountErrorNotifyInterval = time.Minute
	// accountErrorNotifyNotifiedKey 保存已通知且仍处于错误状态的账号 ID（JSON 数组），重启后不重复通知。
	accountErrorNotifyNotifiedKey = "account_error_notify_notified_ids"
	// accountErrorNotifyMaxMessageLen 邮件中单个账号错误信息的最大长度。
	accountErrorNotifyMaxMessageLen = 300
)

// AccountErrorNotifyService 定时扫描进入错误状态的上游账号，并邮件通知配置的邮箱。
type AccountErrorNotifyService struct {
	accountRepo              AccountRepository
	settingRepo              SettingRepository
	emailService             *EmailService
	notificationEmailService *NotificationEmailService
	stopCh                   chan struct{}
	stopOnce                 sync.Once
	wg                       sync.WaitGroup
}

// NewAccountErrorNotifyService 创建账号故障通知服务。
func NewAccountErrorNotifyService(accountRepo AccountRepository, settingRepo SettingRepository, emailService *EmailService, notificationEmailService *NotificationEmailService) *AccountErrorNotifyService {
	return &AccountErrorNotifyService{
		accountRepo:              accountRepo,
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

	// 2. 查询当前处于错误状态的账号。
	accounts, err := s.accountRepo.ListAllWithFilters(ctx, "", "", StatusError, "", 0, "")
	if err != nil {
		slog.Error("account error notify: list error accounts failed", "error", err)
		return
	}

	// 3. 与已通知集合比较，找出新进入错误状态的账号。
	notified := s.loadNotifiedIDs(ctx)
	current := make([]int64, 0, len(accounts))
	newAccounts := make([]Account, 0)
	for _, account := range accounts {
		current = append(current, account.ID)
		if !notified[account.ID] {
			newAccounts = append(newAccounts, account)
		}
	}

	// 4. 有新故障账号时发送一封汇总邮件；全部发送失败则不更新已通知集合，下一轮重试。
	if len(newAccounts) > 0 && !s.sendAlert(ctx, recipients, newAccounts) {
		return
	}

	// 5. 用当前错误账号覆盖已通知集合：已恢复的账号被移除，下次再出错会重新通知。
	sort.Slice(current, func(i, j int) bool { return current[i] < current[j] })
	payload, _ := json.Marshal(current)
	if err := s.settingRepo.Set(ctx, accountErrorNotifyNotifiedKey, string(payload)); err != nil {
		slog.Error("account error notify: save notified ids failed", "error", err)
	}
}

// loadNotifiedIDs 读取已通知的账号 ID 集合。
func (s *AccountErrorNotifyService) loadNotifiedIDs(ctx context.Context) map[int64]bool {
	result := map[int64]bool{}
	raw, err := s.settingRepo.GetValue(ctx, accountErrorNotifyNotifiedKey)
	if err != nil || strings.TrimSpace(raw) == "" {
		return result
	}
	var ids []int64
	if err := json.Unmarshal([]byte(raw), &ids); err != nil {
		return result
	}
	for _, id := range ids {
		result[id] = true
	}
	return result
}

// sendAlert 给每个收件人发送故障汇总邮件，至少一封成功返回 true。
func (s *AccountErrorNotifyService) sendAlert(ctx context.Context, recipients []string, accounts []Account) bool {
	// 1. 组装账号列表文本，每个账号一段：ID、名称、平台和错误信息。
	lines := make([]string, 0, len(accounts))
	for _, account := range accounts {
		message := strings.TrimSpace(account.ErrorMessage)
		if len([]rune(message)) > accountErrorNotifyMaxMessageLen {
			message = string([]rune(message)[:accountErrorNotifyMaxMessageLen]) + "..."
		}
		lines = append(lines, fmt.Sprintf("#%d %s (%s)\n%s", account.ID, account.Name, account.Platform, message))
	}
	accountList := strings.Join(lines, "\n\n")
	triggeredAt := time.Now().UTC().Format(time.RFC3339)
	count := strconv.Itoa(len(accounts))

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
		subject := fmt.Sprintf("账号故障告警 / Account Error Alert - %s", count)
		body := fmt.Sprintf(`<p>%s 个上游账号进入错误状态 / %s account(s) entered the error state.</p><p>%s</p><pre style="white-space:pre-wrap;">%s</pre>`,
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
		slog.Info("account error notify: alert sent", "accounts", len(accounts), "recipients", len(recipients))
	}
	return anySent
}
