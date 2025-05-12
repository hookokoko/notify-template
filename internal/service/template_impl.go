package service

import (
	"context"
	"errors"
	"fmt"
	"notify-template/internal/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	ErrTemplateNotFound      = errors.New("template not found")
	ErrVersionNotFound       = errors.New("version not found")
	ErrInvalidVersionName    = errors.New("invalid version name")
	ErrInvalidOperation      = errors.New("invalid operation")
	ErrTemplateAlreadyExists = errors.New("template already exists")
)

type templateServiceImpl struct {
	db *gorm.DB
}

// NewTemplateService 创建模板服务实例
func NewTemplateService(db *gorm.DB) TemplateService {
	return &templateServiceImpl{db: db}
}

// CreateTemplate 创建模板
func (s *templateServiceImpl) CreateTemplate(ctx context.Context, template *model.Template) error {
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()
	return s.db.Create(template).Error
}

// UpdateTemplate 更新模板
func (s *templateServiceImpl) UpdateTemplate(ctx context.Context, template *model.Template) error {
	template.UpdatedAt = time.Now()
	return s.db.Save(template).Error
}

// DeleteTemplate 删除模板
func (s *templateServiceImpl) DeleteTemplate(ctx context.Context, id int64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除模板版本
		if err := tx.Where("channel_template_id = ?", id).Delete(&model.TemplateVersion{}).Error; err != nil {
			return err
		}
		// 删除模板
		return tx.Delete(&model.Template{}, id).Error
	})
}

// GetTemplate 获取模板
func (s *templateServiceImpl) GetTemplate(ctx context.Context, id int64) (*model.Template, error) {
	var template model.Template
	if err := s.db.First(&template, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTemplateNotFound
		}
		return nil, err
	}
	return &template, nil
}

// ListTemplates 获取模板列表
func (s *templateServiceImpl) ListTemplates(ctx context.Context, query *TemplateQuery) ([]*model.Template, int64, error) {
	var templates []*model.Template
	var total int64

	db := s.db.Model(&model.Template{})

	if query.OwnerID != 0 {
		db = db.Where("owner_id = ?", query.OwnerID)
	}
	if query.OwnerType != "" {
		db = db.Where("owner_type = ?", query.OwnerType)
	}
	if query.DeliveryType != "" {
		db = db.Where("delivery_type = ?", query.DeliveryType)
	}
	if query.BusinessType != 0 {
		db = db.Where("business_type = ?", query.BusinessType)
	}
	if query.Department != "" {
		db = db.Where("department = ?", query.Department)
	}
	if query.Scenario != "" {
		db = db.Where("scenario = ?", query.Scenario)
	}
	if query.Country != "" {
		db = db.Where("country = ?", query.Country)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}

	if err := db.Find(&templates).Error; err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}

// generateNextVersionName 生成下一个版本名称
func (s *templateServiceImpl) generateNextVersionName(ctx context.Context, templateID int64, currentVersion string) (string, error) {
	parts := strings.Split(currentVersion, ".")
	if len(parts) < 3 {
		return "", ErrInvalidVersionName
	}

	major := parts[0]
	minor := parts[1]
	patch := parts[2]

	// 如果是子版本（如v0.0.1.1），增加子版本号
	if len(parts) > 3 {
		subVersion := parts[3]
		subVersionNum := 0
		fmt.Sscanf(subVersion, "%d", &subVersionNum)
		return fmt.Sprintf("%s.%s.%s.%d", major, minor, patch, subVersionNum+1), nil
	}

	// 否则增加补丁版本号
	patchNum := 0
	fmt.Sscanf(patch, "%d", &patchNum)
	return fmt.Sprintf("%s.%s.%d", major, minor, patchNum+1), nil
}

// CreateTemplateVersion 创建模板版本
func (s *templateServiceImpl) CreateTemplateVersion(ctx context.Context, version *model.TemplateVersion) error {
	// 获取当前模板的最新版本
	var latestVersion model.TemplateVersion
	err := s.db.Where("channel_template_id = ?", version.ChannelTemplateID).
		Order("id DESC").
		First(&latestVersion).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 生成新版本号
	if errors.Is(err, gorm.ErrRecordNotFound) {
		version.Name = "v0.0.1"
	} else {
		version.Name, err = s.generateNextVersionName(ctx, version.ChannelTemplateID, latestVersion.Name)
		if err != nil {
			return err
		}
	}

	version.CreatedAt = time.Now()
	version.UpdatedAt = time.Now()
	version.AuditStatus = "PENDING"

	return s.db.Create(version).Error
}

