package api

import (
	"net/http"
	"notify-template/internal/model"
	"notify-template/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TemplateHandler struct {
	templateService service.TemplateService
}

func NewTemplateHandler(templateService service.TemplateService) *TemplateHandler {
	return &TemplateHandler{
		templateService: templateService,
	}
}

// CreateTemplate 创建模板
func (h *TemplateHandler) CreateTemplate(c *gin.Context) {
	var template model.Template
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.templateService.CreateTemplate(c.Request.Context(), &template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

// UpdateTemplate 更新模板
func (h *TemplateHandler) UpdateTemplate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
		return
	}

	var template model.Template
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template.ID = id
	if err := h.templateService.UpdateTemplate(c.Request.Context(), &template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

// DeleteTemplate 删除模板
func (h *TemplateHandler) DeleteTemplate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
		return
	}

	if err := h.templateService.DeleteTemplate(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "template deleted successfully"})
}

// GetTemplate 获取模板
func (h *TemplateHandler) GetTemplate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
		return
	}

	template, err := h.templateService.GetTemplate(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrTemplateNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

// ListTemplates 获取模板列表
func (h *TemplateHandler) ListTemplates(c *gin.Context) {
	query := &service.TemplateQuery{
		Page:     1,
		PageSize: 10,
	}

	// 解析查询参数
	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		query.Page = page
	}
	if pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10")); err == nil {
		query.PageSize = pageSize
	}
	if ownerID, err := strconv.ParseInt(c.Query("owner_id"), 10, 64); err == nil {
		query.OwnerID = ownerID
	}
	if ownerType := c.Query("owner_type"); ownerType != "" {
		query.OwnerType = model.OwnerType(ownerType)
	}
	if deliveryType := c.Query("delivery_type"); deliveryType != "" {
		query.DeliveryType = model.DeliveryType(deliveryType)
	}
	if businessType, err := strconv.ParseInt(c.Query("business_type"), 10, 64); err == nil {
		query.BusinessType = businessType
	}
	query.Department = c.Query("department")
	query.Scenario = c.Query("scenario")
	query.Country = c.Query("country")

	templates, total, err := h.templateService.ListTemplates(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"templates": templates,
	})
}

// CreateTemplateVersion 创建模板版本
func (h *TemplateHandler) CreateTemplateVersion(c *gin.Context) {
	var version model.TemplateVersion
	if err := c.ShouldBindJSON(&version); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.templateService.CreateTemplateVersion(c.Request.Context(), &version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, version)
}

// UpdateTemplateVersion 更新模板版本
func (h *TemplateHandler) UpdateTemplateVersion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version id"})
		return
	}

	var version model.TemplateVersion
	if err := c.ShouldBindJSON(&version); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	version.ID = id
	if err := h.templateService.UpdateTemplateVersion(c.Request.Context(), &version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, version)
}

// GetTemplateVersion 获取模板版本
func (h *TemplateHandler) GetTemplateVersion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version id"})
		return
	}

	version, err := h.templateService.GetTemplateVersion(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrVersionNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, version)
}

// ListTemplateVersions 获取模板版本列表
func (h *TemplateHandler) ListTemplateVersions(c *gin.Context) {
	templateID, err := strconv.ParseInt(c.Param("template_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
		return
	}

	versions, err := h.templateService.ListTemplateVersions(c.Request.Context(), templateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, versions)
}

// SubmitForReview 提交审核
func (h *TemplateHandler) SubmitForReview(c *gin.Context) {
	versionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version id"})
		return
	}

	var req struct {
		NeedVendorReview bool `json:"need_vendor_review"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.templateService.SubmitForReview(c.Request.Context(), versionID, req.NeedVendorReview); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "submitted for review successfully"})
}

// ApproveTemplate 审核通过
func (h *TemplateHandler) ApproveTemplate(c *gin.Context) {
	versionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version id"})
		return
	}

	var req struct {
		AuditorID int64 `json:"auditor_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.templateService.ApproveTemplate(c.Request.Context(), versionID, req.AuditorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "template approved successfully"})
}

// RejectTemplate 审核拒绝
func (h *TemplateHandler) RejectTemplate(c *gin.Context) {
	versionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version id"})
		return
	}

	var req struct {
		AuditorID    int64  `json:"auditor_id"`
		RejectReason string `json:"reject_reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.templateService.RejectTemplate(c.Request.Context(), versionID, req.AuditorID, req.RejectReason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "template rejected successfully"})
}

// RollbackTemplate 回滚模板版本
func (h *TemplateHandler) RollbackTemplate(c *gin.Context) {
	templateID, err := strconv.ParseInt(c.Param("template_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
		return
	}

	versionID, err := strconv.ParseInt(c.Param("version_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version id"})
		return
	}

	if err := h.templateService.RollbackTemplate(c.Request.Context(), templateID, versionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "template rolled back successfully"})
}

// GetTemplateContent 获取模板内容
func (h *TemplateHandler) GetTemplateContent(c *gin.Context) {
	templateID, err := strconv.ParseInt(c.Param("template_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
		return
	}

	language := c.Query("language")
	if language == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "language is required"})
		return
	}

	content, err := h.templateService.GetTemplateContent(c.Request.Context(), templateID, language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"content": content})
}
