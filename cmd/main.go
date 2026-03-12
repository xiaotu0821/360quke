package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/xuri/excelize/v2"
)

type Config struct {
	APIUrl string `json:"api_url"`
	APIKey string `json:"api_key"`
}

type AssetRow struct {
	IP       string
	Port     string
	Protocol string
	Title    string
	Country  string
	City     string
	ASN      string
	Org      string
}

var config Config
var configPath string
var results []AssetRow

func main() {
	configDir := filepath.Join(os.Getenv("APPDATA"), "QuakeGUI")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Printf("创建配置目录失败: %v", err)
	}
	configPath = filepath.Join(configDir, "config.json")

	loadConfig()

	app := tview.NewApplication()

	// 顶部标题
	title := tview.NewTextView()
	title.SetText(" ========== Quake 资产测绘工具 ========== ")
	title.SetTextAlign(tview.AlignCenter)
	title.SetTextColor(tcell.ColorGreen)

	// 模式选择
	cookieBtn := tview.NewButton(" [1] Cookie模式 ")
	apiBtn := tview.NewButton(" [2] API模式 ")

	// Cookie输入
	cookieLabel := tview.NewTextView()
	cookieLabel.SetText(" Cookie: ")
	cookieInput := tview.NewInputField()
	cookieInput.SetPlaceholder("粘贴Cookie...")

	// API输入
	apiLabel := tview.NewTextView()
	apiLabel.SetText(" API地址: ")
	apiUrlInput := tview.NewInputField()
	apiUrlInput.SetPlaceholder("https://quake.chaitin.com")
	apiUrlInput.SetText(config.APIUrl)

	apiKeyLabel := tview.NewTextView()
	apiKeyLabel.SetText(" API Key: ")
	apiKeyInput := tview.NewInputField()
	apiKeyInput.SetPlaceholder("输入API Key")
	apiKeyInput.SetText(config.APIKey)
	apiKeyInput.SetMaskCharacter('*')

	// 功能按钮
	testBtn := tview.NewButton(" [T] 测试连接 ")
	saveBtn := tview.NewButton(" [S] 保存配置 ")

	// 查询
	queryLabel := tview.NewTextView()
	queryLabel.SetText(" 查询语句: ")
	queryInput := tview.NewInputField()
	queryInput.SetPlaceholder("输入Quake查询语法，如: port: 80")
	searchBtn := tview.NewButton(" [回车] 搜索 ")

	// 结果表格
	resultTable := tview.NewTable()
	resultTable.SetBorder(true).SetTitle("查询结果")

	headers := []string{"IP", "端口", "协议", "标题", "国家", "城市", "ASN", "组织"}
	for col, header := range headers {
		cell := tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).
			SetSelectable(false)
		resultTable.SetCell(0, col, cell)
	}
	resultTable.SetSelectable(true, false)

	// 导出按钮
	exportCSVBtn := tview.NewButton(" [C] 导出CSV ")
	exportJSONBtn := tview.NewButton(" [J] 导出JSON ")
	exportExcelBtn := tview.NewButton(" [E] 导出Excel ")

	// 状态栏
	statusBar := tview.NewTextView()
	statusBar.SetText(" 提示: 使用 Tab 切换焦点 | Esc 退出 ")
	statusBar.SetTextColor(tcell.ColorDarkGray)

	// 切换模式
	showAPI := false

	cookieModeView := tview.NewFlex().SetDirection(tview.FlexRow)
	cookieModeView.AddItem(cookieLabel, 1, 0, false)
	cookieModeView.AddItem(cookieInput, 1, 0, true)

	apiModeView := tview.NewFlex().SetDirection(tview.FlexRow)
	apiModeView.AddItem(apiLabel, 1, 0, false)
	apiModeView.AddItem(apiUrlInput, 1, 0, true)
	apiModeView.AddItem(apiKeyLabel, 1, 0, false)
	apiModeView.AddItem(apiKeyInput, 1, 0, true)

	settingsView := tview.NewPages()
	settingsView.AddPage("cookie", cookieModeView, true, true)
	settingsView.AddPage("api", apiModeView, true, false)

	// 按钮行
	btnRow := tview.NewFlex()
	btnRow.AddItem(cookieBtn, 0, 1, false)
	btnRow.AddItem(apiBtn, 0, 1, false)
	btnRow.AddItem(testBtn, 0, 1, false)
	btnRow.AddItem(saveBtn, 0, 1, false)

	// 查询行
	queryRow := tview.NewFlex()
	queryRow.AddItem(queryLabel, 10, 0, false)
	queryRow.AddItem(queryInput, 0, 1, true)
	queryRow.AddItem(searchBtn, 0, 1, false)

	// 导出行
	exportRow := tview.NewFlex()
	exportRow.AddItem(exportCSVBtn, 0, 1, false)
	exportRow.AddItem(exportJSONBtn, 0, 1, false)
	exportRow.AddItem(exportExcelBtn, 0, 1, false)

	// 主布局
	mainView := tview.NewFlex().SetDirection(tview.FlexRow)
	mainView.AddItem(title, 1, 0, false)
	mainView.AddItem(btnRow, 1, 0, false)
	mainView.AddItem(settingsView, 4, 0, false)
	mainView.AddItem(queryRow, 1, 0, false)
	mainView.AddItem(resultTable, 0, 1, true)
	mainView.AddItem(exportRow, 1, 0, false)
	mainView.AddItem(statusBar, 1, 0, false)

	// 按钮事件
	cookieBtn.SetSelectedFunc(func() {
		settingsView.SwitchToPage("cookie")
		showAPI = false
		setStatus(statusBar, "当前模式: Cookie模式")
	})

	apiBtn.SetSelectedFunc(func() {
		settingsView.SwitchToPage("api")
		showAPI = true
		setStatus(statusBar, "当前模式: API模式")
	})

	testBtn.SetSelectedFunc(func() {
		if !showAPI {
			if cookieInput.GetText() == "" {
				setStatus(statusBar, "错误: 请输入Cookie")
				return
			}
			setStatus(statusBar, "Cookie模式: 连接成功(模拟)")
		} else {
			if apiUrlInput.GetText() == "" || apiKeyInput.GetText() == "" {
				setStatus(statusBar, "错误: 请输入API地址和Key")
				return
			}
			setStatus(statusBar, "API模式: 连接成功(模拟)")
		}
	})

	saveBtn.SetSelectedFunc(func() {
		config.APIUrl = apiUrlInput.GetText()
		config.APIKey = apiKeyInput.GetText()
		data, _ := json.Marshal(config)
		os.WriteFile(configPath, data, 0644)
		setStatus(statusBar, "配置已保存到: "+configPath)
	})

	searchBtn.SetSelectedFunc(func() {
		query := queryInput.GetText()
		if query == "" {
			setStatus(statusBar, "错误: 请输入查询语句")
			return
		}

		apiUrl := apiUrlInput.GetText()
		cookie := cookieInput.GetText()

		if apiUrl == "" && cookie == "" {
			setStatus(statusBar, "错误: 请先配置API地址或Cookie")
			return
		}

		// 模拟搜索结果
		results = []AssetRow{
			{IP: "192.168.1.1", Port: "80", Protocol: "HTTP", Title: "Test Web Server", Country: "CN", City: "Beijing", ASN: "AS4808", Org: "China Unicom"},
			{IP: "192.168.1.2", Port: "443", Protocol: "HTTPS", Title: "Nginx Server", Country: "CN", City: "Shanghai", ASN: "AS4134", Org: "China Telecom"},
			{IP: "192.168.1.3", Port: "22", Protocol: "SSH", Title: "SSH Service", Country: "US", City: "New York", ASN: "AS7922", Org: "Comcast"},
			{IP: "10.0.0.1", Port: "3389", Protocol: "RDP", Title: "Windows RDP", Country: "CN", City: "Guangzhou", ASN: "AS58453", Org: "China Mobile"},
			{IP: "172.16.0.1", Port: "3306", Protocol: "MySQL", Title: "Database", Country: "US", City: "Los Angeles", ASN: "AS21501", Org: "Cloudflare"},
		}

		// 清除旧数据
		for row := resultTable.GetRowCount() - 1; row > 0; row-- {
			resultTable.RemoveRow(row)
		}

		// 填充新数据
		for i, asset := range results {
			r := i + 1
			cells := []string{asset.IP, asset.Port, asset.Protocol, asset.Title, asset.Country, asset.City, asset.ASN, asset.Org}
			for col, text := range cells {
				cell := tview.NewTableCell(text).SetAlign(tview.AlignLeft)
				resultTable.SetCell(r, col, cell)
			}
		}

		setStatus(statusBar, fmt.Sprintf("查询成功，共 %d 条结果", len(results)))
	})

	// 导出功能
	exportCSVBtn.SetSelectedFunc(func() {
		if len(results) == 0 {
			setStatus(statusBar, "错误: 没有可导出的数据")
			return
		}
		err := exportToCSV(results, "quake_results.csv")
		if err != nil {
			setStatus(statusBar, fmt.Sprintf("导出失败: %v", err))
		} else {
			setStatus(statusBar, "已导出到: quake_results.csv")
		}
	})

	exportJSONBtn.SetSelectedFunc(func() {
		if len(results) == 0 {
			setStatus(statusBar, "错误: 没有可导出的数据")
			return
		}
		err := exportToJSON(results, "quake_results.json")
		if err != nil {
			setStatus(statusBar, fmt.Sprintf("导出失败: %v", err))
		} else {
			setStatus(statusBar, "已导出到: quake_results.json")
		}
	})

	exportExcelBtn.SetSelectedFunc(func() {
		if len(results) == 0 {
			setStatus(statusBar, "错误: 没有可导出的数据")
			return
		}
		err := exportToExcel(results, "quake_results.xlsx")
		if err != nil {
			setStatus(statusBar, fmt.Sprintf("导出失败: %v", err))
		} else {
			setStatus(statusBar, "已导出到: quake_results.xlsx")
		}
	})

	app.SetRoot(mainView, true)

	// 键盘事件
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			app.Stop()
		}
		return event
	})

	setStatus(statusBar, "就绪 - 请选择模式并输入查询语句")

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func loadConfig() {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return
	}
	json.Unmarshal(data, &config)
}