// UpdateTemplateVersion 更新模板版本
func (s *templateServiceImpl) UpdateTemplateVersion(ctx context.Context, version *model.TemplateVersion) error {
	version.UpdatedAt = time.Now()
	return s.db.Save(version).Error
}

// GetTemplateVersion 获取模板版本
func (s *templateServiceImpl) GetTemplateVersion(ctx context.Context, id int64) (*model.TemplateVersion, error) {
	var version model.TemplateVersion
	if err := s.db.First(&version, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVersionNotFound
		}
		return nil, err
	}
	return &version, nil
}

// ListTemplateVersions 获取模板版本列表
func (s *templateServiceImpl) ListTemplateVersions(ctx context.Context, templateID int64) ([]*model.TemplateVersion, error) {
	var versions []*model.TemplateVersion
	if err := s.db.Where("channel_template_id = ?", templateID).
		Order("id DESC").
		Find(&versions).Error; err != nil {
		return nil, err
	}
	return versions, nil
}

// SubmitForReview 提交审核
func (s *templateServiceImpl) SubmitForReview(ctx context.Context, versionID int64, needVendorReview bool) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var version model.TemplateVersion
		if err := tx.First(&version, versionID).Error; err != nil {
			return err
		}

		now := time.Now()
		version.LastReviewSubmissionTime = &now
		version.AuditStatus = "IN_REVIEW"

		return tx.Save(&version).Error
	})
}

// ApproveTemplate 审核通过
func (s *templateServiceImpl) ApproveTemplate(ctx context.Context, versionID int64, auditorID int64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var version model.TemplateVersion
		if err := tx.First(&version, versionID).Error; err != nil {
			return err
		}

		now := time.Now()
		version.AuditTime = &now
		version.AuditorID = auditorID
		version.AuditStatus = "APPROVED"

		// 更新模板的活跃版本
		if err := tx.Model(&model.Template{}).
			Where("id = ?", version.ChannelTemplateID).
			Updates(map[string]interface{}{
				"active_version_id":   version.ID,
				"active_version_name": version.Name,
				"updated_at":          now,
			}).Error; err != nil {
			return err
		}

		return tx.Save(&version).Error
	})
}

// RejectTemplate 审核拒绝
func (s *templateServiceImpl) RejectTemplate(ctx context.Context, versionID int64, auditorID int64, reason string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var version model.TemplateVersion
		if err := tx.First(&version, versionID).Error; err != nil {
			return err
		}

		now := time.Now()
		version.AuditTime = &now
		version.AuditorID = auditorID
		version.AuditStatus = "REJECTED"
		version.RejectReason = reason

		return tx.Save(&version).Error
	})
}

// RollbackTemplate 回滚模板版本
func (s *templateServiceImpl) RollbackTemplate(ctx context.Context, templateID int64, versionID int64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var version model.TemplateVersion
		if err := tx.First(&version, versionID).Error; err != nil {
			return err
		}

		if version.ChannelTemplateID != templateID {
			return ErrInvalidOperation
		}

		now := time.Now()
		return tx.Model(&model.Template{}).
			Where("id = ?", templateID).
			Updates(map[string]interface{}{
				"active_version_id":   version.ID,
				"active_version_name": version.Name,
				"updated_at":          now,
			}).Error
	})
}

// GetTemplateContent 获取模板内容
func (s *templateServiceImpl) GetTemplateContent(ctx context.Context, templateID int64, language string) (string, error) {
	var template model.Template
	if err := s.db.First(&template, templateID).Error; err != nil {
		return "", err
	}

	var version model.TemplateVersion
	if err := s.db.Where("channel_template_id = ? AND language = ?", templateID, language).
		First(&version).Error; err != nil {
		return "", err
	}

	return version.Content, nil
}
