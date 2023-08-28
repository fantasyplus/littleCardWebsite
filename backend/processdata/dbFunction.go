package processdata

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func ConnectDB() *gorm.DB {
	dsn := "root:yxdbc2008@tcp(127.0.0.1:3306)/non_commercial_test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("failed to connect database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic("获取数据库连接失败")
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	return db
}

func CreateTable(db *gorm.DB) {
	// 使用AutoMigrate方法创建表
	db.AutoMigrate(&PersonInfo{}, &CardInfo{}, &CardIndex{})
}

func InsertPersonInfoTable(db *gorm.DB) (
	map[uint]([]map[string]float64),
	map[string]string,
	map[string][]map[uint]string) {
	// 定义插入数据的闭包
	insertPersonInfoTable := func(cn, qq string) (uint, error) {
		//如果已经有数据，会自动更新；如果没有数据，会自动插入
		//where查找，assign定义更新或插入的字段，firstorcreate执行更新或插入并返回结果
		var personinfo PersonInfo = PersonInfo{
			Model: gorm.Model{},
			CN:    cn,
			QQ:    qq,
		}
		result := db.Where(&PersonInfo{CN: cn, QQ: qq}).Assign(personinfo).FirstOrCreate(&personinfo)
		if result.Error != nil {
			return 0, result.Error
		}

		// 获取刚插入数据的主键
		id := personinfo.ID

		return id, nil
	}

	var selldata = ReadSellData()
	// person_id为key， value为一个map数组，map数组的元素的key为card_id，value为card_num
	// 相当于一个人对应了多少个谷子，每个谷子又有多少数量
	var person_id2card_ids2card_num = make(map[uint]([]map[string]float64))
	// 谷子id和谷子名字的映射
	var card_id2card_name = make(map[string]string)
	// 当某个cn和qq对应的status有信息时，用card_id存储map数组，每个map包含person_id和status的映射
	// 根据card_id查找cardNo表，根据person_id查找对应行，更新status
	var card_id2person_ids2status = make(map[string][]map[uint]string)
	var person_id uint
	var card_name, card_id string
	var card_num float64
	var cn, qq, card_status string
	for key, items := range selldata {
		card_id = key
		for _, item := range items {
			//标题行
			if len(item) <= 3 {
				var ok1 bool
				card_name, ok1 = item[0].(string)

				match1 := regexp.MustCompile(`\d+`).FindStringSubmatch(card_name)
				match2 := regexp.MustCompile(`\d+_\d+`).FindStringSubmatch(card_name)

				if match1 != nil && match2 == nil {
					cardID := match1[0]
					card_name = card_name[len(cardID):]
					// fmt.Printf("Card ID: %s, Card Name: %s\n", cardID, card_name)
					continue
				} else if match2 != nil {
					cardID := match2[0]
					card_name = card_name[len(cardID):]
					// fmt.Printf("Card ID: %s, Card Name: %s\n", cardID, card_name)
					continue
				}

				if ok1 {
					// fmt.Printf("title:%s\n", card_name)
				} else {
					fmt.Printf("read card_name error\n")
				}
				//数据行
			} else if len(item) == 4 {
				var ok1, ok2, ok3, ok4 bool
				cn, ok1 = item[0].(string)
				qq, ok2 = item[1].(string)
				card_num, ok3 = item[2].(float64)
				card_status, ok4 = item[3].(string)
				if ok1 && ok2 && ok3 && ok4 {
					// fmt.Printf("cn:%s, qq:%s, amount:%f, card_status:%s\n", cn, qq, card_num, card_status)
				} else {
					fmt.Printf("read sell data error on card_id:%s, card_name:%s,cn:%s,qq:%s\n", card_id, card_name, cn, qq)
				}
			}
			// 调用插入数据的闭包
			var err error
			person_id, err = insertPersonInfoTable(cn, qq)
			if err != nil {
				panic("failed to insert data: " + err.Error())
			}

			// 为每个card_id创建一个map数组，数组每个map包含person_id和status的映射
			card_id2person_ids2status[card_id] = append(card_id2person_ids2status[card_id], map[uint]string{person_id: card_status})

			// 创建新的嵌套card_id2card_num并插入数据，对应person_id
			person_id2card_ids2card_num[person_id] = append(person_id2card_ids2card_num[person_id], map[string]float64{card_id: card_num})
			// if person_id == 20 {
			// 	fmt.Println(cn, person_id, person_id2card_ids2card_num[person_id])
			// }

			//card_id和card_name对应关系
			card_id2card_name[card_id] = card_name
			// fmt.Println(card_id, card_id2card_name[card_id])
		}
	}
	return person_id2card_ids2card_num, card_id2card_name, card_id2person_ids2status
}

