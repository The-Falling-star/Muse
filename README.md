# Muse

🎭 一款精简高效的 AI 角色扮演聊天工具，兼容 SillyTavern 角色卡格式。

## ✨ 特性

- 🎴 **角色卡兼容** - 完整支持 SillyTavern 角色卡（PNG/JSON）、预设、世界书导入
- 🔄 **正则脚本** - 强大的文本处理能力，支持多种作用域
- 🌐 **多 API 支持** - OpenAI、Claude、Gemini 三大主流 AI 协议
- 🚀 **高性能** - Vue 3 + Go 架构，数据库存储，虚拟滚动
- 🎨 **现代 UI** - 精美的深色/浅色主题，流畅的交互体验
- 🔒 **代理支持** - 灵活的网络代理配置

## 🛠️ 技术栈

### 前端
- Vue 3 + TypeScript
- Vite
- Element Plus
- Pinia
- Vue Router

### 后端
- Go 1.24
- Gin
- GORM
- MySQL

## 🚀 快速开始

### 前端开发

```bash
cd frontend
npm install
npm run dev
```

### 后端开发

```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

## 📁 项目结构

```
Muse/
├── frontend/                # 前端项目
│   ├── src/
│   │   ├── api/            # API 请求
│   │   ├── assets/         # 静态资源
│   │   ├── components/     # 通用组件
│   │   ├── composables/    # 组合式函数
│   │   ├── layouts/        # 布局组件
│   │   ├── router/         # 路由配置
│   │   ├── stores/         # Pinia 状态管理
│   │   ├── types/          # TypeScript 类型定义
│   │   ├── utils/          # 工具函数
│   │   └── views/          # 页面组件
│   └── ...
├── backend/                 # 后端项目
│   ├── cmd/server/         # 主入口
│   ├── internal/           # 内部包
│   │   ├── config/         # 配置
│   │   ├── handler/        # HTTP 处理器
│   │   ├── middleware/     # 中间件
│   │   ├── model/          # 数据模型
│   │   ├── repository/     # 数据访问层
│   │   ├── router/         # 路由
│   │   └── service/        # 业务逻辑
│   └── pkg/                # 公共包
└── README.md
```

## 📝 开发规范

### 前端
- 组件式架构，单文件组件
- 统一使用小驼峰命名法
- 优先使用 Element Plus 组件
- TypeScript 严格模式

### 后端
- 遵循 Google Go 代码规范
- 分层架构：Handler -> Service -> Repository
- 统一错误处理
- 完善的注释文档

## 📄 License

MIT
