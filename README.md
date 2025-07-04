# Kuanzhan CLI 工具

Kuanzhan CLI 是一个快站管理的命令行工具，用于快速创建、管理和部署快站站点。

## 功能特性

- 🚀 **站点管理**: 创建、列表、升级、更新站点信息
- 📄 **页面管理**: 创建、更新、删除页面
- 🔄 **批量操作**: 支持批量创建站点和页面
- 📊 **表格展示**: 美观的表格显示站点和页面信息
- 🌐 **内容上传**: 从源站URL抓取内容并上传到快站
- 📱 **套餐管理**: 自动开通商业套餐
- 🔧 **域名管理**: 随机更换站点域名

## 安装

### 从源码编译

```bash
git clone <repository-url>
cd kuanzhan
go build -o kuanzhan ./cmd/kuanzhan
```

### 使用 Go 安装

```bash
go install pkg.blksails.net/kuanzhan/cmd/kuanzhan@latest
```

## 配置

在项目根目录创建 `kuanzhan.yaml` 配置文件：

```yaml
app_key: "your_app_key"
app_secret: "your_app_secret"
```

## 使用方法

### 基本语法

```bash
kuanzhan [command] [flags]
```

### 全局参数

- `-d, --debug`: 开启调试模式

## 命令说明

### 1. 创建站点

创建指定数量的站点并自动开通商业套餐。

```bash
kuanzhan create-site [flags]
```

**参数**:
- `-s, --size`: 创建站点数量 (必需)
- `-n, --name`: 站点名称 (必需)
- `-t, --type`: 站点类型 (默认: "FAST")
- `-b, --business-type`: 商业套餐类型 (默认: "SITE_EXCLUSIVE_YEAR")

**示例**:
```bash
kuanzhan create-site --size 5 --name "我的站点" --type FAST --business-type SITE_EXCLUSIVE_YEAR
```

### 2. 站点列表

查看站点列表及其页面信息。显示的表格列包括：站点ID、站点名称、站点域名、套餐类型、站点状态、页面ID、页面名称、页面URL。

```bash
kuanzhan list [flags]
```

**参数**:
- `-o, --only-site`: 只显示站点信息，不显示页面信息

**示例**:
```bash
# 显示所有站点和页面信息
kuanzhan list

# 只显示站点信息
kuanzhan list --only-site
```

### 3. 上传站点

从源站URL抓取内容并上传到指定站点。可以指定现有页面ID或创建新页面。

```bash
kuanzhan upload [flags]
```

**参数**:
- `-s, --source-url`: 源站URL (必需)
- `-i, --site-ids`: 站点ID列表 (必需)
- `-p, --page`: 创建页面数量 (默认: 1)
- `-t, --tpl`: 页面模板 (默认: "WHITE")
- `-n, --name`: 页面名称 (必需)
- `-g, --page-ids`: 指定页面ID列表，如果指定则更新现有页面而不是创建新页面

**示例**:
```bash
# 上传到新页面
kuanzhan upload \
  --source-url "https://example.com" \
  --site-ids 123,456,789 \
  --page 3 \
  --tpl WHITE \
  --name "首页"

# 上传到指定的现有页面
kuanzhan upload \
  --source-url "https://example.com" \
  --site-ids 123,456,789 \
  --page-ids 111,222,333 \
  --name "首页"
```

### 4. 更新页面

更新指定页面的名称。

```bash
kuanzhan update [flags]
```

**参数**:
- `-n, --name`: 页面名称 (必需)
- `-i, --page-ids`: 页面ID列表

**示例**:
```bash
kuanzhan update --name "新页面名称" --page-ids 123,456,789
```

### 5. 删除页面

删除指定的页面。

```bash
kuanzhan delete-page [flags]
```

**参数**:
- `-i, --page-ids`: 页面ID列表

**示例**:
```bash
kuanzhan delete-page --page-ids 123,456,789
```

### 6. 升级站点

升级站点套餐。