func InsertCardIndexTable(
	db *gorm.DB,
	person_id2card_ids2card_num map[uint]([]map[string]float64)) {

	for person_id, card_id2card_num := range person_id2card_ids2card_num {
		temp_ids := []string{}
		for _, item := range card_id2card_num {
			for card_id := range item {
				temp_ids = append(temp_ids, card_id)
			}
		}

		//去重，把一张表里一个人买了多次的谷子数量合并成一条记录
		temp_ids = RemoveDuplicates(temp_ids)
		card_ids := strings.Join(temp_ids, ",")

		//如果已经有数据，会自动更新；如果没有数据，会自动插入
		//where查找，assign定义更新或插入的字段，firstorcreate执行更新或插入并返回结果
		var cardindex CardIndex = CardIndex{
			Model:    gorm.Model{},
			PersonID: int(person_id),
			CardIDs:  card_ids,
		}
		result := db.Where(&CardIndex{PersonID: int(person_id)}).Assign(cardindex).FirstOrCreate(&cardindex)

		if result.Error != nil {
			fmt.Println("insertIntoCardIndexTable error:", result.Error)
		}

	}
}

func InsertCardInfoTable(db *gorm.DB) {
	cardData := ReadCardData()
	for _, item := range cardData.Data {
		card_id, card_name, card_character, card_type, card_condition, other := item[0], item[1], item[2], item[3], item[4], item[5]
		// fmt.Println(card_id, card_name, card_character, card_type, card_condition, other)

		//如果已经有数据，会自动更新；如果没有数据，会自动插入
		//where查找，assign定义更新或插入的字段，firstorcreate执行更新或插入并返回结果
		var cardinfo CardInfo
		var assign_cardinfo CardInfo = CardInfo{
			Model:         gorm.Model{},
			CardID:        card_id.(string),
			CardName:      card_name.(string),
			CardCharacter: card_character.(string),
			CardType:      card_type.(string),
			CardCondition: card_condition.(string),
			Other:         other.(string),
		}
		result := db.Where(&CardInfo{CardID: card_id.(string)}).Assign(assign_cardinfo).FirstOrCreate(&cardinfo)

		if result.Error != nil {
			fmt.Println("insertIntoCardInfoTable error:", result.Error)
		}
	}
}

// 获取某个人对应的某个谷子的状态
func getStatus(
	card_id2person_ids2status map[string][]map[uint]string,
	card_id string,
	person_id uint) string {

	for _, item := range card_id2person_ids2status[card_id] {
		for key, status := range item {
			if key == person_id {
				// fmt.Printf("card_id:%s,person_id:%d, status:%s\n", card_id, person_id, status)
				return status
			}
		}
	}
	return "none"
}

func InsertCardNoTable(
	db *gorm.DB,
	person_id2card_ids2card_num map[uint]([]map[string]float64),
	card_id2card_name map[string]string,
	card_id2person_ids2status map[string][]map[uint]string) {

	for card_id := range card_id2card_name {
		// Dynamically set the table name based on cardID
		tableName := fmt.Sprintf("cardNo%s", card_id)

		// Use TableName to set the table name for CardNo model
		db.Table(tableName).AutoMigrate(&CardNo{})
		// fmt.Printf("Created table %s\n", tableName)
	}

	// 根据每个person_id对应的card_id2card_num字典，插入或更新cardNo表
	for person_id, card_id2card_num := range person_id2card_ids2card_num {
		// 把相同card_id的card_num加起来
		card_id2card_num = MergeMap(card_id2card_num)

		//得到一个人对应有多少个card_id，以及每个card_id对应的card_num
		for _, item := range card_id2card_num {
			for card_id := range item {

				var cardno CardNo = CardNo{
					Model:    gorm.Model{},
					PersonID: int(person_id),
					CardName: card_id2card_name[card_id],
					CardNum:  item[card_id],
					Status:   getStatus(card_id2person_ids2status, card_id, person_id),
				}

				tableName := fmt.Sprintf("cardNo%s", card_id)
				//如果已经有数据，会自动更新；如果没有数据，会自动插入
				//where查找，assign定义更新或插入的字段，firstorcreate执行更新或插入并返回结果
				result := db.Table(tableName).Where(&CardNo{PersonID: int(person_id)}).Assign(cardno).FirstOrCreate(&cardno)
				if result.Error != nil {
					fmt.Printf("card_id:%s", card_id)
					fmt.Println("insertIntoCardNoTable error:", result.Error)
				}
			}
		}
	}
}

