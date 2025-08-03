# Kuanzhan 自动升级功能

## 概述

Kuanzhan 现在支持自动升级功能，可以自动检查并下载最新版本。

## 使用方法

### 检查并升级到最新版本

```bash
kuanzhan upgrade
```

### 功能说明

1. **版本检查**: 自动从 GitHub Releases 获取最新版本信息
2. **系统兼容性**: 自动识别当前系统和 CPU 架构，查找匹配的二进制文件
3. **自动下载**: 下载适合当前系统的二进制文件
4. **自动解压**: 支持 `.tar.gz` 和 `.zip` 格式的压缩文件自动解压
5. **自动更新**: 使用 `github.com/minio/selfupdate` 库进行安全的自我更新

### 支持的系统和架构

- **操作系统**: Linux, macOS (darwin), Windows
- **CPU 架构**: amd64, arm64, arm, 386 等

### 二进制文件命名规则

系统会自动查找以下命名模式的二进制文件：

- `kuanzhan-{os}-{arch}` (例如: `kuanzhan-darwin-amd64`)
- `kuanzhan_{os}_{arch}` (例如: `kuanzhan_linux_arm64`)
- `kuanzhan-{os}-{arch}.exe` (Windows 系统)
- `kuanzhan_{os}_{arch}.exe` (Windows 系统)
- 支持压缩格式: `.tar.gz`, `.zip`

### 注意事项

1. **开发版本**: 如果当前版本是 "dev"，升级功能会被跳过
2. **权限要求**: 更新需要写入当前可执行文件的权限
3. **网络连接**: 需要网络连接来访问 GitHub API 和下载文件
4. **版本比较**: 只有当发现更新版本时才会进行升级

### 错误处理

- 如果找不到适合当前系统的二进制文件，会显示相应错误信息
- 如果下载或更新过程中出现错误，会显示详细的错误信息
- 如果当前版本已经是最新版本，会提示无需升级

### 示例输出

#### 直接二进制文件更新
```bash
$ kuanzhan upgrade
Latest version: v1.0.0
New version found: v1.0.0
Found compatible binary: kuanzhan-darwin-amd64
Downloading from: https://github.com/blksails/kuanzhan/releases/download/v1.0.0/kuanzhan-darwin-amd64
Updating executable: /usr/local/bin/kuanzhan
Successfully updated to version: v1.0.0
```

#### 压缩文件更新
```bash
$ kuanzhan upgrade
Latest version: v1.0.0
New version found: v1.0.0
Found compatible binary: kuanzhan-darwin-amd64.tar.gz
Downloading from: https://github.com/blksails/kuanzhan/releases/download/v1.0.0/kuanzhan-darwin-amd64.tar.gz
Detected tar.gz file, extracting...
Found binary in tar.gz: kuanzhan-darwin-amd64
Updating executable: /usr/local/bin/kuanzhan
Successfully updated to version: v1.0.0
```

或者如果已经是最新版本：

```bash
$ kuanzhan upgrade
Latest version: v0.0.12
You are using the latest version
```

## 技术实现

- 使用 GitHub Releases API 获取版本信息
- 使用 `runtime.GOOS` 和 `runtime.GOARCH` 获取系统信息
- 使用 `github.com/minio/selfupdate` 进行安全的自我更新
- 支持多种二进制文件命名约定
- 自动解压 `.tar.gz` 和 `.zip` 格式的压缩文件
- 智能识别压缩文件中的二进制文件 