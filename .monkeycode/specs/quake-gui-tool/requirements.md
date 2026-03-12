# Requirements Document

## Introduction

开发一个基于Go语言的GUI工具，用于调用360 Quake资产测绘平台。支持两种模式：Cookie模式（模拟浏览器操作）和API模式（直接调用API）。

## Glossary

- **Quake**: 360旗下的网络空间资产测绘平台
- **Cookie模式**: 通过导入浏览器Cookie模拟用户登录进行查询
- **API模式**: 使用API Key直接调用Quake服务
- **GIUI**: Go语言的跨平台UI框架

## Requirements

### Requirement 1: 应用主窗口

**User Story:** 作为用户，我需要一个图形化界面来操作Quake工具

#### Acceptance Criteria

1. WHEN 用户启动应用, THE 应用 SHALL 显示主窗口
2. THE 主窗口 SHALL 包含标签页切换功能
3. THE 主窗口 SHALL 包含Cookie模式和API模式两个标签页

### Requirement 2: Cookie模式界面

**User Story:** 作为用户，我需要通过导入Cookie的方式查询Quake

#### Acceptance Criteria

1. THE Cookie输入区域 SHALL 提供文本框用于输入Cookie
2. THE Cookie输入区域 SHALL 提供文件导入按钮用于导入Cookie文件
3. WHEN 用户点击"测试连接"按钮, THE 应用 SHALL 验证Cookie有效性
4. THE 查询区域 SHALL 提供查询输入框（支持Quake查询语法）
5. THE 查询区域 SHALL 提供分页控件（每页数量、当前页码）
6. WHEN 用户点击"搜索"按钮, THE 应用 SHALL 执行查询并显示结果
7. THE 结果展示区域 SHALL 以表格形式显示结果
8. THE 结果展示区域 SHALL 支持滚动浏览
9. THE 导出功能 SHALL 支持导出为CSV格式
10. THE 导出功能 SHALL 支持导出为JSON格式
11. THE 导出功能 SHALL 支持导出为Excel格式

### Requirement 3: API模式界面

**User Story:** 作为用户，我需要通过API Key直接调用Quake

#### Acceptance Criteria

1. THE API设置区域 SHALL 提供API地址输入框
2. THE API设置区域 SHALL 提供API Key输入框
3. WHEN 用户点击"保存设置"按钮, THE 应用 SHALL 保存API配置到本地
4. WHEN 用户点击"测试连接"按钮, THE 应用 SHALL 验证API连接
5. THE 查询功能 SHALL 与Cookie模式一致
6. THE 导出功能 SHALL 与Cookie模式一致

### Requirement 4: 查询结果展示

**User Story:** 作为用户，我需要清晰地查看查询结果

#### Acceptance Criteria

1. THE 结果表格 SHALL 显示关键字段（IP、端口、协议、标题等）
2. THE 结果表格 SHALL 支持点击行查看详情
3. THE 详情面板 SHALL 显示完整的资产信息
4. THE 结果 SHALL 支持复制到剪贴板

### Requirement 5: 系统集成

**User Story:** 作为用户，我需要应用与系统更好地集成

#### Acceptance Criteria

1. THE 应用 SHALL 支持窗口最大化/最小化
2. THE 应用 SHALL 支持关闭按钮
3. THE 应用 SHALL 在窗口标题显示当前模式状态

## Technical Notes

- 使用gioui作为GUI框架
- 使用fyne作为备选GUI框架
- API调用使用HTTP请求
- 本地配置使用JSON文件存储
