package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/orders-service/app"
)

func getOrders(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Orders")
}

func getOrder(c *gin.Context) {
	orderId := c.Param("id")
	c.IndentedJSON(http.StatusOK, orderId)
}

func createOrder(c *gin.Context) {}

func updateOrder(c *gin.Context) {}

func deleteOrder(c *gin.Context) {}

/*
healthPing func used checkhealth of the service

HEAD method is typically used to retrieve only the headers of a resource without fetching the entire response body, which can be useful for checking the availability or status of an endpoint.
*/
func healthPing(c *gin.Context) {}

/*
setLogger() configures the gin.log file where all the HTTP requests and application logs will be written.
We are currently writing the logs in console as well as in .log file
*/
func setLogger() {

	f, err := os.Create("./log/gin.log")
	if err != nil {
		log.Printf("Unable to create a log  file. Error : %v", err)
	} else {
		gin.DefaultWriter = io.MultiWriter(f)
		log.SetOutput(gin.DefaultWriter)
		log.Println("Logger is setup for the microservice")
	}
}

/*
setRouters() creates a default gin router with appropriate handlers for multiple REST API endpoints for Order service.
We are currently using version v1 as its initial version of API and keeping all routers in block of code {} for readability and maintainibility
*/
// initRouters
func setRouters() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/v1/order")
	{

		v1.HEAD("/health", healthPing)
		v1.GET("/", getOrders)
		v1.POST("/create", createOrder)
		v1.GET("/:id", getOrder)
		v1.PUT("/:id", updateOrder)
		v1.DELETE("/:id", deleteOrder)

	}
	return router
}

func main() {
	setLogger()
	router := setRouters()
	fmt.Println("configgg ", app.GetConfiguration())

	// Start the server with port as 8080
	router.Run()

}
