package quakeclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	APIUrl  string `json:"api_url"`
	APIKey  string `json:"api_key"`
	LastDir string `json:"last_dir"`
}

type SearchResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
}

type Data struct {
	Total    int     `json:"total"`
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
	Items    []Asset `json:"items"`
}

type Asset struct {
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	Protocol  string `json:"protocol"`
	Title     string `json:"title"`
	Banner    string `json:"banner"`
	Country   string `json:"country"`
	City      string `json:"city"`
	ASN       string `json:"asn"`
	Org       string `json:"org"`
	UpdatedAt string `json:"updated_at"`
}

type SearchRequest struct {
	Query    string `json:"query"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

type QuakeClient struct {
	baseURL    string
	apiKey     string
	cookie     string
	httpClient *http.Client
}

func NewQuakeClient() *QuakeClient {
	return &QuakeClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *QuakeClient) SetAPIConfig(baseURL, apiKey string) {
	c.baseURL = strings.TrimRight(baseURL, "/")
	c.apiKey = apiKey
}

func (c *QuakeClient) SetCookie(cookie string) {
	c.cookie = cookie
}

func (c *QuakeClient) TestConnection() error {
	if c.baseURL == "" {
		return fmt.Errorf("请先配置API地址")
	}

	req, err := http.NewRequest("GET", c.baseURL+"/api/v3/search/user-info", nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	if c.apiKey != "" {
		req.Header.Set("X-QuakeAIKey", c.apiKey)
	}
	if c.cookie != "" {
		req.Header.Set("Cookie", c.cookie)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("连接失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return fmt.Errorf("认证失败，请检查Cookie或API Key")
	}
	if resp.StatusCode == 403 {
		return fmt.Errorf("权限不足")
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("服务器返回错误: %d", resp.StatusCode)
	}

	return nil
}

func (c *QuakeClient) Search(query string, page, pageSize int) (*SearchResult, error) {
	if c.baseURL == "" {
		return nil, fmt.Errorf("请先配置API地址")
	}

	searchReq := SearchRequest{
		Query:    query,
		Page:     page,
		PageSize: pageSize,
	}

	jsonData, err := json.Marshal(searchReq)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/api/v3/search", strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("X-QuakeAIKey", c.apiKey)
	}
	if c.cookie != "" {
		req.Header.Set("Cookie", c.cookie)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("认证失败，请检查Cookie或API Key")
	}
	if resp.StatusCode == 403 {
		return nil, fmt.Errorf("权限不足")
	}
	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("请求过于频繁，请稍后重试")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("服务器返回错误: %d - %s", resp.StatusCode, string(body))
	}

	var result SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &result, nil
}

func (c *QuakeClient) GetAssetDetail(id string) (*Asset, error) {
	if c.baseURL == "" {
		return nil, fmt.Errorf("请先配置API地址")
	}

	req, err := http.NewRequest("GET", c.baseURL+"/api/v3/assets/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	if c.apiKey != "" {
		req.Header.Set("X-QuakeAIKey", c.apiKey)
	}
	if c.cookie != "" {
		req.Header.Set("Cookie", c.cookie)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var asset Asset
	if err := json.Unmarshal(body, &asset); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &asset, nil
}
