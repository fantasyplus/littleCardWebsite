package processdata

import (
	"encoding/json"
	"fmt"
	"os"
)

func ReadSellData() (selldata map[string][][]interface{}) {
	// 打开JSON文件
	file, err := os.Open("./data/json/selldata.json")
	if err != nil {
		fmt.Println("can't open selldata.json", err)
		return
	}
	defer file.Close()

	// 读取JSON数据并解码到一个map中
	var data map[string][][]interface{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("selldata.json decode fail", err)
		return
	}

	// 输出读取的数据
	// for key, items := range data {
	// 	fmt.Println("id:", key)
	// 	for _, item := range items {
	// 		if len(item) == 4 {
	// 			cn, ok1 := item[0].(string)
	// 			qq, ok2 := item[1].(string)
	// 			amount, ok3 := item[2].(float64)
	// 			if ok1 && ok2 && ok3 {
	// 				fmt.Printf("cn:%s, qq:%s, amount:%f\n", cn, qq, amount)
	// 			}

	// 		}
	// 	}
	// }

	return data
}
