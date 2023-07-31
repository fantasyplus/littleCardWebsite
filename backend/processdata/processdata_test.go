package processdata_test

import (
	"backend/processdata"
	"testing"
	"fmt"
)

func TestCardData(t *testing.T) {
	data:=processdata.ReadCardData()
	fmt.Println("Title:", data.Title)
	fmt.Println("Data:")
	for _, item := range data.Data {
		fmt.Printf("序号: %v, 谷名: %v, 角色: %v, 制品: %v, 状态: %v, None1: %v, None2: %v\n",
			item[0], item[1], item[2], item[3], item[4], item[5], item[6])
	}
}

func TestSellData(t *testing.T) {
	processdata.ReadSellData()
}

func TestDelData(t *testing.T) {
	db := processdata.ConnectDB()
	processdata.DelteTable(db)
}
func TestDb(t *testing.T) {
	db := processdata.ConnectDB()
	processdata.CreateTable(db)
	person_id2card_id2card_num, card_id2card_name := processdata.InsertPersonInfoTable(db)
	processdata.InsertCardIndexTable(db, person_id2card_id2card_num)
	processdata.InsertCardInfoTable(db)
	processdata.InsertCardNoTable(db, person_id2card_id2card_num, card_id2card_name)
}

func TestFind(t *testing.T) {
	db := processdata.ConnectDB()
	processdata.FindCardInfoByCNQQ(db, "heitai", "")
}

func TestMergeMap(t *testing.T) {
	db := processdata.ConnectDB()
	processdata.CreateTable(db)
	person_id2card_id2card_num, _ := processdata.InsertPersonInfoTable(db)
	for person_id, card_id2card_num := range person_id2card_id2card_num {
		if person_id==20{
			processdata.MergeMap(card_id2card_num)
		}
	}
}

func TestRemoveDuplicates(t *testing.T) {
	var temp []string=[]string{"1","2","3","3","3","6","6","6","7","8","9","10","11","12","13","14","15"}
	processdata.RemoveDuplicates(temp)
}