func FindCardInfoByCNQQ(db *gorm.DB, cn, qq string) []CardInfoRes {
	fmt.Printf("-----查找%s的谷子-----\n", cn)

	var cardInfoRes []CardInfoRes
	// 查找 person_id
	var personInfo []PersonInfo
	db.Where("cn = ? OR qq = ?", cn, qq).Find(&personInfo)

	//一个cn有可能对应多个person_id（不同的qq号）
	for _, item := range personInfo {
		fmt.Println(item.CN, item.QQ)

		//查询 card_ids
		var cardindex CardIndex
		db.Where("person_id = ?", item.ID).Find(&cardindex)
		card_ids := strings.Split(cardindex.CardIDs, ",")

		fmt.Println("谷子总数:", len(card_ids))
		for _, card_id := range card_ids {
			//根据 person_id 查询对应cardNo的信息
			var cardno []CardNo
			tableName := fmt.Sprintf("cardNo%s", card_id)
			db.Table(tableName).Where("person_id = ?", item.ID).Find(&cardno)

			//如果一个谷子表里同一个人买了好几次，就会有好几条记录
			var card_num float64
			var card_name string
			for _, cardno := range cardno {
				card_num = cardno.CardNum + card_num
				card_name = cardno.CardName
			}
			// fmt.Printf("序号:%s, 谷子名:%s, 谷子数量:%d ", card_id, card_name, int(card_num))

			//预处理，如果是一对多的情况，card_id为(19_1,1_1形式)，改成19,1
			//因为card_info表里的card_id是不带下划线的
			if strings.Contains(card_id, "_") {
				card_id = strings.Split(card_id, "_")[0]
			}

			var cardinfo CardInfo
			db.Where("card_id = ?", card_id).Find(&cardinfo)
			// fmt.Printf("角色:%s, 制品:%s, 状态:%s, 备注:%s\n",
			// cardinfo.CardCharacter, cardinfo.CardType, cardinfo.CardCondition, cardinfo.Other)

			cardInfoRes = append(cardInfoRes, CardInfoRes{
				CardID:         card_id,
				Card_name:      card_name,
				Card_character: cardinfo.CardCharacter,
				Card_type:      cardinfo.CardType,
				Card_condition: cardinfo.CardCondition,
				Card_num:       strconv.Itoa(int(card_num)),
				Card_deliver:   cardno[0].Status,
				Other:          cardinfo.Other,
			})
		}
	}

	return cardInfoRes
}

// 根据cn和qq和cardInfo更新谷子状态，必须要是cardInfo里已发货的谷子才能更新
func UpdateStatusByCNQQ(db *gorm.DB, cn, qq, status string) {
	fmt.Printf("-----更新%s的谷子状态为:%s-----\n", cn, status)

	// 先把谷子的发货信息保存下来
	var cardinfo []CardInfo
	db.Table("card_info").Find(&cardinfo)

	//保存card_id和card_status的对应关系
	card_id2card_status := map[string]string{}
	for _, item := range cardinfo {
		card_id2card_status[item.CardID] = item.CardCondition
	}
	// fmt.Println(card_id2card_status)

	// 查找 person_id对应的card_ids
	// 先通过person_info表查找person_id
	var personInfo []PersonInfo
	db.Where("cn = ? OR qq = ?", cn, qq).Find(&personInfo)

	//一个cn有可能对应多个person_id（不同的qq号）
	for _, item := range personInfo {
		fmt.Println(item.CN, item.QQ)

		//再根据person_id查找card_ids
		var cardindex CardIndex
		db.Where("person_id = ?", item.ID).Find(&cardindex)
		card_ids := strings.Split(cardindex.CardIDs, ",")

		//根据card_id2card_status的信息更新card_ids的状态
		for _, card_id := range card_ids {
			//预处理，如果是一对多的情况，card_id为(19_1,1_1形式)，改成19,1
			//这里只取第一个，因为card_id2card_status里的card_id是不带_的
			query_card_id := card_id
			if strings.Contains(query_card_id, "_") {
				query_card_id = strings.Split(query_card_id, "_")[0]
			}

			//如果这个谷子是已发货的，更新cardNo表的对应person_id的状态
			if card_id2card_status[query_card_id] != "none" {
				fmt.Printf("更新card_id:%s,person_id:%d的状态为:%s\n", card_id, item.ID, status)

				tableName := fmt.Sprintf("cardNo%s", card_id)
				db.Table(tableName).Where("person_id = ?", item.ID).Update("status", status)
			}
		}
	}
}

// 去重函数，输入一个字符串切片，返回去重后的切片
func RemoveDuplicates(input []string) []string {
	// Helper map to track unique values
	uniqueMap := make(map[string]bool)

	// Result slice to store unique values
	uniqueSlice := []string{}

	// Iterate through the input slice and add unique values to the uniqueMap
	for _, value := range input {
		if _, ok := uniqueMap[value]; !ok {
			// If value is not present in the uniqueMap, add it and set its value to true
			uniqueMap[value] = true
			uniqueSlice = append(uniqueSlice, value)
		}
	}

	// fmt.Println("uniqueSlice:", uniqueSlice)

	return uniqueSlice
}

// Helper map to merge values with the same key
func MergeMap(card_id2card_num []map[string]float64) []map[string]float64 {
	mergedMap := make(map[string]float64)

	// Iterate through card_id2card_num slice and merge values with the same key
	for _, cardMap := range card_id2card_num {
		for cardID, value := range cardMap {
			// Check if the cardID already exists in the mergedMap
			if currentValue, ok := mergedMap[cardID]; ok {
				// If cardID already exists, add the value to the current value
				mergedMap[cardID] = currentValue + value
			} else {
				// If cardID does not exist, add it to the mergedMap with the value
				mergedMap[cardID] = value
			}
		}
	}

	// Convert the mergedMap back to a slice of maps
	mergedSlice := []map[string]float64{}
	for key, value := range mergedMap {
		mergedSlice = append(mergedSlice, map[string]float64{key: value})
	}

	return mergedSlice
}
