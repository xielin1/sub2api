package service

import (
	"context"
	"errors"
)

// GetSalesRecruitment 只供已登录用户读取招募配置，不进入公开配置或 HTML 注入。
func (s *SettingService) GetSalesRecruitment(ctx context.Context) (string, error) {
	// 1. 老站未配置时返回空配置，其他存储错误继续向上传递。
	value, err := s.settingRepo.GetValue(ctx, SettingKeySalesRecruitment)
	if errors.Is(err, ErrSettingNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return value, nil
}
