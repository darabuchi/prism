package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prism/core/internal/service"
)

// SubscriptionHandler 订阅处理器
type SubscriptionHandler struct {
	subscriptionSvc *service.SubscriptionService
}

// NewSubscriptionHandler 创建订阅处理器
func NewSubscriptionHandler(subscriptionSvc *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionSvc: subscriptionSvc,
	}
}

// CreateSubscription 创建订阅
// @Summary 创建订阅
// @Description 创建新的订阅配置
// @Tags 订阅管理
// @Accept json
// @Produce json
// @Param subscription body service.CreateSubscriptionRequest true "订阅信息"
// @Success 200 {object} APIResponse{data=storage.Subscription}
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var req service.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(1001, "参数错误", err.Error()))
		return
	}

	// 设置默认值
	if req.UserAgent == "" {
		req.UserAgent = "clash"
	}
	if req.UpdateInterval == 0 {
		req.UpdateInterval = 3600
	}

	subscription, err := h.subscriptionSvc.CreateSubscription(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(5000, "创建订阅失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse(subscription))
}

// UpdateSubscription 更新订阅
// @Summary 更新订阅
// @Description 更新订阅配置
// @Tags 订阅管理
// @Accept json
// @Produce json
// @Param subscription_id path int true "订阅ID"
// @Param subscription body service.UpdateSubscriptionRequest true "订阅信息"
// @Success 200 {object} APIResponse{data=storage.Subscription}
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /subscriptions/{subscription_id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("subscription_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(1001, "参数错误", "无效的订阅ID"))
		return
	}

	var req service.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(1001, "参数错误", err.Error()))
		return
	}

	subscription, err := h.subscriptionSvc.UpdateSubscription(uint(id), &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, NewErrorResponse(2001, "订阅不存在", err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, NewErrorResponse(5000, "更新订阅失败", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse(subscription))
}

// GetSubscription 获取单个订阅
// @Summary 获取订阅详情
// @Description 根据ID获取订阅详细信息
// @Tags 订阅管理
// @Produce json
// @Param subscription_id path int true "订阅ID"
// @Success 200 {object} APIResponse{data=storage.Subscription}
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Router /subscriptions/{subscription_id} [get]
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("subscription_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(1001, "参数错误", "无效的订阅ID"))
		return
	}

	subscription, err := h.subscriptionSvc.GetSubscription(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, NewErrorResponse(2001, "订阅不存在", err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, NewErrorResponse(5000, "获取订阅失败", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse(subscription))
}

// ListSubscriptions 获取订阅列表
// @Summary 获取订阅列表
// @Description 分页获取订阅列表，支持过滤条件
// @Tags 订阅管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页大小" default(20)
// @Param status query string false "状态过滤" Enums(active,inactive,error)
// @Param auto_update query bool false "自动更新过滤"
// @Success 200 {object} APIResponse{data=service.ListSubscriptionsResponse}
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /subscriptions [get]
func (h *SubscriptionHandler) ListSubscriptions(c *gin.Context) {
	var req service.ListSubscriptionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(1001, "参数错误", err.Error()))
		return
	}

	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 20
	}
	if req.Size > 100 {
		req.Size = 100
	}

	response, err := h.subscriptionSvc.ListSubscriptions(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(5000, "获取订阅列表失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse(response))
}

// DeleteSubscription 删除订阅
// @Summary 删除订阅
// @Description 删除指定的订阅
// @Tags 订阅管理
// @Produce json
// @Param subscription_id path int true "订阅ID"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /subscriptions/{subscription_id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("subscription_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(1001, "参数错误", "无效的订阅ID"))
		return
	}

	err = h.subscriptionSvc.DeleteSubscription(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, NewErrorResponse(2001, "订阅不存在", err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, NewErrorResponse(5000, "删除订阅失败", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse(map[string]interface{}{
		"message": "订阅删除成功",
	}))
}

// UpdateSubscriptionContent 手动更新订阅内容
// @Summary 手动更新订阅
// @Description 立即更新指定订阅的节点内容
// @Tags 订阅管理
// @Produce json
// @Param subscription_id path int true "订阅ID"
// @Success 200 {object} APIResponse{data=service.UpdateResult}
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /subscriptions/{subscription_id}/update [post]
func (h *SubscriptionHandler) UpdateSubscriptionContent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("subscription_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(1001, "参数错误", "无效的订阅ID"))
		return
	}

	result, err := h.subscriptionSvc.UpdateSubscriptionContent(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, NewErrorResponse(2001, "订阅不存在", err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, NewErrorResponse(5000, "更新订阅失败", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse(result))
}

// GetSubscriptionStats 获取订阅统计
// @Summary 获取订阅统计信息
// @Description 获取订阅的详细统计信息
// @Tags 订阅管理
// @Produce json
// @Param subscription_id path int true "订阅ID"
// @Success 200 {object} APIResponse{data=service.SubscriptionStats}
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /subscriptions/{subscription_id}/stats [get]
func (h *SubscriptionHandler) GetSubscriptionStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("subscription_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(1001, "参数错误", "无效的订阅ID"))
		return
	}

	stats, err := h.subscriptionSvc.GetSubscriptionStats(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, NewErrorResponse(2001, "订阅不存在", err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, NewErrorResponse(5000, "获取统计信息失败", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse(stats))
}

// GetSubscriptionLogs 获取订阅日志
// @Summary 获取订阅日志
// @Description 获取订阅更新日志记录
// @Tags 订阅管理
// @Produce json
// @Param subscription_id path int true "订阅ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页大小" default(20)
// @Param success query bool false "成功状态过滤"
// @Param update_type query string false "更新类型过滤" Enums(auto,manual,retry)
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Success 200 {object} APIResponse{data=service.LogsResponse}
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /subscriptions/{subscription_id}/logs [get]
func (h *SubscriptionHandler) GetSubscriptionLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("subscription_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(1001, "参数错误", "无效的订阅ID"))
		return
	}

	var req service.LogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(1001, "参数错误", err.Error()))
		return
	}

	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 20
	}
	if req.Size > 100 {
		req.Size = 100
	}

	response, err := h.subscriptionSvc.GetSubscriptionLogs(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(5000, "获取日志失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse(response))
}

// EnableSubscription 启用订阅
// @Summary 启用订阅
// @Description 启用指定的订阅
// @Tags 订阅管理
// @Produce json
// @Param subscription_id path int true "订阅ID"
// @Success 200 {object} APIResponse{data=storage.Subscription}
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /subscriptions/{subscription_id}/enable [post]
func (h *SubscriptionHandler) EnableSubscription(c *gin.Context) {
	h.toggleSubscriptionStatus(c, "active")
}

// DisableSubscription 禁用订阅
// @Summary 禁用订阅
// @Description 禁用指定的订阅
// @Tags 订阅管理
// @Produce json
// @Param subscription_id path int true "订阅ID"
// @Success 200 {object} APIResponse{data=storage.Subscription}
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /subscriptions/{subscription_id}/disable [post]
func (h *SubscriptionHandler) DisableSubscription(c *gin.Context) {
	h.toggleSubscriptionStatus(c, "inactive")
}

// toggleSubscriptionStatus 切换订阅状态
func (h *SubscriptionHandler) toggleSubscriptionStatus(c *gin.Context, status string) {
	id, err := strconv.ParseUint(c.Param("subscription_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(1001, "参数错误", "无效的订阅ID"))
		return
	}

	req := &service.UpdateSubscriptionRequest{
		Status: &status,
	}

	subscription, err := h.subscriptionSvc.UpdateSubscription(uint(id), req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, NewErrorResponse(2001, "订阅不存在", err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, NewErrorResponse(5000, "更新订阅状态失败", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse(subscription))
}

// ImportSubscription 导入订阅
// @Summary 导入订阅
// @Description 通过订阅链接批量导入订阅
// @Tags 订阅管理
// @Accept json
// @Produce json
// @Param import body ImportSubscriptionRequest true "导入信息"
// @Success 200 {object} APIResponse{data=ImportSubscriptionResponse}
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /subscriptions/import [post]
func (h *SubscriptionHandler) ImportSubscription(c *gin.Context) {
	var req ImportSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(1001, "参数错误", err.Error()))
		return
	}

	response := &ImportSubscriptionResponse{
		Total:   len(req.Subscriptions),
		Success: 0,
		Failed:  0,
		Results: make([]ImportResult, 0),
	}

	// 批量创建订阅
	for _, subReq := range req.Subscriptions {
		// 设置默认值
		if subReq.UserAgent == "" {
			subReq.UserAgent = "clash"
		}
		if subReq.UpdateInterval == 0 {
			subReq.UpdateInterval = 3600
		}

		subscription, err := h.subscriptionSvc.CreateSubscription(&subReq)
		result := ImportResult{
			Name: subReq.Name,
			URL:  subReq.URL,
		}

		if err != nil {
			result.Success = false
			result.Error = err.Error()
			response.Failed++
		} else {
			result.Success = true
			result.SubscriptionID = &subscription.ID
			response.Success++
		}

		response.Results = append(response.Results, result)
	}

	c.JSON(http.StatusOK, NewSuccessResponse(response))
}

// ExportSubscriptions 导出订阅
// @Summary 导出订阅
// @Description 导出所有订阅配置
// @Tags 订阅管理
// @Produce json
// @Param format query string false "导出格式" Enums(json,yaml) default(json)
// @Success 200 {object} APIResponse{data=interface{}}
// @Failure 500 {object} APIResponse
// @Router /subscriptions/export [get]
func (h *SubscriptionHandler) ExportSubscriptions(c *gin.Context) {
	format := c.DefaultQuery("format", "json")

	// 获取所有订阅
	req := &service.ListSubscriptionsRequest{
		Page: 1,
		Size: 1000, // 导出时不分页
	}

	response, err := h.subscriptionSvc.ListSubscriptions(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(5000, "导出订阅失败", err.Error()))
		return
	}

	// 简化订阅信息用于导出
	exportData := make([]map[string]interface{}, len(response.Subscriptions))
	for i, sub := range response.Subscriptions {
		exportData[i] = map[string]interface{}{
			"name":            sub.Name,
			"url":             sub.URL,
			"user_agent":      sub.UserAgent,
			"auto_update":     sub.AutoUpdate,
			"update_interval": sub.UpdateInterval,
			"created_at":      sub.CreatedAt,
		}
	}

	switch format {
	case "yaml":
		c.Header("Content-Type", "application/x-yaml")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"subscriptions_%d.yaml\"", len(exportData)))
	default:
		c.Header("Content-Type", "application/json")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"subscriptions_%d.json\"", len(exportData)))
	}

	c.JSON(http.StatusOK, NewSuccessResponse(map[string]interface{}{
		"format":        format,
		"total":         len(exportData),
		"export_time":   getCurrentTime(),
		"subscriptions": exportData,
	}))
}

// ImportSubscriptionRequest 导入订阅请求
type ImportSubscriptionRequest struct {
	Subscriptions []service.CreateSubscriptionRequest `json:"subscriptions" binding:"required"`
	ReplaceAll    bool                                `json:"replace_all"` // 是否替换所有现有订阅
}

// ImportSubscriptionResponse 导入订阅响应
type ImportSubscriptionResponse struct {
	Total   int            `json:"total"`
	Success int            `json:"success"`
	Failed  int            `json:"failed"`
	Results []ImportResult `json:"results"`
}

// ImportResult 导入结果
type ImportResult struct {
	Name           string `json:"name"`
	URL            string `json:"url"`
	Success        bool   `json:"success"`
	Error          string `json:"error,omitempty"`
	SubscriptionID *uint  `json:"subscription_id,omitempty"`
}
