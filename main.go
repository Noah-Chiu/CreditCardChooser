package main

import (
	"bytes"
	"credit-card-chooser/sql"
	"credit-card-chooser/util"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// AddCard()
	// return
	router := gin.Default()
	router.Use(Cors())
	sql.InitialDB()

	// 主要功能
	router.POST("/", Chooser)

	// 行事曆
	router.POST("/calendar", Schedule)

	// 測試
	router.POST("/test", Test)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 是否存取cookie
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		method := c.Request.Method
		// 放行OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		}
		// 處理請求
		c.Next()
	}
}

func Response(res Reply, token string) {
	if len(res.Messages) > 5 {
		res.Messages = []Message{{
			Type: "text",
			Text: "符合條件過多\n請再試一次",
		}}
	}
	var buf bytes.Buffer

	json.NewEncoder(&buf).Encode(res)

	req, err := http.NewRequest("POST", util.ReplyUrl, &buf)
	if err != nil {
		return
	}

	// Set the request content type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer {%s}", token))

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	sitemap, _ := io.ReadAll(resp.Body)

	fmt.Println("收到的回傳:", string(sitemap))
}
