package handler

import (
	"net/http"
	"strconv"

	"clash-server/internal/model"
	"clash-server/internal/service"
	"clash-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	subService  *service.SubscriptionService
	coreService *service.CoreService
}

func NewSubscriptionHandler() *SubscriptionHandler {
	return &SubscriptionHandler{
		subService:  service.NewSubscriptionService(),
		coreService: service.GetCoreService(),
	}
}

func (h *SubscriptionHandler) List(c *gin.Context) {
	subs, err := h.subService.ListWithCounts()
	if err != nil {
		c.JSON(http.StatusOK, response.InternalError("获取订阅列表失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{"subscriptions": subs}))
}

func (h *SubscriptionHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的订阅ID"))
		return
	}
	sub, customization, err := h.subService.GetWithRelations(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, response.NotFound("订阅不存在"))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{
		"subscription":  sub,
		"customization": customization,
	}))
}

type CreateSubscriptionRequest struct {
	Name       string `json:"name" binding:"required"`
	SourceType string `json:"source_type"`
	URL        string `json:"url"`
	Interval   int    `json:"interval"`
	UseProxy   bool   `json:"use_proxy"`
	UserAgent  string `json:"user_agent"`
	SkipCert   bool   `json:"skip_cert"`
}

func (h *SubscriptionHandler) Create(c *gin.Context) {
	var req CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.BadRequest("请输入有效的订阅名称"))
		return
	}
	sub := &model.Subscription{
		Name:       req.Name,
		SourceType: model.SourceType(req.SourceType),
		URL:        req.URL,
		Interval:   req.Interval,
		UseProxy:   req.UseProxy,
		UserAgent:  req.UserAgent,
		SkipCert:   req.SkipCert,
	}
	if sub.SourceType == "" {
		sub.SourceType = model.SourceTypeRemote
	}
	if sub.Interval == 0 {
		sub.Interval = 60
	}
	if err := h.subService.Create(sub); err != nil {
		c.JSON(http.StatusOK, response.InternalError("创建订阅失败，请稍后重试"))
		return
	}

	if sub.SourceType == model.SourceTypeRemote {
		if _, err := h.subService.Refresh(sub.ID); err != nil {
			c.JSON(http.StatusOK, response.Error(response.CodeOperationFailed, "订阅已创建，但获取内容失败: "+err.Error()))
			return
		}
	}

	if scheduler := GetSubscriptionScheduler(); scheduler != nil {
		scheduler.AddSubscription(sub)
	}

	coreErr := h.applyConfig()
	c.JSON(http.StatusOK, response.SuccessWithCoreError(gin.H{"subscription": sub}, coreErr))
}

type UpdateSubscriptionRequest struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	Interval  int    `json:"interval"`
	UseProxy  bool   `json:"use_proxy"`
	UserAgent string `json:"user_agent"`
	SkipCert  bool   `json:"skip_cert"`
}

func (h *SubscriptionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的订阅ID"))
		return
	}
	var req UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.BadRequest("请输入有效的订阅信息"))
		return
	}
	sub, err := h.subService.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, response.NotFound("订阅不存在"))
		return
	}
	sub.Name = req.Name
	if req.URL != "" {
		sub.URL = req.URL
	}
	if req.Interval > 0 {
		sub.Interval = req.Interval
	}
	sub.UseProxy = req.UseProxy
	sub.UserAgent = req.UserAgent
	sub.SkipCert = req.SkipCert
	if err := h.subService.Update(sub); err != nil {
		c.JSON(http.StatusOK, response.InternalError("更新订阅失败，请稍后重试"))
		return
	}

	if sub.SourceType == model.SourceTypeRemote && sub.URL != "" {
		if _, err := h.subService.Refresh(sub.ID); err != nil {
			c.JSON(http.StatusOK, response.Error(response.CodeOperationFailed, "订阅已更新，但获取内容失败: "+err.Error()))
			return
		}
	}

	if scheduler := GetSubscriptionScheduler(); scheduler != nil {
		scheduler.UpdateSubscription(sub)
	}

	coreErr := h.applyConfig()
	c.JSON(http.StatusOK, response.SuccessWithCoreError(gin.H{"subscription": sub}, coreErr))
}

func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的订阅ID"))
		return
	}

	sub, err := h.subService.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, response.NotFound("订阅不存在"))
		return
	}

	wasActive := sub.Active

	if scheduler := GetSubscriptionScheduler(); scheduler != nil {
		scheduler.RemoveSubscription(uint(id))
	}

	if err := h.subService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusOK, response.InternalError("删除订阅失败，请稍后重试"))
		return
	}

	if wasActive {
		subs, _ := h.subService.List()
		if len(subs) > 0 {
			h.subService.Activate(subs[0].ID)
		}
	}

	coreErr := h.applyConfig()
	c.JSON(http.StatusOK, response.SuccessWithCoreError(nil, coreErr))
}

func (h *SubscriptionHandler) Activate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的订阅ID"))
		return
	}

	if _, err := h.subService.Get(uint(id)); err != nil {
		c.JSON(http.StatusOK, response.NotFound("订阅不存在"))
		return
	}

	if err := h.subService.Activate(uint(id)); err != nil {
		c.JSON(http.StatusOK, response.InternalError("激活订阅失败，请稍后重试"))
		return
	}

	coreErr := h.applyConfig()
	c.JSON(http.StatusOK, response.SuccessWithCoreError(nil, coreErr))
}

func (h *SubscriptionHandler) Refresh(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的订阅ID"))
		return
	}
	result, err := h.subService.Refresh(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, response.InternalError("刷新订阅失败，请检查订阅地址是否正确"))
		return
	}
	if result.Error != "" {
		c.JSON(http.StatusOK, response.Error(response.CodeOperationFailed, "刷新订阅失败: "+result.Error))
		return
	}
	coreErr := h.applyConfig()
	c.JSON(http.StatusOK, response.SuccessWithCoreError(gin.H{
		"subscription": result.Subscription,
	}, coreErr))
}

func (h *SubscriptionHandler) GetMerged(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的订阅ID"))
		return
	}
	config, yaml, err := h.subService.GetMerged(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, response.InternalError("获取合并配置失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{
		"config": config,
		"yaml":   yaml,
	}))
}

func (h *SubscriptionHandler) GetContent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的订阅ID"))
		return
	}
	sub, err := h.subService.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, response.NotFound("订阅不存在"))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{
		"content": sub.Content,
	}))
}

type UpdateContentRequest struct {
	Content string `json:"content" binding:"required"`
}

func (h *SubscriptionHandler) UpdateContent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的订阅ID"))
		return
	}
	var req UpdateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.BadRequest("请提供订阅内容"))
		return
	}
	if err := h.subService.UpdateContent(uint(id), req.Content); err != nil {
		c.JSON(http.StatusOK, response.InternalError("更新订阅内容失败"))
		return
	}
	coreErr := h.applyConfig()
	c.JSON(http.StatusOK, response.SuccessWithCoreError(nil, coreErr))
}

func (h *SubscriptionHandler) applyConfig() string {
	result := h.coreService.ApplyConfig()
	if result.Error != nil {
		return result.Error.Error()
	}
	return ""
}
