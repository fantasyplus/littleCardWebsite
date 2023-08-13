package processdata

import (
	_ "fmt"
	"os"

	"github.com/xuri/excelize/v2"
)

//source=target，用于Sheet的覆盖
func CoverSheet(f *excelize.File, source_name string, target_name string) {
	SwapSheet(f, source_name, target_name)
	f.DeleteSheet(source_name)
}

// 交换source和target
func SwapSheet(f *excelize.File, source_name string, target_name string) {
	source_index, _ := f.GetSheetIndex(source_name)
	target_index, _ := f.GetSheetIndex(target_name)

	// 交换sheet内容
	f.NewSheet("temp_data")
	index_temp_data, _ := f.GetSheetIndex("temp_data")
	f.CopySheet(target_index, index_temp_data)
	f.CopySheet(source_index, target_index)
	f.CopySheet(index_temp_data, source_index)
	f.DeleteSheet("temp_data")

	// 交换sheet名
	f.SetSheetName(target_name, "temp_name")
	f.SetSheetName(source_name, target_name)
	f.SetSheetName("temp_name", source_name)
}

func ExeSortSheet() {
	f := excelize.NewFile()
	f.NewSheet("1")
	f.SetCellValue("1", "A1", "1")
	f.NewSheet("2")
	f.SetCellValue("2", "A1", "2")
	f.SetCellValue("2", "A2", "3")
	f.NewSheet("3")
	f.SetCellValue("3", "A1", "3")
	f.SetCellValue("3", "A2", "4")

	SwapSheet(f, "1", "3")
	SwapSheet(f, "2", "1")

	CoverSheet(f, "1", "2")//1=2
	CoverSheet(f, "2", "3")//2=3

	currentDir, _ := os.Getwd()
	var savePath = currentDir + "/./data/output/sort.xlsx"
	f.SaveAs(savePath)
}