```bash
kuanzhan upgrade [flags]
```

**参数**:
- `-b, --business-type`: 商业套餐类型 (默认: "SITE_EXCLUSIVE_YEAR")
- `-i, --site-ids`: 站点ID列表 (必需)

**示例**:
```bash
kuanzhan upgrade --business-type SITE_EXCLUSIVE_YEAR --site-ids 123,456,789
```

### 7. 更换域名

为指定站点更换随机域名。

```bash
kuanzhan change-domain [flags]
```

**参数**:
- `-i, --site-ids`: 站点ID列表 (必需)

**示例**:
```bash
kuanzhan change-domain --site-ids 123,456,789
```

### 8. 更新站点信息

更新站点名称等信息。

```bash
kuanzhan update-site [flags]
```

**参数**:
- `-n, --name`: 站点名称 (必需)
- `-i, --site-ids`: 站点ID列表 (必需)

**示例**:
```bash
kuanzhan update-site --name "新站点名称" --site-ids 123,456,789
```

## 使用示例

### 完整工作流程

1. **创建站点**:
```bash
kuanzhan create-site --size 3 --name "测试站点"
```

2. **查看站点列表**:
```bash
kuanzhan list
```

3. **上传内容到站点**:
```bash
kuanzhan upload \
  --source-url "https://example.com" \
  --site-ids 123,456,789 \
  --page 2 \
  --name "首页"
```

4. **更新页面名称**:
```bash
kuanzhan update --name "新首页" --page-ids 111,222,333
```

5. **更新站点信息**:
```bash
kuanzhan update-site --name "新站点名称" --site-ids 123,456,789
```

6. **更换域名**:
```bash
kuanzhan change-domain --site-ids 123,456,789
```

### 批量操作示例

```bash
# 批量创建10个站点
kuanzhan create-site --size 10 --name "批量站点"

# 批量上传内容到多个站点
kuanzhan upload \
  --source-url "https://template.com" \
  --site-ids 1,2,3,4,5,6,7,8,9,10 \
  --page 5 \
  --name "模板页面"

# 批量更新站点名称
kuanzhan update-site --name "统一站点名称" --site-ids 1,2,3,4,5,6,7,8,9,10

# 批量更换域名
kuanzhan change-domain --site-ids 1,2,3,4,5,6,7,8,9,10
```

### 高级用法

#### 更新现有页面内容

如果你知道页面ID，可以直接更新现有页面而不是创建新页面：

```bash
kuanzhan upload \
  --source-url "https://newcontent.com" \
  --site-ids 123,456 \
  --page-ids 111,222,333,444 \
  --name "更新后的页面"
```

#### 只显示站点信息

当你只需要查看站点基本信息而不需要页面详情时：

```bash
kuanzhan list --only-site
```

## 错误处理

- 确保配置文件 `kuanzhan.yaml` 存在且包含正确的 `app_key` 和 `app_secret`
- 检查网络连接是否正常
- 验证站点ID和页面ID是否有效
- 使用 `--debug` 参数查看详细的调试信息

## 开发

### 项目结构

```
kuanzhan/
├── cmd/kuanzhan/          # CLI 入口
│   └── main.go
├── client.go              # 快站 API 客户端
├── client_test.go         # 客户端测试
├── impl.go                # API 实现
├── go.mod                 # Go 模块文件
├── go.sum                 # Go 依赖锁文件
└── README.md              # 项目文档
```

### 依赖

- `github.com/spf13/cobra`: 命令行框架
- `github.com/spf13/viper`: 配置管理
- `github.com/olekukonko/tablewriter`: 表格输出
- `golang.org/x/net/html`: HTML 解析
- `github.com/go-viper/mapstructure/v2`: 数据结构映射

### 构建

```bash
go build -o kuanzhan ./cmd/kuanzhan
```

### 测试

```bash
go test ./...
```

## 许可证

[许可证信息]

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目。

## 联系方式

如有问题，请联系开发团队或提交 Issue。 