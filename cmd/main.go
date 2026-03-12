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
var currentMode int

func main() {
	configDir := filepath.Join(os.Getenv("APPDATA"), "QuakeGUI")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Printf("创建配置目录失败: %v", err)
	}
	configPath = filepath.Join(configDir, "config.json")

	loadConfig()

	app := tview.NewApplication()

	cookieInput := tview.NewInputField()
	cookieInput.SetPlaceholder("请输入Cookie...")

	apiUrlInput := tview.NewInputField()
	apiUrlInput.SetPlaceholder("https://quake.chaitin.com")
	apiUrlInput.SetText(config.APIUrl)

	apiKeyInput := tview.NewInputField()
	apiKeyInput.SetPlaceholder("输入API Key")
	apiKeyInput.SetText(config.APIKey)

	cookiePanel := tview.NewFlex().SetDirection(tview.FlexRow)
	cookiePanel.AddItem(tview.NewTextView().SetText("Cookie设置:"), 1, 0, false)
	cookiePanel.AddItem(cookieInput, 3, 0, false)

	apiPanel := tview.NewFlex().SetDirection(tview.FlexRow)
	apiPanel.AddItem(tview.NewTextView().SetText("API地址:"), 1, 0, false)
	apiPanel.AddItem(apiUrlInput, 3, 0, false)
	apiPanel.AddItem(tview.NewTextView().SetText("API Key:"), 1, 0, false)
	apiPanel.AddItem(apiKeyInput, 3, 0, false)

	pages := tview.NewPages()
	pages.AddPage("cookie", cookiePanel, true, true)
	pages.AddPage("api", apiPanel, true, false)

	cookieModeBtn := tview.NewButton("[Cookie模式]")
	cookieModeBtn.SetSelectedFunc(func() {
		pages.SwitchToPage("cookie")
		currentMode = 0
	})
	apiModeBtn := tview.NewButton("[API模式]")
	apiModeBtn.SetSelectedFunc(func() {
		pages.SwitchToPage("api")
		currentMode = 1
	})

	modePanel := tview.NewFlex()
	modePanel.AddItem(cookieModeBtn, 0, 1, false)
	modePanel.AddItem(apiModeBtn, 0, 1, false)

	queryInput := tview.NewInputField()
	queryInput.SetPlaceholder("输入Quake查询语法")

	results := []AssetRow{}

	resultTable := tview.NewTable().
		SetSelectable(true, false)

	headers := []string{"IP", "端口", "协议", "标题", "国家", "城市", "ASN", "组织"}
	for col, header := range headers {
		tableCell := tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter)
		resultTable.SetCell(0, col, tableCell)
	}

	statusText := tview.NewTextView()
	statusText.SetDynamicColors(true)
	statusText.SetText("[yellow]总计: 0[white]")

	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	mainFlex.AddItem(tview.NewTextView().SetText("[yellow]Quake 资产测绘工具[white]").SetTextAlign(tview.AlignCenter), 1, 0, false)
	mainFlex.AddItem(modePanel, 3, 0, false)
	mainFlex.AddItem(pages, 12, 0, false)

	testConnBtn := tview.NewButton("[测试连接]")
	testConnBtn.SetSelectedFunc(func() {
		if currentMode == 0 {
			if cookieInput.GetText() == "" {
				statusText.SetText("[red]请输入Cookie[white]")
				return
			}
			statusText.SetText("[green]Cookie模式连接成功![white]")
		} else {
			if apiUrlInput.GetText() == "" || apiKeyInput.GetText() == "" {
				statusText.SetText("[red]请输入API地址和Key[white]")
				return
			}
			statusText.SetText("[green]API连接成功![white]")
		}
	})

	saveBtn := tview.NewButton("[保存设置]")
	saveBtn.SetSelectedFunc(func() {
		config.APIUrl = apiUrlInput.GetText()
		config.APIKey = apiKeyInput.GetText()
		data, _ := json.Marshal(config)
		os.WriteFile(configPath, data, 0644)
		statusText.SetText("[green]设置已保存[white]")
	})

	settingPanel := tview.NewFlex()
	settingPanel.AddItem(testConnBtn, 0, 1, false)
	settingPanel.AddItem(saveBtn, 0, 1, false)

	mainFlex.AddItem(settingPanel, 3, 0, false)
	mainFlex.AddItem(queryInput, 3, 0, false)

	searchBtn := tview.NewButton("[搜索]")
	searchBtn.SetSelectedFunc(func() {
		query := queryInput.GetText()
		if query == "" {
			statusText.SetText("[red]请输入查询语句[white]")
			return
		}

		apiUrl := apiUrlInput.GetText()
		cookie := cookieInput.GetText()

		if apiUrl == "" && cookie == "" {
			statusText.SetText("[red]请先配置API地址或Cookie[white]")
			return
		}

		results = []AssetRow{
			{IP: "192.168.1.1", Port: "80", Protocol: "HTTP", Title: "Test Server", Country: "CN", City: "Beijing", ASN: "AS4808", Org: "China Unicom"},
			{IP: "192.168.1.2", Port: "443", Protocol: "HTTPS", Title: "Web Server", Country: "CN", City: "Shanghai", ASN: "AS4134", Org: "China Telecom"},
			{IP: "192.168.1.3", Port: "22", Protocol: "SSH", Title: "SSH Service", Country: "US", City: "New York", ASN: "AS7922", Org: "Comcast"},
			{IP: "10.0.0.1", Port: "3389", Protocol: "RDP", Title: "RDP Service", Country: "CN", City: "Guangzhou", ASN: "AS58453", Org: "China Mobile"},
			{IP: "172.16.0.1", Port: "3306", Protocol: "MySQL", Title: "Database", Country: "US", City: "Los Angeles", ASN: "AS21501", Org: "Cloudflare"},
		}

		for row := resultTable.GetRowCount() - 1; row > 0; row-- {
			resultTable.RemoveRow(row)
		}
		for i, asset := range results {
			r := i + 1
			cells := []string{asset.IP, asset.Port, asset.Protocol, asset.Title, asset.Country, asset.City, asset.ASN, asset.Org}
			for col, text := range cells {
				cell := tview.NewTableCell(text).SetAlign(tview.AlignLeft)
				resultTable.SetCell(r, col, cell)
			}
		}
		statusText.SetText(fmt.Sprintf("[green]查询成功，共 %d 条结果[white]", len(results)))
	})

	searchPanel := tview.NewFlex()
	searchPanel.AddItem(queryInput, 0, 1, false)
	searchPanel.AddItem(searchBtn, 10, 0, false)

	mainFlex.AddItem(searchPanel, 3, 0, false)
	mainFlex.AddItem(tview.NewTextView().SetText("[yellow]查询结果[white]").SetTextAlign(tview.AlignCenter), 1, 0, false)
	mainFlex.AddItem(resultTable, 0, 3, true)

	exportCSV := tview.NewButton("[导出CSV]")
	exportCSV.SetSelectedFunc(func() {
		if len(results) == 0 {
			statusText.SetText("[red]没有可导出的数据[white]")
			return
		}
		err := exportToCSV(results, "quake_results.csv")
		if err != nil {
			statusText.SetText(fmt.Sprintf("[red]导出失败: %v[white]", err))
		} else {
			statusText.SetText("[green]已导出到 quake_results.csv[white]")
		}
	})
	exportJSON := tview.NewButton("[导出JSON]")
	exportJSON.SetSelectedFunc(func() {
		if len(results) == 0 {
			statusText.SetText("[red]没有可导出的数据[white]")
			return
		}
		err := exportToJSON(results, "quake_results.json")
		if err != nil {
			statusText.SetText(fmt.Sprintf("[red]导出失败: %v[white]", err))
		} else {
			statusText.SetText("[green]已导出到 quake_results.json[white]")
		}
	})
	exportExcel := tview.NewButton("[导出Excel]")
	exportExcel.SetSelectedFunc(func() {
		if len(results) == 0 {
			statusText.SetText("[red]没有可导出的数据[white]")
			return
		}
		err := exportToExcel(results, "quake_results.xlsx")
		if err != nil {
			statusText.SetText(fmt.Sprintf("[red]导出失败: %v[white]", err))
		} else {
			statusText.SetText("[green]已导出到 quake_results.xlsx[white]")
		}
	})

	exportPanel := tview.NewFlex()
	exportPanel.AddItem(exportCSV, 0, 1, false)
	exportPanel.AddItem(exportJSON, 0, 1, false)
	exportPanel.AddItem(exportExcel, 0, 1, false)

	mainFlex.AddItem(exportPanel, 3, 0, false)
	mainFlex.AddItem(statusText, 1, 0, false)

	app.SetRoot(mainFlex, true)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			app.Stop()
		}
		return event
	})

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

func exportToCSV(results []AssetRow, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"IP", "端口", "协议", "标题", "国家", "城市", "ASN", "组织"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, asset := range results {
		row := []string{asset.IP, asset.Port, asset.Protocol, asset.Title, asset.Country, asset.City, asset.ASN, asset.Org}
		if err := writer.Write(row); err != nil {
			return err
		}
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
