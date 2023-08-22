package http

import (
	"backend/processdata"
	_"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


type TestController struct {
	db *gorm.DB
}

func NewTestController(db *gorm.DB) *TestController {
	return &TestController{db: db}
}

func (tc *TestController) AddUser(c *gin.Context) {
	var data processdata.CardInfo

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "add error",
			"data":    gin.H{},
			"code":    http.StatusBadRequest,
		})
		return
	}

	tc.db.Create(&data)

	c.JSON(http.StatusOK, gin.H{
		"message": "add stccess",
		"data":    data,
		"code":    http.StatusOK,
	})
}

func (tc *TestController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var data processdata.CardInfo

	result := tc.db.Where("id = ?", id).First(&data)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "delete error",
			"data":    gin.H{},
			"code":    http.StatusBadRequest,
		})
		return
	}

	tc.db.Delete(&data)

	c.JSON(http.StatusOK, gin.H{
		"message": "delete stccess",
		"data":    data,
		"code":    http.StatusOK,
	})
}

func (tc *TestController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var data processdata.CardInfo

	result := tc.db.Select("id").Where("id = ?", id).First(&data)

	if result.RowsAffected == 0 {
		message := "can't find id " + id
		c.JSON(http.StatusOK, gin.H{
			"message": message,
			"code":    http.StatusBadRequest,
		})
		return
	}

	err := c.ShouldBindJSON(&data)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "update fail",
			"code":    http.StatusBadRequest,
		})
		return
	}

	tc.db.Where("id = ?", id).Updates(&data)

	c.JSON(http.StatusOK, gin.H{
		"message": "update stccess",
		"code":    http.StatusOK,
	})
}

func (tc *TestController) GetUserByName(c *gin.Context) {
	name := c.Param("name")
	var dataList []processdata.CardInfo
	
	tc.db.Where("card_name = ?", name).Find(&dataList)

	if len(dataList) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "select error",
			"code":    http.StatusBadRequest,
			"data":    gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "select stccess",
		"code":    http.StatusOK,
		"data":    dataList,
	})
}

func (tc *TestController) GetAllUsers(c *gin.Context) {
	var dataList []processdata.CardInfo
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "-1"))
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "-1"))

	var totalItems int64

	if pageSize > 0 && pageNum > 0 {
		offsetVal := (pageNum - 1) * pageSize
		tc.db.Model(&processdata.CardInfo{}).Count(&totalItems).Limit(pageSize).Offset(offsetVal).Find(&dataList)
	} else {
		tc.db.Find(&dataList)
		totalItems = int64(len(dataList))
	}

	if len(dataList) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "no data found",
			"code":    http.StatusBadRequest,
			"data":    gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "query stccess",
		"code":    http.StatusOK,
		"data": gin.H{
			"data":       dataList,
			"totalItems": totalItems,
			"pageNum":    pageNum,
			"pageSize":   pageSize,
		},
	})
}

func (tc *TestController) DownloadFile(c *gin.Context) {
	filePath := "/home/web/web/web-project/backend/processdata/data/test_excel/selldata/selldata_2023_08_14_7.xlsx" // 修改为实际文件路径
	c.File(filePath)
}