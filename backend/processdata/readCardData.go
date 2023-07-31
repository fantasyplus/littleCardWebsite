package processdata

import (
	"encoding/json"
	"fmt"
	"os"
)

type JSONData struct {
	Title []string    `json:"title"`
	Data  [][]interface{} `json:"data"`
}

func ReadCardData() (cardData JSONData){
	// 打开 JSON 文件
	file, err := os.Open("./data/json/carddata.json")
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return
	}
	defer file.Close()

	// 创建 JSON 解码器
	decoder := json.NewDecoder(file)

	// 定义变量用于存储解析后的数据
	var jsonData JSONData

	// 解码 JSON 数据
	err = decoder.Decode(&jsonData)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// 输出解析后的数据
	// fmt.Println("Title:", jsonData.Title)
	// fmt.Println("Data:")
	// for _, item := range jsonData.Data {
	// 	fmt.Printf("序号: %v, 谷名: %v, 角色: %v, 制品: %v, 状态: %v, None1: %v, None2: %v\n",
	// 		item[0], item[1], item[2], item[3], item[4], item[5], item[6])
	// }
	return jsonData
}
