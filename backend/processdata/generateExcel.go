package processdata

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	mapset "github.com/deckarep/golang-set"
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
func GenerateSingleTypeSheet(db *gorm.DB, card_data []CardNo, full_name string, f *excelize.File) bool {
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

	return false
}

func GetCharacterName(item [][]CardNo) mapset.Set {
	character_list := mapset.NewSet()
	for _, datas := range item {
		for _, data := range datas {
			//获取当前数据表对应的角色名
			match := regexp.MustCompile(`-`).FindAllStringIndex(data.CardName, -1)
			character_name := data.CardName[match[0][1]:]
			character_list.Add(character_name)
		}
	}
	return character_list
}

func GetCharacterCell(character_list mapset.Set, card_name string, row int) string {
	var cell_name string

	match := regexp.MustCompile(`-`).FindAllStringIndex(card_name, -1)
	character_name := card_name[match[0][1]:]
	temp := character_list.ToSlice()
	for index := range temp {
		if temp[index] == character_name {
			cell_name = fmt.Sprintf("%c%d", 'A'+index+1, row+2)
			break
		}
	}

	return cell_name
}

// 生成一个谷子对应多个角色的表格
func GenerateMultiTypeSheet(db *gorm.DB, multi_character_infos [][]interface{}, f *excelize.File) bool {
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

		//设置sheet名
		f.NewSheet(full_name_map[card_id])

		//设置表头
		SetCellValueWithErrHandle(f, full_name_map[card_id], "A1", full_name_map[card_id])

		//设置角色名（数量列表头）
		for index, character_name := range character_list.ToSlice() {
			cell_name := fmt.Sprintf("%c1", 'A'+index+1)
			SetCellValueWithErrHandle(f, full_name_map[card_id], cell_name, character_name)
		}

		//状态列表头
		SetCellValueWithErrHandle(f, full_name_map[card_id], fmt.Sprintf("%c1", 'A'+character_list.Cardinality()+1), "状态")

		//设置数据
		is_personid_shown := map[int]string{}
		for _, datas := range item {
			for row, data := range datas {
				cn_qq := getCNQQ(db, data.PersonID)
				num := data.CardNum
				status := data.Status

				var cell_name_character string
				if is_personid_shown[data.PersonID] != "" {
					//如果该person_id已经设置过cell_name，那么就直接用之前的cell_name
					cell_name_character = is_personid_shown[data.PersonID]
					
					//同时只需要设置数量就可以了，因为cnqq和状态已经设置过了
					SetCellValueWithErrHandle(f, full_name_map[card_id], cell_name_character, num)

					// fmt.Println(data.CardName,cn_qq,cell_name_character)
				} else {
					//否则即为第一次设置该person_id的数据

					//获得当前角色名对应的位置
					cell_name_character = GetCharacterCell(character_list, data.CardName, row)
					//设置cnqq
					SetCellValueWithErrHandle(f, full_name_map[card_id], fmt.Sprintf("A%d", row+2), cn_qq)
					//设置数量
					SetCellValueWithErrHandle(f, full_name_map[card_id], cell_name_character, num)
					//设置状态
					SetCellValueWithErrHandle(f, full_name_map[card_id], fmt.Sprintf("%c%d", 'A'+character_list.Cardinality()+1, row+2), status)

					fmt.Println(data.CardName,cn_qq,cell_name_character)
				}


				is_personid_shown[data.PersonID] = cell_name_character
			}
		}
	}

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
