# Changelog

所有对此项目的重要更改都将记录在此文件中。

此格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
并且本项目遵循 [Semantic Versioning](https://semver.org/spec/v2.0.0.html)。

## [Unreleased]

### Added
- 快站管理命令行工具初始版本
- 支持创建、管理和部署快站站点
- 站点管理功能：创建、列表、升级、更新站点信息
- 页面管理功能：创建、更新、删除页面
- 批量操作支持：批量创建站点和页面
- 内容上传功能：从源站URL抓取内容并上传到快站
- 域名管理功能：随机更换站点域名
- 美观的表格展示站点和页面信息
- 自动开通商业套餐功能
- 支持多平台构建（Linux、Windows、macOS）
- Docker 镜像支持
- Homebrew 包管理器支持
- Linux 包管理器支持（deb、rpm、apk）

### Commands
- `kuanzhan create-site` - 创建站点
- `kuanzhan list` - 查看站点列表
- `kuanzhan upload` - 上传站点内容
- `kuanzhan update` - 更新页面名称
- `kuanzhan delete-page` - 删除页面
- `kuanzhan upgrade` - 升级站点套餐
- `kuanzhan change-domain` - 更换域名
- `kuanzhan update-site` - 更新站点信息

### Dependencies
- github.com/spf13/cobra - 命令行框架
- github.com/spf13/viper - 配置管理
- github.com/olekukonko/tablewriter - 表格输出
- golang.org/x/net/html - HTML 解析

<!-- 
## [1.0.0] - 2024-XX-XX

### Added
- 初始版本发布

-->

[Unreleased]: https://github.com/your-org/kuanzhan/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/your-org/kuanzhan/releases/tag/v1.0.0 