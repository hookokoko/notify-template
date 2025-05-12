# Notify Template

这是一个基于 Gin 框架的 Go Web 服务模板项目。

## 项目结构

```
.
├── cmd/                    # 主要的应用程序入口
│   └── api/               # API 服务入口
├── config/                # 配置文件目录
├── internal/              # 私有应用程序和库代码
│   ├── handler/          # HTTP 处理器
│   ├── middleware/       # HTTP 中间件
│   ├── model/           # 数据模型
│   ├── repository/      # 数据访问层
│   └── service/         # 业务逻辑层
├── pkg/                  # 可以被外部应用程序使用的库代码
│   ├── database/        # 数据库相关
│   ├── logger/          # 日志相关
│   └── utils/           # 工具函数
└── docs/                # 文档目录
```

## 快速开始

1. 安装依赖：
```bash
go mod download
```

2. 运行服务：
```bash
go run cmd/api/main.go
```

服务将在 http://localhost:8080 启动

## API 测试

测试服务是否正常运行：
```bash
curl http://localhost:8080/ping
```

# notify-template
通知模板，支持email、sms、push、whatapp，和其他通知方式模板内容的自定义扩展
