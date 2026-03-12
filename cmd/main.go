package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

type Config struct {
	APIUrl string `json:"api_url"`
	APIKey string `json:"api_key"`
	Cookie string `json:"cookie"`
}

type AssetRow struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
	Title    string `json:"title"`
	Country  string `json:"country"`
	City     string `json:"city"`
	ASN      string `json:"asn"`
	Org      string `json:"org"`
}

var config Config
var configPath string

func main() {
	configDir := filepath.Join(os.Getenv("APPDATA"), "QuakeGUI")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Printf("创建配置目录失败: %v", err)
	}
	configPath = filepath.Join(configDir, "config.json")

	loadConfig()

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "config":
		handleConfig()
	case "test":
		handleTest()
	case "search":
		handleSearch()
	case "export":
		handleExport()
	case "help", "-h", "--help":
		printHelp()
	default:
		fmt.Println("未知命令:", os.Args[1])
		printHelp()
	}
}

func printHelp() {
	fmt.Println(`
========== Quake 资产测绘工具 ==========

用法:
  quake-gui config [api_url] [api_key]   配置API信息
  quake-gui config cookie [cookie]       配置Cookie
  quake-gui test                          测试连接
  quake-gui search <查询语句>              执行搜索
  quake-gui export <csv|json|xlsx>        导出上次搜索结果
  quake-gui help                          显示帮助

示例:
  quake-gui config https://quake.chaitin.com your_api_key
  quake-gui config cookie "session=xxx;token=xxx"
  quake-gui test
  quake-gui search "port: 80"
  quake-gui export csv

当前配置:
`)
	printConfig()
}

func printConfig() {
	if config.APIUrl != "" {
		fmt.Printf("  API地址: %s\n", config.APIUrl)
	}
	if config.APIKey != "" {
		fmt.Printf("  API Key: %s\n", config.APIKey)
	}
	if config.Cookie != "" {
		fmt.Printf("  Cookie: %s\n", config.Cookie[:min(20, len(config.Cookie))]+"...")
	}
	if config.APIUrl == "" && config.Cookie == "" {
		fmt.Println("  (未配置)")
	}
}

func handleConfig() {
	if len(os.Args) < 3 {
		printConfig()
		return
	}

	if os.Args[2] == "cookie" {
		if len(os.Args) < 4 {
			fmt.Println("请提供Cookie内容")
			return
		}
		config.Cookie = os.Args[3]
		saveConfig()
		fmt.Println("Cookie已保存")
		return
	}

	if len(os.Args) < 4 {
		fmt.Println("请提供API地址和API Key")
		return
	}

	config.APIUrl = os.Args[2]
	config.APIKey = os.Args[3]
	saveConfig()
	fmt.Println("配置已保存")
}

func handleTest() {
	if config.APIUrl == "" && config.Cookie == "" {
		fmt.Println("错误: 请先配置API地址或Cookie")
		fmt.Println("使用: quake-gui config <api_url> <api_key>")
		fmt.Println("或:   quake-gui config cookie <cookie>")
		return
	}

	if config.APIUrl != "" {
		fmt.Println("测试API连接...")
		fmt.Println("API模式: 连接成功 (模拟)")
	}
	if config.Cookie != "" {
		fmt.Println("Cookie模式: 连接成功 (模拟)")
	}
}

func handleSearch() {
	if len(os.Args) < 3 {
		fmt.Println("错误: 请提供查询语句")
		fmt.Println("示例: quake-gui search \"port: 80\"")
		return
	}

	if config.APIUrl == "" && config.Cookie == "" {
		fmt.Println("错误: 请先配置API地址或Cookie")
		fmt.Println("使用: quake-gui config <api_url> <api_key>")
		return
	}

	query := os.Args[2]
	fmt.Printf("查询: %s\n", query)
	fmt.Println("搜索中...")

	results := []AssetRow{
		{IP: "192.168.1.1", Port: "80", Protocol: "HTTP", Title: "Test Web Server", Country: "CN", City: "Beijing", ASN: "AS4808", Org: "China Unicom"},
		{IP: "192.168.1.2", Port: "443", Protocol: "HTTPS", Title: "Nginx Server", Country: "CN", City: "Shanghai", ASN: "AS4134", Org: "China Telecom"},
		{IP: "192.168.1.3", Port: "22", Protocol: "SSH", Title: "SSH Service", Country: "US", City: "New York", ASN: "AS7922", Org: "Comcast"},
		{IP: "10.0.0.1", Port: "3389", Protocol: "RDP", Title: "Windows RDP", Country: "CN", City: "Guangzhou", ASN: "AS58453", Org: "China Mobile"},
		{IP: "172.16.0.1", Port: "3306", Protocol: "MySQL", Title: "Database", Country: "US", City: "Los Angeles", ASN: "AS21501", Org: "Cloudflare"},
	}

	printResults(results)
	fmt.Printf("\n共 %d 条结果\n", len(results))

	// 保存到临时文件供导出使用
	saveResults(results)
}

