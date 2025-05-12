package service

import (
	"context"
	"notify-template/internal/model"
)

// TemplateService 模板服务接口
type TemplateService interface {
	// CreateTemplate 创建模板
	CreateTemplate(ctx context.Context, template *model.Template) error

	// UpdateTemplate 更新模板
	UpdateTemplate(ctx context.Context, template *model.Template) error

	// DeleteTemplate 删除模板
	DeleteTemplate(ctx context.Context, id int64) error

	// GetTemplate 获取模板
	GetTemplate(ctx context.Context, id int64) (*model.Template, error)

	// ListTemplates 获取模板列表
	ListTemplates(ctx context.Context, query *TemplateQuery) ([]*model.Template, int64, error)

	// CreateTemplateVersion 创建模板版本
	CreateTemplateVersion(ctx context.Context, version *model.TemplateVersion) error

	// UpdateTemplateVersion 更新模板版本
	UpdateTemplateVersion(ctx context.Context, version *model.TemplateVersion) error

	// GetTemplateVersion 获取模板版本
	GetTemplateVersion(ctx context.Context, id int64) (*model.TemplateVersion, error)

	// ListTemplateVersions 获取模板版本列表
	ListTemplateVersions(ctx context.Context, templateID int64) ([]*model.TemplateVersion, error)

	// SubmitForReview 提交审核
	SubmitForReview(ctx context.Context, versionID int64, needVendorReview bool) error

	// ApproveTemplate 审核通过
	ApproveTemplate(ctx context.Context, versionID int64, auditorID int64) error

	// RejectTemplate 审核拒绝
	RejectTemplate(ctx context.Context, versionID int64, auditorID int64, reason string) error

	// RollbackTemplate 回滚模板版本
	RollbackTemplate(ctx context.Context, templateID int64, versionID int64) error

	// GetTemplateContent 获取模板内容
	GetTemplateContent(ctx context.Context, templateID int64, language string) (string, error)
}

// TemplateQuery 模板查询参数
type TemplateQuery struct {
	OwnerID      int64
	OwnerType    model.OwnerType
	DeliveryType model.DeliveryType
	BusinessType int64
	Department   string
	Scenario     string
	Country      string
	Page         int
	PageSize     int
}
