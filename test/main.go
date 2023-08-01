package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xuri/excelize/v2"
)

type User struct {
	ID       string
	Name     string
	State      string
	City    string
	Address string
}

func main() {
	// 连接到 MySQL 数据库
	dsn := "root:yxdbc2008@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("failed to connect database")
	}


	// 查询数据
	var people []User
	db.Find(&people)

	// 创建一个新的 Excel 文件
	f := excelize.NewFile()
	f.NewSheet("Sheet1")

	// 写入表头
	f.SetCellValue("Sheet1", "A1", "ID")
	f.SetCellValue("Sheet1", "B1", "Name")
	f.SetCellValue("Sheet1", "C1", "State")
	f.SetCellValue("Sheet1", "D1", "City")
	f.SetCellValue("Sheet1", "E1", "Address")

	// 从数据库读取数据，并将数据写入 Excel 表格
	for i, person := range people {
		row := i + 2
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), person.ID)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), person.Name)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), person.State)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), person.City)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), person.Address)
	}

	// 将数据保存到 Excel 文件
	err = f.SaveAs("output.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Excel 文件生成成功！")
}
