# Quake 资产测绘工具

基于 Go 语言开发的终端 UI (TUI) 工具，用于调用 360 Quake 资产测绘平台。

## 功能特性

- **双模式查询**：支持 Cookie 模式和 API 模式
- **结果展示**：表格形式展示查询结果
- **数据导出**：支持 CSV、JSON、Excel 三种格式
- **配置保存**：API 配置本地持久化

## 环境要求

- Go 1.18+
- Windows/Linux/macOS

## 安装

```bash
# 克隆项目
git clone https://github.com/xiaotu0821/360quke.git

# 进入目录
cd 360quke

# 编译
go build -o quake-gui ./cmd
```

## 使用方法

### 1. 运行程序

```bash
./quake-gui
```

### 2. 模式切换

- 使用 `Cookie模式` 按钮切换到 Cookie 模式
- 使用 `API模式` 按钮切换到 API 模式

### 3. Cookie 模式

1. 从浏览器开发者工具中复制 Cookie
2. 粘贴到 Cookie 输入框
3. 点击「测试连接」验证
4. 输入查询语句并点击「搜索」

### 4. API 模式

1. 在「API地址」输入 Quake API 地址（如 `https://quake.chaitin.com`）
2. 在「API Key」输入您的 API Key
3. 点击「保存设置」保存配置（可选）
4. 点击「测试连接」验证
5. 输入查询语句并点击「搜索」

### 5. 查询语法

支持 Quake 查询语法，常用示例：

```
# 查询特定端口
port: 80

# 查询特定协议
protocol: http

# 查询IP段
ip: 192.168.1.0/24

# 组合查询
port: 443 AND protocol: https
```

### 6. 导出数据

搜索完成后，可使用底部按钮导出数据：
- 「导出CSV」
- 「导出JSON」
- 「导出Excel」

## 快捷键

- `Esc` - 退出程序

## 配置说明

API 配置保存在：
- Windows: `%APPDATA%\QuakeGUI\config.json`
- Linux/macOS: `~/.config/QuakeGUI/config.json`

## 项目结构

```
.
├── cmd/
│   └── main.go          # 主程序入口
├── quakeclient/
│   ├── client.go        # Quake API 客户端
│   └── exporter.go     # 导出功能
├── go.mod              # 依赖管理
└── quake-gui          # 编译后的可执行文件
```

## 依赖

- [tview](https://github.com/rivo/tview) - 终端 UI 库
- [tcell](https://github.com/gdamore/tcell) - 终端渲染库
- [excelize](https://github.com/xuri/excelize) - Excel 操作库

## 许可证

MIT License
