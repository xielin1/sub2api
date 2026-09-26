package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// AnnouncementHandler handles user announcement operations
type AnnouncementHandler struct {
	announcementService *service.AnnouncementService
}

// NewAnnouncementHandler creates a new user announcement handler
func NewAnnouncementHandler(announcementService *service.AnnouncementService) *AnnouncementHandler {
	return &AnnouncementHandler{
		announcementService: announcementService,
	}
}

// List handles listing announcements visible to current user
// GET /api/v1/announcements
func (h *AnnouncementHandler) List(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}

	unreadOnly := parseBoolQuery(c.Query("unread_only"))

	items, err := h.announcementService.ListForUser(c.Request.Context(), subject.UserID, unreadOnly)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	out := make([]dto.UserAnnouncement, 0, len(items))
	for i := range items {
		out = append(out, *dto.UserAnnouncementFromService(&items[i]))
	}
	response.Success(c, out)
}

// MarkRead marks an announcement as read for current user
// POST /api/v1/announcements/:id/read
func (h *AnnouncementHandler) MarkRead(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}

	announcementID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || announcementID <= 0 {
		response.BadRequest(c, "Invalid announcement ID")
		return
	}

	if err := h.announcementService.MarkRead(c.Request.Context(), subject.UserID, announcementID); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"message": "ok"})
}

func parseBoolQuery(v string) bool {
	switch strings.TrimSpace(strings.ToLower(v)) {
	case "1", "true", "yes", "y", "on":
		return true
	default:
		return false
	}
}

// publicAnnouncement 首页公开公告条目（白名单字段）。
type publicAnnouncement struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// ListPublic 返回匿名可见的最近公告
// GET /api/v1/public/announcements
func (h *AnnouncementHandler) ListPublic(c *gin.Context) {
	// 1. 取最近 5 条面向所有用户的生效公告
	items, err := h.announcementService.ListPublic(c.Request.Context(), 5)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	// 2. 映射为白名单字段返回
	out := make([]publicAnnouncement, 0, len(items))
	for i := range items {
		out = append(out, publicAnnouncement{
			ID:        items[i].ID,
			Title:     items[i].Title,
			Content:   items[i].Content,
			CreatedAt: items[i].CreatedAt,
		})
	}
	response.Success(c, out)
}
