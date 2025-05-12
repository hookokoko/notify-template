package model

import "time"

// TemplateStatus 模板状态
type TemplateStatus string

const (
	StatusEditing         TemplateStatus = "EDITING"          // 编辑中
	StatusFeishuReviewing TemplateStatus = "FEISHU_REVIEWING" // 飞书审核中
	StatusFeishuApproved  TemplateStatus = "FEISHU_APPROVED"  // 飞书审核通过
	StatusFeishuRejected  TemplateStatus = "FEISHU_REJECTED"  // 飞书审核拒绝
	StatusVendorReviewing TemplateStatus = "VENDOR_REVIEWING" // 服务商审核中
	StatusVendorApproved  TemplateStatus = "VENDOR_APPROVED"  // 服务商审核通过
	StatusVendorRejected  TemplateStatus = "VENDOR_REJECTED"  // 服务商审核拒绝
	StatusOnline          TemplateStatus = "ONLINE"           // 已上线
	StatusStopped         TemplateStatus = "STOPPED"          // 已停止
)

// DeliveryType 发送类型
type DeliveryType string

const (
	DeliveryTypeSMS   DeliveryType = "SMS"
	DeliveryTypeEmail DeliveryType = "EMAIL"
	DeliveryTypeInApp DeliveryType = "IN_APP"
)

// OwnerType 所有者类型
type OwnerType string

const (
	OwnerTypePerson       OwnerType = "person"
	OwnerTypeOrganization OwnerType = "organization"
)

// Template 模板模型
type Template struct {
	ID                int64        `json:"id" gorm:"primaryKey;autoIncrement"`
	OwnerID           int64        `json:"owner_id"`
	OwnerType         OwnerType    `json:"owner_type"`
	Name              string       `json:"name"`
	Description       string       `json:"description"`
	DeliveryType      DeliveryType `json:"delivery_type"`
	BusinessType      int64        `json:"business_type"`
	Department        string       `json:"department"`
	Scenario          string       `json:"scenario"`
	Country           string       `json:"country"`
	Priority          string       `json:"priority"`
	ActiveVersionID   int64        `json:"active_version_id"`
	ActiveVersionName string       `json:"active_version_name"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}

// TemplateVersion 模板版本模型
type TemplateVersion struct {
	ID                       int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	ChannelTemplateID        int64      `json:"channel_template_id"`
	Name                     string     `json:"name"`
	Signature                string     `json:"signature"`
	Language                 string     `json:"language"`
	Content                  string     `json:"content"`
	Remark                   string     `json:"remark"`
	AuditID                  int64      `json:"audit_id"`
	AuditorID                int64      `json:"auditor_id"`
	AuditTime                *time.Time `json:"audit_time"`
	AuditStatus              string     `json:"audit_status"`
	RejectReason             string     `json:"reject_reason"`
	LastReviewSubmissionTime *time.Time `json:"last_review_submission_time"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
}

// TableName 指定表名
func (Template) TableName() string {
	return "t_templates"
}

// TableName 指定表名
func (TemplateVersion) TableName() string {
	return "t_template_versions"
}
