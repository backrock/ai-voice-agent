# AI Voice Agent - 智能语音聊天系统

一个完整的语音聊天系统，支持多种大模型（OpenAI、Claude、通义千问、文心一言等），提供Web H5客户端和Golang服务端。

## 📋 项目特性

### 客户端 (Vue3 + Vite + UioTos)
- 🎤 实时语音聊天
- 📝 支持多个聊天会话管理
- 🎨 微信小程序风格UI
- 📱 完全响应式设计
- 🔄 实时消息流式输出
- 🎯 简洁直观的用户界面

### 服务端 (Go)
- 🚀 高性能Go服务
- 🔌 Provider模式支持多种大模型
- ⚙️ 后台配置管理系统
- 🔐 API密钥管理
- 💾 会话数据持久化
- 📊 使用统计和日志

### 支持的模型
- **云端大模型**
  - OpenAI (GPT-4, GPT-3.5-turbo)
  - Anthropic Claude
  - 阿里通义千问
  - 百度文心一言
  - 讯飞星火
  - 智谱GLM
  - 本地Ollama

## 🚀 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+
- Docker (可选)

### 服务端启动

```bash
cd backend
go mod download
go run cmd/main.go
```

服务器将在 `http://localhost:8080` 启动

### 客户端启动

```bash
cd frontend
npm install
npm run dev
```

客户端将在 `http://localhost:5173` 启动

## 📦 项目结构

```
ai-voice-agent/
├── backend/              # Go服务端
│   ├── cmd/
│   │   └── main.go      # 主程序入口
│   ├── internal/
│   │   ├── api/         # API处理器
│   │   ├── models/      # 数据模型
│   │   ├── providers/   # LLM提供商
│   │   ├── service/     # 业务逻辑
│   │   └── storage/     # 数据存储
│   ├── config/          # 配置文件
│   ├── go.mod
│   └── Dockerfile
├── frontend/            # Vue3客户端
│   ├── src/
│   │   ├── components/  # 组件
│   │   ├── pages/       # 页面
│   │   ├── stores/      # 状态管理
│   │   ├── utils/       # 工具函数
│   │   └── App.vue
│   ├── vite.config.ts
│   └── package.json
├── deploy/              # 部署脚本
│   ├── docker-compose.yml
│   ├── nginx.conf
│   └── startup.sh
├── docs/                # 文档
│   ├── API.md          # API文档
│   ├── DEPLOY.md       # 部署指南
│   ├── DEVELOPMENT.md  # 开发指南
│   └── ARCHITECTURE.md # 架构设计
└── README.md
```

## 📖 文档

- [API文档](./docs/API.md) - 完整的API接口文档
- [部署指南](./docs/DEPLOY.md) - Docker和服务器部署指南
- [开发指南](./docs/DEVELOPMENT.md) - 开发环境设置和贡献指南
- [架构设计](./docs/ARCHITECTURE.md) - 系统架构和设计模式说明

## 🔧 配置

### 服务端配置文件 (`backend/config/app.yaml`)

```yaml
server:
  port: 8080
  host: 0.0.0.0

database:
  type: sqlite
  sqlite:
    path: ./data/app.db

providers:
  openai:
    enabled: true
    api_key: your_api_key_here
    base_url: https://api.openai.com/v1
  
  claude:
    enabled: true
    api_key: your_api_key_here
  
  ollama:
    enabled: true
    base_url: http://localhost:11434

logger:
  level: info
  format: json
```

## 🎯 API端点

### Chat API
- `POST /api/v1/chat/sessions` - 创建新聊天会话
- `GET /api/v1/chat/sessions` - 获取会话列表
- `DELETE /api/v1/chat/sessions/:id` - 删除会话
- `POST /api/v1/chat/messages` - 发送消息
- `GET /api/v1/chat/messages/:sessionId` - 获取消息历史

### Configuration API
- `GET /api/v1/admin/models` - 获取可用模型列表
- `GET /api/v1/admin/config` - 获取系统配置
- `PUT /api/v1/admin/config` - 更新系统配置
- `POST /api/v1/admin/providers/test` - 测试提供商连接

## 🐳 Docker部署

```bash
cd deploy
docker-compose up -d
```

## 📱 客户端使用

1. 打开应用后自动连接到服务器
2. 点击底部按钮开始录音
3. 说出你的问题
4. 等待AI响应
5. 点击左上角可以创建新的聊天会话

## 🔐 安全性

- API密钥加密存储
- JWT令牌认证
- CORS跨域配置
- 输入验证和XSS防护

## 📝 许可证

MIT License

## 👥 贡献

欢迎提交Issue和Pull Request！

## 📞 支持

如有问题，请提交Issue或参考[开发文档](./docs/DEVELOPMENT.md)
