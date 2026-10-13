-- Migration: 246_channel_monitor_iq_mode
-- 渠道监控新增智商雷达检测方式 iq：发送糖果推理题，回答出现独立的 21 判为通过。
ALTER TABLE channel_monitors
    DROP CONSTRAINT IF EXISTS channel_monitors_check_mode_check;

ALTER TABLE channel_monitors
    ADD CONSTRAINT channel_monitors_check_mode_check
    CHECK (check_mode IN ('probe', 'quota', 'quota_probe', 'iq'));

COMMENT ON COLUMN channel_monitors.check_mode IS
    'probe = LLM 探活（默认）；quota = 仅查关联账号用量；quota_probe = 探活 + 配额；iq = 智商雷达糖果题';
