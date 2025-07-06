# 敏感词功能实现说明

## 功能概述

本项目实现了基于AC自动机的敏感词检测和管理系统，包含以下功能：

1. **敏感词初始化**: 从文件加载敏感词到数据库
2. **敏感词管理**: 实时添加、删除敏感词
3. **敏感词检测**: 检测文本中是否包含敏感词
4. **敏感词替换**: 将敏感词替换为指定字符
5. **令牌桶限流**: 对验证码发送接口进行限流保护

## 实现细节

### 1. 数据库模型

**文件**: `server/internal/domain/model/sensitive_word.go`

```go
type SensitiveWord struct {
    ID        uint      `json:"id" gorm:"primarykey"`
    Word      string    `json:"word" gorm:"type:varchar(100);not null;uniqueIndex;comment:敏感词"`
    Status    uint      `json:"status" gorm:"type:tinyint;default:1;comment:状态 1:启用 0:禁用"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 2. AC自动机实现

**文件**: `server/internal/initialize/sensitive.go`

- 实现了完整的AC自动机算法
- 支持敏感词的添加、删除、检测、替换
- 使用Trie树构建敏感词字典
- 支持失败指针优化

### 3. 服务层

**文件**: `server/internal/service/sensitive_word.go`

主要功能：
- `InitSensitiveFilter()`: 初始化敏感词过滤器
- `LoadSensitiveWordsFromFile()`: 从文件加载敏感词到数据库
- `AddSensitiveWord()`: 添加敏感词
- `DeleteSensitiveWord()`: 删除敏感词
- `GetSensitiveWordList()`: 获取敏感词列表
- `CheckSensitiveWord()`: 检测敏感词
- `ReplaceSensitiveWord()`: 替换敏感词

### 4. API接口

**文件**: `server/internal/api/v1/sensitive_word.go`

提供以下REST API：

- `POST /api/v1/sensitive-word/add` - 添加敏感词
- `DELETE /api/v1/sensitive-word/delete` - 删除敏感词
- `GET /api/v1/sensitive-word/list` - 获取敏感词列表
- `POST /api/v1/sensitive-word/check` - 检测敏感词
- `POST /api/v1/sensitive-word/replace` - 替换敏感词

### 5. 令牌桶限流

**文件**: `server/internal/middleware/rate_limit.go`

实现了令牌桶算法限流：
- 支持全局限流配置
- 针对邮箱验证码发送的专门限流
- 可配置的容量、速率、超时时间

### 6. 配置扩展

**文件**: `server/internal/config/security.go`

扩展了Security配置结构：

```go
type RateLimit struct {
    Enabled  bool  `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
    Capacity int   `mapstructure:"capacity" json:"capacity" yaml:"capacity"`
    Rate     int   `mapstructure:"rate" json:"rate" yaml:"rate"`
    Timeout  int64 `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
}
```

## 使用示例

### 1. 检测敏感词

```bash
curl -X POST http://localhost:9000/api/v1/sensitive-word/check \
  -H "Content-Type: application/json" \
  -d '{"text": "这是一个测试文本"}'
```

响应：
```json
{
  "code": 200,
  "data": {
    "has_sensitive": true,
    "words": ["测试"],
    "first_word": "测试"
  }
}
```

### 2. 替换敏感词

```bash
curl -X POST http://localhost:9000/api/v1/sensitive-word/replace \
  -H "Content-Type: application/json" \
  -d '{"text": "这是一个测试文本", "replace_by": "*"}'
```

响应：
```json
{
  "code": 200,
  "data": {
    "original_text": "这是一个测试文本",
    "replaced_text": "这是一个**文本",
    "replaced_by": "*"
  }
}
```

### 3. 添加敏感词

```bash
curl -X POST http://localhost:9000/api/v1/sensitive-word/add \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"word": "新敏感词"}'
```

### 4. 获取敏感词列表

```bash
curl -X GET "http://localhost:9000/api/v1/sensitive-word/list?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 配置文件

在 `server/conf/application.dev.yaml` 中添加限流配置：

```yaml
security:
  rate_limit:
    enabled: true
    capacity: 60
    rate: 1
    timeout: 60
```

## 初始化流程

1. 应用启动时，`main.go` 会调用 `service.InitSensitiveFilter()` 初始化敏感词过滤器
2. 从数据库加载所有启用的敏感词
3. 构建AC自动机的Trie树和失败指针
4. 从文件 `data/sensitive/dic.txt` 加载敏感词到数据库（如果不存在）

## 性能特点

- **时间复杂度**: O(n+m+k)，其中n是文本长度，m是敏感词总长度，k是匹配的敏感词数量
- **空间复杂度**: O(m)，其中m是敏感词总长度
- **实时性**: 支持动态添加和删除敏感词，无需重启服务
- **准确性**: 基于AC自动机算法，支持多模式匹配

## 安全特性

- 敏感词管理接口需要管理员权限
- 检测和替换接口为公开接口，无需认证
- 邮箱验证码发送接口集成限流保护
- 支持敏感词的启用/禁用状态管理 