package processdata_test

import (
	"backend/processdata"
	"fmt"
	"testing"
)

func TestCardData(t *testing.T) {
	data := processdata.ReadCardData()
	fmt.Println("Title:", data.Title)
	fmt.Println("Data:")
	for _, item := range data.Data {
		fmt.Printf("序号: %v, 谷名: %v, 角色: %v, 制品: %v, 状态: %v, None1: %v, None2: %v\n",
			item[0], item[1], item[2], item[3], item[4], item[5], item[6])
	}
}

func TestSellData(t *testing.T) {
	data := processdata.ReadSellData()
	for _, items := range data {
		// fmt.Println("id:", key)
		for _, item := range items {
			if len(item) <= 3 {
				fmt.Println(item)
			}
			if len(item) == 4 {
				cn, ok1 := item[0].(string)
				qq, ok2 := item[1].(string)
				amount, ok3 := item[2].(float64)
				status, ok4 := item[3].(string)
				if ok1 && ok2 && ok3 && ok4 {
					fmt.Printf("cn:%s, qq:%s, amount:%f,status:%s\n", cn, qq, amount, status)
				}

			}
		}
	}
}

// 更新数据库
func TestDb(t *testing.T) {
	//执行python脚本，从网站下载数据并读取数据，更新json文件
	processdata.DownloadAndRead("DownloadAndRead.py")
 
	//根据excel文件更新数据库
	processdata.UpdateDataBase()
}

// 根据cn和qq查找谷子信息
func TestFind(t *testing.T) {
	db := processdata.ConnectDB()
	processdata.FindCardInfoByCNQQ(db, "银河", "")
}

// 指定cn和qq更新status
func TestUpdateStatusByCNQQ(t *testing.T) {
	db := processdata.ConnectDB()
	processdata.UpdateStatusByCNQQ(db, "银河", "", "none")
}

func TestGenerateExcel(t *testing.T) {
	db := processdata.ConnectDB()
	processdata.GenerateSellExcel(db, "/home/web/web/web-project/backend/processdata/data/test_excel/selldata/selldata_2023_08_14_23.xlsx")
}

func TestLqbz(t *testing.T) {
	db := processdata.ConnectDB()

	//执行python脚本，从网站下载数据并读取数据，更新json文件
	_, sell_data_path := processdata.DownloadAndRead("DownloadAndRead.py")
 
	//根据json文件更新数据库
	processdata.UpdateDataBase()

	//更新某人的状态
	// processdata.UpdateStatusByCNQQ(db, "银河", "", "none")

	//在这之前要更新数据库
	processdata.GenerateSellExcel(db, sell_data_path)
}
