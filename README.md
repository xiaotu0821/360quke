# Quake 资产测绘工具

基于 Go 语言开发的命令行工具，用于调用 360 Quake 资产测绘平台。

## 功能特性

- **双模式查询**：支持 Cookie 模式和 API 模式
- **命令行操作**：简单易用的 CLI 界面，Windows/Linux/macOS 完美兼容
- **数据导出**：支持 CSV、JSON、Excel 三种格式

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

### 查看帮助

```bash
quake-gui help
```

### 配置

```bash
# 配置API模式
quake-gui config https://quake.chaitin.com your_api_key

# 配置Cookie模式
quake-gui config cookie "session=xxx;token=xxx"

# 查看当前配置
quake-gui config
```

### 测试连接

```bash
quake-gui test
```

### 搜索

```bash
quake-gui search "port: 80"
quake-gui search "protocol: https"
quake-gui search "ip: 192.168.1.0/24"
```

### 导出数据

搜索后可以导出结果：

```bash
quake-gui export csv     # 导出为CSV
quake-gui export json    # 导出为JSON
quake-gui export xlsx   # 导出为Excel
```

## 查询语法

支持 Quake 查询语法，常用示例：

```
# 查询特定端口
port: 80
port: 443

# 查询特定协议
protocol: http
protocol: ssh

# 查询IP段
ip: 192.168.1.0/24

# 组合查询
port: 443 AND protocol: https

# 查询标题
title: nginx

# 查询国家
country: CN
```

## 配置说明

配置文件保存在：
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
└── quake-gui           # 编译后的可执行文件
```

## 依赖

- [excelize](https://github.com/xuri/excelize) - Excel 操作库

## 许可证

MIT License