func setStatus(statusBar *tview.TextView, text string) {
	statusBar.SetText(" " + text + " ")
	statusBar.SetTextColor(tcell.ColorWhite)
}

func exportToCSV(results []AssetRow, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"IP", "端口", "协议", "标题", "国家", "城市", "ASN", "组织"}
	writer.Write(header)

	for _, asset := range results {
		row := []string{asset.IP, asset.Port, asset.Protocol, asset.Title, asset.Country, asset.City, asset.ASN, asset.Org}
		writer.Write(row)
	}
	return nil
}

func exportToJSON(results []AssetRow, filename string) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func exportToExcel(results []AssetRow, filename string) error {
	f := excelize.NewFile()
	defer f.Close()

	headers := []string{"IP", "端口", "协议", "标题", "国家", "城市", "ASN", "组织"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue("Sheet1", cell, header)
	}

	for idx, asset := range results {
		row := idx + 2
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), asset.IP)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), asset.Port)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), asset.Protocol)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), asset.Title)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), asset.Country)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", row), asset.City)
		f.SetCellValue("Sheet1", fmt.Sprintf("G%d", row), asset.ASN)
		f.SetCellValue("Sheet1", fmt.Sprintf("H%d", row), asset.Org)
	}

	return f.SaveAs(filename)
}
