package wire

import (
	"fmt"
	"notify-template/internal/model"
	"notify-template/internal/service"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitializeAPI 初始化API依赖
func InitializeAPI() (service.NotifyService, error) {
	// 初始化数据库连接
	db, err := initDB()
	if err != nil {
		return nil, err
	}

	// 初始化模板服务
	templateService := service.NewTemplateService(db)

	// 初始化通知服务
	notifyService := service.NewNotifyService(templateService)

	return notifyService, nil
}

// initDB 初始化数据库连接
func initDB() (*gorm.DB, error) {
	dsn := "root:password@tcp(127.0.0.1:3306)/notify_template?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// 自动迁移数据库表
	if err := db.AutoMigrate(
		&model.Template{},
		&model.TemplateVersion{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	return db, nil
}
