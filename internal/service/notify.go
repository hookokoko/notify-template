package service

// NotifyService 通知服务接口
type NotifyService interface {
	// SendNotification 发送通知
	SendNotification(message string) error

	// GetTemplateService 获取模板服务
	GetTemplateService() TemplateService
}

type notifyServiceImpl struct {
	templateService TemplateService
}

// NewNotifyService 创建通知服务实例
func NewNotifyService(templateService TemplateService) NotifyService {
	return &notifyServiceImpl{
		templateService: templateService,
	}
}

// SendNotification 发送通知
func (s *notifyServiceImpl) SendNotification(message string) error {
	// TODO: 实现通知发送逻辑
	return nil
}

// GetTemplateService 获取模板服务
func (s *notifyServiceImpl) GetTemplateService() TemplateService {
	return s.templateService
}
