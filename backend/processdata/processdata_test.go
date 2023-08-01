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

func TestDelData(t *testing.T) {
	db := processdata.ConnectDB()
	processdata.DelteTable(db)
}
func TestDb(t *testing.T) {
	db := processdata.ConnectDB()
	processdata.CreateTable(db)
	person_id2card_ids2card_num, card_id2card_name, card_id2person_ids2status := processdata.InsertPersonInfoTable(db)
	processdata.InsertCardIndexTable(db, person_id2card_ids2card_num)
	processdata.InsertCardInfoTable(db)
	processdata.InsertCardNoTable(db, person_id2card_ids2card_num, card_id2card_name, card_id2person_ids2status)
}

func TestFind(t *testing.T) {
	db := processdata.ConnectDB()
	processdata.FindCardInfoByCNQQ(db, "糖", "")
}

func TestMergeMap(t *testing.T) {
	db := processdata.ConnectDB()
	processdata.CreateTable(db)
	person_id2card_ids2card_num, _, _ := processdata.InsertPersonInfoTable(db)
	for person_id, card_id2card_num := range person_id2card_ids2card_num {
		if person_id == 20 {
			processdata.MergeMap(card_id2card_num)
		}
	}
}

func TestRemoveDuplicates(t *testing.T) {
	var temp []string = []string{"1", "2", "3", "3", "3", "6", "6", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15"}
	processdata.RemoveDuplicates(temp)
}

func TestUpdateStatusByCNQQ(t *testing.T) {
	db := processdata.ConnectDB()
	processdata.UpdateStatusByCNQQ(db, "糖", "", "什么鬼")
}
