package processdata_test

import(
	"testing"
	"backend/processdata"
)

func TestCardData(t *testing.T){
	processdata.ReadCardData()
}

func TestSellData(t *testing.T){
	processdata.ReadSellData()
}