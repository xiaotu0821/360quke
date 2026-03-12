package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
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
	os.MkdirAll(configDir, 0755)
	configPath = filepath.Join(configDir, "config.json")
	loadConfig()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/api/config", handleConfig)
	mux.HandleFunc("/api/save-config", handleSaveConfig)
	mux.HandleFunc("/api/test", handleTest)
	mux.HandleFunc("/api/search", handleSearch)
	mux.HandleFunc("/api/export-csv", handleExportCSV)
	mux.HandleFunc("/api/export-json", handleExportJSON)
	mux.HandleFunc("/api/export-excel", handleExportExcel)

	server := &http.Server{Addr: fmt.Sprintf("127.0.0.1:%d", port), Handler: mux}
	go server.ListenAndServe()

	url := fmt.Sprintf("http://127.0.0.1:%d", port)
	log.Printf("服务器启动: %s", url)

	openBrowser(url)
	log.Println("按 Ctrl+C 停止服务器")
	<-make(chan struct{})
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(htmlContent))
}

func handleConfig(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(config)
}

func handleSaveConfig(w http.ResponseWriter, r *http.Request) {
	var req struct {
		APIUrl string `json:"apiUrl"`
		APIKey string `json:"apiKey"`
		Cookie string `json:"cookie"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	config.APIUrl = req.APIUrl
	config.APIKey = req.APIKey
	config.Cookie = req.Cookie
	data, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(configPath, data, 0644)
	w.Write([]byte(`{"status":"ok"}`))
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		APIUrl string `json:"apiUrl"`
		APIKey string `json:"apiKey"`
		Cookie string `json:"cookie"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if req.APIUrl == "" && req.Cookie == "" {
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "请先配置API地址或Cookie"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "连接成功（模拟）"})
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Query  string `json:"query"`
		APIUrl string `json:"apiUrl"`
		APIKey string `json:"apiKey"`
		Cookie string `json:"cookie"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if req.Query == "" {
		json.NewEncoder(w).Encode([]AssetRow{})
		return
	}
	if req.APIUrl == "" && req.Cookie == "" {
		json.NewEncoder(w).Encode([]AssetRow{})
		return
	}

	results := []AssetRow{
		{IP: "192.168.1.1", Port: "80", Protocol: "HTTP", Title: "Test Web Server", Country: "CN", City: "Beijing", ASN: "AS4808", Org: "China Unicom"},
		{IP: "192.168.1.2", Port: "443", Protocol: "HTTPS", Title: "Nginx Server", Country: "CN", City: "Shanghai", ASN: "AS4134", Org: "China Telecom"},
		{IP: "192.168.1.3", Port: "22", Protocol: "SSH", Title: "SSH Service", Country: "US", City: "New York", ASN: "AS7922", Org: "Comcast"},
		{IP: "10.0.0.1", Port: "3389", Protocol: "RDP", Title: "Windows RDP", Country: "CN", City: "Guangzhou", ASN: "AS58453", Org: "China Mobile"},
		{IP: "172.16.0.1", Port: "3306", Protocol: "MySQL", Title: "Database", Country: "US", City: "Los Angeles", ASN: "AS21501", Org: "Cloudflare"},
	}

	json.NewEncoder(w).Encode(results)
}

func handleExportCSV(w http.ResponseWriter, r *http.Request) {
	var data []AssetRow
	json.NewDecoder(r.Body).Decode(&data)

	filename := getSavePath("csv")
	if filename == "" {
		w.Write([]byte(`{"status":"cancel"}`))
		return
	}

	file, _ := os.Create(filename)
	defer file.Close()
	fmt.Fprintln(file, "IP,端口,协议,标题,国家,城市,ASN,组织")
	for _, r := range data {
		fmt.Fprintf(file, "%s,%s,%s,%s,%s,%s,%s,%s\n", r.IP, r.Port, r.Protocol, r.Title, r.Country, r.City, r.ASN, r.Org)
	}
	w.Write([]byte(`{"status":"ok","filename":"` + filename + `"}`))
}

func handleExportJSON(w http.ResponseWriter, r *http.Request) {
	var data []AssetRow
	json.NewDecoder(r.Body).Decode(&data)

	filename := getSavePath("json")
	if filename == "" {
		w.Write([]byte(`{"status":"cancel"}`))
		return
	}

	jsonBytes, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile(filename, jsonBytes, 0644)
	w.Write([]byte(`{"status":"ok","filename":"` + filename + `"}`))
}

func handleExportExcel(w http.ResponseWriter, r *http.Request) {
	var data []AssetRow
	json.NewDecoder(r.Body).Decode(&data)

	filename := getSavePath("xlsx")
	if filename == "" {
		w.Write([]byte(`{"status":"cancel"}`))
		return
	}

	file, _ := os.Create(filename)
	defer file.Close()
	fmt.Fprintln(file, "IP,端口,协议,标题,国家,城市,ASN,组织")
	for _, r := range data {
		fmt.Fprintf(file, "%s,%s,%s,%s,%s,%s,%s,%s\n", r.IP, r.Port, r.Protocol, r.Title, r.Country, r.City, r.ASN, r.Org)
	}
	w.Write([]byte(`{"status":"ok","filename":"` + filename + `"}`))
}

func getSavePath(ext string) string {
	currentDir, _ := os.Getwd()
	defaultName := fmt.Sprintf("quake_results.%s", ext)
	filename := filepath.Join(currentDir, defaultName)
	return filename
}

func loadConfig() {
	data, _ := os.ReadFile(configPath)
	json.Unmarshal(data, &config)
}

func openBrowser(url string) {
	time.Sleep(500 * time.Millisecond)
	switch runtime.GOOS {
	case "windows":
		exec.Command("cmd", "/c", "start", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	default:
		exec.Command("xdg-open", url).Start()
	}
}

var htmlContent = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Quake 资产测绘工具</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif; background: #f0f2f5; min-height: 100vh; }
        .container { max-width: 1200px; margin: 0 auto; padding: 20px; }
        .header { background: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); margin-bottom: 20px; }
        .header h1 { color: #1890ff; font-size: 24px; }
        .tabs { display: flex; gap: 10px; margin-bottom: 20px; }
        .tab { padding: 10px 20px; background: #fff; border: 1px solid #d9d9d9; border-radius: 4px; cursor: pointer; transition: all 0.3s; }
        .tab:hover { border-color: #1890ff; color: #1890ff; }
        .tab.active { background: #1890ff; color: #fff; border-color: #1890ff; }
        .panel { display: none; background: #fff; padding: 24px; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
        .panel.active { display: block; }
        .form-group { margin-bottom: 16px; }
        .form-group label { display: block; margin-bottom: 8px; color: #333; font-weight: 500; }
        .form-group input, .form-group textarea { width: 100%; padding: 8px 12px; border: 1px solid #d9d9d9; border-radius: 4px; font-size: 14px; transition: border-color 0.3s; }
        .form-group input:focus, .form-group textarea:focus { outline: none; border-color: #40a9ff; }
        .form-group textarea { min-height: 80px; resize: vertical; }
        .btn { padding: 8px 20px; border: none; border-radius: 4px; cursor: pointer; font-size: 14px; transition: all 0.3s; }
        .btn-primary { background: #1890ff; color: #fff; }
        .btn-primary:hover { background: #40a9ff; }
        .btn-success { background: #52c41a; color: #fff; }
        .btn-success:hover { background: #73d13d; }
        .btn-default { background: #fff; border: 1px solid #d9d9d9; color: #333; }
        .btn-default:hover { border-color: #1890ff; color: #1890ff; }
        .btn-group { display: flex; gap: 12px; margin-top: 16px; }
        .query-box { background: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); margin-bottom: 20px; }
        .query-row { display: flex; gap: 12px; align-items: center; }
        .query-box input { flex: 1; }
        .results-box { background: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
        .results-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
        .results-header h3 { color: #333; }
        .results-table { width: 100%; border-collapse: collapse; }
        .results-table th, .results-table td { padding: 12px; text-align: left; border-bottom: 1px solid #f0f0f0; }
        .results-table th { background: #fafafa; font-weight: 500; color: #333; }
        .results-table tr:hover { background: #f5f5f5; }
        .results-table td { font-size: 13px; }
        .empty { text-align: center; padding: 40px; color: #999; }
        .toast { position: fixed; top: 20px; right: 20px; background: #52c41a; color: #fff; padding: 12px 24px; border-radius: 4px; display: none; z-index: 1000; }
        .toast.error { background: #ff4d4f; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Quake 资产测绘工具</h1>
        </div>
        
        <div class="tabs">
            <div class="tab active" onclick="switchTab('cookie')">Cookie 模式</div>
            <div class="tab" onclick="switchTab('api')">API 模式</div>
        </div>
        
        <div id="cookie-panel" class="panel active">
            <div class="form-group">
                <label>Cookie (从浏览器开发者工具复制)</label>
                <textarea id="cookie-input" placeholder="粘贴Cookie内容..."></textarea>
            </div>
            <div class="btn-group">
                <button class="btn btn-primary" onclick="testConnection()">测试连接</button>
            </div>
        </div>
        
        <div id="api-panel" class="panel">
            <div class="form-group">
                <label>API地址</label>
                <input type="text" id="api-url" placeholder="https://quake.chaitin.com">
            </div>
            <div class="form-group">
                <label>API Key</label>
                <input type="password" id="api-key" placeholder="输入API Key">
            </div>
            <div class="btn-group">
                <button class="btn btn-primary" onclick="testConnection()">测试连接</button>
                <button class="btn btn-success" onclick="saveConfig()">保存配置</button>
            </div>
        </div>
        
        <div class="query-box">
            <div class="query-row">
                <input type="text" id="query-input" placeholder="输入Quake查询语法，如: port: 80" onkeydown="if(event.key==='Enter')search()">
                <button class="btn btn-primary" onclick="search()">搜索</button>
            </div>
        </div>
        
        <div class="results-box">
            <div class="results-header">
                <h3>查询结果 <span id="result-count"></span></h3>
                <div class="btn-group" style="margin-top:0">
                    <button class="btn btn-default" onclick="exportData('csv')">导出CSV</button>
                    <button class="btn btn-default" onclick="exportData('json')">导出JSON</button>
                    <button class="btn btn-default" onclick="exportData('excel')">导出Excel</button>
                </div>
            </div>
            <div class="table-container" style="overflow-x:auto">
                <table class="results-table" id="results-table">
                    <thead>
                        <tr>
                            <th>IP</th>
                            <th>端口</th>
                            <th>协议</th>
                            <th>标题</th>
                            <th>国家</th>
                            <th>城市</th>
                            <th>ASN</th>
                            <th>组织</th>
                        </tr>
                    </thead>
                    <tbody id="results-body">
                    </tbody>
                </table>
            </div>
            <div id="empty-msg" class="empty">暂无数据</div>
        </div>
    </div>
    
    <div id="toast" class="toast"></div>
    
    <script>
        let currentTab = 'cookie';
        let searchResults = [];
        
        window.onload = function() {
            fetch('/api/config').then(r=>r.json()).then(data=>{
                document.getElementById('api-url').value = data.api_url || '';
                document.getElementById('api-key').value = data.api_key || '';
                document.getElementById('cookie-input').value = data.cookie || '';
            });
        };
        
        function switchTab(tab) {
            currentTab = tab;
            document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
            document.querySelectorAll('.panel').forEach(p => p.classList.remove('active'));
            event.target.classList.add('active');
            document.getElementById(tab + '-panel').classList.add('active');
        }
        
        function showToast(msg, isError=false) {
            const toast = document.getElementById('toast');
            toast.textContent = msg;
            toast.className = 'toast' + (isError?' error':'');
            toast.style.display = 'block';
            setTimeout(() => toast.style.display = 'none', 3000);
        }
        
        function testConnection() {
            let apiUrl = document.getElementById('api-url').value;
            let apiKey = document.getElementById('api-key').value;
            let cookie = document.getElementById('cookie-input').value;
            
            if (currentTab === 'cookie' && !cookie) {
                showToast('请输入Cookie', true);
                return;
            }
            if (currentTab === 'api' && (!apiUrl || !apiKey)) {
                showToast('请输入API地址和Key', true);
                return;
            }
            
            fetch('/api/test', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({apiUrl, apiKey, cookie})
            }).then(r=>r.json()).then(data=>{
                if (data.status === 'ok') {
                    showToast(data.message);
                } else {
                    showToast(data.message, true);
                }
            });
        }
        
        function saveConfig() {
            let apiUrl = document.getElementById('api-url').value;
            let apiKey = document.getElementById('api-key').value;
            let cookie = document.getElementById('cookie-input').value;
            
            fetch('/api/save-config', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({apiUrl, apiKey, cookie})
            }).then(()=>showToast('配置已保存'));
        }
        
        function search() {
            let query = document.getElementById('query-input').value;
            let apiUrl = document.getElementById('api-url').value;
            let apiKey = document.getElementById('api-key').value;
            let cookie = document.getElementById('cookie-input').value;
            
            if (!query) {
                showToast('请输入查询语句', true);
                return;
            }
            
            fetch('/api/search', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({query, apiUrl, apiKey, cookie})
            }).then(r=>r.json()).then(data=>{
                searchResults = data;
                renderResults(data);
            });
        }
        
        function renderResults(items) {
            let tbody = document.getElementById('results-body');
            let emptyMsg = document.getElementById('empty-msg');
            let countSpan = document.getElementById('result-count');
            
            countSpan.textContent = items.length ? '(' + items.length + ')' : '';
            
            if (!items || items.length === 0) {
                tbody.innerHTML = '';
                emptyMsg.style.display = 'block';
                return;
            }
            
            emptyMsg.style.display = 'none';
            tbody.innerHTML = items.map(item => 
                '<tr>' +
                    '<td>' + (item.ip || '') + '</td>' +
                    '<td>' + (item.port || '') + '</td>' +
                    '<td>' + (item.protocol || '') + '</td>' +
                    '<td>' + (item.title || '') + '</td>' +
                    '<td>' + (item.country || '') + '</td>' +
                    '<td>' + (item.city || '') + '</td>' +
                    '<td>' + (item.asn || '') + '</td>' +
                    '<td>' + (item.org || '') + '</td>' +
                '</tr>'
            ).join('');
        }
        
        function exportData(format) {
            if (searchResults.length === 0) {
                showToast('没有可导出的数据', true);
                return;
            }
            
            let endpoint = format === 'csv' ? '/api/export-csv' : 
                          format === 'json' ? '/api/export-json' : '/api/export-excel';
            
            fetch(endpoint, {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify(searchResults)
            }).then(r=>r.json()).then(data=>{
                if (data.status === 'ok') {
                    showToast('已保存到: ' + data.filename);
                } else if (data.status !== 'cancel') {
                    showToast('导出失败', true);
                }
            });
        }
    </script>
</body>
</html>`
