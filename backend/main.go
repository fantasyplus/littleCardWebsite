package main

import (
	"github.com/gin-gonic/gin"

	"backend/http"
	"backend/processdata"
)


func main() {
	// 连接数据库

	db := processdata.ConnectDB()
	TestController := http.NewTestController(db)
	DataController := http.NewDataController(db)

	r := gin.Default()

	v1 := r.Group("/v1")
	{
		testGroup := v1.Group("/test")
		{
			testGroup.POST("/add", TestController.AddUser)
			testGroup.DELETE("/delete/:id", TestController.DeleteUser)
			testGroup.PUT("/update/:id", TestController.UpdateUser)
			testGroup.GET("/list/:name", TestController.GetUserByName)
			testGroup.GET("/list", TestController.GetAllUsers)
			testGroup.GET("/download", TestController.DownloadFile) 
		}

		dataGroup := v1.Group("/data")
		{
			dataGroup.GET("/search", DataController.SearchCardInfos)
			dataGroup.GET("/list", DataController.GetAllCardInfos)
		}
	}

	port := ":3001"
	r.Run(port)
}
