package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"os"
)

type Request struct {
	URL string `json:"url"`
}

var channelAmqp *amqp.Channel

func init() {
	amqpConnection, err := amqp.Dial(os.Getenv("RABBITMQ_URI"))
	if err != nil {
		log.Fatal(err)
	}

	channelAmqp, _ = amqpConnection.Channel()
}

func main() {
	router := gin.Default()
	router.POST("/parse", ParserHandler)
	log.Fatal(router.Run(":6000"))
}

func ParserHandler(c *gin.Context) {
	var request Request

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	data, _ := json.Marshal(request)
	err := channelAmqp.Publish(
		"",
		os.Getenv("RABBITMQ_QUEUE"),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
