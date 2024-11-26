package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendSms(c *gin.Context) {
	log.Println("sendSms")
	log.Println(c.Request.Header)
	log.Println(c.Request.Body)
	log.Println(c.Request.URL.Query())
	c.JSON(http.StatusOK, gin.H{"message": "sendSms", "assert": true})
}
