package processdata

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xuri/excelize/v2"
)

type User struct {
	ID      string
	Name    string
	State   string
	City    string
	Address string
}

func GenerateExcelTest() {
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

//生成一个谷子对应一个角色的表格
func GenerateSingleTypeSheet(card_data CardNo, full_name string, f *excelize.File) bool {
	//设置sheet名
	f.NewSheet(full_name)

	f.SetCellValue(full_name, "A1", full_name)
	f.SetCellValue(full_name, "B1", "数量")



	// f.SaveAs()
	return false
}

//思路：
/*
1.先读取所有CardNo表，拆分出表序号，加上card_name可以作为谷子名和sheet名
2.第一列的“cn+群内qq:行时2986454288”，可以通过person_id查找到person_info表，得到cn和qq
3.数量和状态列就直接读取
*/
func GenerateSellExcel(db *gorm.DB) {

	// 创建一个新的 Excel 文件
	f := excelize.NewFile()

	var tableNames []string
	tableNames, _ = db.Migrator().GetTables()

	//一个谷子多个角色对应,每一组包含谷子全名和cardno数据
	multi_character_infos := []interface{}{}

	for _, tableName := range tableNames {
		if strings.Contains(tableName, "cardNo") {
			var cardno CardNo
			db.Table(tableName).Find(&cardno)

			//保存谷子的全名（包括序号）
			card_id := regexp.MustCompile(`\d+`).FindStringSubmatch(tableName)
			var full_card_name string = card_id[0] + cardno.CardName

			//如果是一个谷子多个角色的情况，先存下来，之后单独处理
			if strings.Contains(tableName, "_") {
				multi_character_infos = append(multi_character_infos, cardno, full_card_name)
			}

			//处理正常的情况（一对一）
			GenerateSingleTypeSheet(cardno,full_card_name,f)
		}
	}

}
