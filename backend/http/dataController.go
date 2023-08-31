package http

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend/processdata"
)

type DataController struct {
	db *gorm.DB
}

func NewDataController(db *gorm.DB) *DataController {
	return &DataController{db: db}
}

// /search?cn=xxx&qq=xxx
func (dc *DataController) SearchCardInfos(c *gin.Context) {
	cn := c.DefaultQuery("cn", "-1")
	qq := c.DefaultQuery("qq", "-1")
	fmt.Println(cn, qq)
	decodedCN, _ := url.QueryUnescape(cn) // URL解码
	decodedQQ, _ := url.QueryUnescape(qq) // URL解码

	res := processdata.FindCardInfoByCNQQ(dc.db, decodedCN, decodedQQ)

	// fmt.Println(res)
	if len(res) == 0 {
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
		"data":    res,
	})

}

// "/v1/data/list?pageSize=100&pageNum=1"
func (dc *DataController) GetAllCardInfos(c *gin.Context) {
	var dataList []processdata.CardInfo
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "-1"))
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "-1"))

	var totalItems int64

	if pageSize > 0 && pageNum > 0 {
		offsetVal := (pageNum - 1) * pageSize
		dc.db.Model(&processdata.CardInfo{}).Count(&totalItems).Limit(pageSize).Offset(offsetVal).Find(&dataList)
	} else {
		dc.db.Find(&dataList)
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

// "/v1/data/add"
/*
input:cn,qq,card_id,card_name,card_character
card_type,card_condition,card_num,card_deliver,other_info

1.通过cn和qq获得对应的person_id
2.通过card_id和card数据插入对应的cardNo表
*/
func (dc *DataController) AddCardInfoByCNOrQQ(c *gin.Context) {
	type temp_info struct {
		CN            string `json:"cn"`
		QQ            string `json:"qq"`
		CardID        string `json:"card_id"`
		CardName      string `json:"card_name"`
		CardCharacter string `json:"card_character"`
		CardType      string `json:"card_type"`
		CardCondition string `json:"card_condition"`
		CardNum       string `json:"card_num"`
		CardDeliver   string `json:"card_deliver"`
		CardOtherInfo string `json:"other_info"`
	}
	
	var tempInfo temp_info
	if err := c.ShouldBindJSON(&tempInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "add error",
			"data":    gin.H{},
			"code":    http.StatusBadRequest,
		})
		return
	}
	
	fmt.Println(tempInfo)
	
}
