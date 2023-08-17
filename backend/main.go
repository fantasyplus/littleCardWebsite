package main

import (
	_"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	_"backend/processdata"
)

// 定义表示您的表的模型结构
type User struct {
	gorm.Model
	Name    string `gorm:"type:varchar(255); not null" json:"name" binding:"required"`
	State   string `gorm:"type:varchar(255); not null" json:"state" binding:"required"`
	City    string `gorm:"type:varchar(255); not null" json:"city" binding:"required"`
	Address string `gorm:"type:varchar(255); not null" json:"address" binding:"required"`
}

//relate to Table cardNo
type SingleCardData struct{
	gorm.Model
	Person_id uint `gorm:"type:int; not null" json:"person_id" binding:"required"`
	Card_name string `gorm:"type:varchar(255); not null" json:"card_name" binding:"required"`
	Card_num uint `gorm:"type:int; not null" json:"card_num" binding:"required"`
}

type CardInfo struct{
	gorm.Model
	Card_id string `gorm:"type:varchar(255); not null" json:"card_id" binding:"required"`
	Card_name string `gorm:"type:varchar(255); not null" json:"card_name" binding:"required"`
	Card_character string `gorm:"type:varchar(255); not null" json:"card_character" binding:"required"`
	Card_type string `gorm:"type:varchar(255); not null" json:"card_type" binding:"required"`
	Card_condition string `gorm:"type:varchar(255); not null" json:"card_condition" binding:"required"`
	Other string `gorm:"type:varchar(255); not null" json:"other" binding:"required"`
}

/* 注意点 :
结构体里面的变量 (Name) 必须是首字符大写
gorm 指定类型
json 表示json接受的时候的名称
binding required 表示必须传入
*/

func main() {
	// 将 "your_username"、"your_password" 和 "your_host:port" 替换为您实际的 MySQL 凭据
	dsn := "root:yxdbc2008@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

	// 连接到 MySQL 数据库服务器，不指定数据库名称
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("连接数据库服务器失败")
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

	db.AutoMigrate(&User{})

	r := gin.Default()

	/*业务码约定
	正确 : 200
	错误 : 400
	*/

	// add
	r.POST("/user/add", func(c *gin.Context) {
		var data User

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(200, gin.H{
				"message": "add error",
				"data":    gin.H{},
				"code":    400,
			})
		} else {
			db.Create(&data)

			c.JSON(200, gin.H{
				"message": "add success",
				"data":    data,
				"code":    200,
			})
		}
	})

	//delete
	// 1. 找到对应的 id 所对应的条目
	// 2. 判断 id 是否存在
	// 3. 从数据库中删除
	// 3. 返回, id 没有找到
	r.DELETE("/user/delete/:id", func(c *gin.Context) {
		var data []User

		id := c.Param("id")

		db.Where("id = ?", id).Find(&data)

		if len(data) == 0 {
			c.JSON(200, gin.H{
				"message": "delete error",
				"data":    gin.H{},
				"code":    400,
			})
		} else {
			db.Delete(&data)

			c.JSON(200, gin.H{
				"message": "delete success",
				"data":    data,
				"code":    200,
			})
		}
	})

	//update
	r.PUT("/user/update/:id", func(c *gin.Context) {

		// 1. 找到对应的 id 所对应的条目
		// 2. 判断 id 是否存在
		// 3. 修改对应条目
		// 4. 返回 id,没有找到

		var data User

		// 接受 id
		id := c.Param("id")

		// 判断 id 是否存在
		db.Select("id").Where("id = ? ", id).Find(&data)

		// 判断 id 是否存在
		if data.ID == 0 {
			var message string = "can't find id" + strconv.FormatUint(uint64(data.ID), 10)
			c.JSON(200, gin.H{
				"message": message,
				"code":    400,
			})
		} else {
			err := c.ShouldBindJSON(&data)

			if err != nil {
				c.JSON(200, gin.H{
					"message": "update fail",
					"code":    400,
				})
			} else {
				// db 修改数据库内容
				db.Where("id = ?", id).Updates(&data)

				c.JSON(200, gin.H{
					"message": "update success",
					"code":    200,
				})
			}
		}
	})

	// 查 ( 条件查询 , 全部查询 / 分页查询)

	// 条件查询
	r.GET("/user/list/:name", func(c *gin.Context) {

		// 获取路径参数
		name := c.Param("name")

		var dataList []User

		// 查询数据库
		db.Where("name = ? ", name).Find(&dataList)

		// 判断是否查询到数据
		if len(dataList) == 0 {
			c.JSON(200, gin.H{
				"message": "select error",
				"code":    400,
				"data":    gin.H{},
			})
		} else {
			c.JSON(200, gin.H{
				"message": "select success",
				"code":    200,
				"data":    dataList,
			})
		}
	})

	// 全部查询
	r.GET("/user/list", func(c *gin.Context) {

		var dataList []User

		// 1. 查询全部数据,  查询分页数据
		pageSize, _ := strconv.Atoi(c.Query("pageSize"))
		pageNum, _ := strconv.Atoi(c.Query("pageNum"))

		// 判断是否需要分页
		if pageSize == 0 {
			pageSize = -1
		}
		if pageNum == 0 {
			pageNum = -1
		}

		offsetVal := (pageNum - 1) * pageSize
		if pageNum == -1 && pageSize == -1 {
			offsetVal = -1
		}

		// 返回一个总数
		var totalitems int64
		// 查询数据库
		db.Model(dataList).Count(&totalitems).Limit(pageSize).Offset(offsetVal).Find(&dataList)

		if len(dataList) == 0 {
			c.JSON(200, gin.H{
				"message": "没有查询到数据",
				"code":    400,
				"data":    gin.H{},
			})
		} else {
			c.JSON(200, gin.H{
				"message": "查询成功",
				"code":    200,
				"data": gin.H{
					"data":       dataList,
					"totalitems": totalitems,
					"pageNum":    pageNum,
					"pageSize":   pageSize,
				},
			})
		}
	})

	// 路由处理函数，用于处理文件下载请求
	r.GET("/user/download", func(c *gin.Context) {
		filePath := "/home/web/web/web-project/backend/processdata/data/test_excel/selldata/selldata_2023_08_14_7.xlsx" // 修改为实际文件路径
		c.File(filePath)
	}) 

	port := ":3001"
	r.Run(port)
}
