package processdata

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func getCNQQ(db *gorm.DB, person_id int) string {
	var person_info PersonInfo
	db.Where("id = ?", person_id).Find(&person_info)
	prefix := "cn+群内qq:"
	return prefix + person_info.CN + person_info.QQ
}

func SetCellValueWithErrHandle(f *excelize.File, sheet_name string, cell_name string, value interface{}) {
	err := f.SetCellValue(sheet_name, cell_name, value)
	if err != nil {
		fmt.Println(sheet_name, cell_name, err)
	}
}

// 保存文件
func SaveExcelFile(f *excelize.File) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("generate excel error when save,", err)
	}

	var savePath = currentDir + "/./data/output/test.xlsx"
	save_err := f.SaveAs(savePath)
	if save_err != nil {
		fmt.Println(save_err)
	}
}

// 生成一个谷子对应一个角色的表格
func GenerateSingleTypeSheet(db *gorm.DB, card_data []CardNo, full_name string, f *excelize.File) {
	//设置sheet名
	f.NewSheet(full_name)

	//设置表头
	SetCellValueWithErrHandle(f, full_name, "A1", full_name)
	//数量列表头
	SetCellValueWithErrHandle(f, full_name, "B1", "数量")
	//状态列表头
	SetCellValueWithErrHandle(f, full_name, "C1", "状态")

	//设置数据
	for index, data := range card_data {
		cn_qq := getCNQQ(db, data.PersonID)
		num := data.CardNum
		status := data.Status

		SetCellValueWithErrHandle(f, full_name, fmt.Sprintf("A%d", index+2), cn_qq)
		SetCellValueWithErrHandle(f, full_name, fmt.Sprintf("B%d", index+2), num)
		SetCellValueWithErrHandle(f, full_name, fmt.Sprintf("C%d", index+2), status)
	}
}

func GetCharacterName(item [][]CardNo) []string {
	character_list := []string{}
	for _, datas := range item {
		for _, data := range datas {
			//获取当前数据表对应的角色名
			match := regexp.MustCompile(`-`).FindAllStringIndex(data.CardName, -1)
			character_name := data.CardName[match[0][1]:]

			//去重,判断角色名是否已经存在
			if !strings.Contains(strings.Join(character_list, ","), character_name) {
				character_list = append(character_list, character_name)
			}
		}
	}
	return character_list
}

func GetCharacterCell(character_list []string, card_name string, row int) string {
	var cell_name string

	match := regexp.MustCompile(`-`).FindAllStringIndex(card_name, -1)
	character_name := card_name[match[0][1]:]

	for index := range character_list {
		if character_list[index] == character_name {
			cell_name = fmt.Sprintf("%c%d", 'A'+index+1, row)
			break
		}
	}

	return cell_name
}

func SetMultiTypeData(db *gorm.DB, item [][]CardNo, sheet_name string, character_list []string, f *excelize.File) {
	//设置数据
	is_personid_shown := map[int]string{}
	//excel表数据部分的总行数,从第二行开始
	total_row := 2
	for _, datas := range item {
		for _, data := range datas {
			cn_qq := getCNQQ(db, data.PersonID)
			num := data.CardNum
			status := data.Status

			var cell_name_character string
			if is_personid_shown[data.PersonID] != "" {
				//如果该person_id已经设置过cell_name，那么就直接用之前的cell_name
				cell_name_character = is_personid_shown[data.PersonID]
				//分离出之前的row，再重新通过角色名获得新的位置
				match := regexp.MustCompile(`\d+`).FindStringSubmatch(cell_name_character)
				old_row, _ := strconv.Atoi(match[0])
				cell_name_character = GetCharacterCell(character_list, data.CardName, old_row)

				//同时只需要设置数量就可以了，因为cnqq和状态已经设置过了
				SetCellValueWithErrHandle(f, sheet_name, cell_name_character, num)

				fmt.Println("two", data.CardName, cn_qq, cell_name_character, num)
			} else {
				//否则即为第一次设置该person_id的数据

				//获得当前角色名对应的位置
				cell_name_character = GetCharacterCell(character_list, data.CardName, total_row)
				//设置cnqq
				SetCellValueWithErrHandle(f, sheet_name, fmt.Sprintf("A%d", total_row), cn_qq)
				//设置数量
				SetCellValueWithErrHandle(f, sheet_name, cell_name_character, num)
				//设置状态
				SetCellValueWithErrHandle(f, sheet_name, fmt.Sprintf("%c%d", 'A'+len(character_list)+1, total_row), status)

				fmt.Println("one", data.CardName, cn_qq, cell_name_character, num)

				is_personid_shown[data.PersonID] = cell_name_character

				//行数++
				total_row++
			}
		}
	}
}

// 生成一个谷子对应多个角色的表格
func GenerateMultiTypeSheet(db *gorm.DB, multi_character_infos [][]interface{}, f *excelize.File) {
	merged_data := map[string]([][]CardNo){}
	full_name_map := map[string]string{}

	for _, datas := range multi_character_infos {
		card_data := datas[0].([]CardNo)
		full_card_name := datas[1].(string)

		//把同属一个谷子的角色数据放在一起（比如下午茶甜点的三种，合并到一起）
		//这样merged_data中每个card_id对应的数据就是同一个谷子的所有角色数据
		match1 := regexp.MustCompile(`\d+`).FindStringSubmatch(full_card_name)
		merged_data[match1[0]] = append(merged_data[match1[0]], card_data)

		//获取纯粹的谷子名，去除角色名
		match2 := regexp.MustCompile(`-`).FindAllStringIndex(full_card_name, -1)
		//保存谷子的全名（包括序号,不包括角色名）
		full_card_name = full_card_name[:match2[0][0]]
		full_name_map[match1[0]] = full_card_name
	}

	//开始生成excel
	for card_id, item := range merged_data {
		//存放某一种谷子的所有角色名
		character_list := GetCharacterName(item)

		sheet_name := full_name_map[card_id]
		//设置sheet名
		f.NewSheet(sheet_name)
		f.SetColWidth(sheet_name, "A", "A", 30)

		//设置表头
		SetCellValueWithErrHandle(f, sheet_name, "A1", sheet_name)

		//设置角色名（数量列表头）
		for index, character_name := range character_list {
			cell_name := fmt.Sprintf("%c1", 'A'+index+1)
			SetCellValueWithErrHandle(f, sheet_name, cell_name, character_name)
		}

		//状态列表头
		SetCellValueWithErrHandle(f, sheet_name, fmt.Sprintf("%c1", 'A'+len(character_list)+1), "状态")

		//设置数据部分
		SetMultiTypeData(db, item, sheet_name, character_list, f)
	}
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
	multi_character_infos := [][]interface{}{}

	for _, tableName := range tableNames {
		if strings.Contains(tableName, "cardNo") {
			var cardno []CardNo
			db.Table(tableName).Find(&cardno)

			//保存谷子的全名（包括序号）
			card_id := regexp.MustCompile(`\d+`).FindStringSubmatch(tableName)
			var full_card_name string = card_id[0] + cardno[0].CardName

			//如果是一个谷子多个角色的情况，先存下来，之后单独处理
			if strings.Contains(tableName, "_") {
				multi_character_infos = append(multi_character_infos, []interface{}{cardno, full_card_name})
				continue
			}

			//处理正常的情况（一对一）
			// GenerateSingleTypeSheet(db, cardno, full_card_name, f)
		}
	}

	//处理一个谷子多个角色的情况
	GenerateMultiTypeSheet(db, multi_character_infos, f)

	//保存文件
	SaveExcelFile(f)
}