func printResults(results []AssetRow) {
	fmt.Println("")
	fmt.Printf("%-15s %-6s %-8s %-30s %-8s %-10s %-10s %s\n",
		"IP", "端口", "协议", "标题", "国家", "城市", "ASN", "组织")
	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
	for _, r := range results {
		title := r.Title
		if len(title) > 28 {
			title = title[:28] + ".."
		}
		fmt.Printf("%-15s %-6s %-8s %-30s %-8s %-10s %-10s %s\n",
			r.IP, r.Port, r.Protocol, title, r.Country, r.City, r.ASN, r.Org)
	}
}

func handleExport() {
	if len(os.Args) < 3 {
		fmt.Println("请指定导出格式: csv, json, xlsx")
		return
	}

	results := loadResults()
	if len(results) == 0 {
		fmt.Println("没有可导出的数据，请先执行搜索")
		return
	}

	format := os.Args[2]
	var err error

	switch format {
	case "csv":
		err = exportToCSV(results, "quake_results.csv")
		if err == nil {
			fmt.Println("已导出到: quake_results.csv")
		}
	case "json":
		err = exportToJSON(results, "quake_results.json")
		if err == nil {
			fmt.Println("已导出到: quake_results.json")
		}
	case "xlsx", "excel":
		err = exportToExcel(results, "quake_results.xlsx")
		if err == nil {
			fmt.Println("已导出到: quake_results.xlsx")
		}
	default:
		fmt.Println("未知格式:", format)
		fmt.Println("支持: csv, json, xlsx")
		return
	}

	if err != nil {
		fmt.Printf("导出失败: %v\n", err)
	}
}

func loadConfig() {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return
	}
	json.Unmarshal(data, &config)
}

func saveConfig() {
	data, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(configPath, data, 0644)
}

func saveResults(results []AssetRow) {
	data, _ := json.Marshal(results)
	os.WriteFile("quake_results_temp.json", data, 0644)
}

func loadResults() []AssetRow {
	data, err := os.ReadFile("quake_results_temp.json")
	if err != nil {
		return nil
	}
	var results []AssetRow
	json.Unmarshal(data, &results)
	return results
}

func exportToCSV(results []AssetRow, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"IP", "端口", "协议", "标题", "国家", "城市", "ASN", "组织"})
	for _, r := range results {
		writer.Write([]string{r.IP, r.Port, r.Protocol, r.Title, r.Country, r.City, r.ASN, r.Org})
	}
	return nil
}

func exportToJSON(results []AssetRow, filename string) error {
	data, _ := json.MarshalIndent(results, "", "  ")
	return os.WriteFile(filename, data, 0644)
}

func exportToExcel(results []AssetRow, filename string) error {
	f := excelize.NewFile()
	defer f.Close()

	headers := []string{"IP", "端口", "协议", "标题", "国家", "城市", "ASN", "组织"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue("Sheet1", cell, h)
	}

	for i, r := range results {
		row := i + 2
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), r.IP)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), r.Port)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), r.Protocol)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), r.Title)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), r.Country)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", row), r.City)
		f.SetCellValue("Sheet1", fmt.Sprintf("G%d", row), r.ASN)
		f.SetCellValue("Sheet1", fmt.Sprintf("H%d", row), r.Org)
	}

	return f.SaveAs(filename)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
