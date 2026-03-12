package quakeclient

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
)

func ExportCSV(results []Asset, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"IP", "Port", "Protocol", "Title", "Country", "City", "ASN", "Org", "UpdatedAt"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("写入表头失败: %v", err)
	}

	for _, asset := range results {
		row := []string{
			asset.IP,
			fmt.Sprintf("%d", asset.Port),
			asset.Protocol,
			asset.Title,
			asset.Country,
			asset.City,
			asset.ASN,
			asset.Org,
			asset.UpdatedAt,
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("写入数据失败: %v", err)
		}
	}

	return nil
}

func ExportJSON(results []Asset, filename string) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化失败: %v", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

func ExportExcel(results []Asset, filename string) error {
	f := excelize.NewFile()
	defer f.Close()

	headers := []string{"IP", "Port", "Protocol", "Title", "Country", "City", "ASN", "Org", "UpdatedAt"}
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
		f.SetCellValue("Sheet1", fmt.Sprintf("I%d", row), asset.UpdatedAt)
	}

	if err := f.SaveAs(filename); err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}

	return nil
}
